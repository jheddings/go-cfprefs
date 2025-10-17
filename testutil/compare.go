package testutil

import (
	"math"
	"reflect"
	"time"
)

// valuesEqualApprox compares two values with tolerance for type conversions and floating point precision
func ValuesEqualApprox(expected, actual any) bool {
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
				if !ValuesEqualApprox(expSlice[i], actSlice[i]) {
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
				if !exists || !ValuesEqualApprox(expVal, actVal) {
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
