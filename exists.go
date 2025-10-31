package cfprefs

import (
	"github.com/go-openapi/jsonpointer"
	"github.com/jheddings/go-cfprefs/internal"
)

// Exists checks if a preference key exists for the given application ID.
//
// The keypath can be a simple name or include a JSON Pointer path (e.g.,
// "config/server/port") to access nested values within the preference.
//
// Example usage:
//
//	// Check if a simple key exists
//	exists, err := Exists("com.example.app", "username")
//
//	// Check if a nested value exists
//	exists, err := Exists("com.example.app", "config/server/port")
//
// Returns true if the key exists, false otherwise.
func Exists(appID, keypath string) (bool, error) {
	kp, err := parseKeypath(keypath)
	if err != nil {
		return false, NewKeyPathError().Wrap(err).WithMsgF("invalid keypath: %s", keypath)
	}

	exists, err := internal.Exists(appID, kp.Key)
	if err != nil {
		return false, NewInternalError().Wrap(err).WithMsgF("failed to check: %s", kp.Key)
	}

	// look for a quick exit
	if !exists || kp.Path == "" {
		return exists, nil
	}

	// rebuild the pointer path
	ptr, err := jsonpointer.New(kp.Path)
	if err != nil {
		return false, NewKeyPathError().Wrap(err).WithMsgF("invalid pointer: %s", kp.Path)
	}

	// get the preference value
	val, err := internal.Get(appID, kp.Key)
	if err != nil {
		return false, NewInternalError().Wrap(err).WithMsgF("failed to get value: %s", kp)
	}

	_, _, err = ptr.Get(val)
	return err == nil, nil
}
