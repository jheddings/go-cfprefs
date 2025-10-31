package cfprefs

import (
	"reflect"
	"testing"

	"github.com/jheddings/go-cfprefs/testutil"
)

// Note: Test helpers are defined in get_test.go since they're shared

func TestSet(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Clean up any previous test data
	Delete(appID, "test-key")

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

func TestSetQ(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Clean up any previous test data
	Delete(appID, "nested-test")

	// Test setting a value in a nested path that doesn't exist yet
	err := SetQ(appID, "nested-test", "$.level1.level1_1.value", "hello from nested path")
	testutil.AssertNoError(t, err, "set nested path")
	defer Delete(appID, "nested-test")

	// Add a sibling branch
	err = SetQ(appID, "nested-test", "$.neighbor.level1_1.value", "hello from neighbor")
	testutil.AssertNoError(t, err, "set sibling branch")

	// Verify the value was set correctly using GetQ
	value, err := Get(appID, "nested-test/level1/level1_1/value")
	testutil.AssertNoError(t, err, "get nested value")

	strValue, ok := value.(string)
	if !ok {
		t.Fatalf("value is not a string: got %T", value)
	}
	if strValue != "hello from nested path" {
		t.Fatalf("value does not match: expected 'hello from nested path', got '%s'", strValue)
	}

	// Verify the intermediate dictionaries were created
	rootValue, err := Get(appID, "nested-test")
	testutil.AssertNoError(t, err, "get root dictionary")

	rootDict, ok := rootValue.(map[string]any)
	if !ok {
		t.Fatalf("root value is not a dictionary: got %T", rootValue)
	}

	// Check level1 exists
	level1Value, ok := rootDict["level1"]
	if !ok {
		t.Fatal("level1 dictionary not found")
	}

	level1Dict, ok := level1Value.(map[string]any)
	if !ok {
		t.Fatalf("level1 is not a dictionary: got %T", level1Value)
	}

	// Check level1_1 exists
	level2Value, ok := level1Dict["level1_1"]
	if !ok {
		t.Fatal("level1_1 dictionary not found")
	}

	level2Dict, ok := level2Value.(map[string]any)
	if !ok {
		t.Fatalf("level1_1 is not a dictionary: got %T", level2Value)
	}

	// Check the final value
	finalValue, ok := level2Dict["value"]
	if !ok {
		t.Fatal("final value not found in level1_1 dictionary")
	}

	finalStr, ok := finalValue.(string)
	if !ok {
		t.Fatalf("final value is not a string: got %T", finalValue)
	}
	if finalStr != "hello from nested path" {
		t.Fatalf("final value does not match: expected 'hello from nested path', got '%s'", finalStr)
	}

	// Test adding another value to the same nested structure
	err = SetQ(appID, "nested-test", "$.level1.level1_1.another", int64(42))
	testutil.AssertNoError(t, err, "set another value in existing path")

	// Verify both values exist
	value1, err := Get(appID, "nested-test/level1/level1_1/value")
	testutil.AssertNoError(t, err, "get first value")
	if value1.(string) != "hello from nested path" {
		t.Fatal("first value was modified when setting second value")
	}

	value2, err := Get(appID, "nested-test/level1/level1_1/another")
	testutil.AssertNoError(t, err, "get second value")
	if value2.(int64) != 42 {
		t.Fatalf("second value does not match: expected 42, got %d", value2.(int64))
	}
}

func TestSetQArrayOperations(t *testing.T) {
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
	err = SetQ(appID, "array-test", "$.items[1]", "updated-second")
	testutil.AssertNoError(t, err, "update array element")

	value, err := Get(appID, "array-test/items/1")
	testutil.AssertNoError(t, err, "get updated array element")
	if value.(string) != "updated-second" {
		t.Fatalf("array element not updated: expected 'updated-second', got '%s'", value.(string))
	}

	// Test appending to array
	err = SetQ(appID, "array-test", "$.items[]", "fourth")
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
	err = SetQ(appID, "array-test", "$.pets[].name", "Fido")
	testutil.AssertNoError(t, err, "append new element")

	err = SetQ(appID, "array-test", "$.pets[].name", "Spot")
	testutil.AssertNoError(t, err, "append new element")

	items, err = GetSlice(appID, "array-test/pets")
	testutil.AssertNoError(t, err, "get array after append")

	if !reflect.DeepEqual(items, []any{map[string]any{"name": "Fido"}, map[string]any{"name": "Spot"}}) {
		t.Fatalf("array does not match expected: found %v", items)
	}

	// test appending to a new array
	err = SetQ(appID, "array-test", "$.deep-array[].children[]", "child-1")
	testutil.AssertNoError(t, err, "append to new array")

	err = SetQ(appID, "array-test", "$.deep-array[0].children[]", "child-2")
	testutil.AssertNoError(t, err, "append child to new array")

	// NOTE: the brackets are required when querying names with dashes
	items, err = GetSlice(appID, "array-test/deep-array")
	testutil.AssertNoError(t, err, "get deep-array")
	if !reflect.DeepEqual(items, []any{map[string]any{"children": []any{"child-1", "child-2"}}}) {
		t.Fatalf("array does not match expected: found %v", items)
	}

	// Test array index out of bounds
	err = SetQ(appID, "array-test", "$.items[10]", "should fail")
	testutil.AssertError(t, err, "array index out of bounds")
	err = SetQ(appID, "array-test", "$.pets[3]", "should fail")
	testutil.AssertError(t, err, "array index out of bounds")
}

func TestSetQReplaceRoot(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Clean up any previous test data
	Delete(appID, "replace-test")
	defer Delete(appID, "replace-test")

	// Set initial value
	err := Set(appID, "replace-test", map[string]any{"old": "value"})
	testutil.AssertNoError(t, err, "set initial value")

	// Replace entire root with $ or empty query
	newValue := map[string]any{"new": "value"}
	err = SetQ(appID, "replace-test", "$", newValue)
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

func TestSetQNonObjectSegment(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Set a simple string value
	cleanup := setupTest(t, appID, "simple-string", "hello")
	defer cleanup()

	// Try to set a nested value through a non-object segment
	err := SetQ(appID, "simple-string", "$.nested.value", "should fail")
	testutil.AssertError(t, err, "setting through non-object segment")
}
