package cfprefs

import (
	"testing"
)

// Common test constants
const testAppID = "com.jheddings.cfprefs.testing"

// setupTest sets a value in preferences and returns a cleanup function.
// The cleanup function will delete the key when called, typically via defer.
// Uses testing.TB interface to work with both *testing.T and *testing.B.
func setupTest(tb testing.TB, appID, key string, value any) func() {
	tb.Helper()
	
	err := Set(appID, key, value)
	if err != nil {
		tb.Fatalf("failed to set test value: %v", err)
	}

	return func() {
		if err := Delete(appID, key); err != nil {
			tb.Errorf("failed to cleanup test key: %v", err)
		}
	}
}

// assertKeyExists verifies that a key exists or doesn't exist as expected.
// It fails the test if the existence check returns an error or if the result
// doesn't match the expected value.
func assertKeyExists(t *testing.T, appID, key string, expected bool) {
	t.Helper()

	exists, err := Exists(appID, key)
	if err != nil {
		t.Fatalf("failed to check if key exists: %v", err)
	}

	if exists != expected {
		t.Fatalf("expected key existence to be %t, got %t", expected, exists)
	}
}
