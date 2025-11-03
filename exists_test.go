package cfprefs

import (
	"testing"

	"github.com/jheddings/go-cfprefs/testutil"
)

func TestExistsBasic(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// if this fails, we need to manually cleanup
	assertKeyExists(t, appID, "exists-basic-test", false)

	cleanup := setupTest(t, appID, "exists-basic-test", "test value")
	defer cleanup()

	assertKeyExists(t, appID, "exists-basic-test", true)
	assertKeyExists(t, appID, "nonexistent-key", false)
}

func TestExistsNested(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

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

	testKey := "exists-nested-test"
	cleanup := setupTest(t, appID, testKey, nestedData)
	defer cleanup()

	t.Run("Root object exists", func(t *testing.T) {
		exists, err := Exists(appID, "exists-nested-test")
		testutil.AssertNoError(t, err, "check root exists")
		if !exists {
			t.Fatal("expected root to exist")
		}
	})

	t.Run("Empty query checks root", func(t *testing.T) {
		exists, err := Exists(appID, "exists-nested-test/")
		testutil.AssertNoError(t, err, "check empty query")
		if !exists {
			t.Fatal("expected root to exist with empty query")
		}
	})

	t.Run("Nested field exists", func(t *testing.T) {
		exists, err := Exists(appID, "exists-nested-test/user/name")
		testutil.AssertNoError(t, err, "check user.name exists")
		if !exists {
			t.Fatal("expected user.name to exist")
		}
	})

	t.Run("Deeply nested field exists", func(t *testing.T) {
		exists, err := Exists(appID, "exists-nested-test/user/address/city")
		testutil.AssertNoError(t, err, "check user.address.city exists")
		if !exists {
			t.Fatal("expected user.address.city to exist")
		}
	})

	t.Run("Array element exists", func(t *testing.T) {
		exists, err := Exists(appID, "exists-nested-test/items/0")
		testutil.AssertNoError(t, err, "check items[0] exists")
		if !exists {
			t.Fatal("expected items[0] to exist")
		}
	})

	t.Run("Array field exists", func(t *testing.T) {
		exists, err := Exists(appID, "exists-nested-test/items/1/id")
		testutil.AssertNoError(t, err, "check items[1].id exists")
		if !exists {
			t.Fatal("expected items[1].id to exist")
		}
	})

	t.Run("Non-existent field returns false", func(t *testing.T) {
		exists, err := Exists(appID, "exists-nested-test/user/nonexistent")
		testutil.AssertNoError(t, err, "check non-existent field")
		if exists {
			t.Fatal("expected user.nonexistent to not exist")
		}
	})

	t.Run("Non-existent nested path returns false", func(t *testing.T) {
		exists, err := Exists(appID, "exists-nested-test/user/address/country")
		testutil.AssertNoError(t, err, "check non-existent nested path")
		if exists {
			t.Fatal("expected user.address.country to not exist")
		}
	})

	t.Run("Out of bounds array index returns false", func(t *testing.T) {
		exists, err := Exists(appID, "exists-nested-test/items/999")
		testutil.AssertNoError(t, err, "check out of bounds index")
		if exists {
			t.Fatal("expected items[999] to not exist")
		}
	})

	t.Run("Non-existent root key returns false", func(t *testing.T) {
		exists, err := Exists(appID, "nonexistent-root/user/name")
		testutil.AssertNoError(t, err, "check non-existent root key")
		if exists {
			t.Fatal("expected non-existent root key to return false")
		}
	})
}
