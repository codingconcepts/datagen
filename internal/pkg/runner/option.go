package runner

import "github.com/ejdem86/datagen/internal/pkg/random"

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
