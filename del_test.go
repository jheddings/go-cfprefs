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
