package cfprefs

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/jheddings/go-cfprefs/testutil"
)

// setupTest sets a value and returns a cleanup function
func setupTest(t *testing.T, appID, key string, value any) func() {
	t.Helper()
	err := Set(appID, key, value)

	if err != nil {
		t.Fatalf("failed to set test value: %v", err)
	}

	return func() {
		if err := Delete(appID, key); err != nil {
			t.Errorf("failed to cleanup test key: %v", err)
		}
	}
}

func TestGetKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testValue := map[string]any{
		"string": "hello",
		"number": int64(42),
		"float":  3.14,
		"bool":   true,
	}

	cleanup := setupTest(t, appID, "map-test", testValue)
	defer cleanup()

	// retrieve a nested value using keypath
	value, err := Get(appID, "map-test/string")
	testutil.AssertNoError(t, err, "get nested string")

	strValue, ok := value.(string)
	if !ok {
		t.Fatalf("value is not a string: got %T", value)
	}
	if strValue != "hello" {
		t.Fatalf("value does not match: expected 'hello', got '%s'", strValue)
	}

	// retrieve another nested value
	value, err = Get(appID, "map-test/number")
	testutil.AssertNoError(t, err, "get nested number")

	numValue, ok := value.(int64)
	if !ok {
		t.Fatalf("value is not an int64: got %T", value)
	}
	if numValue != 42 {
		t.Fatalf("value does not match: expected 42, got %d", numValue)
	}

	// error case: non-existent key in path
	_, err = Get(appID, "map-test/nonexistent")
	testutil.AssertError(t, err, "non-existent key in path")

	// retrieve the whole map without keypath (backward compatibility)
	value, err = Get(appID, "map-test")
	testutil.AssertNoError(t, err, "get whole map")

	if !reflect.DeepEqual(value, testValue) {
		t.Fatalf("map does not match: expected %v, got %v", testValue, value)
	}
}

func TestGetStr(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"
	testValue := "hello"

	cleanup := setupTest(t, appID, "str-test", testValue)
	defer cleanup()

	assertKeyExists(t, appID, "str-test", true)

	// Type mismatch - should error
	_, err := GetInt(appID, "str-test")
	testutil.AssertError(t, err, "non-int value")

	value, err := GetStr(appID, "str-test")
	testutil.AssertNoError(t, err, "GetStr")
	if value != testValue {
		t.Fatalf("expected %s, got %s", testValue, value)
	}
}

func TestGetInt(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"
	testValue := int64(42)

	cleanup := setupTest(t, appID, "int-test", testValue)
	defer cleanup()

	assertKeyExists(t, appID, "int-test", true)

	// Type mismatch - should error
	_, err := GetBool(appID, "int-test")
	testutil.AssertError(t, err, "non-bool value")

	value, err := GetInt(appID, "int-test")
	testutil.AssertNoError(t, err, "GetInt")
	if value != testValue {
		t.Fatalf("expected %d, got %d", testValue, value)
	}
}

func TestGetFloat(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"
	testValue := 3.14159

	cleanup := setupTest(t, appID, "float-test", testValue)
	defer cleanup()

	assertKeyExists(t, appID, "float-test", true)

	// Type mismatch - should error
	_, err := GetData(appID, "float-test")
	testutil.AssertError(t, err, "non-data value")

	value, err := GetFloat(appID, "float-test")
	testutil.AssertNoError(t, err, "GetFloat")
	if math.Abs(value-testValue) > 1e-10 {
		t.Fatalf("expected %f, got %f", testValue, value)
	}
}

func TestGetBool(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"
	testValue := true

	cleanup := setupTest(t, appID, "bool-test", testValue)
	defer cleanup()

	assertKeyExists(t, appID, "bool-test", true)

	// Type mismatch - should error
	_, err := GetDate(appID, "bool-test")
	testutil.AssertError(t, err, "non-date value")

	value, err := GetBool(appID, "bool-test")
	testutil.AssertNoError(t, err, "GetBool")
	if value != testValue {
		t.Fatalf("expected %t, got %t", testValue, value)
	}
}

func TestGetDate(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"
	// Use a fixed time for deterministic testing
	testValue := time.Date(2024, 10, 15, 12, 30, 45, 123456000, time.UTC)

	cleanup := setupTest(t, appID, "date-test", testValue)
	defer cleanup()

	assertKeyExists(t, appID, "date-test", true)

	// Type mismatch - should error
	_, err := GetMap(appID, "date-test")
	testutil.AssertError(t, err, "non-map value")

	value, err := GetDate(appID, "date-test")
	testutil.AssertNoError(t, err, "GetDate")

	if !testutil.ValuesEqualApprox(testValue, value) {
		t.Fatalf("expected %v, got %v", testValue, value)
	}
}

func TestGetSlice(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"
	testTime := time.Date(2024, 10, 15, 12, 30, 45, 123456000, time.UTC)
	testValue := []any{int64(123), 3.14159, true, testTime}

	cleanup := setupTest(t, appID, "slice-test", testValue)
	defer cleanup()

	assertKeyExists(t, appID, "slice-test", true)

	// Type mismatch - should error
	_, err := GetStr(appID, "slice-test")
	testutil.AssertError(t, err, "non-string value")

	value, err := GetSlice(appID, "slice-test")
	testutil.AssertNoError(t, err, "GetSlice")

	if !testutil.ValuesEqualApprox(testValue, value) {
		t.Fatalf("expected %v, got %v", testValue, value)
	}
}

func TestGetData(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"
	testValue := []byte("hello world")

	cleanup := setupTest(t, appID, "data-test", testValue)
	defer cleanup()

	assertKeyExists(t, appID, "data-test", true)

	// Type mismatch - should error
	_, err := GetFloat(appID, "data-test")
	testutil.AssertError(t, err, "non-float value")

	value, err := GetData(appID, "data-test")
	testutil.AssertNoError(t, err, "GetData")
	if !reflect.DeepEqual(value, testValue) {
		t.Fatalf("expected %v, got %v", testValue, value)
	}
}

func TestGetMap(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"
	testTime := time.Date(2024, 10, 15, 12, 30, 45, 123456000, time.UTC)
	testValue := map[string]any{
		"string": "hello",
		"number": int64(456),
		"float":  2.71828,
		"bool":   false,
		"time":   testTime,
	}

	cleanup := setupTest(t, appID, "map-test", testValue)
	defer cleanup()

	assertKeyExists(t, appID, "map-test", true)

	// Type mismatch - should error
	_, err := GetSlice(appID, "map-test")
	testutil.AssertError(t, err, "non-slice value")

	value, err := GetMap(appID, "map-test")
	testutil.AssertNoError(t, err, "GetMap")

	if !testutil.ValuesEqualApprox(testValue, value) {
		t.Fatalf("expected %v, got %v", testValue, value)
	}
}

func TestGetMissingKey(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test getting a key that doesn't exist
	_, err := Get(appID, "this-key-does-not-exist")
	testutil.AssertError(t, err, "missing key")
}

func TestGetEmptyKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test empty keypath
	_, err := Get(appID, "")
	testutil.AssertError(t, err, "empty keypath")
}

func TestGetKeypathOnlySlashes(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Test keypath with only slashes
	_, err := Get(appID, "///")
	testutil.AssertError(t, err, "keypath with only slashes")
}

func TestGetKeypathNonDictSegment(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// Set a simple string value
	cleanup := setupTest(t, appID, "simple-string", "hello")
	defer cleanup()

	// Try to traverse through it as if it were a dictionary
	_, err := Get(appID, "simple-string/nested")
	testutil.AssertError(t, err, "traversing non-dict segment")
}
