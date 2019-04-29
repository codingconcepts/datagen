package runner

import (
	"database/sql/driver"
	"reflect"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/codingconcepts/datagen/internal/pkg/parse"
	"github.com/codingconcepts/datagen/internal/pkg/test"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name     string
		b        parse.Block
		expError bool
	}{
		{
			name: "empty template to simulate db error",
			b: parse.Block{
				Repeat: 1,
				Name:   "owner",
				Body:   ``,
			},
			expError: true,
		},
		{
			name: "invalid template",
			b: parse.Block{
				Repeat: 1,
				Name:   "owner",
				Body:   `{{range $i, $e := $.times_1000 }}`,
			},
			expError: true,
		},
		{
			name: "valid block",
			b: parse.Block{
				Repeat: 1,
				Name:   "owner",
				Body:   `insert into "owner" ("name") values ("Alice") returning "id", "name", "date_of_birth"`,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resetMock()
			r := New(db)

			id, name, dob := 123, "Alice", time.Date(2019, time.January, 2, 3, 4, 5, 0, time.UTC)

			if !c.expError {
				rows := []driver.Value{id, name, dob}
				mock.ExpectQuery(`insert into "owner" (.*) values (.*) returning "id", "name", "date_of_birth"`).WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "date_of_birth"}).AddRow(rows...))
			}

			err := r.Run(c.b)
			test.ErrorExists(t, c.expError, err)
			if err != nil {
				return
			}

			// Check the values committed to context, doing a string
			// comparison, as we're operating against reflect.Values.
			//
			// Note that no error expectation cases are being set up,
			// as we expect there to be values in these cases.
			actID, err := r.store.reference(c.b.Name, "id")
			test.ErrorExists(t, false, err)
			test.StringEquals(t, id, actID)

			actName, err := r.store.reference(c.b.Name, "name")
			test.ErrorExists(t, false, err)
			test.StringEquals(t, name, actName)

			actDob, err := r.store.reference(c.b.Name, "date_of_birth")
			test.ErrorExists(t, false, err)
			test.StringEquals(t, dob, actDob)
		})
	}
}

func TestReference(t *testing.T) {
	s := newStore()
	s.set("owner", map[string]interface{}{
		"id":   123,
		"name": "Alice",
	})

	cases := []struct {
		name     string
		key      string
		column   string
		expValue interface{}
		expError bool
	}{
		{name: "id found", key: "owner", column: "id", expValue: 123},
		{name: "name found", key: "owner", column: "name", expValue: "Alice"},
		{name: "key not found", key: "invalid", column: "name", expError: true},
		{name: "column not found", key: "owner", column: "invalid", expError: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			act, err := s.reference(c.key, c.column)
			test.ErrorExists(t, c.expError, err)
			test.Equals(t, c.expValue, act)
		})
	}
}

func TestRow(t *testing.T) {
	s := newStore()
	s.set("owner", map[string]interface{}{
		"id":   123,
		"name": "Alice",
	})

	cases := []struct {
		name     string
		key      string
		group    int
		lookups  map[string]interface{}
		expError bool
	}{
		{name: "id found", key: "owner", group: 1, lookups: map[string]interface{}{"id": 123}},
		{name: "name found", key: "owner", group: 2, lookups: map[string]interface{}{"name": "Alice"}},
		{name: "columns found", key: "owner", group: 3, lookups: map[string]interface{}{"id": 123, "name": "Alice"}},
		{name: "column not found for new group", group: 4, key: "owner", lookups: map[string]interface{}{"invalid": nil}, expError: true},
		{name: "column not found for existing group", group: 3, key: "owner", lookups: map[string]interface{}{"invalid": nil}, expError: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for lk, lv := range c.lookups {
				act, err := s.row(c.key, lk, c.group)
				test.ErrorExists(t, c.expError, err)
				test.Equals(t, lv, act)
			}
		})
	}
}

func TestPrepareValue(t *testing.T) {
	r := New(db, WithDateFormat("20060102"))

	cases := []struct {
		name  string
		value interface{}
		exp   interface{}
	}{
		{
			name:  "string",
			value: "Alice",
			exp:   "Alice",
		},
		{
			name:  "time.Time",
			value: time.Date(2019, time.July, 8, 9, 0, 1, 0, time.UTC),
			exp:   "20190708",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			act := r.prepareValue(reflect.ValueOf(c.value))
			test.StringEquals(t, c.exp, act)
		})
	}
}
