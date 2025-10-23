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

func TestDeleteQ(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Create a nested structure
	nestedData := map[string]any{
		"level1": map[string]any{
			"value1": "first value",
			"value2": "second value",
		},
		"level2": map[string]any{
			"nested": int64(42),
		},
	}

	testKey := "deleteq-test"
	cleanup := setupTest(t, appID, testKey, nestedData)
	defer cleanup()

	// Verify all values exist
	exists, _ := ExistsQ(appID, testKey, "$.level1.value1")
	if !exists {
		t.Fatal("value1 should exist before deletion")
	}

	// Delete one nested value
	err := DeleteQ(appID, testKey, "$.level1.value1")
	testutil.AssertNoError(t, err, "delete nested field")

	// Verify it was deleted
	exists, _ = ExistsQ(appID, testKey, "$.level1.value1")
	if exists {
		t.Fatal("value1 should not exist after deletion")
	}

	// Verify sibling value still exists
	value, err := GetQ(appID, testKey, "$.level1.value2")
	testutil.AssertNoError(t, err, "get sibling value")
	if value.(string) != "second value" {
		t.Fatalf("sibling value was modified: expected 'second value', got '%s'", value.(string))
	}

	// Verify parent dictionary still exists
	exists, _ = ExistsQ(appID, testKey, "$.level1")
	if !exists {
		t.Fatal("parent level1 should still exist")
	}

	// Verify other branch still exists
	value, err = GetQ(appID, testKey, "$.level2.nested")
	testutil.AssertNoError(t, err, "get other branch value")
	if value.(int64) != int64(42) {
		t.Fatalf("other branch was modified: expected 42, got %v", value)
	}
}
func TestDeleteQArrayElement(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Create a structure with an array
	arrayData := map[string]any{
		"items": []any{
			map[string]any{"id": int64(1), "name": "first"},
			map[string]any{"id": int64(2), "name": "second"},
			map[string]any{"id": int64(3), "name": "third"},
		},
	}

	testKey := "deleteq-array-test"
	cleanup := setupTest(t, appID, testKey, arrayData)
	defer cleanup()

	// Verify array has 3 items
	items, err := GetSliceQ(appID, testKey, "$.items")
	testutil.AssertNoError(t, err, "get items array")
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	// Delete the second item (index 1)
	err = DeleteQ(appID, testKey, "$.items[1]")
	testutil.AssertNoError(t, err, "delete array element")

	// Verify array now has 2 items
	items, err = GetSliceQ(appID, testKey, "$.items")
	testutil.AssertNoError(t, err, "get items array after deletion")
	if len(items) != 2 {
		t.Fatalf("expected 2 items after deletion, got %d", len(items))
	}

	// Verify the remaining items are correct
	firstItem := items[0].(map[string]any)
	if firstItem["id"].(int64) != int64(1) {
		t.Fatalf("first item should have id 1, got %v", firstItem["id"])
	}

	secondItem := items[1].(map[string]any)
	if secondItem["id"].(int64) != int64(3) {
		t.Fatalf("second item should have id 3 (was third), got %v", secondItem["id"])
	}
}

func TestDeleteQNestedInArray(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Create a structure with nested data in an array
	arrayData := map[string]any{
		"items": []any{
			map[string]any{"id": int64(1), "name": "first"},
			map[string]any{"id": int64(2), "name": "second"},
		},
	}

	testKey := "deleteq-nested-array-test"
	cleanup := setupTest(t, appID, testKey, arrayData)
	defer cleanup()

	// Delete a field from an array element
	err := DeleteQ(appID, testKey, "$.items[0].name")
	testutil.AssertNoError(t, err, "delete field from array element")

	// Verify the name field was deleted
	exists, _ := ExistsQ(appID, testKey, "$.items[0].name")
	if exists {
		t.Fatal("items[0].name should not exist after deletion")
	}

	// Verify the id field still exists
	value, err := GetIntQ(appID, testKey, "$.items[0].id")
	testutil.AssertNoError(t, err, "get items[0].id")
	if value != int64(1) {
		t.Fatalf("items[0].id should be 1, got %v", value)
	}

	// Verify the second item is unchanged
	value2, err := GetQ(appID, testKey, "$.items[1].name")
	testutil.AssertNoError(t, err, "get items[1].name")
	if value2.(string) != "second" {
		t.Fatalf("items[1].name should be 'second', got %v", value2)
	}
}

func TestDeleteQRootKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testKey := "deleteq-root-test"
	cleanup := setupTest(t, appID, testKey, "test value")
	defer cleanup()

	// Verify key exists
	assertKeyExists(t, appID, testKey, true)

	// Delete using empty query (deletes entire root key)
	err := DeleteQ(appID, testKey, "")
	testutil.AssertNoError(t, err, "delete with empty query")

	// Verify key was deleted
	assertKeyExists(t, appID, testKey, false)
}

func TestDeleteQWithDollarRoot(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testKey := "deleteq-dollar-test"
	cleanup := setupTest(t, appID, testKey, map[string]any{"field": "value"})
	defer cleanup()

	// Verify key exists
	assertKeyExists(t, appID, testKey, true)

	// Delete using "$" query (deletes entire root key)
	err := DeleteQ(appID, testKey, "$")
	testutil.AssertNoError(t, err, "delete with $ query")

	// Verify key was deleted
	assertKeyExists(t, appID, testKey, false)
}

func TestDeleteQNonExistentPath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	nestedData := map[string]any{
		"user": map[string]any{
			"name": "John",
		},
	}

	testKey := "deleteq-nonexistent-test"
	cleanup := setupTest(t, appID, testKey, nestedData)
	defer cleanup()

	// Delete a non-existent path should be idempotent (no error)
	err := DeleteQ(appID, testKey, "$.user.age")
	testutil.AssertNoError(t, err, "delete non-existent path should be idempotent")

	// Verify the existing data is unchanged
	value, err := GetStrQ(appID, testKey, "$.user.name")
	testutil.AssertNoError(t, err, "get existing field")
	if value != "John" {
		t.Fatalf("existing field should be unchanged, got %v", value)
	}
}

func TestDeleteQNonExistentRootKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Delete from a non-existent root key should be idempotent (no error)
	err := DeleteQ(appID, "nonexistent-root-key", "$.field")
	testutil.AssertNoError(t, err, "delete from non-existent root should be idempotent")
}

func TestDeleteMissingKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Deleting a non-existent key should not error (idempotent)
	err := Delete(appID, "this-key-does-not-exist")
	testutil.AssertNoError(t, err, "delete missing key should be idempotent")
}
