package internal

// This file contains the public API operations for CoreFoundation preferences.

// TODO: fail gracefully if not running on macOS

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
	appIDRef, err := createCFStringRef(appID)
	if err != nil {
		return nil, CFRefError().Wrap(err).WithMsg("failed to create CFString for appID")
	}
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef, err := createCFStringRef(key)
	if err != nil {
		return nil, CFRefError().Wrap(err).WithMsg("failed to create CFString for key")
	}
	defer C.CFRelease(C.CFTypeRef(keyRef))

	// https://developer.apple.com/documentation/corefoundation/cfpreferencescopyappvalue(_:_:)
	value := C.CFPreferencesCopyAppValue(keyRef, appIDRef)
	if value == nilCFType {
		return nil, CFLookupError().WithMsgF("key '%s' not found in app '%s'", key, appID)
	}
	defer C.CFRelease(value)

	goValue, err := convertCFTypeToGo(value)
	if err != nil {
		return nil, CFTypeError().Wrap(err).WithMsg("failed to convert preference value")
	}

	return goValue, nil
}

// GetKeys retrieves all keys for the given appID.
func GetKeys(appID string) ([]string, error) {
	appIDRef, err := createCFStringRef(appID)
	if err != nil {
		return nil, CFRefError().Wrap(err).WithMsg("failed to create CFString for appID")
	}
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	// https://developer.apple.com/documentation/corefoundation/cfpreferencescopykeylist(_:_:_:)
	keysCF := C.CFPreferencesCopyKeyList(appIDRef, C.kCFPreferencesCurrentUser, C.kCFPreferencesAnyHost)
	if keysCF == nilCFArray {
		return nil, NewCFError("get keys", nil).WithMsgF("app not found: %s", appID)
	}
	defer C.CFRelease(C.CFTypeRef(keysCF))

	keys, err := convertCFArrayToGoStr(keysCF)
	if err != nil {
		return nil, CFTypeError().Wrap(err).WithMsg("failed to convert keys array")
	}

	return keys, nil
}

// Set updates a preference value for the given key and appID.
func Set(appID, key string, value any) error {
	appIDRef, err := createCFStringRef(appID)
	if err != nil {
		return CFRefError().Wrap(err).WithMsg("failed to create CFString for appID")
	}
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef, err := createCFStringRef(key)
	if err != nil {
		return CFRefError().Wrap(err).WithMsg("failed to create CFString for key")
	}
	defer C.CFRelease(C.CFTypeRef(keyRef))

	valueRef, err := convertGoToCFType(value)
	if err != nil {
		return CFTypeError().Wrap(err).WithMsg("failed to convert value")
	}
	defer func() {
		safeCFRelease(valueRef)
	}()

	// https://developer.apple.com/documentation/corefoundation/cfpreferencessetappvalue(_:_:_:)
	C.CFPreferencesSetAppValue(keyRef, valueRef, appIDRef)

	// https://developer.apple.com/documentation/corefoundation/cfpreferencesappsynchronize(_:)
	success := C.CFPreferencesAppSynchronize(appIDRef)
	if success == 0 {
		return CFSyncError().WithMsg("failed to synchronize preferences")
	}

	return nil
}

// Delete removes a preference value for the given key and appID.
func Delete(appID, key string) error {
	appIDRef, err := createCFStringRef(appID)
	if err != nil {
		return CFRefError().Wrap(err).WithMsg("failed to create CFString for appID")
	}
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef, err := createCFStringRef(key)
	if err != nil {
		return CFRefError().Wrap(err).WithMsg("failed to create CFString for key")
	}
	defer C.CFRelease(C.CFTypeRef(keyRef))

	// https://developer.apple.com/documentation/corefoundation/cfpreferencessetappvalue(_:_:_:)
	C.CFPreferencesSetAppValue(keyRef, nilCFType, appIDRef)

	// https://developer.apple.com/documentation/corefoundation/cfpreferencesappsynchronize(_:)
	success := C.CFPreferencesAppSynchronize(appIDRef)
	if success == 0 {
		return CFSyncError().WithMsg("failed to synchronize preferences")
	}

	return nil
}

// Exists checks if a key exists for the given appID.
func Exists(appID, key string) (bool, error) {
	appIDRef, err := createCFStringRef(appID)
	if err != nil {
		return false, CFRefError().Wrap(err).WithMsg("failed to create CFString for appID")
	}
	defer C.CFRelease(C.CFTypeRef(appIDRef))

	keyRef, err := createCFStringRef(key)
	if err != nil {
		return false, CFRefError().Wrap(err).WithMsg("failed to create CFString for key")
	}
	defer C.CFRelease(C.CFTypeRef(keyRef))

	// https://developer.apple.com/documentation/corefoundation/cfpreferencescopyappvalue(_:_:)
	value := C.CFPreferencesCopyAppValue(keyRef, appIDRef)
	if value == nilCFType {
		return false, nil
	}
	defer C.CFRelease(value)

	return true, nil
}
