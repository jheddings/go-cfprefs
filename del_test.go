package cfprefs

import (
	"reflect"
	"testing"

	"github.com/jheddings/go-cfprefs/testutil"
)

// TestDeleteOperations tests various Delete operations
func TestDeleteOperations(t *testing.T) {
	t.Run("basic delete", func(t *testing.T) {
		key := "test-basic-delete"
		
		// Set a value
		err := Set(testAppID, key, "test value")
		testutil.AssertNoError(t, err, "set key")
		
		// Verify it exists
		assertKeyExists(t, testAppID, key, true)
		
		// Delete the key
		err = Delete(testAppID, key)
		testutil.AssertNoError(t, err, "delete key")
		
		// Verify it no longer exists
		assertKeyExists(t, testAppID, key, false)
	})

	t.Run("delete from nested structure", func(t *testing.T) {
		rootKey := "test-delete-nested"
		nestedData := map[string]any{
			"level1": map[string]any{
				"value1": "first value",
				"value2": "second value",
			},
			"level2": map[string]any{
				"nested": int64(42),
			},
		}
		
		cleanup := setupTest(t, testAppID, rootKey, nestedData)
		defer cleanup()
		
		// Verify value exists before deletion
		exists, _ := Exists(testAppID, rootKey+"/level1/value1")
		if !exists {
			t.Fatal("value1 should exist before deletion")
		}
		
		// Delete one nested value
		err := Delete(testAppID, rootKey+"/level1/value1")
		testutil.AssertNoError(t, err, "delete nested field")
		
		// Verify it was deleted
		exists, _ = Exists(testAppID, rootKey+"/level1/value1")
		if exists {
			t.Fatal("value1 should not exist after deletion")
		}
		
		// Verify sibling value still exists
		value, err := Get(testAppID, rootKey+"/level1/value2")
		testutil.AssertNoError(t, err, "get sibling value")
		if value.(string) != "second value" {
			t.Errorf("sibling value was modified: expected 'second value', got '%s'", value.(string))
		}
		
		// Verify parent dictionary still exists
		exists, _ = Exists(testAppID, rootKey+"/level1")
		if !exists {
			t.Fatal("parent level1 should still exist")
		}
		
		// Verify other branch still exists
		value, err = Get(testAppID, rootKey+"/level2/nested")
		testutil.AssertNoError(t, err, "get other branch value")
		if value.(int64) != int64(42) {
			t.Errorf("other branch was modified: expected 42, got %v", value)
		}
	})

	t.Run("delete entire root key", func(t *testing.T) {
		rootKey := "test-delete-root"
		cleanup := setupTest(t, testAppID, rootKey, "test value")
		defer cleanup()
		
		assertKeyExists(t, testAppID, rootKey, true)
		
		// Delete using empty path (deletes entire root key)
		err := Delete(testAppID, rootKey+"/")
		testutil.AssertNoError(t, err, "delete with empty path")
		
		assertKeyExists(t, testAppID, rootKey, false)
	})
}

// TestDeleteArrayOperations tests array-specific Delete operations
func TestDeleteArrayOperations(t *testing.T) {
	t.Run("delete array element", func(t *testing.T) {
		rootKey := "test-delete-array-element"
		arrayData := map[string]any{
			"items": []any{
				map[string]any{"id": int64(1), "name": "first"},
				map[string]any{"id": int64(2), "name": "second"},
				map[string]any{"id": int64(3), "name": "third"},
			},
		}
		
		cleanup := setupTest(t, testAppID, rootKey, arrayData)
		defer cleanup()
		
		// Delete the second item (index 1)
		err := Delete(testAppID, rootKey+"/items/1")
		testutil.AssertNoError(t, err, "delete array element")
		
		// Verify the modified array
		items, err := Get(testAppID, rootKey+"/items")
		testutil.AssertNoError(t, err, "get items array after deletion")
		
		expected := []any{
			map[string]any{"id": int64(1), "name": "first"},
			map[string]any{"id": int64(3), "name": "third"},
		}
		
		if !reflect.DeepEqual(items, expected) {
			t.Errorf("items should be %v, got %v", expected, items)
		}
	})

	t.Run("delete field from array element", func(t *testing.T) {
		rootKey := "test-delete-array-field"
		arrayData := map[string]any{
			"items": []any{
				map[string]any{"id": int64(1), "name": "first"},
				map[string]any{"id": int64(2), "name": "second"},
			},
		}
		
		cleanup := setupTest(t, testAppID, rootKey, arrayData)
		defer cleanup()
		
		// Delete a field from an array element
		err := Delete(testAppID, rootKey+"/items/0/name")
		testutil.AssertNoError(t, err, "delete field from array element")
		
		// Verify the name field was deleted
		exists, _ := Exists(testAppID, rootKey+"/items/0/name")
		if exists {
			t.Fatal("items[0].name should not exist after deletion")
		}
		
		// Verify the id field still exists
		value, err := GetInt(testAppID, rootKey+"/items/0/id")
		testutil.AssertNoError(t, err, "get items[0].id")
		if value != int64(1) {
			t.Errorf("items[0].id should be 1, got %v", value)
		}
		
		// Verify the second item is unchanged
		value2, err := GetStr(testAppID, rootKey+"/items/1/name")
		testutil.AssertNoError(t, err, "get items[1].name")
		if value2 != "second" {
			t.Errorf("items[1].name should be 'second', got '%s'", value2)
		}
	})
}

// TestDeleteIdempotency tests that delete operations are idempotent
func TestDeleteIdempotency(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (key string, cleanup func())
		path  string
	}{
		{
			name: "non-existent path in existing structure",
			setup: func() (string, func()) {
				key := "test-delete-idempotent-path"
				data := map[string]any{
					"user": map[string]any{
						"name": "John",
					},
				}
				cleanup := setupTest(t, testAppID, key, data)
				return key, cleanup
			},
			path: "/user/age",
		},
		{
			name: "non-existent root key",
			setup: func() (string, func()) {
				return "non-existent-root-key", func() {}
			},
			path: "/field",
		},
		{
			name: "already deleted key",
			setup: func() (string, func()) {
				key := "test-already-deleted"
				cleanup := setupTest(t, testAppID, key, "value")
				// Delete it first
				Delete(testAppID, key)
				return key, cleanup
			},
			path: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, cleanup := tt.setup()
			defer cleanup()

			// Delete operation should not error on non-existent paths
			err := Delete(testAppID, key+tt.path)
			testutil.AssertNoError(t, err, "delete should be idempotent")

			// Verify any existing data is unchanged
			if tt.name == "non-existent path in existing structure" {
				value, err := GetStr(testAppID, key+"/user/name")
				testutil.AssertNoError(t, err, "get existing field")
				if value != "John" {
					t.Errorf("existing field should be unchanged, got %v", value)
				}
			}
		})
	}
}