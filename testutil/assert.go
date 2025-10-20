package testutil

import (
	"testing"
)

// AssertError verifies that an error occurred
func AssertError(t *testing.T, err error, context string) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error for %s, got nil", context)
	}
}

// AssertNoError verifies that no error occurred
func AssertNoError(t *testing.T, err error, context string) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error for %s: %v", context, err)
	}
}
