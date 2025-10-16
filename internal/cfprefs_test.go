package internal

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// valuesEqualApprox compares two values with tolerance for type conversions and floating point precision
func valuesEqualApprox(expected, actual any) bool {
	// handle time.Time specially
	if expTime, ok := expected.(time.Time); ok {
		if actTime, ok := actual.(time.Time); ok {
			diff := expTime.Sub(actTime)
			return diff >= -time.Microsecond && diff <= time.Microsecond
		}
		return false
	}

	// handle int -> int64 conversion
	if expInt, ok := expected.(int); ok {
		if actInt, ok := actual.(int64); ok {
			return int64(expInt) == actInt
		}
	}

	// handle float comparison with epsilon tolerance
	if expFloat, ok := expected.(float64); ok {
		if actFloat, ok := actual.(float64); ok {
			return math.Abs(expFloat-actFloat) < 1e-10
		}
	}

	// handle slices (arrays) with recursive comparison
	if expSlice, ok := expected.([]any); ok {
		if actSlice, ok := actual.([]any); ok {
			if len(expSlice) != len(actSlice) {
				return false
			}
			for i := range expSlice {
				if !valuesEqualApprox(expSlice[i], actSlice[i]) {
					return false
				}
			}
			return true
		}
		return false
	}

	// handle maps with recursive comparison
	if expMap, ok := expected.(map[string]any); ok {
		if actMap, ok := actual.(map[string]any); ok {
			if len(expMap) != len(actMap) {
				return false
			}
			for k, expVal := range expMap {
				actVal, exists := actMap[k]
				if !exists || !valuesEqualApprox(expVal, actVal) {
					return false
				}
			}
			return true
		}
		return false
	}

	// default to reflect.DeepEqual for other types
	return reflect.DeepEqual(expected, actual)
}

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

			if !valuesEqualApprox(tc.val, readVal) {
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
