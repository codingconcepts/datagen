package runner

import (
	"testing"
	"time"

	"github.com/codingconcepts/datagen/internal/pkg/test"
)

func TestWithDateFormat(t *testing.T) {
	r := New(db, WithDateFormat(time.RFC3339))
	test.Equals(t, time.RFC3339, r.dateFormat)
}
