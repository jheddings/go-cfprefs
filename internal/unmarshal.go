package internal

// This file contains functions to convert CoreFoundation types to Go types.

import (
	"encoding/json"
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

// getCFNumberAsInt64WithCheck safely converts a CFNumber to int64_t with error checking
Boolean getCFNumberAsInt64(CFNumberRef numRef, int64_t *outValue) {
    if (numRef == NULL || outValue == NULL) return false;
    return CFNumberGetValue(numRef, kCFNumberLongLongType, outValue);
}

// getCFNumberAsFloat64WithCheck safely converts a CFNumber to double with error checking
Boolean getCFNumberAsFloat64(CFNumberRef numRef, double *outValue) {
    if (numRef == NULL || outValue == NULL) return false;
    return CFNumberGetValue(numRef, kCFNumberDoubleType, outValue);
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

Boolean getCFNumberAsInt8(CFNumberRef numRef, int8_t *outValue) {
    if (numRef == NULL || outValue == NULL) return false;
    return CFNumberGetValue(numRef, kCFNumberSInt8Type, outValue);
}

Boolean getCFNumberAsInt16(CFNumberRef numRef, int16_t *outValue) {
    if (numRef == NULL || outValue == NULL) return false;
    return CFNumberGetValue(numRef, kCFNumberSInt16Type, outValue);
}

Boolean getCFNumberAsInt32(CFNumberRef numRef, int32_t *outValue) {
    if (numRef == NULL || outValue == NULL) return false;
    return CFNumberGetValue(numRef, kCFNumberSInt32Type, outValue);
}

Boolean getCFNumberAsFloat32(CFNumberRef numRef, float *outValue) {
    if (numRef == NULL || outValue == NULL) return false;
    return CFNumberGetValue(numRef, kCFNumberFloat32Type, outValue);
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
	if cfValue == nilCFType {
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
	}

	return nil, CFTypeError().WithMsgF("unsupported CFType: %v", typeID)
}

// converts a CFBooleanRef to a Go bool
func convertCFBooleanToGo(boolRef C.CFBooleanRef) bool {
	return C.getCFBoolean(boolRef) != 0
}

// converts a CFStringRef to a Go string
func convertCFStringToGo(strRef C.CFStringRef) (string, error) {
	cStr := C.cfStringToC(strRef)
	if cStr == nil {
		return "", CFTypeError().WithMsg("failed to convert CFString")
	}
	defer C.free(unsafe.Pointer(cStr))
	return C.GoString(cStr), nil
}

// converts a CFNumberRef to the appropriate Go numeric type
func convertCFNumberToGo(numRef C.CFNumberRef) (any, error) {
	numberType := C.getCFNumberType(numRef)

	switch numberType {
	case C.kCFNumberSInt8Type, C.kCFNumberCharType:
		var value C.int8_t
		if C.getCFNumberAsInt8(numRef, &value) == C.false {
			return nil, CFTypeError().WithMsg("failed to convert CFNumber to int8")
		}
		return int8(value), nil

	case C.kCFNumberSInt16Type, C.kCFNumberShortType:
		var value C.int16_t
		if C.getCFNumberAsInt16(numRef, &value) == C.false {
			return nil, CFTypeError().WithMsg("failed to convert CFNumber to int16")
		}
		return int16(value), nil

	case C.kCFNumberSInt32Type, C.kCFNumberIntType:
		var value C.int32_t
		if C.getCFNumberAsInt32(numRef, &value) == C.false {
			return nil, CFTypeError().WithMsg("failed to convert CFNumber to int32")
		}
		return int32(value), nil

	case C.kCFNumberSInt64Type, C.kCFNumberLongLongType:
		var value C.int64_t
		if C.getCFNumberAsInt64(numRef, &value) == C.false {
			return nil, CFTypeError().WithMsg("failed to convert CFNumber to int64")
		}
		return int64(value), nil

	case C.kCFNumberFloat32Type, C.kCFNumberFloatType:
		var value C.float
		if C.getCFNumberAsFloat32(numRef, &value) == C.false {
			return nil, CFTypeError().WithMsg("failed to convert CFNumber to float32")
		}
		return float32(value), nil

	case C.kCFNumberFloat64Type, C.kCFNumberDoubleType, C.kCFNumberCGFloatType:
		var value C.double
		if C.getCFNumberAsFloat64(numRef, &value) == C.false {
			return nil, CFTypeError().WithMsg("failed to convert CFNumber to float64")
		}
		return float64(value), nil
	}

	return nil, CFTypeError().WithMsgF("unsupported CFNumber type: %v", numberType)
}

// converts a CFDateRef to a Go time.Time
func convertCFDateToGo(dateRef C.CFDateRef) time.Time {
	absoluteTime := float64(C.getCFDateAbsoluteTime(dateRef))
	unixTime := absoluteTime + CFAbsoluteTimeIntervalSince1970

	seconds := int64(unixTime)
	nanoseconds := int64((unixTime - float64(seconds)) * 1e9)

	return time.Unix(seconds, nanoseconds)
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
	if plist := C.tryDeserializePlist(dataRef); plist != nilCFType {
		defer C.CFRelease(C.CFTypeRef(plist))
		if value, err := convertCFTypeToGo(C.CFTypeRef(plist)); err == nil {
			return NewBinaryPlist(data, value)
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

// converts a CFArrayRef to a Go slice
func convertCFArrayToGo(arrRef C.CFArrayRef) ([]any, error) {
	count := int(C.getCFArrayCount(arrRef))
	result := make([]any, count)

	for i := range count {
		cfValue := C.getCFArrayValueAtIndex(arrRef, C.CFIndex(i))
		value, err := convertCFTypeToGo(cfValue)
		if err != nil {
			return nil, CFTypeError().Wrap(err).WithMsgF("failed to convert array element %d", i)
		}
		result[i] = value
	}

	return result, nil
}

// converts a CFArrayRef to a Go slice of strings
func convertCFArrayToGoStr(arrRef C.CFArrayRef) ([]string, error) {
	count := int(C.getCFArrayCount(arrRef))
	result := make([]string, count)

	for i := range count {
		cfValue := C.getCFArrayValueAtIndex(arrRef, C.CFIndex(i))
		value, err := convertCFStringToGo(C.CFStringRef(cfValue))
		if err != nil {
			return nil, CFTypeError().Wrap(err).WithMsgF("failed to convert array element %d", i)
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
			return nil, CFTypeError().Wrap(err).WithMsg("failed to convert dictionary key to string")
		}

		valueRef := C.getCFDictionaryValue(dictRef, keyRef)
		value, err := convertCFTypeToGo(valueRef)
		if err != nil {
			return nil, CFTypeError().Wrap(err).WithMsgF("failed to convert dictionary value for key '%s'", keyStr)
		}

		result[keyStr] = value
	}

	return result, nil
}
