package runner

import (
	"testing"

	"github.com/codingconcepts/datagen/internal/pkg/test"
)

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

func TestEach(t *testing.T) {
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
			s := newStore()
			s.set("owner", map[string]interface{}{
				"id":   123,
				"name": "Alice",
			})

			for lk, lv := range c.lookups {
				act, err := s.each(c.key, lk, c.group)
				test.ErrorExists(t, c.expError, err)
				test.Equals(t, lv, act)
			}
		})
	}
}
