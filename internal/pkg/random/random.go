package random

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	utcNow = func() time.Time { return time.Now().UTC() }

	ascii = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	verbPattern = regexp.MustCompile("%[sd]{1}")
)

// String returns a random string between two lengths.
func String(min, max int64, prefix string, set string) string {
	if int64(len(prefix)) >= max {
		return prefix
	}

	var length int64
	if min == max {
		length = min
	} else {
		length = between64(min, max) - int64(len(prefix))
	}

	runes := ascii
	if set != "" {
		runes = []rune(set)
	}

	output := []rune{}
	for i := 0; i < int(length-int64(len(prefix))); i++ {
		output = append(output, runes[rand.Intn(len(runes))])
	}

	return prefix + string(output)
}

// StringF returns a random string built around a format string.
func StringF(d StringFDefaults) func(format string, args ...interface{}) (string, error) {
	return func(format string, args ...interface{}) (string, error) {
		fargs := []interface{}{}

		verbs := verbPattern.FindAllString(format, -1)

		min, max, pattern, argIndex := int64(0), int64(0), "", 0
		var err error

		for _, v := range verbs {
			switch v[1:] {
			case "d":
				min, max, argIndex, err = intArgs(argIndex, d, args...)
				if err != nil {
					return "", errors.Wrap(err, "generating integer placeholder")
				}
				fargs = append(fargs, Int(min, max))
			case "s":
				min, max, pattern, argIndex, err = stringArgs(argIndex, d, args...)
				if err != nil {
					return "", errors.Wrap(err, "generating string placeholder")
				}
				fargs = append(fargs, String(min, max, "", pattern))
			}
		}

		return fmt.Sprintf(format, fargs...), nil
	}
}

// intArgs returns the minimum and maximum values to generate between,
// the next index to use from the arguments provided by the user, and
// any error that occurred parsing the parameters.
func intArgs(i int, d StringFDefaults, args ...interface{}) (int64, int64, int, error) {
	if len(args) <= i {
		return d.IntMinDefault, d.IntMaxDefault, i, nil
	}

	// The next 2 args should be integers.
	min, ok := args[i].(int)
	if !ok {
		return 0, 0, 0, fmt.Errorf("argument for min: %v is not an integer", args[i])
	}
	max, ok := args[i+1].(int)
	if !ok {
		return 0, 0, 0, fmt.Errorf("argument for max: %v is not an integer", args[i])
	}

	return int64(min), int64(max), i + 2, nil
}

// stringArgs returns the minimum and maximum length values to generate
// between, the character set to use when generating the random string,
// the next index to use from the arguments provided by the user, and
// any error that occurred parsing the parameters.
func stringArgs(i int, d StringFDefaults, args ...interface{}) (int64, int64, string, int, error) {
	if len(args) <= i {
		return d.StringMinDefault, d.StringMaxDefault, "", i, nil
	}

	// The next 2 args should be integers.
	min, ok := args[i].(int)
	if !ok {
		return 0, 0, "", 0, fmt.Errorf("argument for min: %v is not an integer", args[i])
	}

	max, ok := args[i+1].(int)
	if !ok {
		return 0, 0, "", 0, fmt.Errorf("argument for max: %v is not an integer", args[i])
	}

	// If there's a next argument, it might be a pattern, or for the next verb.
	if len(args) <= i+2 {
		return int64(min), int64(max), "", i + 2, nil
	}

	if s, ok := args[i+2].(string); ok {
		return int64(min), int64(max), s, i + 3, nil
	}

	return int64(min), int64(max), "", i + 2, nil
}

// Int returns a random 64 integer between a minimum and maximum.
func Int(min, max int64) int64 {
	return between64(min, max)
}

// Date returns a random date between two dates and formats it
// as a string provided by Runner.  It can optionally accept a
// format string to override the Runner's format. Leave empty
// to use the default.
func Date(dateFormat string) func(minStr, maxStr, format string) (string, error) {
	return func(minStr, maxStr, format string) (string, error) {
		if format == "" {
			format = dateFormat
		}

		min, err := parseDate(format, minStr)
		if err != nil {
			return "", errors.Wrap(err, "parsing min date")
		}

		max, err := parseDate(format, maxStr)
		if err != nil {
			return "", errors.Wrap(err, "parsing max date")
		}

		if min == max {
			return min.UTC().Format(format), nil
		}

		diff := between64(min.Unix(), max.Unix())

		return time.Unix(diff, 0).UTC().Format(format), nil
	}
}

// Float returns a random 64 bit float between a minimum and maximum.
func Float(min, max float64) float64 {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}

	return min + rand.Float64()*(max-min)
}

// Set returns a random item from a set
func Set(set ...interface{}) interface{} {
	i := between64(0, int64(len(set)))
	return set[i]
}

func between64(min, max int64) int64 {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return rand.Int63n(max-min) + min
}

func parseDate(format, input string) (time.Time, error) {
	if strings.EqualFold(input, "now") {
		return utcNow(), nil
	}

	return time.Parse(format, input)
}
