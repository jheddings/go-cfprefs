package internal

// This file contains helper functions for common CGO operations.

import (
	"fmt"
	"time"
	"unsafe"
)

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
#include <stdlib.h>

extern CFStringRef createCFString(const char *str);
*/
import "C"

// nilCFType is a nil CFTypeRef.
var nilCFType = C.CFTypeRef(unsafe.Pointer(nil))

// nilCFString is a nil CFStringRef.
var nilCFString = C.CFStringRef(unsafe.Pointer(nil))

// nilCFArray is a nil CFArrayRef.
var nilCFArray = C.CFArrayRef(unsafe.Pointer(nil))

// CFAbsoluteTimeIntervalSince1970 is the offset between CoreFoundation's epoch
// (Jan 1, 2001 00:00:00 GMT) and Unix epoch (Jan 1, 1970 00:00:00 GMT).
// This is calculated as the number of seconds between these two dates.
var CFAbsoluteTimeIntervalSince1970 = calculateCFEpochOffset()

// createCFStringRef creates a CFStringRef from a Go string.
// It handles C memory allocation and returns an error if creation fails.
// The caller is responsible for releasing the returned CFStringRef with CFRelease.
func createCFStringRef(str string) (C.CFStringRef, error) {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	strRef := C.createCFString(cStr)
	if strRef == nilCFString {
		return nilCFString, fmt.Errorf("failed to create CFString from: %s", str)
	}

	return strRef, nil
}

// calculateCFEpochOffset computes the offset between CF epoch and Unix epoch.
func calculateCFEpochOffset() float64 {
	// CoreFoundation epoch: Jan 1, 2001 00:00:00 GMT
	cfEpoch := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)

	// Unix epoch: Jan 1, 1970 00:00:00 GMT
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	// return the difference in seconds
	return cfEpoch.Sub(unixEpoch).Seconds()
}
