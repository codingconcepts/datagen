package runner

import "github.com/codingconcepts/datagen/internal/pkg/random"

// Option allows the Runner to be configured by the user.
type Option func(*Runner)

// WithDateFormat sets the default date format for the Runner.
func WithDateFormat(f string) Option {
	return func(r *Runner) {
		r.dateFormat = f
	}
}

// WithStringFDefaults sets the default format min and max values
// for the Runner.
func WithStringFDefaults(d random.StringFDefaults) Option {
	return func(r *Runner) {
		r.stringFdefaults = d
	}
}

// WithDebug puts the Runner in debug mode, meaning nothing will be
// written to a database.
func WithDebug(d bool) Option {
	return func(r *Runner) {
		r.debug = d
	}
}
