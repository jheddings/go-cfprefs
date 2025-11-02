package cfprefs

import (
	"reflect"
	"testing"

	"github.com/jheddings/go-cfprefs/testutil"
)

// TestSetOperations tests various Set operations
func TestSetOperations(t *testing.T) {
	t.Run("basic set and get", func(t *testing.T) {
		key := "test-basic-set"
		err := Set(testAppID, key, "hello world")
		testutil.AssertNoError(t, err, "set simple value")
		defer Delete(testAppID, key)

		// Verify the value was set correctly
		value, err := Get(testAppID, key)
		testutil.AssertNoError(t, err, "get value")

		strValue, ok := value.(string)
		if !ok {
			t.Fatalf("value is not a string: got %T", value)
		}
		if strValue != "hello world" {
			t.Errorf("value does not match: expected 'hello world', got '%s'", strValue)
		}
	})

	t.Run("replace existing value", func(t *testing.T) {
		key := "test-replace"
		// Set initial value
		err := Set(testAppID, key, "initial value")
		testutil.AssertNoError(t, err, "set initial value")
		defer Delete(testAppID, key)

		// Replace with different type
		err = Set(testAppID, key, int64(42))
		testutil.AssertNoError(t, err, "replace value")

		value, err := Get(testAppID, key)
		testutil.AssertNoError(t, err, "get replaced value")

		intValue, ok := value.(int64)
		if !ok {
			t.Fatalf("value is not an int64: got %T", value)
		}
		if intValue != 42 {
			t.Errorf("value does not match: expected 42, got %d", intValue)
		}
	})

	t.Run("nested path creation", func(t *testing.T) {
		rootKey := "test-nested-creation"
		defer Delete(testAppID, rootKey)

		// Set value in non-existent nested path
		err := Set(testAppID, rootKey+"/level1/level2/value", "deeply nested")
		testutil.AssertNoError(t, err, "set deeply nested path")

		// Add sibling branch
		err = Set(testAppID, rootKey+"/sibling/level2/value", "sibling value")
		testutil.AssertNoError(t, err, "set sibling branch")

		// Verify structure
		expected := map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"value": "deeply nested",
				},
			},
			"sibling": map[string]any{
				"level2": map[string]any{
					"value": "sibling value",
				},
			},
		}

		root, err := GetMap(testAppID, rootKey)
		testutil.AssertNoError(t, err, "get root map")

		if !reflect.DeepEqual(root, expected) {
			t.Errorf("structure does not match expected: expected %v, got %v", expected, root)
		}
	})

	t.Run("add to existing nested structure", func(t *testing.T) {
		rootKey := "test-nested-add"
		defer Delete(testAppID, rootKey)

		// Create initial structure
		initial := map[string]any{
			"existing": map[string]any{
				"field": "value",
			},
		}
		err := Set(testAppID, rootKey, initial)
		testutil.AssertNoError(t, err, "set initial structure")

		// Add new field to existing object
		err = Set(testAppID, rootKey+"/existing/newField", int64(123))
		testutil.AssertNoError(t, err, "add field to existing object")

		// Verify both fields exist
		existing, err := GetMap(testAppID, rootKey+"/existing")
		testutil.AssertNoError(t, err, "get existing object")

		expected := map[string]any{
			"field":    "value",
			"newField": int64(123),
		}
		if !reflect.DeepEqual(existing, expected) {
			t.Errorf("object does not match expected: expected %v, got %v", expected, existing)
		}
	})
}

// TestSetArrayOperations tests array-specific Set operations
func TestSetArrayOperations(t *testing.T) {
	t.Run("update array element", func(t *testing.T) {
		rootKey := "test-array-update"
		defer Delete(testAppID, rootKey)

		// Create initial array
		initial := map[string]any{
			"items": []any{"first", "second", "third"},
		}
		err := Set(testAppID, rootKey, initial)
		testutil.AssertNoError(t, err, "set initial array")

		// Update middle element
		err = Set(testAppID, rootKey+"/items/1", "updated-second")
		testutil.AssertNoError(t, err, "update array element")

		// Verify update
		value, err := Get(testAppID, rootKey+"/items/1")
		testutil.AssertNoError(t, err, "get updated element")
		if value.(string) != "updated-second" {
			t.Errorf("element not updated: expected 'updated-second', got '%s'", value.(string))
		}
	})

	t.Run("append to existing array", func(t *testing.T) {
		rootKey := "test-array-append"
		defer Delete(testAppID, rootKey)

		// Create initial array
		initial := map[string]any{
			"items": []any{"one", "two"},
		}
		err := Set(testAppID, rootKey, initial)
		testutil.AssertNoError(t, err, "set initial array")

		// Append new element
		err = Set(testAppID, rootKey+"/items/~]", "three")
		testutil.AssertNoError(t, err, "append to array")

		// Verify length and content
		items, err := GetSlice(testAppID, rootKey+"/items")
		testutil.AssertNoError(t, err, "get array after append")
		
		if len(items) != 3 {
			t.Errorf("array length incorrect: expected 3, got %d", len(items))
		}
		if items[2].(string) != "three" {
			t.Errorf("appended value incorrect: expected 'three', got '%s'", items[2].(string))
		}
	})

	t.Run("append object to array", func(t *testing.T) {
		rootKey := "test-array-append-object"
		defer Delete(testAppID, rootKey)

		// Start with empty preferences
		err := Set(testAppID, rootKey, map[string]any{})
		testutil.AssertNoError(t, err, "set empty root")

		// Append objects with fields
		err = Set(testAppID, rootKey+"/pets/~]/name", "Fido")
		testutil.AssertNoError(t, err, "append first pet")

		err = Set(testAppID, rootKey+"/pets/~]/name", "Spot")
		testutil.AssertNoError(t, err, "append second pet")

		// Verify result
		pets, err := GetSlice(testAppID, rootKey+"/pets")
		testutil.AssertNoError(t, err, "get pets array")

		expected := []any{
			map[string]any{"name": "Fido"},
			map[string]any{"name": "Spot"},
		}
		if !reflect.DeepEqual(pets, expected) {
			t.Errorf("pets array does not match expected: expected %v, got %v", expected, pets)
		}
	})

	t.Run("nested array append", func(t *testing.T) {
		rootKey := "test-nested-array-append"
		defer Delete(testAppID, rootKey)

		// Create nested array structure through append
		err := Set(testAppID, rootKey+"/groups/~]/children/~]", "child-1")
		testutil.AssertNoError(t, err, "append to nested array")

		err = Set(testAppID, rootKey+"/groups/0/children/~]", "child-2")
		testutil.AssertNoError(t, err, "append to existing nested array")

		// Verify structure
		groups, err := GetSlice(testAppID, rootKey+"/groups")
		testutil.AssertNoError(t, err, "get groups")

		expected := []any{
			map[string]any{
				"children": []any{"child-1", "child-2"},
			},
		}
		if !reflect.DeepEqual(groups, expected) {
			t.Errorf("nested array does not match expected: expected %v, got %v", expected, groups)
		}
	})

	t.Run("array index out of bounds", func(t *testing.T) {
		rootKey := "test-array-bounds"
		defer Delete(testAppID, rootKey)

		// Create array with 3 elements
		initial := map[string]any{
			"items": []any{"a", "b", "c"},
		}
		err := Set(testAppID, rootKey, initial)
		testutil.AssertNoError(t, err, "set initial array")

		// Try to set beyond array bounds
		err = Set(testAppID, rootKey+"/items/10", "should fail")
		testutil.AssertError(t, err, "array index out of bounds")
	})
}

// TestSetRootOperations tests operations on root keys
func TestSetRootOperations(t *testing.T) {
	t.Run("replace entire root", func(t *testing.T) {
		rootKey := "test-root-replace"
		defer Delete(testAppID, rootKey)

		// Set initial value
		initial := map[string]any{"old": "value", "count": int64(1)}
		err := Set(testAppID, rootKey, initial)
		testutil.AssertNoError(t, err, "set initial value")

		// Replace entire root
		replacement := map[string]any{"new": "value", "count": int64(2)}
		err = Set(testAppID, rootKey, replacement)
		testutil.AssertNoError(t, err, "replace root")

		// Verify replacement
		value, err := Get(testAppID, rootKey)
		testutil.AssertNoError(t, err, "get replaced value")

		mapValue := value.(map[string]any)
		if _, hasOld := mapValue["old"]; hasOld {
			t.Error("old field still exists after replacement")
		}
		if mapValue["new"] != "value" {
			t.Errorf("new field incorrect: expected 'value', got '%v'", mapValue["new"])
		}
		if mapValue["count"] != int64(2) {
			t.Errorf("count incorrect: expected 2, got %v", mapValue["count"])
		}
	})
}

// TestSetErrors tests error conditions for Set operations
func TestSetErrors(t *testing.T) {
	t.Run("set through non-object segment", func(t *testing.T) {
		rootKey := "test-non-object"
		cleanup := setupTest(t, testAppID, rootKey, "string value")
		defer cleanup()

		// Try to set nested value through string
		err := Set(testAppID, rootKey+"/nested/value", "should fail")
		testutil.AssertError(t, err, "setting through non-object segment")
	})

	t.Run("invalid array operations", func(t *testing.T) {
		rootKey := "test-invalid-array"
		cleanup := setupTest(t, testAppID, rootKey, []any{"a", "b"})
		defer cleanup()

		// Try to set with invalid array index
		err := Set(testAppID, rootKey+"/not-a-number", "should fail")
		testutil.AssertError(t, err, "invalid array index")
	})
}