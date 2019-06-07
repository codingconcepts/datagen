package random

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ejdem86/datagen/internal/pkg/test"
)

func TestString(t *testing.T) {
	cases := []struct {
		name   string
		min    int64
		max    int64
		prefix string
		set    string
	}{
		{name: "length 1 without prefix", min: 1, max: 1, prefix: ""},
		{name: "length 1 with prefix", min: 1, max: 1, prefix: "a"},
		{name: "length 2 without prefix", min: 2, max: 2, prefix: ""},
		{name: "length 2 with prefix", min: 2, max: 2, prefix: "aa"},
		{name: "different lengths 2 without prefix", min: 1, max: 10, prefix: ""},
		{name: "different lengths 2 with prefix", min: 1, max: 10, prefix: "a"},
		{name: "min > max without prefix", min: 10, max: 1, prefix: ""},
		{name: "min > max with prefix", min: 10, max: 1, prefix: "a"},
		{name: "custom set without prefix", min: 10, max: 10, prefix: "", set: "ab"},
		{name: "custom set with prefix", min: 10, max: 10, prefix: "c", set: "ab"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := String(c.min, c.max, c.prefix, c.set)

			if c.min > c.max {
				c.min, c.max = c.max, c.min
			}

			test.Assert(t, int64(len(s)) >= c.min)
			test.Assert(t, int64(len(s)) <= c.max)

			if c.prefix != "" {
				test.Assert(t, strings.HasPrefix(s, c.prefix))
			}

			if c.set != "" {
				runesInSet(t, []rune(c.set), []rune(strings.TrimPrefix(s, c.prefix)))
			}
		})
	}
}

func BenchmarkString(b *testing.B) {
	cases := []struct {
		name   string
		min    int64
		max    int64
		prefix string
		set    string
	}{
		{name: "1 1 no prefix", min: 1, max: 1, prefix: ""},
		{name: "1 1 prefix", min: 1, max: 1, prefix: "a"},
		{name: "10 10 no prefix", min: 10, max: 10, prefix: ""},
		{name: "10 10 prefix", min: 10, max: 10, prefix: "a"},
		{name: "1 10 no prefix", min: 1, max: 10, prefix: ""},
		{name: "1 10 prefix", min: 1, max: 10, prefix: "a"},
		{name: "1 10 not prefix set", min: 1, max: 10, prefix: "", set: "abcABC"},
		{name: "1 10 prefix set", min: 1, max: 10, prefix: "a", set: "abcABC"},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				String(c.min, c.max, c.prefix, c.set)
			}
		})
	}
}

func TestInt(t *testing.T) {
	cases := []struct {
		name string
		min  int64
		max  int64
	}{
		{name: "min eq max", min: 1, max: 1},
		{name: "min lt max", min: 1, max: 10},
		{name: "min gt max", min: 10, max: 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := Int(c.min, c.max)

			if c.min > c.max {
				c.min, c.max = c.max, c.min
			}

			test.Assert(t, i >= c.min)
			test.Assert(t, i <= c.max)
		})
	}
}

func TestStringF(t *testing.T) {
	defaults := StringFDefaults{
		StringMinDefault: 10,
		StringMaxDefault: 10,
		IntMinDefault:    1000,
		IntMaxDefault:    1000,
	}

	cases := []struct {
		name     string
		format   string
		args     []interface{}
		assert   func(string) error
		expError bool
	}{
		{
			name:   "string without arguments",
			format: "%s",
			assert: func(s string) error {
				if int64(len(s)) != defaults.StringMaxDefault {
					return fmt.Errorf("%q length is not equal to %d", s, defaults.StringMaxDefault)
				}
				return nil
			},
		},
		{
			name:   "string with length arguments",
			format: "%s",
			args:   []interface{}{20, 20},
			assert: func(s string) error {
				if int64(len(s)) != 20 {
					return fmt.Errorf("%q length is not equal to %d", s, 20)
				}
				return nil
			},
		},
		{
			name:   "string with all arguments",
			format: "%s",
			args:   []interface{}{20, 20, "abc"},
			assert: func(s string) error {
				if int64(len(s)) != 20 {
					return fmt.Errorf("%q length is not equal to %d", s, 20)
				}
				return nil
			},
		},
		{
			name:   "string with another placeholder's arguments",
			format: "%s",
			args:   []interface{}{20, 20, 30, 30},
			assert: func(s string) error {
				if int64(len(s)) != 20 {
					return fmt.Errorf("%q length is not equal to %d", s, 20)
				}
				return nil
			},
		},
		{
			name:     "string with invalid min argument",
			format:   "%s",
			args:     []interface{}{"hello", 20},
			expError: true,
		},
		{
			name:     "string with invalid max argument",
			format:   "%s",
			args:     []interface{}{20, "hello"},
			expError: true,
		},
		{
			name:   "int without arguments",
			format: "%d",
			assert: func(s string) error {
				i, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return err
				}

				if i != defaults.IntMaxDefault {
					return fmt.Errorf("%q is not equal to %d", s, defaults.IntMaxDefault)
				}
				return nil
			},
		},
		{
			name:   "int with arguments",
			format: "%d",
			args:   []interface{}{2000, 2000},
			assert: func(s string) error {
				i, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return err
				}

				if i != 2000 {
					return fmt.Errorf("%q is not equal to %d", s, 2000)
				}
				return nil
			},
		},
		{
			name:     "int with invalid min argument",
			format:   "%d",
			args:     []interface{}{"hello", 2000},
			expError: true,
		},
		{
			name:     "int with invalid max argument",
			format:   "%d",
			args:     []interface{}{2000, "hello"},
			expError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := StringF(defaults)
			act, err := s(c.format, c.args...)

			test.ErrorExists(t, c.expError, err)
			if c.expError {
				return
			}

			test.ErrorExists(t, false, c.assert(act))
		})
	}
}

func BenchmarkStringF(b *testing.B) {
	defaults := StringFDefaults{
		StringMinDefault: 10,
		StringMaxDefault: 10,
		IntMinDefault:    1000,
		IntMaxDefault:    1000,
	}

	cases := []struct {
		name   string
		format string
		args   []interface{}
	}{
		{
			name:   "string without arguments",
			format: "%s",
		},
		{
			name:   "string with length arguments",
			format: "%s",
			args:   []interface{}{20, 20},
		},
		{
			name:   "string with all arguments",
			format: "%s",
			args:   []interface{}{20, 20, "abc"},
		},
		{
			name:   "string with another placeholder's arguments",
			format: "%s",
			args:   []interface{}{20, 20, 30, 30},
		},
		{
			name:   "int without arguments",
			format: "%d",
		},
		{
			name:   "int with arguments",
			format: "%d",
			args:   []interface{}{2000, 2000},
		},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			s := StringF(defaults)

			for i := 0; i < b.N; i++ {
				s(c.format, c.args...)
			}
		})
	}
}

func BenchmarkInt(b *testing.B) {
	cases := []struct {
		name string
		min  int64
		max  int64
	}{
		{name: "1 1", min: 1, max: 1},
		{name: "10 10", min: 10, max: 10},
		{name: "1 10", min: 1, max: 10},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Int(c.min, c.max)
			}
		})
	}
}

func TestDate(t *testing.T) {
	date := time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC)

	origUTCNow := utcNow
	utcNow = func() time.Time { return date }
	defer func() { utcNow = origUTCNow }()

	cases := []struct {
		name           string
		min            string
		max            string
		format         string
		overrideFormat string
		expError       bool
	}{
		{name: "min eq max", min: "2019-04-23", max: "2019-04-23", format: "2006-01-02"},
		{name: "min lt max", min: "2018-04-23", max: "2019-04-23", format: "2006-01-02"},
		{name: "min gt max", min: "2019-04-23", max: "2018-04-23", format: "2006-01-02"},

		{name: "min parse failure", min: "2019-13-32", expError: true},
		{name: "max parse failure", min: "2019-04-23", max: "2019-13-32", format: "2006-01-02", expError: true},
		{name: "max parse failure", min: "2019-04-23", max: "2019-04-23", format: "1006-01-02", expError: true},

		{name: "override format", min: "20190423", max: "20180423", format: "2006-01-02", overrideFormat: "20060102"},
		{name: "override format failure", min: "20190423", max: "20180423", format: "2006-01-02", overrideFormat: "20170102", expError: true},

		{name: "override format failure", min: "20190423", max: "20180423", format: "2006-01-02", overrideFormat: "20170102", expError: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			df := Date(c.format)
			d, err := df(c.min, c.max, c.overrideFormat)
			test.ErrorExists(t, c.expError, err)

			// Don't continue if we expect an error.
			if c.expError {
				return
			}

			format := c.format
			if c.overrideFormat != "" {
				format = c.overrideFormat
			}

			minD, err := time.Parse(format, c.min)
			if err != nil {
				t.Fatalf("invalid min format: %v", err)
			}

			maxD, err := time.Parse(format, c.max)
			if err != nil {
				t.Fatalf("invalid max format: %v", err)
			}

			if minD.Unix() > maxD.Unix() {
				minD, maxD = maxD, minD
			}

			dD, err := time.Parse(format, d)
			if err != nil {
				t.Fatalf("invalid d format: %v", err)
			}

			test.Assert(t, dD.Unix() >= minD.Unix())
			test.Assert(t, dD.Unix() <= maxD.Unix())
		})
	}
}

func BenchmarkDate(b *testing.B) {
	cases := []struct {
		name           string
		min            string
		max            string
		format         string
		overrideFormat string
	}{
		{name: "min eq max date", min: "2019-01-02", max: "2019-01-02", format: "2006-01-02"},
		{name: "min eq max date and time", min: "2019-01-02 03:04:05", max: "2019-01-02 03:04:05", format: "2006-01-02 15:04:05"},
		{name: "min lt max date", min: "2018-01-02", max: "2019-01-02", format: "2006-01-02"},
		{name: "min lt max date and time", min: "2018-01-02 03:04:05", max: "2019-01-02 03:04:05", format: "2006-01-02 15:04:05"},
		{name: "min gt max date", min: "2019-01-02", max: "2018-01-02", format: "2006-01-02"},
		{name: "min gt max date and time", min: "2019-01-02 03:04:05", max: "2018-01-02 03:04:05", format: "2006-01-02 15:04:05"},

		{name: "override format", min: "2019-01-02 03:04:05", max: "2018-01-02 03:04:05", format: "2006-01-02 15:04:05", overrideFormat: "20060102"},
		{name: "override format failure", min: "2019-01-02 03:04:05", max: "2018-01-02 03:04:05", format: "2006-01-02 15:04:05", overrideFormat: "20170102"},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			d := Date(c.format)
			for i := 0; i < b.N; i++ {
				d(c.min, c.max, c.overrideFormat)
			}
		})
	}
}

func TestFloat(t *testing.T) {
	cases := []struct {
		name string
		min  float64
		max  float64
	}{
		{name: "min eq max", min: 1, max: 1},
		{name: "min lt max", min: 1, max: 10},
		{name: "min gt max", min: 10, max: 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := Float(c.min, c.max)

			if c.min > c.max {
				c.min, c.max = c.max, c.min
			}

			test.Assert(t, i >= c.min)
			test.Assert(t, i <= c.max)
		})
	}
}

func BenchmarkFloat(b *testing.B) {
	cases := []struct {
		name string
		min  float64
		max  float64
	}{
		{name: "1.23 1.23", min: 1.23, max: 1.23},
		{name: "10.23 10.23", min: 10.23, max: 10.23},
		{name: "1.23 10.23", min: 1.23, max: 10.23},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Float(c.min, c.max)
			}
		})
	}
}

func TestSet(t *testing.T) {
	cases := []struct {
		name string
		set  []string
	}{
		{name: "one item", set: []string{"a"}},
		{name: "multiple items", set: []string{"a", "b"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Set(c.set...)
		})
	}
}

func BenchmarkSet(b *testing.B) {
	cases := []struct {
		name  string
		items []string
	}{
		{name: "one item", items: []string{"a"}},
		{name: "multiple items", items: []string{"a", "b", "c"}},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Set(c.items...)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	date := time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC)

	origUTCNow := utcNow
	utcNow = func() time.Time { return date }
	defer func() { utcNow = origUTCNow }()

	cases := []struct {
		name   string
		format string
		exp    string
	}{
		{name: "default format", format: "", exp: "2000-01-02"},
		{name: "custom format", format: time.RFC3339, exp: "2000-01-02T03:04:05Z"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := Date("2006-01-02")
			act, err := d("now", "now", c.format)
			test.ErrorExists(t, false, err)
			test.Equals(t, c.exp, act)
		})
	}
}

func runesInSet(t *testing.T, exp, act []rune) {
	for _, a := range act {
		runeInSet(t, a, exp)
	}
}

func runeInSet(t *testing.T, r rune, set []rune) {
	for _, s := range set {
		if r == s {
			return
		}
	}
	t.Fatalf("rune %v not found in set", r)
}
