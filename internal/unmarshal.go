package internal

// This file contains functions to convert CoreFoundation types to Go types.

import (
	"encoding/json"
	"fmt"
	"time"
	"unsafe"
)

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
#include <stdlib.h>

extern char* cfStringToC(CFStringRef str);

Boolean getCFBoolean(CFBooleanRef boolRef) {
    return CFBooleanGetValue(boolRef);
}

int64_t getCFNumberAsInt64(CFNumberRef numRef) {
    int64_t value = 0;
    CFNumberGetValue(numRef, kCFNumberLongLongType, &value);
    return value;
}

double getCFNumberAsFloat64(CFNumberRef numRef) {
    double value = 0;
    CFNumberGetValue(numRef, kCFNumberDoubleType, &value);
    return value;
}

CFNumberType getCFNumberType(CFNumberRef numRef) {
    return CFNumberGetType(numRef);
}

Boolean isCFNumberFloat(CFNumberRef numRef) {
    CFNumberType type = CFNumberGetType(numRef);
    return (type == kCFNumberFloatType ||
            type == kCFNumberDoubleType ||
            type == kCFNumberFloat32Type ||
            type == kCFNumberFloat64Type);
}

int8_t getCFNumberAsInt8(CFNumberRef numRef) {
    int8_t value = 0;
    CFNumberGetValue(numRef, kCFNumberSInt8Type, &value);
    return value;
}

int16_t getCFNumberAsInt16(CFNumberRef numRef) {
    int16_t value = 0;
    CFNumberGetValue(numRef, kCFNumberSInt16Type, &value);
    return value;
}

int32_t getCFNumberAsInt32(CFNumberRef numRef) {
    int32_t value = 0;
    CFNumberGetValue(numRef, kCFNumberSInt32Type, &value);
    return value;
}

float getCFNumberAsFloat32(CFNumberRef numRef) {
    float value = 0;
    CFNumberGetValue(numRef, kCFNumberFloat32Type, &value);
    return value;
}

CFIndex getCFArrayCount(CFArrayRef arr) {
    return CFArrayGetCount(arr);
}

CFTypeRef getCFArrayValueAtIndex(CFArrayRef arr, CFIndex idx) {
    return CFArrayGetValueAtIndex(arr, idx);
}

CFIndex getCFDictionaryCount(CFDictionaryRef dict) {
    return CFDictionaryGetCount(dict);
}

void getCFDictionaryKeys(CFDictionaryRef dict, const void **keys) {
    CFDictionaryGetKeysAndValues(dict, keys, NULL);
}

CFTypeRef getCFDictionaryValue(CFDictionaryRef dict, CFStringRef key) {
    return CFDictionaryGetValue(dict, key);
}

CFAbsoluteTime getCFDateAbsoluteTime(CFDateRef date) {
    return CFDateGetAbsoluteTime(date);
}

const uint8_t* getCFDataBytes(CFDataRef data) {
    return CFDataGetBytePtr(data);
}

CFIndex getCFDataLength(CFDataRef data) {
    return CFDataGetLength(data);
}

CFPropertyListRef tryDeserializePlist(CFDataRef data) {
    CFErrorRef error = NULL;

    CFPropertyListRef plist = CFPropertyListCreateWithData(
        kCFAllocatorDefault,
        data,
        kCFPropertyListImmutable,
        NULL,
        &error
    );

    if (error) {
        CFRelease(error);
        return NULL;
    }

    return plist;
}
*/
import "C"

// converts a CFTypeRef to a native Go type
func convertCFTypeToGo(cfValue C.CFTypeRef) (any, error) {
	if unsafe.Pointer(cfValue) == nil {
		return nil, nil
	}

	typeID := C.CFGetTypeID(cfValue)

	switch typeID {
	case C.CFStringGetTypeID():
		return convertCFStringToGo(C.CFStringRef(cfValue))

	case C.CFNumberGetTypeID():
		return convertCFNumberToGo(C.CFNumberRef(cfValue))

	case C.CFBooleanGetTypeID():
		return convertCFBooleanToGo(C.CFBooleanRef(cfValue)), nil

	case C.CFArrayGetTypeID():
		return convertCFArrayToGo(C.CFArrayRef(cfValue))

	case C.CFDictionaryGetTypeID():
		return convertCFDictionaryToGo(C.CFDictionaryRef(cfValue))

	case C.CFDateGetTypeID():
		return convertCFDateToGo(C.CFDateRef(cfValue)), nil

	case C.CFDataGetTypeID():
		return convertCFDataToGo(C.CFDataRef(cfValue)), nil

	default:
		return nil, fmt.Errorf("unsupported CFType: %v", typeID)
	}
}

// converts a CFStringRef to a Go string
func convertCFStringToGo(strRef C.CFStringRef) (string, error) {
	cStr := C.cfStringToC(strRef)
	if cStr == nil {
		return "", fmt.Errorf("failed to convert CFString")
	}
	defer C.free(unsafe.Pointer(cStr))
	return C.GoString(cStr), nil
}

// converts a CFNumberRef to the appropriate Go numeric type
func convertCFNumberToGo(numRef C.CFNumberRef) (any, error) {
	numberType := C.getCFNumberType(numRef)

	switch numberType {
	case C.kCFNumberSInt8Type, C.kCFNumberCharType:
		return int8(C.getCFNumberAsInt8(numRef)), nil

	case C.kCFNumberSInt16Type, C.kCFNumberShortType:
		return int16(C.getCFNumberAsInt16(numRef)), nil

	case C.kCFNumberSInt32Type, C.kCFNumberIntType:
		return int32(C.getCFNumberAsInt32(numRef)), nil

	case C.kCFNumberSInt64Type, C.kCFNumberLongLongType:
		return int64(C.getCFNumberAsInt64(numRef)), nil

	case C.kCFNumberFloat32Type, C.kCFNumberFloatType:
		return float32(C.getCFNumberAsFloat32(numRef)), nil

	case C.kCFNumberFloat64Type, C.kCFNumberDoubleType, C.kCFNumberCGFloatType:
		return float64(C.getCFNumberAsFloat64(numRef)), nil

	default:
		if C.isCFNumberFloat(numRef) != 0 {
			return float64(C.getCFNumberAsFloat64(numRef)), nil
		}
		return int64(C.getCFNumberAsInt64(numRef)), nil
	}
}

// converts a CFBooleanRef to a Go bool
func convertCFBooleanToGo(boolRef C.CFBooleanRef) bool {
	return C.getCFBoolean(boolRef) != 0
}

// converts a CFArrayRef to a Go slice
func convertCFArrayToGo(arrRef C.CFArrayRef) ([]any, error) {
	count := int(C.getCFArrayCount(arrRef))
	result := make([]any, count)

	for i := range count {
		cfValue := C.getCFArrayValueAtIndex(arrRef, C.CFIndex(i))
		value, err := convertCFTypeToGo(cfValue)
		if err != nil {
			return nil, fmt.Errorf("failed to convert array element %d: %w", i, err)
		}
		result[i] = value
	}

	return result, nil
}

// converts a CFDictionaryRef to a Go map
func convertCFDictionaryToGo(dictRef C.CFDictionaryRef) (map[string]any, error) {
	count := int(C.getCFDictionaryCount(dictRef))
	if count == 0 {
		return make(map[string]any), nil
	}

	keys := make([]unsafe.Pointer, count)
	C.getCFDictionaryKeys(dictRef, (*unsafe.Pointer)(unsafe.Pointer(&keys[0])))
	result := make(map[string]any, count)

	for i := range count {
		keyRef := C.CFStringRef(keys[i])
		keyStr, err := convertCFStringToGo(keyRef)
		if err != nil {
			return nil, fmt.Errorf("failed to convert dictionary key: %w", err)
		}

		valueRef := C.getCFDictionaryValue(dictRef, keyRef)
		value, err := convertCFTypeToGo(valueRef)
		if err != nil {
			return nil, fmt.Errorf("failed to convert dictionary value for key '%s': %w", keyStr, err)
		}

		result[keyStr] = value
	}

	return result, nil
}

// converts a complex CFDataRef to a Go value
func convertCFDataToGo(dataRef C.CFDataRef) any {
	length := int(C.getCFDataLength(dataRef))
	if length == 0 {
		return []byte{}
	}

	bytes := C.getCFDataBytes(dataRef)
	data := C.GoBytes(unsafe.Pointer(bytes), C.int(length))

	// first, try to deserialize as property list
	if plist := C.tryDeserializePlist(dataRef); unsafe.Pointer(plist) != nil {
		defer C.CFRelease(C.CFTypeRef(plist))
		if value, err := convertCFTypeToGo(C.CFTypeRef(plist)); err == nil {
			return value
		}
	}

	// next, try to deserialize as JSON
	var jsonValue any
	if err := json.Unmarshal(data, &jsonValue); err == nil {
		return jsonValue
	}

	// neither worked... return as raw bytes
	return data
}

// converts a CFDateRef to a Go time.Time
func convertCFDateToGo(dateRef C.CFDateRef) time.Time {
	// CFAbsoluteTime is seconds since Jan 1, 2001 00:00:00 GMT
	// Unix epoch is Jan 1, 1970 00:00:00 GMT
	// Difference is 31 years = 978307200 seconds
	// FIXME: perform the calculation instead of using a constant
	const cfAbsoluteTimeIntervalSince1970 = 978307200.0

	absoluteTime := float64(C.getCFDateAbsoluteTime(dateRef))
	unixTime := absoluteTime + cfAbsoluteTimeIntervalSince1970

	seconds := int64(unixTime)
	nanoseconds := int64((unixTime - float64(seconds)) * 1e9)

	return time.Unix(seconds, nanoseconds)
}
