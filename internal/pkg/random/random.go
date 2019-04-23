package random

import (
	"log"
	"math/rand"
	"time"
)

var (
	ascii = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

const (
	dateFormat = "2006-01-02 15:04:05"
)

// String returns a random string between two lengths.
func String(min, max int, prefix string) string {
	if len(prefix) >= max {
		return prefix
	}

	var length int
	if min >= max {
		length = min
	} else {
		length = between(min, max) - len(prefix)
	}

	output := []rune{}
	for i := 0; i < length-len(prefix); i++ {
		output = append(output, ascii[rand.Intn(len(ascii))])
	}

	return prefix + string(output)
}

// Int returns a random integer between a minimum and maximum.
func Int(min, max int) int {
	return between(min, max)
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
		log.Fatalf("invalid min date: %v", err)
	}
	max, err := time.Parse(format, maxStr)
	if err != nil {
		log.Fatalf("invalid max date: %v", err)
	}

	if min == max {
		return min.UTC().Format(format)
	}

	diff := between64(min.Unix(), max.Unix())

	// Ensure max is greater than min.
	if diff < 0 {
		diff = -diff
	}

	return time.Unix(diff, 0).UTC().Format(format)
}

// Float32 returns a random float between a minimum and maximum.
func Float32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

// Float64 returns a random float between a minimum and maximum.
func Float64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Set returns a random item from a set
func Set(set ...string) string {
	i := between(0, len(set))
	return set[i]
}

func between(min, max int) int {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min) + min
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
