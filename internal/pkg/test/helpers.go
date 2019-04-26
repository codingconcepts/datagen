package test

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

// StringEquals performs a comparison against two values
// values by comparing their string values and fails if
// they are not the same.
func StringEquals(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(fmt.Sprintf("%v", expected), fmt.Sprintf("%v", actual)) {
		tb.Helper()
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}

// Equals performs a deep equal comparison against two
// values and fails if they are not the same.
func Equals(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		tb.Helper()
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}

// Assert checks the result of a predicate.
func Assert(tb testing.TB, result bool) {
	tb.Helper()
	if !result {
		tb.Fatal("\n\tassertion failed\n")
	}
}

// ErrorExists fails if an error is expected but doesn't
// exist or if an error is exists but is not expected.
// It does not check equality.
func ErrorExists(tb testing.TB, exp bool, err error) {
	tb.Helper()
	if !exp && err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	if exp && err == nil {
		log.Fatal("expect error but didn't get one")
	}
}
