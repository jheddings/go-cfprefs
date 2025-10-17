package internal

import (
	"slices"
	"testing"
	"time"
)

func TestBasicStr(t *testing.T) {
	dateTime := time.Now().Format(time.RFC3339)
	t.Log(dateTime)

	err := Set("com.jheddings.cfprefs.testing", "str", dateTime)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "str")

	if err != nil {
		t.Fatal(err)
	}

	if value != dateTime {
		t.Fatal("value does not match")
	}
}

func TestBasicInt(t *testing.T) {
	number := 1234567890
	t.Log(number)

	err := Set("com.jheddings.cfprefs.testing", "int", number)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "int")

	if err != nil {
		t.Fatal(err)
	}

	if value != int64(number) {
		t.Fatalf("value does not match: expected %d, got %v", number, value)
	}
}

func TestBasicBool(t *testing.T) {
	testValue := true

	err := Set("com.jheddings.cfprefs.testing", "bool", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "bool")

	if err != nil {
		t.Fatal(err)
	}

	if value != testValue {
		t.Fatalf("value does not match: expected %v, got %v", testValue, value)
	}
}

func TestBasicFloat(t *testing.T) {
	testValue := 3.14159265

	err := Set("com.jheddings.cfprefs.testing", "float", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "float")

	if err != nil {
		t.Fatal(err)
	}

	if value != testValue {
		t.Fatalf("value does not match: expected %f, got %v", testValue, value)
	}
}

func TestBasicTime(t *testing.T) {
	testValue := time.Now()

	err := Set("com.jheddings.cfprefs.testing", "time", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "time")

	if err != nil {
		t.Fatal(err)
	}

	retrievedTime, ok := value.(time.Time)
	if !ok {
		t.Fatalf("value is not time.Time: got %T", value)
	}

	// Compare times with some tolerance (within 1 second)
	diff := retrievedTime.Sub(testValue)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Fatalf("time does not match: expected %v, got %v", testValue, retrievedTime)
	}
}

func TestBasicBytes(t *testing.T) {
	testValue := []byte("hello world")

	err := Set("com.jheddings.cfprefs.testing", "bytes", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "bytes")

	if err != nil {
		t.Fatal(err)
	}

	retrievedBytes, ok := value.([]byte)
	if !ok {
		t.Fatalf("value is not []byte: got %T", value)
	}

	if string(retrievedBytes) != string(testValue) {
		t.Fatalf("value does not match: expected %s, got %s", testValue, retrievedBytes)
	}
}

func TestBasicArray(t *testing.T) {
	testValue := []any{"hello", int64(42), 3.14, true}

	err := Set("com.jheddings.cfprefs.testing", "array", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "array")

	if err != nil {
		t.Fatal(err)
	}

	retrievedArray, ok := value.([]any)
	if !ok {
		t.Fatalf("value is not []any: got %T", value)
	}

	if len(retrievedArray) != len(testValue) {
		t.Fatalf("array length does not match: expected %d, got %d", len(testValue), len(retrievedArray))
	}

	for i, v := range testValue {
		if retrievedArray[i] != v {
			t.Fatalf("array element %d does not match: expected %v, got %v", i, v, retrievedArray[i])
		}
	}
}

func TestBasicMap(t *testing.T) {
	testValue := map[string]any{
		"string": "hello",
		"number": int64(42),
		"float":  3.14,
		"bool":   true,
	}

	err := Set("com.jheddings.cfprefs.testing", "map", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "map")

	if err != nil {
		t.Fatal(err)
	}

	retrievedMap, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("value is not map[string]any: got %T", value)
	}

	if len(retrievedMap) != len(testValue) {
		t.Fatalf("map length does not match: expected %d, got %d", len(testValue), len(retrievedMap))
	}

	for k, v := range testValue {
		if retrievedMap[k] != v {
			t.Fatalf("map value for key %s does not match: expected %v, got %v", k, v, retrievedMap[k])
		}
	}
}

func TestGetKeys(t *testing.T) {
	keys := []string{"key1", "key2"}

	for _, key := range keys {
		err := Set("com.jheddings.cfprefs.testing", key, key+"-value")
		if err != nil {
			t.Fatal(err)
		}
		defer Delete("com.jheddings.cfprefs.testing", key)
	}

	prefKeys, err := GetKeys("com.jheddings.cfprefs.testing")
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range keys {
		if !slices.Contains(prefKeys, key) {
			t.Fatalf("key %s not found in keys", key)
		}
	}
}
