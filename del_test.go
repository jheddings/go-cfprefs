package cfprefs

import (
	"testing"

	"github.com/jheddings/go-cfprefs/testutil"
)

// assertKeyExists verifies that a key exists
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

func TestDeleteBasic(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"
	key := "del-test"

	// Set a value
	err := Set(appID, key, "test")
	testutil.AssertNoError(t, err, "set key")

	// Verify it exists
	assertKeyExists(t, appID, key, true)

	// Delete the key
	err = Delete(appID, key)
	testutil.AssertNoError(t, err, "delete key")

	// Verify it no longer exists
	assertKeyExists(t, appID, key, false)
}

func TestDeleteKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Create a nested structure with multiple keys
	err := Set(appID, "delete-test/level1/value1", "first value")
	testutil.AssertNoError(t, err, "set first value")

	err = Set(appID, "delete-test/level1/value2", "second value")
	testutil.AssertNoError(t, err, "set second value")

	err = Set(appID, "delete-test/level2/nested", int64(42))
	testutil.AssertNoError(t, err, "set nested value")
	defer Delete(appID, "delete-test")

	// Verify all values exist
	assertKeyExists(t, appID, "delete-test/level1/value1", true)
	assertKeyExists(t, appID, "delete-test/level1/value2", true)

	// Delete one nested value
	err = Delete(appID, "delete-test/level1/value1")
	testutil.AssertNoError(t, err, "delete nested key")

	// Verify it was deleted
	assertKeyExists(t, appID, "delete-test/level1/value1", false)

	// Verify sibling value still exists
	assertKeyExists(t, appID, "delete-test/level1/value2", true)

	value, err := Get(appID, "delete-test/level1/value2")
	testutil.AssertNoError(t, err, "get sibling value")
	if value.(string) != "second value" {
		t.Fatalf("sibling value was modified: expected 'second value', got '%s'", value.(string))
	}

	// Verify parent dictionary still exists
	assertKeyExists(t, appID, "delete-test/level1", true)

	// Verify other branch still exists
	assertKeyExists(t, appID, "delete-test/level2/nested", true)
}

func TestExistsKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Set up a nested structure
	err := Set(appID, "exists-test/level1/level2/value", "nested value")
	testutil.AssertNoError(t, err, "set nested value")
	defer Delete(appID, "exists-test")

	// Test that full path exists
	assertKeyExists(t, appID, "exists-test/level1/level2/value", true)

	// Test that intermediate paths exist
	assertKeyExists(t, appID, "exists-test", true)
	assertKeyExists(t, appID, "exists-test/level1", true)

	// Test that non-existent paths return false
	assertKeyExists(t, appID, "exists-test/nonexistent", false)
	assertKeyExists(t, appID, "exists-test/level1/wrong/path", false)
}

func TestDeleteEmptyKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test deleting with empty keypath
	err := Delete(appID, "")
	testutil.AssertError(t, err, "empty keypath")
}

func TestDeleteKeypathOnlySlashes(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test keypath with only slashes
	err := Delete(appID, "///")
	testutil.AssertError(t, err, "keypath with only slashes")
}

func TestDeleteKeypathNonDictSegment(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Set a simple string value
	cleanup := setupTest(t, appID, "simple-string", "hello")
	defer cleanup()

	// Try to delete through a non-dict segment - should return error
	err := Delete(appID, "simple-string/nested")
	testutil.AssertError(t, err, "deleting through non-dict segment")
}

func TestDeleteMissingKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Deleting a non-existent key should not error (idempotent)
	err := Delete(appID, "this-key-does-not-exist")
	testutil.AssertNoError(t, err, "delete missing key should be idempotent")
}

func TestExistsEmptyKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test empty keypath
	_, err := Exists(appID, "")
	testutil.AssertError(t, err, "empty keypath")
}

func TestExistsKeypathOnlySlashes(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test keypath with only slashes
	_, err := Exists(appID, "///")
	testutil.AssertError(t, err, "keypath with only slashes")
}
