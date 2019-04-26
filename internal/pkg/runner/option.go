package runner

// Option allows the Runner to be configured by the user.
type Option func(*Runner)

// WithDateFormat sets the default date format for the Runner.
func WithDateFormat(f string) Option {
	return func(r *Runner) {
		r.dateFormat = f
	}
}
