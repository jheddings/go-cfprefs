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

func TestExists(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// if this fails, we need to manually cleanup
	assertKeyExists(t, appID, "exists-test", false)

	cleanup := setupTest(t, appID, "exists-test", "test value")
	defer cleanup()

	assertKeyExists(t, appID, "exists-test", true)
	assertKeyExists(t, appID, "nonexistent-key", false)
}

func TestExistsQ(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	nestedData := map[string]any{
		"user": map[string]any{
			"name": "John Doe",
			"age":  int64(30),
			"address": map[string]any{
				"city":  "Anytown",
				"state": "CA",
			},
		},
		"items": []any{
			map[string]any{"id": int64(1), "active": true},
			map[string]any{"id": int64(2), "active": false},
			map[string]any{"id": int64(3), "active": true},
		},
	}

	testKey := "existsq-test"
	cleanup := setupTest(t, appID, testKey, nestedData)
	defer cleanup()

	t.Run("Root object exists", func(t *testing.T) {
		exists, err := ExistsQ(appID, testKey, "$")
		testutil.AssertNoError(t, err, "check root exists")
		if !exists {
			t.Fatal("expected root to exist")
		}
	})

	t.Run("Empty query checks root", func(t *testing.T) {
		exists, err := ExistsQ(appID, testKey, "")
		testutil.AssertNoError(t, err, "check empty query")
		if !exists {
			t.Fatal("expected root to exist with empty query")
		}
	})

	t.Run("Nested field exists", func(t *testing.T) {
		exists, err := ExistsQ(appID, testKey, "$.user.name")
		testutil.AssertNoError(t, err, "check user.name exists")
		if !exists {
			t.Fatal("expected user.name to exist")
		}
	})

	t.Run("Deeply nested field exists", func(t *testing.T) {
		exists, err := ExistsQ(appID, testKey, "$.user.address.city")
		testutil.AssertNoError(t, err, "check user.address.city exists")
		if !exists {
			t.Fatal("expected user.address.city to exist")
		}
	})

	t.Run("Array element exists", func(t *testing.T) {
		exists, err := ExistsQ(appID, testKey, "$.items[0]")
		testutil.AssertNoError(t, err, "check items[0] exists")
		if !exists {
			t.Fatal("expected items[0] to exist")
		}
	})

	t.Run("Array field exists", func(t *testing.T) {
		exists, err := ExistsQ(appID, testKey, "$.items[1].id")
		testutil.AssertNoError(t, err, "check items[1].id exists")
		if !exists {
			t.Fatal("expected items[1].id to exist")
		}
	})

	t.Run("Non-existent field returns false", func(t *testing.T) {
		exists, err := ExistsQ(appID, testKey, "$.user.nonexistent")
		testutil.AssertNoError(t, err, "check non-existent field")
		if exists {
			t.Fatal("expected user.nonexistent to not exist")
		}
	})

	t.Run("Non-existent nested path returns false", func(t *testing.T) {
		exists, err := ExistsQ(appID, testKey, "$.user.address.country")
		testutil.AssertNoError(t, err, "check non-existent nested path")
		if exists {
			t.Fatal("expected user.address.country to not exist")
		}
	})

	t.Run("Out of bounds array index returns false", func(t *testing.T) {
		exists, err := ExistsQ(appID, testKey, "$.items[999]")
		testutil.AssertNoError(t, err, "check out of bounds index")
		if exists {
			t.Fatal("expected items[999] to not exist")
		}
	})

	t.Run("Non-existent root key returns false", func(t *testing.T) {
		exists, err := ExistsQ(appID, "nonexistent-root", "$.user.name")
		testutil.AssertNoError(t, err, "check non-existent root key")
		if exists {
			t.Fatal("expected non-existent root key to return false")
		}
	})

	t.Run("Invalid JSONPath returns error", func(t *testing.T) {
		_, err := ExistsQ(appID, testKey, "$.user[invalid")
		testutil.AssertError(t, err, "invalid JSONPath should not error")
	})
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
