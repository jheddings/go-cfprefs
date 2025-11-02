package cfprefs

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkGet benchmarks various Get operations
func BenchmarkGet(b *testing.B) {
	// Setup test data
	setupData := map[string]any{
		"string": "benchmark test value",
		"int":    int64(123456),
		"float":  3.14159,
		"bool":   true,
		"date":   time.Now(),
		"array":  []any{"a", "b", "c", "d", "e"},
		"object": map[string]any{
			"field1": "value1",
			"field2": int64(42),
			"nested": map[string]any{
				"deep": "value",
			},
		},
	}

	key := "benchmark-get"
	cleanup := setupTest(b, testAppID, key, setupData)
	defer cleanup()

	b.Run("GetStr", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetStr(testAppID, key+"/string")
		}
	})

	b.Run("GetInt", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetInt(testAppID, key+"/int")
		}
	})

	b.Run("GetFloat", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetFloat(testAppID, key+"/float")
		}
	})

	b.Run("GetBool", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetBool(testAppID, key+"/bool")
		}
	})

	b.Run("GetDate", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetDate(testAppID, key+"/date")
		}
	})

	b.Run("GetSlice", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetSlice(testAppID, key+"/array")
		}
	})

	b.Run("GetMap", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GetMap(testAppID, key+"/object")
		}
	})

	b.Run("GetNested", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = Get(testAppID, key+"/object/nested/deep")
		}
	})
}

// BenchmarkSet benchmarks various Set operations
func BenchmarkSet(b *testing.B) {
	b.Run("SetStr", func(b *testing.B) {
		key := "benchmark-set-str"
		defer Delete(testAppID, key)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Set(testAppID, key, fmt.Sprintf("value-%d", i))
		}
	})

	b.Run("SetInt", func(b *testing.B) {
		key := "benchmark-set-int"
		defer Delete(testAppID, key)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Set(testAppID, key, int64(i))
		}
	})

	b.Run("SetNested", func(b *testing.B) {
		baseKey := "benchmark-set-nested"
		defer Delete(testAppID, baseKey)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			path := fmt.Sprintf("%s/level1/level2/item%d", baseKey, i)
			_ = Set(testAppID, path, fmt.Sprintf("value-%d", i))
		}
	})

	b.Run("SetArray", func(b *testing.B) {
		key := "benchmark-set-array"
		defer Delete(testAppID, key)
		
		array := make([]any, 100)
		for i := range array {
			array[i] = fmt.Sprintf("element-%d", i)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Set(testAppID, key, array)
		}
	})

	b.Run("SetMap", func(b *testing.B) {
		key := "benchmark-set-map"
		defer Delete(testAppID, key)
		
		mapData := make(map[string]any)
		for i := 0; i < 50; i++ {
			mapData[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Set(testAppID, key, mapData)
		}
	})
}

// BenchmarkDelete benchmarks Delete operations
func BenchmarkDelete(b *testing.B) {
	b.Run("DeleteSimple", func(b *testing.B) {
		// Pre-create keys
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark-delete-%d", i)
			_ = Set(testAppID, key, "value")
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark-delete-%d", i)
			_ = Delete(testAppID, key)
		}
	})

	b.Run("DeleteNested", func(b *testing.B) {
		baseKey := "benchmark-delete-nested"
		
		// Pre-create nested structure
		data := make(map[string]any)
		for i := 0; i < 100; i++ {
			data[fmt.Sprintf("field%d", i)] = fmt.Sprintf("value%d", i)
		}
		_ = Set(testAppID, baseKey, data)
		defer Delete(testAppID, baseKey)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fieldNum := i % 100
			path := fmt.Sprintf("%s/field%d", baseKey, fieldNum)
			_ = Delete(testAppID, path)
			// Re-add for next iteration
			_ = Set(testAppID, path, fmt.Sprintf("value%d", fieldNum))
		}
	})
}

// BenchmarkExists benchmarks Exists operations
func BenchmarkExists(b *testing.B) {
	// Setup test data
	key := "benchmark-exists"
	cleanup := setupTest(b, testAppID, key, map[string]any{
		"field1": "value1",
		"field2": int64(123),
		"nested": map[string]any{
			"deep": "value",
		},
	})
	defer cleanup()

	b.Run("ExistsRoot", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = Exists(testAppID, key)
		}
	})

	b.Run("ExistsNested", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = Exists(testAppID, key+"/nested/deep")
		}
	})

	b.Run("ExistsNonExistent", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = Exists(testAppID, "non-existent-key")
		}
	})
}

// BenchmarkComplexOperations benchmarks complex scenarios
func BenchmarkComplexOperations(b *testing.B) {
	b.Run("ReadModifyWrite", func(b *testing.B) {
		key := "benchmark-rmw"
		cleanup := setupTest(b, testAppID, key, map[string]any{
			"counter": int64(0),
			"data":    "initial",
		})
		defer cleanup()
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Read
			data, _ := GetMap(testAppID, key)
			
			// Modify
			if counter, ok := data["counter"].(int64); ok {
				data["counter"] = counter + 1
			}
			data["data"] = fmt.Sprintf("modified-%d", i)
			
			// Write
			_ = Set(testAppID, key, data)
		}
	})

	b.Run("DeepNesting", func(b *testing.B) {
		baseKey := "benchmark-deep"
		defer Delete(testAppID, baseKey)
		
		// Build a path with 10 levels
		path := baseKey
		for j := 0; j < 10; j++ {
			path += fmt.Sprintf("/level%d", j)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Set(testAppID, path+"/value", fmt.Sprintf("deep-value-%d", i))
		}
	})

	b.Run("LargeArray", func(b *testing.B) {
		key := "benchmark-large-array"
		
		// Create array with 1000 elements
		largeArray := make([]any, 1000)
		for i := range largeArray {
			largeArray[i] = map[string]any{
				"id":    int64(i),
				"value": fmt.Sprintf("item-%d", i),
			}
		}
		
		cleanup := setupTest(b, testAppID, key, largeArray)
		defer cleanup()
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Access random element
			idx := i % len(largeArray)
			_, _ = Get(testAppID, fmt.Sprintf("%s/%d/value", key, idx))
		}
	})
}
