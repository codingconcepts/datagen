package builder

import "strings"

// ErrBuilder simplifies the building of strings in a similar way
// to Rob Pike's errWriter: https://blog.golang.org/errors-are-values.
type ErrBuilder struct {
	b   strings.Builder
	err error
}

func NewErrBuilder() *ErrBuilder {
	return &ErrBuilder{b: strings.Builder{}}
}

// WriteStrings writes a given set of strings to the writer unless
// an error has previously occurred, in which case it returns early.
func (b *ErrBuilder) WriteStrings(ss ...string) {
	if b.err != nil {
		return
	}

	for _, s := range ss {
		if _, b.err = b.b.WriteString(s); b.err != nil {
			return
		}
	}
}

// Error returns the collected error value.
func (b *ErrBuilder) Error() error {
	return b.err
}

// String returns the underlying builder's string value.
func (b *ErrBuilder) String() string {
	return b.b.String()
}

// Reset resets the underlying builder.
func (b *ErrBuilder) Reset() {
	b.b.Reset()
}
