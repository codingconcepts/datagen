package random

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/codingconcepts/datagen/internal/pkg/test"
)

func TestString(t *testing.T) {
	cases := []struct {
		name   string
		min    int
		max    int
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
			fmt.Printf("%q, %d, %d, %q\n", s, c.min, c.max, c.prefix)
			test.Assert(t, len(s) >= c.min)
			test.Assert(t, len(s) <= c.max)

			if c.prefix != "" {
				test.Assert(t, strings.HasPrefix(s, c.prefix))
			}
		})
	}
}

func TestInt(t *testing.T) {
	cases := []struct {
		name string
		min  int
		max  int
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

func TestFloat32(t *testing.T) {
	cases := []struct {
		name string
		min  float32
		max  float32
	}{
		{name: "min eq max", min: 1, max: 1},
		{name: "min lt max", min: 1, max: 10},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := Float32(c.min, c.max)
			test.Assert(t, i >= c.min)
			test.Assert(t, i <= c.max)
		})
	}
}

func TestFloat64(t *testing.T) {
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
			i := Float64(c.min, c.max)
			test.Assert(t, i >= c.min)
			test.Assert(t, i <= c.max)
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
