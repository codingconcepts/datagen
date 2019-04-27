package parse

import (
	"errors"
	"strings"
	"testing"

	"github.com/codingconcepts/datagen/internal/pkg/test"
)

func TestBlocks(t *testing.T) {
	input := `-- NAME a
	-- REPEAT 1
	A

	-- NAME b
	-- REPEAT 2
	B

	-- NAME e
	-- REPEAT 4

	-- NAME c
	-- REPEAT 3
	C
	
	-- NAME e
	-- REPEAT 5`

	blocks, err := Blocks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("error parsing blocks: %v", err)
	}

	test.Equals(t, "a", blocks[0].Name)
	test.Equals(t, 1, blocks[0].Repeat)
	test.Equals(t, "A", blocks[0].Body)

	test.Equals(t, "b", blocks[1].Name)
	test.Equals(t, 2, blocks[1].Repeat)
	test.Equals(t, "B", blocks[1].Body)

	test.Equals(t, "c", blocks[2].Name)
	test.Equals(t, 3, blocks[2].Repeat)
	test.Equals(t, "C", blocks[2].Body)

	// Note that nothing is parsed for block d.
	test.Equals(t, 3, len(blocks))
}

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
		{
			name: "returns error for invalid repeat",
			input: `-- REPEAT a
			insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 0,
			exp:      0,
			expError: true,
		},
		{
			name: "returns error for invalid repeat with blank line",
			input: `
			
			-- REPEAT a
			insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 0,
			exp:      0,
			expError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			blocks, err := Blocks(strings.NewReader(c.input))
			test.ErrorExists(t, c.expError, err)
			if err != nil {
				return
			}

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

func TestBlocksEOF(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expCount int
		expError bool
	}{
		{
			name: "one block",
			input: `insert into "t" ("a", "b") values ('a', 'b');

			-- EOF`,
			expCount: 1,
			expError: false,
		},
		{
			name: "two blocks",
			input: `insert into "t" ("a", "b") values ('a', 'b');

			insert into "t" ("a", "b") values ('a', 'b');

			-- EOF`,
			expCount: 2,
			expError: false,
		},
		{
			name: "ignore block",
			input: `insert into "t" ("a", "b") values ('a', 'b');
			
			-- EOF

			insert into "t" ("a", "b") values ('a', 'b');`,
			expCount: 1,
			expError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			blocks, err := Blocks(strings.NewReader(c.input))
			test.ErrorExists(t, c.expError, err)
			test.Equals(t, c.expCount, len(blocks))
		})
	}
}

func TestBlocksScanError(t *testing.T) {
	r := &errReader{err: errors.New("oh noes!")}
	_, err := Blocks(r)
	test.Equals(t, r.err, err)
}

type errReader struct {
	err error
}

func (r *errReader) Read(_ []byte) (int, error) {
	return 0, r.err
}
