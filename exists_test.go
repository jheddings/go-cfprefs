package cfprefs

import (
	"testing"

	"github.com/jheddings/go-cfprefs/testutil"
)

// TestExistsOperations tests various Exists operations
func TestExistsOperations(t *testing.T) {
	t.Run("basic exists check", func(t *testing.T) {
		key := "test-exists-basic"
		
		// Verify key doesn't exist initially
		assertKeyExists(t, testAppID, key, false)
		
		// Create the key
		cleanup := setupTest(t, testAppID, key, "test value")
		defer cleanup()
		
		// Verify it now exists
		assertKeyExists(t, testAppID, key, true)
	})

	t.Run("non-existent key", func(t *testing.T) {
		assertKeyExists(t, testAppID, "nonexistent-key", false)
	})
}

// TestExistsNested tests Exists operations on nested structures
func TestExistsNested(t *testing.T) {
	nestedData := map[string]any{
		"user": map[string]any{
			"name": "John Doe",
			"age":  int64(30),
			"address": map[string]any{
				"city":  "Anytown",
				"state": "CA",
			},
		},
		"items": []any{
			map[string]any{"id": int64(1), "active": true},
			map[string]any{"id": int64(2), "active": false},
			map[string]any{"id": int64(3), "active": true},
		},
	}

	rootKey := "test-exists-nested"
	cleanup := setupTest(t, testAppID, rootKey, nestedData)
	defer cleanup()

	t.Run("root object exists", func(t *testing.T) {
		exists, err := Exists(testAppID, rootKey)
		testutil.AssertNoError(t, err, "check root exists")
		if !exists {
			t.Fatal("expected root to exist")
		}
	})

	t.Run("empty query checks root", func(t *testing.T) {
		exists, err := Exists(testAppID, rootKey+"/")
		testutil.AssertNoError(t, err, "check empty query")
		if !exists {
			t.Fatal("expected root to exist with empty query")
		}
	})

	t.Run("nested field exists", func(t *testing.T) {
		exists, err := Exists(testAppID, rootKey+"/user/name")
		testutil.AssertNoError(t, err, "check user.name exists")
		if !exists {
			t.Fatal("expected user.name to exist")
		}
	})

	t.Run("deeply nested field exists", func(t *testing.T) {
		exists, err := Exists(testAppID, rootKey+"/user/address/city")
		testutil.AssertNoError(t, err, "check user.address.city exists")
		if !exists {
			t.Fatal("expected user.address.city to exist")
		}
	})

	t.Run("array element exists", func(t *testing.T) {
		exists, err := Exists(testAppID, rootKey+"/items/0")
		testutil.AssertNoError(t, err, "check items[0] exists")
		if !exists {
			t.Fatal("expected items[0] to exist")
		}
	})

	t.Run("array field exists", func(t *testing.T) {
		exists, err := Exists(testAppID, rootKey+"/items/1/id")
		testutil.AssertNoError(t, err, "check items[1].id exists")
		if !exists {
			t.Fatal("expected items[1].id to exist")
		}
	})

	t.Run("non-existent field returns false", func(t *testing.T) {
		exists, err := Exists(testAppID, rootKey+"/user/nonexistent")
		testutil.AssertNoError(t, err, "check non-existent field")
		if exists {
			t.Fatal("expected user.nonexistent to not exist")
		}
	})

	t.Run("non-existent nested path returns false", func(t *testing.T) {
		exists, err := Exists(testAppID, rootKey+"/user/address/country")
		testutil.AssertNoError(t, err, "check non-existent nested path")
		if exists {
			t.Fatal("expected user.address.country to not exist")
		}
	})

	t.Run("out of bounds array index returns false", func(t *testing.T) {
		exists, err := Exists(testAppID, rootKey+"/items/999")
		testutil.AssertNoError(t, err, "check out of bounds index")
		if exists {
			t.Fatal("expected items[999] to not exist")
		}
	})

	t.Run("non-existent root key returns false", func(t *testing.T) {
		exists, err := Exists(testAppID, "nonexistent-root/user/name")
		testutil.AssertNoError(t, err, "check non-existent root key")
		if exists {
			t.Fatal("expected non-existent root key to return false")
		}
	})
}