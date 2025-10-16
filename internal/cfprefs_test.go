package internal

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestGetSetTypes(t *testing.T) {
	testCases := []struct {
		name string
		val  any
	}{
		{name: "string", val: time.Now().Format(time.RFC3339)},
		{name: "int", val: rand.Int()},
		{name: "float", val: rand.Float64()},
		{name: "time", val: time.Now()},
		{name: "bytes", val: []byte("hello world")},
		{name: "array", val: []any{
			rand.Int(),
			rand.Float64(),
			true,
			time.Now(),
		}},
		{name: "map", val: map[string]any{
			"string": "hello",
			"number": rand.Int(),
			"float":  rand.Float64(),
			"bool":   true,
			"time":   time.Now(),
			"bytes":  []byte("hello world"),
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Set("com.jheddings.cfprefs.testing", tc.name, tc.val)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err := Delete("com.jheddings.cfprefs.testing", tc.name)
				if err != nil {
					t.Fatal(err)
				}
			}()

			exists, err := Exists("com.jheddings.cfprefs.testing", tc.name)
			if err != nil {
				t.Fatal(err)
			}
			if !exists {
				t.Fatal("expected true for existing key, got false")
			}

			readVal, err := Get("com.jheddings.cfprefs.testing", tc.name)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(readVal, tc.val) {
				t.Fatalf("expected %v [%T], got %v [%T]", tc.val, tc.val, readVal, readVal)
			}
		})
	}
}

func TestMissingGet(t *testing.T) {
	_, err := Get("com.jheddings.cfprefs.testing", "this-key-should-not-exist")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestMissingDelete(t *testing.T) {
	err := Delete("com.jheddings.cfprefs.testing", "this-key-will-not-exist")
	if err != nil {
		t.Fatalf("expected no error for missing key, got '%v'", err)
	}
}

func TestMissingExists(t *testing.T) {
	exists, err := Exists("com.jheddings.cfprefs.testing", "this-key-will-not-exist")
	if err != nil {
		t.Fatalf("expected no error for missing key, got '%v'", err)
	}
	if exists {
		t.Fatal("expected false for missing key, got true")
	}
}
