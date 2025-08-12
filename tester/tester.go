package tester

import (
	"strings"
	"testing"
)

func Assert[T comparable](t testing.TB, want T, got T) {
	if want != got {
		t.Errorf("want: %v | got: %v | want != got", want, got)
	}
}

func AssertContains(t testing.TB, got string, shouldContain string) {
	if !strings.Contains(got, shouldContain) {
		t.Errorf("want: %v | got: %v | got does not contain want", shouldContain, got)
	}
}
