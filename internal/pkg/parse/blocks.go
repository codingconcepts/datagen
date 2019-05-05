package parse

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	commentEOF    = "-- EOF"
	commentRepeat = "-- REPEAT"
	commentName   = "-- NAME"
)

// Block represents an instruction block in a script file.
type Block struct {
	// Repeat tells the application how many times to run the body.
	Repeat int

	// The name of the block can be used to identify the return values
	// from one block execution from another.
	Name string

	// The body of the template.
	Body string
}

// Blocks reads an input reader line by line, parsing blocks than
// can be executed by the Runner.  If a block does not have an
// explicit REPEAT value, a default of 1 will be used.
func Blocks(r io.Reader) ([]Block, error) {
	scanner := bufio.NewScanner(r)
	output := []Block{}

	for {
		ok, block, err := parseBlock(scanner)
		if err != nil {
			return nil, err
		}
		if block.Body != "" {
			output = append(output, block)
		}
		if !ok {
			return output, nil
		}
	}
}

func parseBlock(scanner *bufio.Scanner) (ok bool, block Block, err error) {
	b := strings.Builder{}
	block.Repeat = 1
	for scanner.Scan() {
		t := strings.Trim(scanner.Text(), " \t")

		if strings.HasPrefix(t, commentName) {
			block.Name = parseName(t)
			continue
		}

		if strings.HasPrefix(t, commentRepeat) {
			var err error
			if block.Repeat, err = parseRepeat(t); err != nil {
				return false, Block{}, errors.Wrap(err, "parsing repeat")
			}
			continue
		}

		// We've hit the gap between statements,break out and
		// signal that there could be more blocks to come.
		if t == "" {
			block.Body = b.String()
			return true, block, nil
		}

		// We've git the user-defined EOF, break out and signal
		// that there are no more blocks to come.
		if strings.HasPrefix(t, commentEOF) {
			block.Body = b.String()
			return false, block, nil
		}

		b.WriteString(t)
	}

	block.Body = b.String()
	return false, block, scanner.Err()
}

func parseRepeat(input string) (int, error) {
	clean := strings.Trim(strings.TrimPrefix(input, commentRepeat), " \t")
	return strconv.Atoi(clean)
}

func parseName(input string) string {
	return strings.Trim(strings.TrimPrefix(input, commentName), " \t")
}
