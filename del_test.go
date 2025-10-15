package cfprefs

import "testing"

func TestDeleteBasic(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"
	key := "del-test"

	err := Set(appID, key, "test")
	if err != nil {
		t.Fatalf("failed to set key: %v", err)
	}

	exists, err := Exists(appID, key)
	if err != nil {
		t.Fatalf("failed to check if key exists: %v", err)
	}
	if !exists {
		t.Fatalf("key does not exist")
	}

	err = Delete(appID, key)
	if err != nil {
		t.Fatalf("failed to delete key: %v", err)
	}

	exists, err = Exists(appID, key)
	if err != nil {
		t.Fatalf("failed to check if key exists: %v", err)
	}
	if exists {
		t.Fatalf("key still exists")
	}
}

func TestDeleteKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// create a nested structure with multiple keys
	err := Set(appID, "delete-test/level1/value1", "first value")
	if err != nil {
		t.Fatalf("failed to set first value: %v", err)
	}

	err = Set(appID, "delete-test/level1/value2", "second value")
	if err != nil {
		t.Fatalf("failed to set second value: %v", err)
	}

	err = Set(appID, "delete-test/level2/nested", int64(42))
	if err != nil {
		t.Fatalf("failed to set nested value: %v", err)
	}

	// verify all values exist
	exists, err := Exists(appID, "delete-test/level1/value1")
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Fatal("value1 should exist")
	}

	exists, err = Exists(appID, "delete-test/level1/value2")
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Fatal("value2 should exist")
	}

	// delete one nested value
	err = Delete(appID, "delete-test/level1/value1")
	if err != nil {
		t.Fatalf("failed to delete nested key: %v", err)
	}

	// verify it was deleted
	exists, err = Exists(appID, "delete-test/level1/value1")
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if exists {
		t.Fatal("value1 should not exist after deletion")
	}

	// verify sibling value still exists
	exists, err = Exists(appID, "delete-test/level1/value2")
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Fatal("value2 should still exist")
	}

	value, err := Get(appID, "delete-test/level1/value2")
	if err != nil {
		t.Fatalf("failed to get sibling value: %v", err)
	}
	if value.(string) != "second value" {
		t.Fatalf("sibling value was modified: expected 'second value', got '%s'", value.(string))
	}

	// verify parent dictionary still exists
	exists, err = Exists(appID, "delete-test/level1")
	if err != nil {
		t.Fatalf("failed to check parent existence: %v", err)
	}
	if !exists {
		t.Fatal("parent dictionary should still exist")
	}

	// verify other branch still exists
	exists, err = Exists(appID, "delete-test/level2/nested")
	if err != nil {
		t.Fatalf("failed to check other branch: %v", err)
	}
	if !exists {
		t.Fatal("other branch should still exist")
	}

	// clean up
	err = Delete(appID, "delete-test")
	if err != nil {
		t.Fatalf("failed to clean up: %v", err)
	}
}

func TestExistsKeypath(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	// set up a nested structure
	err := Set(appID, "exists-test/level1/level2/value", "nested value")
	if err != nil {
		t.Fatalf("failed to set nested value: %v", err)
	}

	// test that full path exists
	exists, err := Exists(appID, "exists-test/level1/level2/value")
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Fatal("full keypath should exist")
	}

	// test that intermediate paths exist
	exists, err = Exists(appID, "exists-test")
	if err != nil {
		t.Fatalf("failed to check root existence: %v", err)
	}
	if !exists {
		t.Fatal("root should exist")
	}

	exists, err = Exists(appID, "exists-test/level1")
	if err != nil {
		t.Fatalf("failed to check intermediate existence: %v", err)
	}
	if !exists {
		t.Fatal("intermediate path should exist")
	}

	// test that non-existent paths return false
	exists, err = Exists(appID, "exists-test/nonexistent")
	if err != nil {
		t.Fatalf("failed to check non-existent key: %v", err)
	}
	if exists {
		t.Fatal("non-existent key should not exist")
	}

	exists, err = Exists(appID, "exists-test/level1/wrong/path")
	if err != nil {
		t.Fatalf("failed to check non-existent path: %v", err)
	}
	if exists {
		t.Fatal("non-existent path should not exist")
	}

	// clean up
	err = Delete(appID, "exists-test")
	if err != nil {
		t.Fatalf("failed to clean up: %v", err)
	}
}
