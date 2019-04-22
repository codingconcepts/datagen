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
	dateFormat = "2006-01-02 15:04:05Z07:00"
)

func String(min, max int, prefix string) string {
	var length int
	if min >= max {
		length = min
	} else {
		length = between(min, max) - len(prefix)
	}

	output := make([]rune, length)
	for i := 0; i < length; i++ {
		output[i] = ascii[rand.Intn(len(ascii))]
	}

	return prefix + string(output)
}

func Int(min, max int) int {
	return between(min, max)
}

func Date(minStr, maxStr string) string {
	min, err := time.Parse(time.RFC3339, minStr)
	if err != nil {
		log.Fatalf("invalid min date: %v", err)
	}
	max, err := time.Parse(time.RFC3339, maxStr)
	if err != nil {
		log.Fatalf("invalid max date: %v", err)
	}

	diff := between64(min.Unix(), max.Unix())
	return time.Unix(diff, 0).Format(dateFormat)
}

func Float32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func Float64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
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
