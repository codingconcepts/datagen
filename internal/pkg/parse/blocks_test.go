package parse

import (
	"strings"
	"testing"

	"github.com/codingconcepts/datagen/internal/pkg/test"
)

func TestBlocksRepeat(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expCount int
		exp      int
		expError bool
	}{
		{
			name:     "defaults to 1",
			input:    `insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 1,
			exp:      1,
			expError: false,
		},
		{
			name: "defaults to 1 with blank line",
			input: `

			insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 1,
			exp:      1,
			expError: false,
		},
		{
			name: "sets to 2",
			input: `-- REPEAT 2
			insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 1,
			exp:      2,
			expError: false,
		},
		{
			name: "sets to 2 with blank line",
			input: `
			
			-- REPEAT 2
			insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 1,
			exp:      2,
			expError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			blocks, err := Blocks(strings.NewReader(c.input))
			test.ErrorExists(t, c.expError, err)
			test.Equals(t, c.expCount, len(blocks))

			for _, block := range blocks {
				test.Equals(t, c.exp, block.Repeat)
			}
		})
	}
}

func TestBlocksName(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expCount int
		exp      string
		expError bool
	}{
		{
			name:     "defaults to empty",
			input:    `insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 1,
			exp:      "",
			expError: false,
		},
		{
			name: "defaults to empty with blank line",
			input: `
			
			insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 1,
			exp:      "",
			expError: false,
		},
		{
			name: "sets to hello",
			input: `-- NAME hello
			insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 1,
			exp:      "hello",
			expError: false,
		},
		{
			name: "sets to hello with blank line",
			input: `
			
			-- NAME hello
			insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 1,
			exp:      "hello",
			expError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			blocks, err := Blocks(strings.NewReader(c.input))
			test.ErrorExists(t, c.expError, err)
			test.Equals(t, c.expCount, len(blocks))

			for _, block := range blocks {
				test.Equals(t, c.exp, block.Name)
			}
		})
	}
}
