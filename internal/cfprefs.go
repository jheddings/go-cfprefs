package internal

// TODO fail gracefully if the framework is not available (e.g. non-macOS)

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

CFStringRef createCFString(const char *str) {
    return CFStringCreateWithCString(kCFAllocatorDefault, str, kCFStringEncodingUTF8);
}

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

Boolean isCFNumberFloat(CFNumberRef numRef) {
    CFNumberType type = CFNumberGetType(numRef);
    return (type == kCFNumberFloatType ||
            type == kCFNumberDoubleType ||
            type == kCFNumberFloat32Type ||
            type == kCFNumberFloat64Type);
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

const uint8_t* getCFDataBytes(CFDataRef data) {
    return CFDataGetBytePtr(data);
}

CFIndex getCFDataLength(CFDataRef data) {
    return CFDataGetLength(data);
}

CFAbsoluteTime getCFDateAbsoluteTime(CFDateRef date) {
    return CFDateGetAbsoluteTime(date);
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

// converts a CFTypeRef to a native Go type
func convertCFTypeToGo(cfValue C.CFTypeRef) (any, error) {
	if unsafe.Pointer(cfValue) == nil {
		return nil, nil
	}

	typeID := C.CFGetTypeID(cfValue)

	switch typeID {
	case C.CFStringGetTypeID():
		return convertCFString(C.CFStringRef(cfValue))

	case C.CFNumberGetTypeID():
		return convertCFNumber(C.CFNumberRef(cfValue))

	case C.CFBooleanGetTypeID():
		return convertCFBoolean(C.CFBooleanRef(cfValue)), nil

	case C.CFArrayGetTypeID():
		return convertCFArray(C.CFArrayRef(cfValue))

	case C.CFDictionaryGetTypeID():
		return convertCFDictionary(C.CFDictionaryRef(cfValue))

	case C.CFDateGetTypeID():
		return convertCFDate(C.CFDateRef(cfValue)), nil

	case C.CFDataGetTypeID():
		return convertCFData(C.CFDataRef(cfValue)), nil

	default:
		return nil, fmt.Errorf("unsupported CFType: %v", typeID)
	}
}

// converts a CFStringRef to a Go string
func convertCFString(strRef C.CFStringRef) (string, error) {
	cStr := C.cfStringToC(strRef)
	if cStr == nil {
		return "", fmt.Errorf("failed to convert CFString")
	}
	defer C.free(unsafe.Pointer(cStr))
	return C.GoString(cStr), nil
}

// converts a CFNumberRef to either int64 or float64
func convertCFNumber(numRef C.CFNumberRef) (any, error) {
	if C.isCFNumberFloat(numRef) != 0 {
		return float64(C.getCFNumberAsFloat64(numRef)), nil
	}
	return int64(C.getCFNumberAsInt64(numRef)), nil
}

// converts a CFBooleanRef to a Go bool
func convertCFBoolean(boolRef C.CFBooleanRef) bool {
	return C.getCFBoolean(boolRef) != 0
}

// converts a CFArrayRef to a Go slice
func convertCFArray(arrRef C.CFArrayRef) ([]any, error) {
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
func convertCFDictionary(dictRef C.CFDictionaryRef) (map[string]any, error) {
	count := int(C.getCFDictionaryCount(dictRef))
	if count == 0 {
		return make(map[string]any), nil
	}

	keys := make([]unsafe.Pointer, count)
	C.getCFDictionaryKeys(dictRef, (*unsafe.Pointer)(unsafe.Pointer(&keys[0])))
	result := make(map[string]any, count)

	for i := range count {
		keyRef := C.CFStringRef(keys[i])
		keyStr, err := convertCFString(keyRef)
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
func convertCFData(dataRef C.CFDataRef) any {
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
func convertCFDate(dateRef C.CFDateRef) time.Time {
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

// converts a native Go type to a CFTypeRef
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

// converts a Go string to a CFStringRef
func convertStringToCF(value string) C.CFTypeRef {
	cStr := C.CString(value)
	defer C.free(unsafe.Pointer(cStr))
	strRef := C.createCFString(cStr)
	return C.CFTypeRef(strRef)
}

// converts a Go bool to a CFBooleanRef
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

// converts a Go int64 to a CFNumberRef
func convertInt64ToCF(value int64) C.CFTypeRef {
	numRef := C.createCFNumberInt64(C.int64_t(value))
	return C.CFTypeRef(numRef)
}

// converts a Go float64 to a CFNumberRef
func convertFloat64ToCF(value float64) C.CFTypeRef {
	numRef := C.createCFNumberFloat64(C.double(value))
	return C.CFTypeRef(numRef)
}

// converts a Go time.Time to a CFDateRef
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

// converts a Go []byte to a CFDataRef
func convertBytesToCF(value []byte) C.CFTypeRef {
	if len(value) == 0 {
		dataRef := C.createCFData(nil, 0)
		return C.CFTypeRef(dataRef)
	}

	dataRef := C.createCFData(unsafe.Pointer(&value[0]), C.CFIndex(len(value)))
	return C.CFTypeRef(dataRef)
}

// converts a Go []any to a CFArrayRef
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

// converts a Go map[string]any to a CFDictionaryRef
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

// Get a preference value for the given key, appID.
func Get(appID, key string) (any, error) {
	appIDRef := C.createCFString(C.CString(appID))
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef := C.createCFString(C.CString(key))
	defer C.CFRelease(C.CFTypeRef(keyRef))

	value := C.CFPreferencesCopyAppValue(keyRef, appIDRef)
	if unsafe.Pointer(value) == nil {
		return nil, fmt.Errorf("key not found: %s [%s]", key, appID)
	}
	defer C.CFRelease(value)

	// Convert the CFType to a native Go type
	goValue, err := convertCFTypeToGo(value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert preference value: %w", err)
	}

	return goValue, nil
}

// Set a preference value for the given key, appID.
func Set(appID, key string, value any) error {
	appIDRef := C.createCFString(C.CString(appID))
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef := C.createCFString(C.CString(key))
	defer C.CFRelease(C.CFTypeRef(keyRef))

	valueRef, err := convertGoToCFType(value)
	if err != nil {
		return fmt.Errorf("failed to convert value: %w", err)
	}
	defer func() {
		if unsafe.Pointer(valueRef) != nil {
			C.CFRelease(valueRef)
		}
	}()

	C.CFPreferencesSetAppValue(keyRef, valueRef, appIDRef)

	// Synchronize to ensure the change is written
	success := C.CFPreferencesAppSynchronize(appIDRef)
	if success == 0 {
		return fmt.Errorf("failed to synchronize preferences")
	}

	return nil
}

// Delete a preference value for the given key, appID.
func Delete(appID, key string) error {
	appIDRef := C.createCFString(C.CString(appID))
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef := C.createCFString(C.CString(key))
	defer C.CFRelease(C.CFTypeRef(keyRef))

	// Setting value to nil deletes the preference
	nilRef := C.CFTypeRef(unsafe.Pointer(nil))
	C.CFPreferencesSetAppValue(keyRef, nilRef, appIDRef)

	// Synchronize to ensure the change is written
	success := C.CFPreferencesAppSynchronize(appIDRef)
	if success == 0 {
		return fmt.Errorf("failed to synchronize preferences")
	}

	return nil
}

// Exists checks if a key exists for the given appID.
func Exists(appID, key string) (bool, error) {
	appIDRef := C.createCFString(C.CString(appID))
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef := C.createCFString(C.CString(key))
	defer C.CFRelease(C.CFTypeRef(keyRef))

	value := C.CFPreferencesCopyAppValue(keyRef, appIDRef)
	if unsafe.Pointer(value) == nil {
		return false, nil
	}
	defer C.CFRelease(value)

	return true, nil
}
