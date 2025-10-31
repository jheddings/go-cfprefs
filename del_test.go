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

func TestDeleteStructured(t *testing.T) {
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

	testKey := "deleteq-test"
	cleanup := setupTest(t, appID, testKey, nestedData)
	defer cleanup()

	// basic check that the value exists (using bracket notation)
	exists, _ := Exists(appID, "deleteq-test/level1/value1")
	if !exists {
		t.Fatal("value1 should exist before deletion")
	}

	// delete one nested value
	err := DeleteQ(appID, testKey, "$.level1.value1")
	testutil.AssertNoError(t, err, "delete nested field")

	// verify it was deleted (using dot notation)
	exists, _ = Exists(appID, "deleteq-test/level1/value1")
	if exists {
		t.Fatal("value1 should not exist after deletion")
	}

	// verify sibling value still exists (using bracket notation)
	value, err := Get(appID, "deleteq-test/level1/value2")
	testutil.AssertNoError(t, err, "get sibling value")
	if value.(string) != "second value" {
		t.Fatalf("sibling value was modified: expected 'second value', got '%s'", value.(string))
	}

	// verify parent dictionary still exists
	exists, _ = Exists(appID, "deleteq-test/level1")
	if !exists {
		t.Fatal("parent level1 should still exist")
	}

	// verify other branch still exists
	value, err = Get(appID, "deleteq-test/level2/nested")
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

	testKey := "deleteq-array-test"
	cleanup := setupTest(t, appID, testKey, arrayData)
	defer cleanup()

	// delete the second item (index 1)
	err := DeleteQ(appID, testKey, "$.items[1]")
	testutil.AssertNoError(t, err, "delete array element")

	// verify the modified array
	items, err := Get(appID, "deleteq-array-test/items")
	testutil.AssertNoError(t, err, "get items array after deletion")

	expected := []any{
		map[string]any{"id": int64(1), "name": "first"},
		map[string]any{"id": int64(3), "name": "third"},
	}

	if !reflect.DeepEqual(items, expected) {
		t.Fatalf("items should be %v, got %v", expected, items)
	}
}

func TestDeleteQNestedInArray(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	arrayData := map[string]any{
		"items": []any{
			map[string]any{"id": int64(1), "name": "first"},
			map[string]any{"id": int64(2), "name": "second"},
		},
	}

	testKey := "deleteq-nested-array-test"
	cleanup := setupTest(t, appID, testKey, arrayData)
	defer cleanup()

	// delete a field from an array element (using bracket notation)
	err := DeleteQ(appID, testKey, "$['items'][0]['name']")
	testutil.AssertNoError(t, err, "delete field from array element")

	// verify the name field was deleted
	exists, _ := Exists(appID, "deleteq-nested-array-test/items/0/name")
	if exists {
		t.Fatal("items[0].name should not exist after deletion")
	}

	// verify the id field still exists
	value, err := GetInt(appID, "deleteq-nested-array-test/items/0/id")
	testutil.AssertNoError(t, err, "get items[0].id")
	if value != int64(1) {
		t.Fatalf("items[0].id should be 1, got %v", value)
	}

	// verify the second item is unchanged
	value2, err := GetStr(appID, "deleteq-nested-array-test/items/1/name")
	testutil.AssertNoError(t, err, "get items[1].name")
	if value2 != "second" {
		t.Fatalf("items[1].name should be 'second', got '%s'", value2)
	}
}

func TestDeleteQRootKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testKey := "deleteq-root-test"
	cleanup := setupTest(t, appID, testKey, "test value")
	defer cleanup()

	assertKeyExists(t, appID, testKey, true)

	// delete using empty query (deletes entire root key)
	err := DeleteQ(appID, testKey, "")
	testutil.AssertNoError(t, err, "delete with empty query")

	assertKeyExists(t, appID, testKey, false)
}

func TestDeleteQWithDollarRoot(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testKey := "deleteq-dollar-test"
	cleanup := setupTest(t, appID, testKey, map[string]any{"field": "value"})
	defer cleanup()

	assertKeyExists(t, appID, testKey, true)

	// delete using "$" query (deletes entire root key)
	err := DeleteQ(appID, testKey, "$")
	testutil.AssertNoError(t, err, "delete with $ query")

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

	// delete a non-existent path should be idempotent (no error)
	err := DeleteQ(appID, testKey, "$.user.age")
	testutil.AssertNoError(t, err, "delete non-existent path should be idempotent")

	value, err := GetStr(appID, "deleteq-nonexistent-test/user/name")
	testutil.AssertNoError(t, err, "get existing field")
	if value != "John" {
		t.Fatalf("existing field should be unchanged, got %v", value)
	}
}

func TestDeleteQNonExistentRootKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// delete from a non-existent root key should be idempotent (no error)
	err := DeleteQ(appID, "nonexistent-root-key", "$['field']")
	testutil.AssertNoError(t, err, "delete from non-existent root should be idempotent")
}

func TestDeleteMissingKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// deleting a non-existent key should not error (idempotent)
	err := Delete(appID, "delete-missing-key")
	testutil.AssertNoError(t, err, "delete missing key should be idempotent")
}

func TestDeleteMultipleValues(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	data := map[string]any{
		"users": []any{
			map[string]any{"id": int64(1), "name": "Alice", "active": true},
			map[string]any{"id": int64(2), "name": "Bob", "active": false},
			map[string]any{"id": int64(3), "name": "Charlie", "active": true},
		},
	}

	testKey := "deleteq-multiple-test"
	cleanup := setupTest(t, appID, testKey, data)
	defer cleanup()

	// delete all 'active' fields from users
	err := DeleteQ(appID, testKey, "$.users[*].active")
	testutil.AssertNoError(t, err, "delete multiple fields")

	users, err := GetSlice(appID, "deleteq-multiple-test/users")
	testutil.AssertNoError(t, err, "get users")

	expected := []any{
		map[string]any{"id": int64(1), "name": "Alice"},
		map[string]any{"id": int64(2), "name": "Bob"},
		map[string]any{"id": int64(3), "name": "Charlie"},
	}

	if !reflect.DeepEqual(users, expected) {
		t.Fatalf("users should be %v, got %v", expected, users)
	}
}

func TestDeleteMultipleArrayElements(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	data := map[string]any{
		"data": map[string]any{
			"items": []any{
				map[string]any{"type": "A", "value": int64(1)},
				map[string]any{"type": "B", "value": int64(2)},
				map[string]any{"type": "A", "value": int64(3)},
				map[string]any{"type": "C", "value": int64(4)},
			},
		},
	}

	testKey := "deleteq-multiple-array-test"
	cleanup := setupTest(t, appID, testKey, data)
	defer cleanup()

	// delete all items with type "A" using a filter
	err := DeleteQ(appID, testKey, "$.data.items[?(@.type=='A')]")
	if err != nil {
		t.Fatalf("failed to delete items with type A: %v", err)
	}

	// verify remaining items
	items, err := GetSlice(appID, "deleteq-multiple-array-test/data/items")
	testutil.AssertNoError(t, err, "get remaining items")

	expected := []any{
		map[string]any{"type": "B", "value": int64(2)},
		map[string]any{"type": "C", "value": int64(4)},
	}

	if !reflect.DeepEqual(items, expected) {
		t.Fatalf("items should be %v, got %v", expected, items)
	}
}

func TestDeleteWildcardPaths(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	data := map[string]any{
		"departments": map[string]any{
			"sales": map[string]any{
				"budget":    int64(100000),
				"employees": int64(10),
			},
			"marketing": map[string]any{
				"budget":    int64(50000),
				"employees": int64(5),
			},
			"engineering": map[string]any{
				"budget":    int64(200000),
				"employees": int64(20),
			},
		},
	}

	testKey := "deleteq-wildcard-test"
	cleanup := setupTest(t, appID, testKey, data)
	defer cleanup()

	// delete all budget fields using wildcard
	err := DeleteQ(appID, testKey, "$.departments.*.budget")
	testutil.AssertNoError(t, err, "delete with wildcard")

	// verify all budget fields were deleted
	departments, err := GetMap(appID, "deleteq-wildcard-test/departments")
	testutil.AssertNoError(t, err, "get departments")

	expected := map[string]any{
		"sales": map[string]any{
			"employees": int64(10),
		},
		"marketing": map[string]any{
			"employees": int64(5),
		},
		"engineering": map[string]any{
			"employees": int64(20),
		},
	}
	if !reflect.DeepEqual(departments, expected) {
		t.Fatalf("departments should be %v, got %v", expected, departments)
	}
}
