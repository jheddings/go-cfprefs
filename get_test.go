package cfprefs

import (
	"math/rand/v2"
	"reflect"
	"testing"
	"time"
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

	if !reflect.DeepEqual(value, testValue) {
		t.Fatalf("map does not match: expected %v, got %v", testValue, value)
	}
}

func TestGetStr(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testValue := "hello"
	err := Set(appID, "str-test", testValue)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := Delete(appID, "str-test")
		if err != nil {
			t.Fatal(err)
		}
	}()

	exists, err := Exists(appID, "str-test")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected true for existing key, got false")
	}

	if _, err = GetInt(appID, "str-test"); err == nil {
		t.Fatal("expected error for non-int value, got nil")
	}

	value, err := GetStr(appID, "str-test")
	if err != nil {
		t.Fatal(err)
	}
	if value != testValue {
		t.Fatalf("expected %s, got %s", testValue, value)
	}
}

func TestGetInt(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testValue := rand.Int64()
	err := Set(appID, "int-test", testValue)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := Delete(appID, "int-test")
		if err != nil {
			t.Fatal(err)
		}
	}()

	exists, err := Exists(appID, "int-test")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected true for existing key, got false")
	}

	if _, err = GetBool(appID, "int-test"); err == nil {
		t.Fatal("expected error for non-bool value, got nil")
	}

	value, err := GetInt(appID, "int-test")
	if err != nil {
		t.Fatal(err)
	}
	if value != testValue {
		t.Fatalf("expected %d, got %d", testValue, value)
	}
}

func TestGetFloat(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testValue := rand.Float64()
	err := Set(appID, "float-test", testValue)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := Delete(appID, "float-test")
		if err != nil {
			t.Fatal(err)
		}
	}()

	exists, err := Exists(appID, "float-test")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected true for existing key, got false")
	}

	if _, err = GetData(appID, "float-test"); err == nil {
		t.Fatal("expected error for non-data value, got nil")
	}

	value, err := GetFloat(appID, "float-test")
	if err != nil {
		t.Fatal(err)
	}
	if value != testValue {
		t.Fatalf("expected %f, got %f", testValue, value)
	}
}

func TestGetBool(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testValue := true
	err := Set(appID, "bool-test", testValue)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := Delete(appID, "bool-test")
		if err != nil {
			t.Fatal(err)
		}
	}()

	exists, err := Exists(appID, "bool-test")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected true for existing key, got false")
	}

	if _, err = GetDate(appID, "bool-test"); err == nil {
		t.Fatal("expected error for non-date value, got nil")
	}

	value, err := GetBool(appID, "bool-test")
	if err != nil {
		t.Fatal(err)
	}
	if value != testValue {
		t.Fatalf("expected %t, got %t", testValue, value)
	}
}

func TestGetDate(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testValue := time.Now()
	err := Set(appID, "date-test", testValue)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := Delete(appID, "date-test")
		if err != nil {
			t.Fatal(err)
		}
	}()

	exists, err := Exists(appID, "date-test")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected true for existing key, got false")
	}

	if _, err = GetMap(appID, "date-test"); err == nil {
		t.Fatal("expected error for non-slice value, got nil")
	}

	value, err := GetDate(appID, "date-test")
	if err != nil {
		t.Fatal(err)
	}
	if value != testValue {
		t.Fatalf("expected %s, got %s", testValue.Format(time.RFC3339), value.Format(time.RFC3339))
	}
}

func TestGetSlice(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testValue := []any{rand.Int(), rand.Float64(), true, time.Now()}
	err := Set(appID, "slice-test", testValue)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := Delete(appID, "slice-test")
		if err != nil {
			t.Fatal(err)
		}
	}()

	exists, err := Exists(appID, "slice-test")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected true for existing key, got false")
	}

	if _, err = GetStr(appID, "slice-test"); err == nil {
		t.Fatal("expected error for non-map value, got nil")
	}

	value, err := GetSlice(appID, "slice-test")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, testValue) {
		t.Fatalf("expected %v, got %v", testValue, value)
	}
}

func TestGetData(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testValue := []byte("hello world")
	err := Set(appID, "data-test", testValue)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := Delete(appID, "data-test")
		if err != nil {
			t.Fatal(err)
		}
	}()

	exists, err := Exists(appID, "data-test")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected true for existing key, got false")
	}

	if _, err = GetFloat(appID, "data-test"); err == nil {
		t.Fatal("expected error for non-float value, got nil")
	}

	value, err := GetData(appID, "data-test")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, testValue) {
		t.Fatalf("expected %v, got %v", testValue, value)
	}
}

func TestGetMap(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing"

	testValue := map[string]any{
		"string": "hello",
		"number": rand.Int(),
		"float":  rand.Float64(),
		"bool":   false,
		"time":   time.Now(),
	}
	err := Set(appID, "map-test", testValue)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := Delete(appID, "map-test")
		if err != nil {
			t.Fatal(err)
		}
	}()

	exists, err := Exists(appID, "map-test")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected true for existing key, got false")
	}

	if _, err = GetSlice(appID, "map-test"); err == nil {
		t.Fatal("expected error for non-slice value, got nil")
	}

	value, err := GetMap(appID, "map-test")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(value, testValue) {
		t.Fatalf("expected %v, got %v", testValue, value)
	}
}
