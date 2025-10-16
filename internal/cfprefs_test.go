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

	// handle uint -> int64 conversion (CoreFoundation doesn't have unsigned types)
	if expUint, ok := expected.(uint); ok {
		if actInt, ok := actual.(int64); ok {
			return int64(expUint) == actInt
		}
	}

	// handle uint8 -> int16 conversion (to safely represent full uint8 range)
	if expUint8, ok := expected.(uint8); ok {
		if actInt16, ok := actual.(int16); ok {
			return int16(expUint8) == actInt16
		}
	}

	// handle uint16 -> int32 conversion (to safely represent full uint16 range)
	if expUint16, ok := expected.(uint16); ok {
		if actInt32, ok := actual.(int32); ok {
			return int32(expUint16) == actInt32
		}
	}

	// handle uint32 -> int64 conversion (to safely represent full uint32 range)
	if expUint32, ok := expected.(uint32); ok {
		if actInt64, ok := actual.(int64); ok {
			return int64(expUint32) == actInt64
		}
	}

	// handle uint64 -> int64 conversion (may overflow for very large values)
	if expUint64, ok := expected.(uint64); ok {
		if actInt64, ok := actual.(int64); ok {
			return int64(expUint64) == actInt64
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
