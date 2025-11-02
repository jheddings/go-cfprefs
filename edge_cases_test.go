package cfprefs

import (
	"strings"
	"testing"

	"github.com/jheddings/go-cfprefs/testutil"
)

// TestEdgeCaseKeys tests edge cases for key names and paths
func TestEdgeCaseKeys(t *testing.T) {
	t.Run("empty string key", func(t *testing.T) {
		// Empty key should error
		err := Set(testAppID, "", "value")
		testutil.AssertError(t, err, "set with empty key")

		_, err = Get(testAppID, "")
		testutil.AssertError(t, err, "get with empty key")

		_, err = Exists(testAppID, "")
		testutil.AssertError(t, err, "exists with empty key")
	})

	t.Run("special characters in keys", func(t *testing.T) {
		tests := []struct {
			name  string
			key   string
			value any
		}{
			{"key with spaces", "test key with spaces", "value"},
			{"key with dots", "test.key.with.dots", "value"},
			{"key with unicode", "test-κλειδί-🔑", "value"},
			{"key with tabs", "test\tkey\twith\ttabs", "value"},
			{"key with newlines", "test\nkey\nwith\nnewlines", "value"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				key := "edge-case-keys/" + tt.key
				cleanup := setupTest(t, testAppID, key, tt.value)
				defer cleanup()

				// Verify we can read it back
				got, err := Get(testAppID, key)
				testutil.AssertNoError(t, err, "get special character key")
				if got != tt.value {
					t.Errorf("expected %v, got %v", tt.value, got)
				}
			})
		}
	})

	t.Run("very long key names", func(t *testing.T) {
		// Create a very long key name
		longKey := "edge-case-long/" + strings.Repeat("a", 1000)
		cleanup := setupTest(t, testAppID, longKey, "value")
		defer cleanup()

		exists, err := Exists(testAppID, longKey)
		testutil.AssertNoError(t, err, "check long key exists")
		if !exists {
			t.Error("long key should exist")
		}
	})

	t.Run("deeply nested paths", func(t *testing.T) {
		// Create a deeply nested structure
		path := "edge-case-deep"
		for i := 0; i < 50; i++ {
			path += "/level"
		}
		path += "/value"

		err := Set(testAppID, path, "deeply nested")
		testutil.AssertNoError(t, err, "set deeply nested value")
		defer Delete(testAppID, "edge-case-deep")

		got, err := Get(testAppID, path)
		testutil.AssertNoError(t, err, "get deeply nested value")
		if got != "deeply nested" {
			t.Errorf("expected 'deeply nested', got %v", got)
		}
	})
}

// TestEdgeCaseValues tests edge cases for values
func TestEdgeCaseValues(t *testing.T) {
	t.Run("nil and zero values", func(t *testing.T) {
		tests := []struct {
			name  string
			key   string
			value any
		}{
			{"empty string", "test-empty-string", ""},
			{"zero int", "test-zero-int", int64(0)},
			{"zero float", "test-zero-float", 0.0},
			{"false bool", "test-false-bool", false},
			{"empty slice", "test-empty-slice", []any{}},
			{"empty map", "test-empty-map", map[string]any{}},
			{"empty bytes", "test-empty-bytes", []byte{}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cleanup := setupTest(t, testAppID, tt.key, tt.value)
				defer cleanup()

				got, err := Get(testAppID, tt.key)
				testutil.AssertNoError(t, err, "get "+tt.name)
				if !testutil.ValuesEqualApprox(tt.value, got) {
					t.Errorf("expected %v, got %v", tt.value, got)
				}
			})
		}
	})

	t.Run("very large values", func(t *testing.T) {
		t.Run("large string", func(t *testing.T) {
			// 10MB string
			largeString := strings.Repeat("a", 10*1024*1024)
			key := "test-large-string"
			cleanup := setupTest(t, testAppID, key, largeString)
			defer cleanup()

			got, err := GetStr(testAppID, key)
			testutil.AssertNoError(t, err, "get large string")
			if len(got) != len(largeString) {
				t.Errorf("expected string of length %d, got %d", len(largeString), len(got))
			}
		})

		t.Run("large array", func(t *testing.T) {
			// Array with 10000 elements
			largeArray := make([]any, 10000)
			for i := range largeArray {
				largeArray[i] = int64(i)
			}
			key := "test-large-array"
			cleanup := setupTest(t, testAppID, key, largeArray)
			defer cleanup()

			got, err := GetSlice(testAppID, key)
			testutil.AssertNoError(t, err, "get large array")
			if len(got) != len(largeArray) {
				t.Errorf("expected array of length %d, got %d", len(largeArray), len(got))
			}
		})

		t.Run("large nested map", func(t *testing.T) {
			// Create a large nested structure
			largeMap := make(map[string]any)
			for i := 0; i < 1000; i++ {
				largeMap[strings.Repeat("key", i%10)] = map[string]any{
					"index": int64(i),
					"data":  strings.Repeat("value", i%5),
				}
			}
			key := "test-large-map"
			cleanup := setupTest(t, testAppID, key, largeMap)
			defer cleanup()

			got, err := GetMap(testAppID, key)
			testutil.AssertNoError(t, err, "get large map")
			if len(got) != len(largeMap) {
				t.Errorf("expected map of length %d, got %d", len(largeMap), len(got))
			}
		})
	})
}

// TestEdgeCaseArrayOperations tests edge cases for array operations
func TestEdgeCaseArrayOperations(t *testing.T) {
	t.Run("negative array indices", func(t *testing.T) {
		key := "test-negative-index"
		cleanup := setupTest(t, testAppID, key, []any{"a", "b", "c"})
		defer cleanup()

		// JSON Pointer doesn't support negative indices
		_, err := Get(testAppID, key+"/-1")
		testutil.AssertError(t, err, "negative array index")
	})

	t.Run("array with mixed types", func(t *testing.T) {
		key := "test-mixed-array"
		mixedArray := []any{
			"string",
			int64(123),
			3.14,
			true,
			map[string]any{"nested": "object"},
			[]any{"nested", "array"},
		}
		cleanup := setupTest(t, testAppID, key, mixedArray)
		defer cleanup()

		// Verify each element
		for i, expected := range mixedArray {
			got, err := Get(testAppID, key+"/"+string(rune(i+'0')))
			testutil.AssertNoError(t, err, "get array element")
			if !testutil.ValuesEqualApprox(expected, got) {
				t.Errorf("element %d: expected %v, got %v", i, expected, got)
			}
		}
	})

	t.Run("sparse array operations", func(t *testing.T) {
		key := "test-sparse-array"
		// Create array with gaps
		sparseArray := []any{"a", "b", "c", "d", "e"}
		cleanup := setupTest(t, testAppID, key, sparseArray)
		defer cleanup()

		// Delete middle elements to create gaps
		err := Delete(testAppID, key+"/1")
		testutil.AssertNoError(t, err, "delete array element")
		err = Delete(testAppID, key+"/1")  // Delete index 1 again (which is now "c")
		testutil.AssertNoError(t, err, "delete array element")

		// Verify remaining structure
		got, err := GetSlice(testAppID, key)
		testutil.AssertNoError(t, err, "get sparse array")
		expected := []any{"a", "d", "e"}
		if !testutil.ValuesEqualApprox(expected, got) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}

// TestEdgeCaseQueries tests edge cases for JSON Pointer queries
func TestEdgeCaseQueries(t *testing.T) {
	testData := map[string]any{
		"~field":  "tilde field",
		"/field":  "slash field",
		"~/field": "tilde slash field",
		"array": []any{
			"element0",
			map[string]any{"~": "tilde", "/": "slash"},
		},
	}
	
	key := "test-edge-queries"
	cleanup := setupTest(t, testAppID, key, testData)
	defer cleanup()

	t.Run("escaped characters in paths", func(t *testing.T) {
		tests := []struct {
			name     string
			path     string
			expected any
		}{
			{"tilde field", key + "/~0field", "tilde field"},
			{"slash field", key + "/~1field", "slash field"},
			{"tilde slash field", key + "/~0~1field", "tilde slash field"},
			{"array tilde", key + "/array/1/~0", "tilde"},
			{"array slash", key + "/array/1/~1", "slash"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Get(testAppID, tt.path)
				testutil.AssertNoError(t, err, "get with escaped path")
				if got != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, got)
				}
			})
		}
	})
}
