package internal

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/jheddings/go-cfprefs/testutil"
)

func TestGetSetTypes(t *testing.T) {
	testCases := []struct {
		name string
		val  any
	}{
		{name: "string", val: time.Now().Format(time.RFC3339)},
		{name: "int", val: rand.Int()},
		{name: "int8", val: int8(rand.Intn(math.MaxInt8))},
		{name: "int16", val: int16(rand.Intn(math.MaxInt16))},
		{name: "int32", val: rand.Int31()},
		{name: "int64", val: rand.Int63()},
		{name: "uint", val: uint(rand.Uint32())},
		{name: "uint8", val: uint8(rand.Uint32())},
		{name: "uint16", val: uint16(rand.Uint32())},
		{name: "uint32", val: rand.Uint32()},
		{name: "uint64", val: rand.Uint64()},
		{name: "float32", val: rand.Float32()},
		{name: "float64", val: rand.Float64()},
		{name: "bool-true", val: true},
		{name: "bool-false", val: false},
		{name: "date-time", val: time.Now()},
		{name: "bytes", val: []byte("hello world")},
		{name: "array", val: []any{
			rand.Int(),
			rand.Float64(),
			false,
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
		{name: "empty-bytes", val: []byte{}},
		{name: "empty-slice", val: []any{}},
		{name: "empty-map", val: map[string]any{}},
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

			if !testutil.ValuesEqualApprox(tc.val, readVal) {
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

func TestEmptySlice(t *testing.T) {
	err := Set("com.jheddings.cfprefs.testing", "empty-slice", []any{})
	if err != nil {
		t.Fatalf("expected no error for empty slice, got '%v'", err)
	}

	readVal, err := Get("com.jheddings.cfprefs.testing", "empty-slice")
	if err != nil {
		t.Fatalf("expected no error for empty slice, got '%v'", err)
	}

	if !reflect.DeepEqual(readVal, []any{}) {
		t.Fatalf("expected empty slice, got '%v'", readVal)
	}
}
