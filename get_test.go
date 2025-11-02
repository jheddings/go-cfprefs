package cfprefs

import (
	"reflect"
	"testing"
	"time"

	"github.com/jheddings/go-cfprefs/testutil"
)

// TestGetTyped tests all typed getter functions using a table-driven approach
func TestGetTyped(t *testing.T) {
	testTime := time.Date(2024, 10, 15, 12, 30, 45, 123456000, time.UTC)

	tests := []struct {
		name           string
		key            string
		value          any
		getter         func(appID, key string) (any, error)
		wrongGetter    func(appID, key string) (any, error)
		wrongGetterErr string
		checkEqual     func(expected, actual any) bool
	}{
		{
			name:           "string value",
			key:            "test-string",
			value:          "hello world",
			getter:         func(appID, key string) (any, error) { return GetStr(appID, key) },
			wrongGetter:    func(appID, key string) (any, error) { return GetInt(appID, key) },
			wrongGetterErr: "non-int value",
			checkEqual: func(expected, actual any) bool {
				return expected.(string) == actual.(string)
			},
		},
		{
			name:           "int value",
			key:            "test-int",
			value:          int64(42),
			getter:         func(appID, key string) (any, error) { return GetInt(appID, key) },
			wrongGetter:    func(appID, key string) (any, error) { return GetBool(appID, key) },
			wrongGetterErr: "non-bool value",
			checkEqual: func(expected, actual any) bool {
				return expected.(int64) == actual.(int64)
			},
		},
		{
			name:           "float value",
			key:            "test-float",
			value:          3.14159,
			getter:         func(appID, key string) (any, error) { return GetFloat(appID, key) },
			wrongGetter:    func(appID, key string) (any, error) { return GetData(appID, key) },
			wrongGetterErr: "non-data value",
			checkEqual: func(expected, actual any) bool {
				return testutil.ValuesEqualApprox(expected, actual)
			},
		},
		{
			name:           "bool value",
			key:            "test-bool",
			value:          true,
			getter:         func(appID, key string) (any, error) { return GetBool(appID, key) },
			wrongGetter:    func(appID, key string) (any, error) { return GetDate(appID, key) },
			wrongGetterErr: "non-date value",
			checkEqual: func(expected, actual any) bool {
				return expected.(bool) == actual.(bool)
			},
		},
		{
			name:           "date value",
			key:            "test-date",
			value:          testTime,
			getter:         func(appID, key string) (any, error) { return GetDate(appID, key) },
			wrongGetter:    func(appID, key string) (any, error) { return GetMap(appID, key) },
			wrongGetterErr: "non-map value",
			checkEqual: func(expected, actual any) bool {
				return testutil.ValuesEqualApprox(expected, actual)
			},
		},
		{
			name:           "slice value",
			key:            "test-slice",
			value:          []any{int64(123), 3.14159, true, testTime},
			getter:         func(appID, key string) (any, error) { return GetSlice(appID, key) },
			wrongGetter:    func(appID, key string) (any, error) { return GetStr(appID, key) },
			wrongGetterErr: "non-string value",
			checkEqual: func(expected, actual any) bool {
				return testutil.ValuesEqualApprox(expected, actual)
			},
		},
		{
			name:           "data value",
			key:            "test-data",
			value:          []byte("hello world"),
			getter:         func(appID, key string) (any, error) { return GetData(appID, key) },
			wrongGetter:    func(appID, key string) (any, error) { return GetFloat(appID, key) },
			wrongGetterErr: "non-float value",
			checkEqual: func(expected, actual any) bool {
				return reflect.DeepEqual(expected, actual)
			},
		},
		{
			name: "map value",
			key:  "test-map",
			value: map[string]any{
				"string": "hello",
				"number": int64(456),
				"float":  2.71828,
				"bool":   false,
				"time":   testTime,
			},
			getter:         func(appID, key string) (any, error) { return GetMap(appID, key) },
			wrongGetter:    func(appID, key string) (any, error) { return GetSlice(appID, key) },
			wrongGetterErr: "non-slice value",
			checkEqual: func(expected, actual any) bool {
				return testutil.ValuesEqualApprox(expected, actual)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupTest(t, testAppID, tt.key, tt.value)
			defer cleanup()

			// Verify key exists
			assertKeyExists(t, testAppID, tt.key, true)

			// Test correct getter
			got, err := tt.getter(testAppID, tt.key)
			testutil.AssertNoError(t, err, "correct getter")
			if !tt.checkEqual(tt.value, got) {
				t.Errorf("value mismatch: expected %v, got %v", tt.value, got)
			}

			// Test wrong getter (type mismatch)
			_, err = tt.wrongGetter(testAppID, tt.key)
			testutil.AssertError(t, err, tt.wrongGetterErr)
		})
	}
}

// TestGetErrors tests error conditions for Get operations
func TestGetErrors(t *testing.T) {
	t.Run("missing key", func(t *testing.T) {
		_, err := Get(testAppID, "this-key-does-not-exist")
		testutil.AssertError(t, err, "missing key")
	})

	t.Run("empty key name", func(t *testing.T) {
		_, err := Get(testAppID, "")
		testutil.AssertError(t, err, "empty key name")
	})

	t.Run("query errors", func(t *testing.T) {
		t.Run("non-existent field", func(t *testing.T) {
			_, err := Get(testAppID, "userData/nonexistent")
			testutil.AssertError(t, err, "non-existent field should return error")
		})

		t.Run("non-existent root key", func(t *testing.T) {
			_, err := Get(testAppID, "nonexistent/name")
			testutil.AssertError(t, err, "non-existent root key should return error")
		})

		t.Run("invalid JSONPath", func(t *testing.T) {
			_, err := Get(testAppID, "userData/name[0]")
			testutil.AssertError(t, err, "invalid JSONPath should return error")
		})
	})
}

// TestQuery tests JSON Pointer query functionality
func TestQuery(t *testing.T) {
	testData := map[string]any{
		"name": "Jane Doe",
		"age":  int64(30),
		"city": "Anytown",
		"items": []any{
			"first",
			"second",
			"third",
		},
		"pets": []any{
			map[string]any{
				"name": "Fluffy",
				"type": "dog",
			},
			map[string]any{
				"name": "Whiskers",
				"type": "cat",
			},
		},
	}

	cleanup := setupTest(t, testAppID, "bio", testData)
	defer cleanup()

	tests := []struct {
		name     string
		path     string
		expected any
	}{
		{
			name:     "simple field access",
			path:     "bio/name",
			expected: "Jane Doe",
		},
		{
			name:     "numeric field access",
			path:     "bio/age",
			expected: int64(30),
		},
		{
			name:     "root object access",
			path:     "bio",
			expected: testData,
		},
		{
			name:     "array element access",
			path:     "bio/items/0",
			expected: "first",
		},
		{
			name:     "nested array field access",
			path:     "bio/pets/0/name",
			expected: "Fluffy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := Get(testAppID, tt.path)
			testutil.AssertNoError(t, err, "get "+tt.path)
			if !testutil.ValuesEqualApprox(tt.expected, value) {
				t.Errorf("expected %v, got %v", tt.expected, value)
			}
		})
	}
}
