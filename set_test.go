package cfprefs

import (
	"testing"
)

func TestSetKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// clear previous tests
	err := Delete(appID, "nested-test/level1/level1.1/value")
	if err != nil {
		t.Fatalf("failed to delete keypath: %v", err)
	}

	// test setting a value in a nested path that doesn't exist yet
	err = Set(appID, "nested-test/level1/level1.1/value", "hello from nested path")
	err = Set(appID, "nested-test/neighbor/level1.1/value", "hello from nested path")
	if err != nil {
		t.Fatalf("failed to set keypath with missing dictionaries: %v", err)
	}

	// verify the value was set correctly by retrieving it
	value, err := Get(appID, "nested-test/level1/level1.1/value")
	if err != nil {
		t.Fatalf("failed to get keypath: %v", err)
	}

	strValue, ok := value.(string)
	if !ok {
		t.Fatalf("value is not a string: got %T", value)
	}

	if strValue != "hello from nested path" {
		t.Fatalf("value does not match: expected 'hello from nested path', got '%s'", strValue)
	}

	// verify the intermediate dictionaries were created correctly
	rootValue, err := Get(appID, "nested-test")
	if err != nil {
		t.Fatalf("failed to get root dictionary: %v", err)
	}

	rootDict, ok := rootValue.(map[string]any)
	if !ok {
		t.Fatalf("root value is not a dictionary: got %T", rootValue)
	}

	// check level1 exists
	level1Value, ok := rootDict["level1"]
	if !ok {
		t.Fatal("level1 dictionary not found")
	}

	level1Dict, ok := level1Value.(map[string]any)
	if !ok {
		t.Fatalf("level1 is not a dictionary: got %T", level1Value)
	}

	// check level1.1 exists
	level2Value, ok := level1Dict["level1.1"]
	if !ok {
		t.Fatal("level1.1 dictionary not found")
	}

	level2Dict, ok := level2Value.(map[string]any)
	if !ok {
		t.Fatalf("level1.1 is not a dictionary: got %T", level2Value)
	}

	// check the final value
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

	// test adding another value to the same nested structure
	err = Set(appID, "nested-test/level1/level1.1/another", int64(42))
	if err != nil {
		t.Fatalf("failed to set another value in existing path: %v", err)
	}

	// verify both values exist
	value1, err := Get(appID, "nested-test/level1/level1.1/value")
	if err != nil {
		t.Fatalf("failed to get first value after adding second: %v", err)
	}

	if value1.(string) != "hello from nested path" {
		t.Fatal("first value was modified when setting second value")
	}

	value2, err := Get(appID, "nested-test/level1/level1.1/another")
	if err != nil {
		t.Fatalf("failed to get second value: %v", err)
	}

	if value2.(int64) != 42 {
		t.Fatalf("second value does not match: expected 42, got %d", value2.(int64))
	}
}
