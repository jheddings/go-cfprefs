package internal

// This file contains functions to convert Go types to CoreFoundation types.

import (
	"fmt"
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

    if (CFStringGetCString(str, buffer, maxSize, kCFStringEncodingUTF8)) {
        return buffer;
    }

    free(buffer);
    return NULL;
}

CFStringRef createCFString(const char *str) {
    return CFStringCreateWithCString(kCFAllocatorDefault, str, kCFStringEncodingUTF8);
}

CFNumberRef createCFNumberInt64(int64_t value) {
    return CFNumberCreate(kCFAllocatorDefault, kCFNumberLongLongType, &value);
}

CFNumberRef createCFNumberFloat64(double value) {
    return CFNumberCreate(kCFAllocatorDefault, kCFNumberDoubleType, &value);
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
		return C.CFTypeRef(unsafe.Pointer(nil)), nil
	}

	switch v := value.(type) {
	case string:
		return convertStringToCF(v), nil

	case bool:
		return convertBoolToCF(v), nil

	case int:
		return convertInt64ToCF(int64(v)), nil

	case int8:
		return convertInt64ToCF(int64(v)), nil

	case int16:
		return convertInt64ToCF(int64(v)), nil

	case int32:
		return convertInt64ToCF(int64(v)), nil

	case int64:
		return convertInt64ToCF(v), nil

	case uint:
		return convertInt64ToCF(int64(v)), nil

	case uint8:
		return convertInt64ToCF(int64(v)), nil

	case uint16:
		return convertInt64ToCF(int64(v)), nil

	case uint32:
		return convertInt64ToCF(int64(v)), nil

	case uint64:
		return convertInt64ToCF(int64(v)), nil

	case float32:
		return convertFloat64ToCF(float64(v)), nil

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

	default:
		return C.CFTypeRef(unsafe.Pointer(nil)), fmt.Errorf("unsupported Go type: %T", value)
	}
}

// convertStringToCF converts a Go string to a CFStringRef
func convertStringToCF(value string) C.CFTypeRef {
	cStr := C.CString(value)
	defer C.free(unsafe.Pointer(cStr))
	strRef := C.createCFString(cStr)
	return C.CFTypeRef(strRef)
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

// convertInt64ToCF converts a Go int64 to a CFNumberRef
func convertInt64ToCF(value int64) C.CFTypeRef {
	numRef := C.createCFNumberInt64(C.int64_t(value))
	return C.CFTypeRef(numRef)
}

// convertFloat64ToCF converts a Go float64 to a CFNumberRef
func convertFloat64ToCF(value float64) C.CFTypeRef {
	numRef := C.createCFNumberFloat64(C.double(value))
	return C.CFTypeRef(numRef)
}

// convertTimeToCF converts a Go time.Time to a CFDateRef
func convertTimeToCF(value time.Time) C.CFTypeRef {
	// CFAbsoluteTime is seconds since Jan 1, 2001 00:00:00 GMT
	// Unix epoch is Jan 1, 1970 00:00:00 GMT
	// Difference is 31 years = 978307200 seconds
	// FIXME: perform the calculation instead of using a constant
	const cfAbsoluteTimeIntervalSince1970 = 978307200.0

	unixTime := float64(value.Unix()) + float64(value.Nanosecond())/1e9
	absoluteTime := unixTime - cfAbsoluteTimeIntervalSince1970

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
		// Release all the CF objects we created
		for _, ptr := range cfValues {
			if ptr != nil {
				C.CFRelease(C.CFTypeRef(ptr))
			}
		}
	}()

	for i, v := range value {
		cfValue, err := convertGoToCFType(v)
		if err != nil {
			return C.CFTypeRef(unsafe.Pointer(nil)), fmt.Errorf("failed to convert slice element %d: %w", i, err)
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
		// Release all the CF objects we created
		for _, ptr := range cfKeys {
			if ptr != nil {
				C.CFRelease(C.CFTypeRef(ptr))
			}
		}
		for _, ptr := range cfValues {
			if ptr != nil {
				C.CFRelease(C.CFTypeRef(ptr))
			}
		}
	}()

	for k, v := range value {
		keyRef := convertStringToCF(k)
		cfKeys = append(cfKeys, unsafe.Pointer(keyRef))

		valueRef, err := convertGoToCFType(v)
		if err != nil {
			return C.CFTypeRef(unsafe.Pointer(nil)), fmt.Errorf("failed to convert map value for key '%s': %w", k, err)
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
