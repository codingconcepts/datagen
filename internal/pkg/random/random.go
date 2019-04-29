package random

import (
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

var (
	ascii = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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

// Int returns a random 64 integer between a minimum and maximum.
func Int(min, max int64) int64 {
	return between64(min, max)
}

// Date returns a random date between two dates and formats it
// as a string provided by Runner.
func Date(dateFormat string) func(minStr, maxStr string) (string, error) {
	return func(minStr, maxStr string) (string, error) {
		min, err := time.Parse(dateFormat, minStr)
		if err != nil {
			return "", errors.Wrap(err, "parsing min date")
		}

		max, err := time.Parse(dateFormat, maxStr)
		if err != nil {
			return "", errors.Wrap(err, "parsing max date")
		}

		if min == max {
			return min.UTC().Format(dateFormat), nil
		}

		diff := between64(min.Unix(), max.Unix())

		return time.Unix(diff, 0).UTC().Format(dateFormat), nil
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
func Set(set ...string) string {
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
