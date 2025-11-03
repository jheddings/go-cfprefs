package cfprefs

import (
	"reflect"
	"testing"

	"github.com/jheddings/go-cfprefs/testutil"
)

// Note: Test helpers are defined in get_test.go since they're shared

func TestSet(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test setting a simple value
	err := Set(appID, "test-key", "hello world")
	testutil.AssertNoError(t, err, "set simple value")
	defer Delete(appID, "test-key")

	// Verify the value was set correctly
	value, err := Get(appID, "test-key")
	testutil.AssertNoError(t, err, "get value")

	strValue, ok := value.(string)
	if !ok {
		t.Fatalf("value is not a string: got %T", value)
	}
	if strValue != "hello world" {
		t.Fatalf("value does not match: expected 'hello world', got '%s'", strValue)
	}

	// Test replacing an existing value
	err = Set(appID, "test-key", int64(42))
	testutil.AssertNoError(t, err, "replace value")

	value, err = Get(appID, "test-key")
	testutil.AssertNoError(t, err, "get replaced value")

	intValue, ok := value.(int64)
	if !ok {
		t.Fatalf("value is not an int64: got %T", value)
	}
	if intValue != 42 {
		t.Fatalf("value does not match: expected 42, got %d", intValue)
	}
}

func TestSetNested(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test setting a value in a nested path that doesn't exist yet
	err := Set(appID, "set-nested-test/level1/level1_1/value", "hello from nested path")
	testutil.AssertNoError(t, err, "set nested path")
	defer Delete(appID, "set-nested-test")

	// Add a sibling branch
	err = Set(appID, "set-nested-test/neighbor/level1_1/value", "hello from neighbor")
	testutil.AssertNoError(t, err, "set sibling branch")

	// Verify the expected structure was created
	expected := map[string]any{
		"level1": map[string]any{
			"level1_1": map[string]any{
				"value": "hello from nested path",
			},
		},
		"neighbor": map[string]any{
			"level1_1": map[string]any{
				"value": "hello from neighbor",
			},
		},
	}

	root, err := GetMap(appID, "set-nested-test")
	testutil.AssertNoError(t, err, "get root dictionary")

	if !reflect.DeepEqual(root, expected) {
		t.Fatalf("root value does not match expected: expected %v, got %v", expected, root)
	}

	// Test adding another value to the same nested structure
	err = Set(appID, "set-nested-test/level1/level1_1/another", int64(42))
	testutil.AssertNoError(t, err, "set another value in existing path")

	// Verify both values exist
	expected = map[string]any{
		"level1": map[string]any{
			"level1_1": map[string]any{
				"value":   "hello from nested path",
				"another": int64(42),
			},
		},
		"neighbor": map[string]any{
			"level1_1": map[string]any{
				"value": "hello from neighbor",
			},
		},
	}

	root, err = GetMap(appID, "set-nested-test")
	testutil.AssertNoError(t, err, "get root dictionary")

	if !reflect.DeepEqual(root, expected) {
		t.Fatalf("root value does not match expected: expected %v, got %v", expected, root)
	}
}

func TestSetArrayOperations(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Clean up any previous test data
	Delete(appID, "array-test")
	defer Delete(appID, "array-test")

	// Create an array with some initial items
	err := Set(appID, "array-test", map[string]any{
		"items": []any{"first", "second", "third"},
	})
	testutil.AssertNoError(t, err, "set initial array")

	// Test updating an array element
	err = Set(appID, "array-test/items/1", "updated-second")
	testutil.AssertNoError(t, err, "update array element")

	value, err := Get(appID, "array-test/items/1")
	testutil.AssertNoError(t, err, "get updated array element")
	if value.(string) != "updated-second" {
		t.Fatalf("array element not updated: expected 'updated-second', got '%s'", value.(string))
	}

	// Test appending to array
	err = Set(appID, "array-test/items/~]", "fourth")
	testutil.AssertNoError(t, err, "append to array")

	items, err := GetSlice(appID, "array-test/items")
	testutil.AssertNoError(t, err, "get array after append")
	if len(items) != 4 {
		t.Fatalf("array length incorrect: expected 4, got %d", len(items))
	}
	if items[3].(string) != "fourth" {
		t.Fatalf("appended value incorrect: expected 'fourth', got '%s'", items[3].(string))
	}

	// test appending a new element with a new field
	err = Set(appID, "array-test/pets/~]/name", "Fido")
	testutil.AssertNoError(t, err, "append new element")

	err = Set(appID, "array-test/pets/~]/name", "Spot")
	testutil.AssertNoError(t, err, "append new element")

	items, err = GetSlice(appID, "array-test/pets")
	testutil.AssertNoError(t, err, "get array after append")

	if !reflect.DeepEqual(items, []any{map[string]any{"name": "Fido"}, map[string]any{"name": "Spot"}}) {
		t.Fatalf("array does not match expected: found %v", items)
	}

	// test appending to a new array
	err = Set(appID, "array-test/deep-array/~]/children/~]", "child-1")
	testutil.AssertNoError(t, err, "append to new array")

	err = Set(appID, "array-test/deep-array/0/children/~]", "child-2")
	testutil.AssertNoError(t, err, "append child to new array")

	// NOTE: the brackets are required when querying names with dashes
	items, err = GetSlice(appID, "array-test/deep-array")
	testutil.AssertNoError(t, err, "get deep-array")
	if !reflect.DeepEqual(items, []any{map[string]any{"children": []any{"child-1", "child-2"}}}) {
		t.Fatalf("array does not match expected: found %v", items)
	}

	// Test array index out of bounds
	err = Set(appID, "array-test/items/10", "should fail")
	testutil.AssertError(t, err, "array index out of bounds")
	err = Set(appID, "array-test/pets/3", "should fail")
	testutil.AssertError(t, err, "array index out of bounds")
}

func TestSetReplaceRoot(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Clean up any previous test data
	Delete(appID, "replace-test")
	defer Delete(appID, "replace-test")

	// Set initial value
	err := Set(appID, "replace-test", map[string]any{"old": "value"})
	testutil.AssertNoError(t, err, "set initial value")

	// Replace entire root with $ or empty query
	newValue := map[string]any{"new": "value"}
	err = Set(appID, "replace-test", newValue)
	testutil.AssertNoError(t, err, "replace root with $")

	value, err := Get(appID, "replace-test")
	testutil.AssertNoError(t, err, "get replaced value")

	mapValue, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("value is not a map: got %T", value)
	}
	if _, hasOld := mapValue["old"]; hasOld {
		t.Fatal("old value still exists after replacement")
	}
	if mapValue["new"] != "value" {
		t.Fatalf("new value incorrect: expected 'value', got '%v'", mapValue["new"])
	}
}

func TestSetNonObjectSegment(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Set a simple string value
	cleanup := setupTest(t, appID, "simple-string", "hello")
	defer cleanup()

	// Try to set a nested value through a non-object segment
	err := Set(appID, "simple-string/nested/value", "should fail")
	testutil.AssertError(t, err, "setting through non-object segment")
}
