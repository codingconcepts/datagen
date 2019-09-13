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
				Body:   `{{range $i, $e := ntimes 10 }}`,
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
