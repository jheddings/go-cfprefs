package cfprefs

import (
	"testing"
)

func TestGetKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testValue := map[string]any{
		"string": "hello",
		"number": int64(42),
		"float":  3.14,
		"bool":   true,
	}

	err := Set(appID, "map-test", testValue)
	if err != nil {
		t.Fatal(err)
	}

	// retrieve a nested value using keypath
	value, err := Get(appID, "map-test/string")
	if err != nil {
		t.Fatalf("failed to get keypath: %v", err)
	}

	strValue, ok := value.(string)
	if !ok {
		t.Fatalf("value is not a string: got %T", value)
	}

	if strValue != "hello" {
		t.Fatalf("value does not match: expected 'hello', got '%s'", strValue)
	}

	// retrieve another nested value
	value, err = Get(appID, "map-test/number")
	if err != nil {
		t.Fatalf("failed to get keypath: %v", err)
	}

	numValue, ok := value.(int64)
	if !ok {
		t.Fatalf("value is not an int64: got %T", value)
	}

	if numValue != 42 {
		t.Fatalf("value does not match: expected 42, got %d", numValue)
	}

	// error case: non-existent key in path
	_, err = Get(appID, "map-test/nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent key, got nil")
	}

	// retrieve the whole map without keypath (backward compatibility)
	value, err = Get(appID, "map-test")
	if err != nil {
		t.Fatalf("failed to get map: %v", err)
	}

	retrievedMap, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("value is not map[string]any: got %T", value)
	}

	if len(retrievedMap) != len(testValue) {
		t.Fatalf("map length does not match: expected %d, got %d", len(testValue), len(retrievedMap))
	}
}
