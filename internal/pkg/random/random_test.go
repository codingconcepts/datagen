package random

import (
	"strings"
	"testing"
	"time"

	"github.com/codingconcepts/datagen/internal/pkg/test"
)

func TestString(t *testing.T) {
	cases := []struct {
		name   string
		min    int64
		max    int64
		prefix string
	}{
		{name: "length 1 without prefix", min: 1, max: 1, prefix: ""},
		{name: "length 1 with prefix", min: 1, max: 1, prefix: "a"},
		{name: "length 2 without prefix", min: 2, max: 2, prefix: ""},
		{name: "length 2 with prefix", min: 2, max: 2, prefix: "aa"},
		{name: "different lengths 2 without prefix", min: 1, max: 10, prefix: ""},
		{name: "different lengths 2 with prefix", min: 1, max: 10, prefix: "a"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := String(c.min, c.max, c.prefix)
			test.Assert(t, int64(len(s)) >= c.min)
			test.Assert(t, int64(len(s)) <= c.max)

			if c.prefix != "" {
				test.Assert(t, strings.HasPrefix(s, c.prefix))
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
	}{
		{name: "1 1 no prefix", min: 1, max: 1, prefix: ""},
		{name: "1 1 prefix", min: 1, max: 1, prefix: "a"},
		{name: "10 10 no prefix", min: 10, max: 10, prefix: ""},
		{name: "10 10 prefix", min: 10, max: 10, prefix: "a"},
		{name: "1 10 no prefix", min: 1, max: 10, prefix: ""},
		{name: "1 10 prefix", min: 1, max: 10, prefix: "a"},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				String(c.min, c.max, c.prefix)
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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := Int(c.min, c.max)
			test.Assert(t, i >= c.min)
			test.Assert(t, i <= c.max)
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
	cases := []struct {
		name   string
		min    string
		max    string
		format string
	}{
		{name: "min eq max", min: "2019-04-23", max: "2019-04-23", format: "2006-01-02"},
		{name: "min lt max", min: "2018-04-23", max: "2019-04-23", format: "2006-01-02"},
		{name: "min gt max", min: "2019-04-23", max: "2018-04-23", format: "2006-01-02"},
		{name: "min gt max without format", min: "2019-04-23 01:02:03", max: "2018-04-23 01:02:03", format: ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := Date(c.min, c.max, c.format)

			format := c.format
			if format == "" {
				format = dateFormat
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
		name   string
		min    string
		max    string
		format string
	}{
		{name: "min eq max date", min: "2019-01-02", max: "2019-01-02", format: "2006-01-02"},
		{name: "min eq max date and time", min: "2019-01-02 03:04:05", max: "2019-01-02 03:04:05", format: "2006-01-02 15:04:05"},
		{name: "min lt max date", min: "2018-01-02", max: "2019-01-02", format: "2006-01-02"},
		{name: "min lt max date and time", min: "2018-01-02 03:04:05", max: "2019-01-02 03:04:05", format: "2006-01-02 15:04:05"},
		{name: "min gt max date", min: "2019-01-02", max: "2018-01-02", format: "2006-01-02"},
		{name: "min gt max date and time", min: "2019-01-02 03:04:05", max: "2018-01-02 03:04:05", format: "2006-01-02 15:04:05"},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Date(c.min, c.max, c.format)
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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := Float(c.min, c.max)
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
