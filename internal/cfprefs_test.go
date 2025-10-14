package internal

import (
	"testing"
	"time"
)

const testingDoCleanup = false

func TestBasicStr(t *testing.T) {
	dateTime := time.Now().Format(time.RFC3339)
	t.Log(dateTime)

	err := Set("com.jheddings.cfprefs.testing", "str-test", dateTime)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "str-test")

	if err != nil {
		t.Fatal(err)
	}

	if value != dateTime {
		t.Fatal("value does not match")
	}

	if testingDoCleanup {
		err = Delete("com.jheddings.cfprefs.testing", "str-test")

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBasicInt(t *testing.T) {
	number := 1234567890
	t.Log(number)

	err := Set("com.jheddings.cfprefs.testing", "int-test", number)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "int-test")

	if err != nil {
		t.Fatal(err)
	}

	if value != int64(number) {
		t.Fatalf("value does not match: expected %d, got %v", number, value)
	}

	if testingDoCleanup {
		err = Delete("com.jheddings.cfprefs.testing", "int-test")

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBasicBool(t *testing.T) {
	testValue := true

	err := Set("com.jheddings.cfprefs.testing", "bool-test", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "bool-test")

	if err != nil {
		t.Fatal(err)
	}

	if value != testValue {
		t.Fatalf("value does not match: expected %v, got %v", testValue, value)
	}

	if testingDoCleanup {
		err = Delete("com.jheddings.cfprefs.testing", "bool-test")

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBasicFloat(t *testing.T) {
	testValue := 3.14159265

	err := Set("com.jheddings.cfprefs.testing", "float-test", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "float-test")

	if err != nil {
		t.Fatal(err)
	}

	if value != testValue {
		t.Fatalf("value does not match: expected %f, got %v", testValue, value)
	}

	if testingDoCleanup {
		err = Delete("com.jheddings.cfprefs.testing", "float-test")

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBasicTime(t *testing.T) {
	testValue := time.Now()

	err := Set("com.jheddings.cfprefs.testing", "time-test", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "time-test")

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

	if testingDoCleanup {
		err = Delete("com.jheddings.cfprefs.testing", "time-test")

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBasicBytes(t *testing.T) {
	testValue := []byte("hello world")

	err := Set("com.jheddings.cfprefs.testing", "bytes-test", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "bytes-test")

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

	if testingDoCleanup {
		err = Delete("com.jheddings.cfprefs.testing", "bytes-test")

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBasicArray(t *testing.T) {
	testValue := []any{"hello", int64(42), 3.14, true}

	err := Set("com.jheddings.cfprefs.testing", "array-test", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "array-test")

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

	if testingDoCleanup {
		err = Delete("com.jheddings.cfprefs.testing", "array-test")

		if err != nil {
			t.Fatal(err)
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

	err := Set("com.jheddings.cfprefs.testing", "map-test", testValue)

	if err != nil {
		t.Fatal(err)
	}

	value, err := Get("com.jheddings.cfprefs.testing", "map-test")

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

	if testingDoCleanup {
		err = Delete("com.jheddings.cfprefs.testing", "map-test")

		if err != nil {
			t.Fatal(err)
		}
	}
}
