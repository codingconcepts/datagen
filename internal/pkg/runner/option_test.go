package runner

import (
	"testing"
	"time"

	"github.com/codingconcepts/datagen/internal/pkg/random"

	"github.com/codingconcepts/datagen/internal/pkg/test"
)

func TestWithDateFormat(t *testing.T) {
	r := New(db, WithDateFormat(time.RFC3339))

	test.Equals(t, time.RFC3339, r.dateFormat)
}

func TestWithStringFDefaults(t *testing.T) {
	r := New(db, WithStringFDefaults(random.StringFDefaults{
		IntMinDefault:    1,
		IntMaxDefault:    2,
		StringMinDefault: 3,
		StringMaxDefault: 4,
	}))

	test.Equals(t, int64(1), r.stringFdefaults.IntMinDefault)
	test.Equals(t, int64(2), r.stringFdefaults.IntMaxDefault)
	test.Equals(t, int64(3), r.stringFdefaults.StringMinDefault)
	test.Equals(t, int64(4), r.stringFdefaults.StringMaxDefault)
}
