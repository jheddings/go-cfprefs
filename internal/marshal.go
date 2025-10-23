package internal

// This file contains functions to convert Go types to CoreFoundation types.

import (
	"time"
	"unsafe"
)

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
#include <stdlib.h>

char* cfStringToC(CFStringRef str) {
    if (str == NULL) return NULL;

    CFIndex length = CFStringGetLength(str);
    CFIndex maxSize = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
    char *buffer = (char *)malloc(maxSize);

    if (buffer == NULL) {
        // memory allocation failed
        return NULL;
    }

    if (CFStringGetCString(str, buffer, maxSize, kCFStringEncodingUTF8)) {
        return buffer;
    }

    free(buffer);
    return NULL;
}

CFStringRef createCFString(const char *str) {
    return CFStringCreateWithCString(kCFAllocatorDefault, str, kCFStringEncodingUTF8);
}

CFNumberRef createCFNumberInt8(int8_t value) {
    return CFNumberCreate(kCFAllocatorDefault, kCFNumberSInt8Type, &value);
}

CFNumberRef createCFNumberInt16(int16_t value) {
    return CFNumberCreate(kCFAllocatorDefault, kCFNumberSInt16Type, &value);
}

CFNumberRef createCFNumberInt32(int32_t value) {
    return CFNumberCreate(kCFAllocatorDefault, kCFNumberSInt32Type, &value);
}

CFNumberRef createCFNumberInt64(int64_t value) {
    return CFNumberCreate(kCFAllocatorDefault, kCFNumberSInt64Type, &value);
}

CFNumberRef createCFNumberFloat32(float value) {
    return CFNumberCreate(kCFAllocatorDefault, kCFNumberFloat32Type, &value);
}

CFNumberRef createCFNumberFloat64(double value) {
    return CFNumberCreate(kCFAllocatorDefault, kCFNumberFloat64Type, &value);
}

CFBooleanRef createCFBoolean(Boolean value) {
    return value ? kCFBooleanTrue : kCFBooleanFalse;
}

CFDateRef createCFDate(CFAbsoluteTime time) {
    return CFDateCreate(kCFAllocatorDefault, time);
}

CFDataRef createCFData(const void* bytes, CFIndex length) {
    return CFDataCreate(kCFAllocatorDefault, bytes, length);
}

CFArrayRef createCFArray(const void** values, CFIndex count) {
    return CFArrayCreate(kCFAllocatorDefault, values, count, &kCFTypeArrayCallBacks);
}

CFDictionaryRef createCFDictionary(const void** keys, const void** values, CFIndex count) {
    return CFDictionaryCreate(kCFAllocatorDefault, keys, values, count, &kCFTypeDictionaryKeyCallBacks, &kCFTypeDictionaryValueCallBacks);
}
*/
import "C"

// convertGoToCFType converts a native Go type to a CFTypeRef
func convertGoToCFType(value any) (C.CFTypeRef, error) {
	if value == nil {
		return nilCFType, nil
	}

	switch v := value.(type) {
	case string:
		return convertStringToCF(v)

	case bool:
		return convertBoolToCF(v), nil

	case int:
		// int is platform-dependent (32 or 64 bit), use int64 for safety
		return convertInt64ToCF(int64(v)), nil

	case int8:
		return convertInt8ToCF(v), nil

	case int16:
		return convertInt16ToCF(v), nil

	case int32:
		return convertInt32ToCF(v), nil

	case int64:
		return convertInt64ToCF(v), nil

	case uint:
		// uint is platform-dependent, use int64 for safety
		return convertInt64ToCF(int64(v)), nil

	case uint8:
		// Use int16 to safely represent full uint8 range (0-255)
		return convertInt16ToCF(int16(v)), nil

	case uint16:
		// Use int32 to safely represent full uint16 range (0-65535)
		return convertInt32ToCF(int32(v)), nil

	case uint32:
		// Use int64 to safely represent full uint32 range (0-4294967295)
		return convertInt64ToCF(int64(v)), nil

	case uint64:
		// Use int64 (may overflow for values > 2^63-1)
		return convertInt64ToCF(int64(v)), nil

	case float32:
		return convertFloat32ToCF(v), nil

	case float64:
		return convertFloat64ToCF(v), nil

	case time.Time:
		return convertTimeToCF(v), nil

	case []byte:
		return convertBytesToCF(v), nil

	case []any:
		return convertSliceToCF(v)

	case map[string]any:
		return convertMapToCF(v)
	}

	return nilCFType, CFTypeError().WithMsgF("unsupported Go type: %T", value)
}

// convertStringToCF converts a Go string to a CFStringRef
func convertStringToCF(value string) (C.CFTypeRef, error) {
	strRef, err := createCFStringRef(value)
	if err != nil {
		return nilCFType, CFTypeError().Wrap(err).WithMsg("failed to convert string to CFString")
	}
	return C.CFTypeRef(strRef), nil
}

// convertBoolToCF converts a Go bool to a CFBooleanRef
func convertBoolToCF(value bool) C.CFTypeRef {
	var boolVal C.Boolean
	if value {
		boolVal = 1
	} else {
		boolVal = 0
	}
	boolRef := C.createCFBoolean(boolVal)
	return C.CFTypeRef(boolRef)
}

// convertInt8ToCF converts a Go int8 to a CFNumberRef
func convertInt8ToCF(value int8) C.CFTypeRef {
	numRef := C.createCFNumberInt8(C.int8_t(value))
	return C.CFTypeRef(numRef)
}

// convertInt16ToCF converts a Go int16 to a CFNumberRef
func convertInt16ToCF(value int16) C.CFTypeRef {
	numRef := C.createCFNumberInt16(C.int16_t(value))
	return C.CFTypeRef(numRef)
}

// convertInt32ToCF converts a Go int32 to a CFNumberRef
func convertInt32ToCF(value int32) C.CFTypeRef {
	numRef := C.createCFNumberInt32(C.int32_t(value))
	return C.CFTypeRef(numRef)
}

// convertInt64ToCF converts a Go int64 to a CFNumberRef
func convertInt64ToCF(value int64) C.CFTypeRef {
	numRef := C.createCFNumberInt64(C.int64_t(value))
	return C.CFTypeRef(numRef)
}

// convertFloat32ToCF converts a Go float32 to a CFNumberRef
func convertFloat32ToCF(value float32) C.CFTypeRef {
	numRef := C.createCFNumberFloat32(C.float(value))
	return C.CFTypeRef(numRef)
}

// convertFloat64ToCF converts a Go float64 to a CFNumberRef
func convertFloat64ToCF(value float64) C.CFTypeRef {
	numRef := C.createCFNumberFloat64(C.double(value))
	return C.CFTypeRef(numRef)
}

// convertTimeToCF converts a Go time.Time to a CFDateRef
func convertTimeToCF(value time.Time) C.CFTypeRef {
	unixTime := float64(value.Unix()) + float64(value.Nanosecond())/1e9
	absoluteTime := unixTime - CFAbsoluteTimeIntervalSince1970

	dateRef := C.createCFDate(C.CFAbsoluteTime(absoluteTime))
	return C.CFTypeRef(dateRef)
}

// convertBytesToCF converts a Go []byte to a CFDataRef
func convertBytesToCF(value []byte) C.CFTypeRef {
	if len(value) == 0 {
		dataRef := C.createCFData(nil, 0)
		return C.CFTypeRef(dataRef)
	}

	dataRef := C.createCFData(unsafe.Pointer(&value[0]), C.CFIndex(len(value)))
	return C.CFTypeRef(dataRef)
}

// convertSliceToCF converts a Go []any to a CFArrayRef
func convertSliceToCF(value []any) (C.CFTypeRef, error) {
	if len(value) == 0 {
		arrRef := C.createCFArray(nil, 0)
		return C.CFTypeRef(arrRef), nil
	}

	cfValues := make([]unsafe.Pointer, len(value))
	defer func() {
		for _, ptr := range cfValues {
			safeCFRelease(C.CFTypeRef(ptr))
		}
	}()

	for i, v := range value {
		cfValue, err := convertGoToCFType(v)
		if err != nil {
			return nilCFType, CFTypeError().Wrap(err).WithMsgF("failed to convert slice element %d", i)
		}
		cfValues[i] = unsafe.Pointer(cfValue)
	}

	arrRef := C.createCFArray((*unsafe.Pointer)(unsafe.Pointer(&cfValues[0])), C.CFIndex(len(value)))
	return C.CFTypeRef(arrRef), nil
}

// convertMapToCF converts a Go map[string]any to a CFDictionaryRef
func convertMapToCF(value map[string]any) (C.CFTypeRef, error) {
	if len(value) == 0 {
		dictRef := C.createCFDictionary(nil, nil, 0)
		return C.CFTypeRef(dictRef), nil
	}

	cfKeys := make([]unsafe.Pointer, 0, len(value))
	cfValues := make([]unsafe.Pointer, 0, len(value))

	defer func() {
		for _, ptr := range cfKeys {
			safeCFRelease(C.CFTypeRef(ptr))
		}
		for _, ptr := range cfValues {
			safeCFRelease(C.CFTypeRef(ptr))
		}
	}()

	for k, v := range value {
		keyRef, err := convertStringToCF(k)
		if err != nil {
			return nilCFType, CFTypeError().Wrap(err).WithMsgF("failed to convert map key '%s'", k)
		}
		cfKeys = append(cfKeys, unsafe.Pointer(keyRef))

		valueRef, err := convertGoToCFType(v)
		if err != nil {
			return nilCFType, CFTypeError().Wrap(err).WithMsgF("failed to convert value for key '%s'", k)
		}
		cfValues = append(cfValues, unsafe.Pointer(valueRef))
	}

	dictRef := C.createCFDictionary(
		(*unsafe.Pointer)(unsafe.Pointer(&cfKeys[0])),
		(*unsafe.Pointer)(unsafe.Pointer(&cfValues[0])),
		C.CFIndex(len(value)),
	)

	return C.CFTypeRef(dictRef), nil
}
