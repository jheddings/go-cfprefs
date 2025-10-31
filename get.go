package cfprefs

import (
	"strings"
	"time"

	"github.com/go-openapi/jsonpointer"
	"github.com/jheddings/go-cfprefs/internal"
)

// GetKeys retrieves all keys for the given appID.
// Returns an error if the appID is not found.
func GetKeys(appID string) ([]string, error) {
	return internal.GetKeys(appID)
}

// Get retrieves a preference value for the given key and application ID.
//
// The keypath can be a simple name or include a JSON Pointer path (e.g.,
// "config/server/port") to access nested values within the preference.
//
// Example usage:
//
//	// Get a simple value
//	value, err := Get("com.example.app", "username")
//
//	// Get a nested value
//	value, err := Get("com.example.app", "config/server/port")
//
// Returns the value at the specified path or an error if not found.
func Get(appID, keypath string) (any, error) {
	// look for a base key and optional pointer
	parts := strings.SplitN(keypath, "/", 2)
	pref := parts[0]

	val, err := internal.Get(appID, pref)
	if err != nil {
		return nil, NewKeyNotFoundError(appID, pref).Wrap(err)
	}

	// if there is no pointer, just return the value
	if len(parts) == 1 {
		return val, nil
	}

	path := "/" + parts[1]
	ptr, err := jsonpointer.New(path)
	if err != nil {
		return nil, NewKeyPathError().Wrap(err).WithMsgF("invalid pointer: %s", path)
	}

	result, _, err := ptr.Get(val)
	if err != nil {
		return nil, NewKeyNotFoundError(appID, keypath).Wrap(err)
	}

	return result, nil
}

// GetZ retrieves a preference value for the given key and application ID and
// converts it to the desired type.
//
// The keypath can be a simple name or include a JSON Pointer path (e.g.,
// "config/server/port") to access nested values within the preference.
//
// Example usage:
//
//	// Get a simple value
//	value, err := GetZ[string]("com.example.app", "username")
//
//	// Get a nested value
//	value, err := GetZ[int64]("com.example.app", "config/server/port")
//
// Returns an error if the key doesn't exist or if the value is not of the given type.
func GetZ[T any](appID, keypath string) (T, error) {
	var zero T

	value, err := Get(appID, keypath)
	if err != nil {
		return zero, err
	}

	typedValue, ok := value.(T)
	if !ok {
		return zero, NewTypeMismatchError(zero, value).WithKey(appID, keypath)
	}

	return typedValue, nil
}

// GetStr retrieves a string preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a string.
func GetStr(appID, keypath string) (string, error) {
	value, err := Get(appID, keypath)

	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return "", NewTypeMismatchError(string(""), value).WithKey(appID, keypath)
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
		return false, NewTypeMismatchError(bool(false), value).WithKey(appID, keypath)
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
		return 0, NewTypeMismatchError(int64(0), value).WithKey(appID, keypath)
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
		return 0.0, NewTypeMismatchError(float64(0.0), value).WithKey(appID, keypath)
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
		return time.Time{}, NewTypeMismatchError(time.Time{}, value).WithKey(appID, keypath)
	}

	return dateValue, nil
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
		return nil, NewTypeMismatchError([]byte{}, value).WithKey(appID, keypath)
	}

	return dataValue, nil
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
		return nil, NewTypeMismatchError([]any{}, value).WithKey(appID, keypath)
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
		return nil, NewTypeMismatchError(map[string]any{}, value).WithKey(appID, keypath)
	}

	return mapValue, nil
}
