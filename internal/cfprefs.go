package internal

// This file contains the public API operations for CoreFoundation preferences.

import (
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
#include <stdlib.h>

// Forward declarations of helper functions
extern CFStringRef createCFString(const char *str);
*/
import "C"

// Get retrieves a preference value for the given key and appID.
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

	goValue, err := convertCFTypeToGo(value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert preference value: %w", err)
	}

	return goValue, nil
}

// Set updates a preference value for the given key and appID.
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

	success := C.CFPreferencesAppSynchronize(appIDRef)
	if success == 0 {
		return fmt.Errorf("failed to synchronize preferences")
	}

	return nil
}

// Delete removes a preference value for the given key and appID.
func Delete(appID, key string) error {
	appIDRef := C.createCFString(C.CString(appID))
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef := C.createCFString(C.CString(key))
	defer C.CFRelease(C.CFTypeRef(keyRef))

	nilRef := C.CFTypeRef(unsafe.Pointer(nil))
	C.CFPreferencesSetAppValue(keyRef, nilRef, appIDRef)

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
