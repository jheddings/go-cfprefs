package cfprefs

import (
	"testing"

	"github.com/jheddings/go-cfprefs/testutil"
)

// Note: Test helpers are defined in get_test.go since they're shared

func TestSetKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Clean up any previous test data
	Delete(appID, "nested-test")

	// Test setting a value in a nested path that doesn't exist yet
	err := Set(appID, "nested-test/level1/level1.1/value", "hello from nested path")
	testutil.AssertNoError(t, err, "set nested path")
	defer Delete(appID, "nested-test")

	// Add a sibling branch
	err = Set(appID, "nested-test/neighbor/level1.1/value", "hello from neighbor")
	testutil.AssertNoError(t, err, "set sibling branch")

	// Verify the value was set correctly
	value, err := Get(appID, "nested-test/level1/level1.1/value")
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

	// Check level1.1 exists
	level2Value, ok := level1Dict["level1.1"]
	if !ok {
		t.Fatal("level1.1 dictionary not found")
	}

	level2Dict, ok := level2Value.(map[string]any)
	if !ok {
		t.Fatalf("level1.1 is not a dictionary: got %T", level2Value)
	}

	// Check the final value
	finalValue, ok := level2Dict["value"]
	if !ok {
		t.Fatal("final value not found in level1.1 dictionary")
	}

	finalStr, ok := finalValue.(string)
	if !ok {
		t.Fatalf("final value is not a string: got %T", finalValue)
	}
	if finalStr != "hello from nested path" {
		t.Fatalf("final value does not match: expected 'hello from nested path', got '%s'", finalStr)
	}

	// Test adding another value to the same nested structure
	err = Set(appID, "nested-test/level1/level1.1/another", int64(42))
	testutil.AssertNoError(t, err, "set another value in existing path")

	// Verify both values exist
	value1, err := Get(appID, "nested-test/level1/level1.1/value")
	testutil.AssertNoError(t, err, "get first value")
	if value1.(string) != "hello from nested path" {
		t.Fatal("first value was modified when setting second value")
	}

	value2, err := Get(appID, "nested-test/level1/level1.1/another")
	testutil.AssertNoError(t, err, "get second value")
	if value2.(int64) != 42 {
		t.Fatalf("second value does not match: expected 42, got %d", value2.(int64))
	}
}

// Error path tests

func TestSetEmptyKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test setting with empty keypath
	err := Set(appID, "", "value")
	testutil.AssertError(t, err, "empty keypath")
}

func TestSetKeypathOnlySlashes(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test keypath with only slashes
	err := Set(appID, "///", "value")
	testutil.AssertError(t, err, "keypath with only slashes")
}

func TestSetKeypathNonDictSegment(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Set a simple string value
	cleanup := setupTest(t, appID, "simple-string", "hello")
	defer cleanup()

	// Try to set a nested value through a non-dict segment
	err := Set(appID, "simple-string/nested/value", "should fail")
	testutil.AssertError(t, err, "setting through non-dict segment")
}
