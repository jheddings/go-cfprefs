package cfprefs

import (
	"time"

	"github.com/jheddings/go-cfprefs/internal"
)

// Get retrieves a preference value for the given keypath and application ID.
// Returns the value at the specified path or an error if not found.
func Get(appID, keypath string) (any, error) {
	segments, err := splitKeypath(keypath)
	if err != nil {
		return nil, err
	}

	// start with the root value
	value, err := internal.Get(appID, segments[0])
	if err != nil {
		return nil, err
	}

	// traverse the remaining path segments
	for i := 1; i < len(segments); i++ {
		// current value must be a map to traverse further
		dict, ok := value.(map[string]any)
		if !ok {
			return nil, KeyNotFoundError(appID, keypath).WithMsgF("segment '%s' is not a dictionary", segments[i-1])
		}

		// get the next value from the dictionary
		value, ok = dict[segments[i]]
		if !ok {
			return nil, KeyNotFoundError(appID, keypath).WithMsgF("segment '%s' not found in dictionary", segments[i])
		}
	}

	return value, nil
}

// GetKeys retrieves all keys for the given appID.
// Returns an error if the appID is not found.
func GetKeys(appID string) ([]string, error) {
	return internal.GetKeys(appID)
}

// GetStr retrieves a string preference value for the given keypath and application ID.
// Returns an error if the key doesn't exist or if the value is not a string.
func GetStr(appID, keypath string) (string, error) {
	value, err := Get(appID, keypath)
	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return "", TypeMismatchError(appID, keypath, string(""), value)
	}

	return strValue, nil
}

// GetBool retrieves a boolean preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a boolean.
func GetBool(appID, keypath string) (bool, error) {
	value, err := Get(appID, keypath)
	if err != nil {
		return false, err
	}

	boolValue, ok := value.(bool)
	if !ok {
		return false, TypeMismatchError(appID, keypath, bool(false), value)
	}

	return boolValue, nil
}

// GetInt retrieves an integer preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not an integer.
func GetInt(appID, keypath string) (int64, error) {
	value, err := Get(appID, keypath)
	if err != nil {
		return 0, err
	}

	intValue, ok := value.(int64)
	if !ok {
		return 0, TypeMismatchError(appID, keypath, int64(0), value)
	}

	return intValue, nil
}

// GetFloat retrieves a float preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a float.
func GetFloat(appID, keypath string) (float64, error) {
	value, err := Get(appID, keypath)
	if err != nil {
		return 0.0, err
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0.0, TypeMismatchError(appID, keypath, float64(0.0), value)
	}

	return floatValue, nil
}

// GetDate retrieves a time.Time preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a time.Time.
func GetDate(appID, keypath string) (time.Time, error) {
	value, err := Get(appID, keypath)
	if err != nil {
		return time.Time{}, err
	}

	dateValue, ok := value.(time.Time)
	if !ok {
		return time.Time{}, TypeMismatchError(appID, keypath, time.Time{}, value)
	}

	return dateValue, nil
}

// GetSlice retrieves a []any preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a []any.
func GetSlice(appID, keypath string) ([]any, error) {
	value, err := Get(appID, keypath)
	if err != nil {
		return nil, err
	}

	sliceValue, ok := value.([]any)
	if !ok {
		return nil, TypeMismatchError(appID, keypath, []any{}, value)
	}

	return sliceValue, nil
}

// GetMap retrieves a map[string]any preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a map[string]any.
func GetMap(appID, keypath string) (map[string]any, error) {
	value, err := Get(appID, keypath)
	if err != nil {
		return nil, err
	}

	mapValue, ok := value.(map[string]any)
	if !ok {
		return nil, TypeMismatchError(appID, keypath, map[string]any{}, value)
	}

	return mapValue, nil
}

// GetData retrieves a []byte preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a []byte.
func GetData(appID, keypath string) ([]byte, error) {
	value, err := Get(appID, keypath)
	if err != nil {
		return nil, err
	}

	dataValue, ok := value.([]byte)
	if !ok {
		return nil, TypeMismatchError(appID, keypath, []byte{}, value)
	}

	return dataValue, nil
}
