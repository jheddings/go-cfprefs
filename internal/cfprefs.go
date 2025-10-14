package internal

// TODO fail gracefully if the framework is not available (e.g. non-macOS)

import (
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>

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
*/
import "C"

// Get a preference value for the given key, appID.
func Get(appID, key string) (string, error) {
	appIDRef := C.createCFString(C.CString(appID))
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef := C.createCFString(C.CString(key))
	defer C.CFRelease(C.CFTypeRef(keyRef))

	value := C.CFPreferencesCopyAppValue(keyRef, appIDRef)
	if unsafe.Pointer(value) == nil {
		return "", fmt.Errorf("key not found: %s [%s]", key, appID)
	}
	defer C.CFRelease(value)

	// Check if the value is a string
	if C.CFGetTypeID(value) == C.CFStringGetTypeID() {
		strValue := C.CFStringRef(value)
		cStr := C.cfStringToC(strValue)
		if cStr == nil {
			return "", fmt.Errorf("failed to convert CFString to C string")
		}
		defer C.free(unsafe.Pointer(cStr))
		return C.GoString(cStr), nil
	}

	return "", fmt.Errorf("preference value is not a string")
}

// Set a preference value for the given key, appID.
func Set(appID, key, value string) error {
	appIDRef := C.createCFString(C.CString(appID))
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef := C.createCFString(C.CString(key))
	defer C.CFRelease(C.CFTypeRef(keyRef))

	valueRef := C.createCFString(C.CString(value))
	defer C.CFRelease(C.CFTypeRef(valueRef))

	C.CFPreferencesSetAppValue(keyRef, C.CFTypeRef(valueRef), appIDRef)

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
