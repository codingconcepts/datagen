package random

import (
	"log"
	"math/rand"
	"time"
)

var (
	logFatalf = log.Fatalf
	ascii     = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

const (
	dateFormat = "2006-01-02 15:04:05"
)

// String returns a random string between two lengths.
func String(min, max int64, prefix string) string {
	if int64(len(prefix)) >= max {
		return prefix
	}

	var length int64
	if min == max {
		length = min
	} else {
		length = between64(min, max) - int64(len(prefix))
	}

	output := []rune{}
	for i := 0; i < int(length-int64(len(prefix))); i++ {
		output = append(output, ascii[rand.Intn(len(ascii))])
	}

	return prefix + string(output)
}

// Int returns a random 64 integer between a minimum and maximum.
func Int(min, max int64) int64 {
	return between64(min, max)
}

// Date returns a random date between two dates in format provided
// by the caller, returning a formatted UTC date in the same format.
// Defaults to 2006-01-02 15:04:05, if one isn't provided.
func Date(minStr, maxStr, format string) string {
	if format == "" {
		format = dateFormat
	}

	min, err := time.Parse(format, minStr)
	if err != nil {
		logFatalf("invalid min date: %v", err)
		return "" // Break out early for tests.
	}

	max, err := time.Parse(format, maxStr)
	if err != nil {
		logFatalf("invalid max date: %v", err)
		return "" // Break out early for tests.
	}

	if min == max {
		return min.UTC().Format(format)
	}

	diff := between64(min.Unix(), max.Unix())

	return time.Unix(diff, 0).UTC().Format(format)
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
