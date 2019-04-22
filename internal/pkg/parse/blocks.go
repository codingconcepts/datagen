package parse

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/codingconcepts/datagen/internal/pkg/builder"
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
	scanner := bufio.NewScanner(r)

	b := builder.NewErrBuilder()
	output := []Block{}
	current := Block{Repeat: 1}

	// Function to call whenever we've hit the gap between statements
	// or have reach the end of the file (either through manual EOF,
	// or the actual EOF).
	addAndReset := func(body string) {
		if body != "" {
			current.Body = body
			output = append(output, current)
		}
		b.Reset()
		current = Block{Repeat: 1}
	}

	for scanner.Scan() {
		t := strings.Trim(scanner.Text(), " \t")

		if strings.HasPrefix(t, commentRepeat) {
			var err error
			if current.Repeat, err = parseRepeat(t); err != nil {
				return nil, errors.Wrap(err, "parsing repeat")
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
			addAndReset(b.String())
			continue
		}

		// We've hit the end of the file.
		if strings.HasPrefix(t, commentEOF) {
			addAndReset(b.String())
			break
		}
		b.WriteStrings(t)
	}
	addAndReset(b.String())

	if err := b.Error(); err != nil {
		return nil, err
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func parseRepeat(input string) (int, error) {
	clean := strings.Trim(strings.TrimPrefix(input, commentRepeat), " \t")
	return strconv.Atoi(clean)
}

func parseName(input string) string {
	return strings.Trim(strings.TrimPrefix(input, commentName), " \t")
}
