# internal - CFPreferences CGO Wrapper

This package provides a low-level CGO wrapper around the macOS `CFPreferences` API, bridging between Go types and CoreFoundation types.

## Overview

The internal package handles the complexity of interfacing between Go and CoreFoundation's C-based API, including:

- Memory management across the Go/C boundary
- Type conversions between native Go types and CoreFoundation types
- Error handling for CFPreferences operations

## Core APIs

### cfprefs.go

The main public API for CFPreferences operations:

- **`Get(appID, key string) (any, error)`** - Retrieves a preference value
- **`Set(appID, key string, value any) error`** - Sets a preference value
- **`Delete(appID, key string) error`** - Removes a preference value
- **`Exists(appID, key string) (bool, error)`** - Checks if a preference key exists

All operations use `CFPreferencesCopyAppValue`, `CFPreferencesSetAppValue`, and `CFPreferencesAppSynchronize` from the CoreFoundation framework.

### Type Conversions

| CoreFoundation Type     | Go Type                                       |
|-------------------------|-----------------------------------------------|
| `CFStringRef`           | `string`                                      |
| `CFBooleanRef`          | `bool`                                        |
| `CFNumberRef` (int64)   | `int`, `int8`, `int16`, `int32`, `int64`      |
| `CFNumberRef` (int64)   | `uint`, `uint8`, `uint16`, `uint32`, `uint64` |
| `CFNumberRef` (float64) | `float32`, `float64`                          |
| `CFDateRef`             | `time.Time`                                   |
| `CFDataRef`             | `[]byte`                                      |
| `CFArrayRef`            | `[]any`                                       |
| `CFDictionaryRef`       | `map[string]any`                              |

#### marshal.go

Converts Go types to CoreFoundation types (`convertGoToCFType`):

#### unmarshal.go

Converts CoreFoundation types back to Go types (`convertCFTypeToGo`):

- Uses `CFGetTypeID()` to determine the CoreFoundation type
- Recursively converts nested structures (arrays, dictionaries)
- For `CFDataRef`, attempts to deserialize as property list or JSON before falling back to raw bytes

## Memory Management

### Go to C Boundary

When passing data from Go to C:

1. **CStrings**: Created with `C.CString()` and immediately freed with `defer C.free()`
2. **CFTypes**: Created via C helper functions, released with `defer C.CFRelease()`
3. **Helper function**: `createCFStringRef()` encapsulates string creation and validation

### C to Go Boundary

When receiving data from C:

1. **CFTypes**: Always released with `defer C.CFRelease()` after use
2. **Temporary allocations**: Released before returning (e.g., in array/dictionary conversions)
3. **Nested structures**: Intermediate CF objects are released after conversion to Go types

### Key Principles

- Never return C pointers to Go code without proper cleanup
- Use `defer` for consistent cleanup even in error paths
- Validate all CF object creation for nil results
- Release temporary CF objects in nested conversions
