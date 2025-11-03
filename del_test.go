package cfprefs

import (
	"reflect"
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

func TestDeleteNested(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	nestedData := map[string]any{
		"level1": map[string]any{
			"value1": "first value",
			"value2": "second value",
		},
		"level2": map[string]any{
			"nested": int64(42),
		},
	}

	cleanup := setupTest(t, appID, "delete-nested-test", nestedData)
	defer cleanup()

	// basic check that the value exists (using bracket notation)
	exists, _ := Exists(appID, "delete-nested-test/level1/value1")
	if !exists {
		t.Fatal("value1 should exist before deletion")
	}

	// delete one nested value
	err := Delete(appID, "delete-nested-test/level1/value1")
	testutil.AssertNoError(t, err, "delete nested field")

	// verify it was deleted (using dot notation)
	exists, _ = Exists(appID, "delete-nested-test/level1/value1")
	if exists {
		t.Fatal("value1 should not exist after deletion")
	}

	// verify sibling value still exists (using bracket notation)
	value, err := Get(appID, "delete-nested-test/level1/value2")
	testutil.AssertNoError(t, err, "get sibling value")
	if value.(string) != "second value" {
		t.Fatalf("sibling value was modified: expected 'second value', got '%s'", value.(string))
	}

	// verify parent dictionary still exists
	exists, _ = Exists(appID, "delete-nested-test/level1")
	if !exists {
		t.Fatal("parent level1 should still exist")
	}

	// verify other branch still exists
	value, err = Get(appID, "delete-nested-test/level2/nested")
	testutil.AssertNoError(t, err, "get other branch value")
	if value.(int64) != int64(42) {
		t.Fatalf("other branch was modified: expected 42, got %v", value)
	}
}

func TestDeleteQArrayElement(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	arrayData := map[string]any{
		"items": []any{
			map[string]any{"id": int64(1), "name": "first"},
			map[string]any{"id": int64(2), "name": "second"},
			map[string]any{"id": int64(3), "name": "third"},
		},
	}

	cleanup := setupTest(t, appID, "delete-array-test", arrayData)
	defer cleanup()

	// delete the second item (index 1)
	err := Delete(appID, "delete-array-test/items/1")
	testutil.AssertNoError(t, err, "delete array element")

	// verify the modified array
	items, err := Get(appID, "delete-array-test/items")
	testutil.AssertNoError(t, err, "get items array after deletion")

	expected := []any{
		map[string]any{"id": int64(1), "name": "first"},
		map[string]any{"id": int64(3), "name": "third"},
	}

	if !reflect.DeepEqual(items, expected) {
		t.Fatalf("items should be %v, got %v", expected, items)
	}
}

func TestDeleteNestedInArray(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	arrayData := map[string]any{
		"items": []any{
			map[string]any{"id": int64(1), "name": "first"},
			map[string]any{"id": int64(2), "name": "second"},
		},
	}

	cleanup := setupTest(t, appID, "delete-nested-array-test", arrayData)
	defer cleanup()

	// delete a field from an array element (using bracket notation)
	err := Delete(appID, "delete-nested-array-test/items/0/name")
	testutil.AssertNoError(t, err, "delete field from array element")

	// verify the name field was deleted
	exists, _ := Exists(appID, "delete-nested-array-test/items/0/name")
	if exists {
		t.Fatal("items[0].name should not exist after deletion")
	}

	// verify the id field still exists
	value, err := GetInt(appID, "delete-nested-array-test/items/0/id")
	testutil.AssertNoError(t, err, "get items[0].id")
	if value != int64(1) {
		t.Fatalf("items[0].id should be 1, got %v", value)
	}

	// verify the second item is unchanged
	value2, err := GetStr(appID, "delete-nested-array-test/items/1/name")
	testutil.AssertNoError(t, err, "get items[1].name")
	if value2 != "second" {
		t.Fatalf("items[1].name should be 'second', got '%s'", value2)
	}
}

func TestDeleteRootKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	cleanup := setupTest(t, appID, "delete-root-test", "test value")
	defer cleanup()

	assertKeyExists(t, appID, "delete-root-test", true)

	// delete using empty path (deletes entire root key)
	err := Delete(appID, "delete-root-test/")
	testutil.AssertNoError(t, err, "delete with empty path")

	assertKeyExists(t, appID, "delete-root-test", false)
}
func TestDeleteNonExistentPath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	nestedData := map[string]any{
		"user": map[string]any{
			"name": "John",
		},
	}

	cleanup := setupTest(t, appID, "delete-nonexistent-test", nestedData)
	defer cleanup()

	// delete a non-existent path should be idempotent (no error)
	err := Delete(appID, "delete-nonexistent-test/user/age")
	testutil.AssertNoError(t, err, "delete non-existent path should be idempotent")

	value, err := GetStr(appID, "delete-nonexistent-test/user/name")
	testutil.AssertNoError(t, err, "get existing field")
	if value != "John" {
		t.Fatalf("existing field should be unchanged, got %v", value)
	}
}

func TestDeleteNonExistentRootKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// delete from a non-existent root key should be idempotent (no error)
	err := Delete(appID, "nonexistent-root-key/field")
	testutil.AssertNoError(t, err, "delete from non-existent root should be idempotent")
}

func TestDeleteMissingKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	assertKeyExists(t, appID, "nonexistent-key", false)

	// deleting a non-existent key should not error (idempotent)
	err := Delete(appID, "nonexistent-key")
	testutil.AssertNoError(t, err, "delete missing key should be idempotent")
}
