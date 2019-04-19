package parse

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/codingconcepts/dbgen/internal/pkg/builder"
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

func Blocks(r io.Reader) ([]Block, error) {
	var err error
	scanner := bufio.NewScanner(r)

	b := builder.NewErrBuilder()
	output := []Block{}
	current := Block{}
	for scanner.Scan() {
		t := scanner.Text()

		if strings.HasPrefix(t, commentRepeat) {
			if current.Repeat, err = parseRepeat(t); err != nil {
				return nil, errors.Wrap(err, "parsing repeat comment")
			}
			continue
		}

		if strings.HasPrefix(t, commentName) {
			current.Name = parseName(t)
			continue
		}

		// We've hit the gap between statements, add this block to
		// the output slice and reset.
		if t == "" {
			current.Body = b.String()
			output = append(output, current)
			b.Reset()
		}
		b.WriteStrings(t)

		// If the user has specified an end-of-file, break out.
		if t == commentEOF {
			break
		}
	}

	if err := b.Error(); err != nil {
		return nil, err
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func parseRepeat(input string) (int, error) {
	clean := strings.Trim(strings.TrimPrefix(input, commentRepeat), " ")
	return strconv.Atoi(clean)
}

func parseName(input string) string {
	return strings.Trim(strings.TrimPrefix(input, commentRepeat), " ")
}
