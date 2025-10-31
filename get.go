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
// Returns the value at the specified path or an error if not found.
func Get(appID, key string) (any, error) {
	// look for a base key and optional pointer
	parts := strings.SplitN(key, "/", 2)
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
		return nil, NewKeyNotFoundError(appID, key).Wrap(err)
	}

	return result, nil
}

// GetZ retrieves a preference value for the given key and application ID and converts it to the desired type.
// Returns an error if the key doesn't exist or if the value is not of the given type.
func GetZ[T any](appID, key string) (T, error) {
	var zero T

	value, err := Get(appID, key)
	if err != nil {
		return zero, err
	}

	typedValue, ok := value.(T)
	if !ok {
		return zero, NewTypeMismatchError(zero, value).WithKey(appID, key)
	}

	return typedValue, nil
}

// GetStr retrieves a string preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a string.
func GetStr(appID, key string) (string, error) {
	value, err := Get(appID, key)

	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return "", NewTypeMismatchError(string(""), value).WithKey(appID, key)
	}

	return strValue, nil
}

// GetBool retrieves a boolean preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a boolean.
func GetBool(appID, key string) (bool, error) {
	value, err := Get(appID, key)
	if err != nil {
		return false, err
	}

	boolValue, ok := value.(bool)
	if !ok {
		return false, NewTypeMismatchError(bool(false), value).WithKey(appID, key)
	}

	return boolValue, nil
}

// GetInt retrieves an integer preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not an integer.
func GetInt(appID, key string) (int64, error) {
	value, err := Get(appID, key)
	if err != nil {
		return 0, err
	}

	intValue, ok := value.(int64)
	if !ok {
		return 0, NewTypeMismatchError(int64(0), value).WithKey(appID, key)
	}

	return intValue, nil
}

// GetFloat retrieves a float preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a float.
func GetFloat(appID, key string) (float64, error) {
	value, err := Get(appID, key)
	if err != nil {
		return 0.0, err
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0.0, NewTypeMismatchError(float64(0.0), value).WithKey(appID, key)
	}

	return floatValue, nil
}

// GetDate retrieves a time.Time preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a time.Time.
func GetDate(appID, key string) (time.Time, error) {
	value, err := Get(appID, key)
	if err != nil {
		return time.Time{}, err
	}

	dateValue, ok := value.(time.Time)
	if !ok {
		return time.Time{}, NewTypeMismatchError(time.Time{}, value).WithKey(appID, key)
	}

	return dateValue, nil
}

// GetData retrieves a []byte preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a []byte.
func GetData(appID, key string) ([]byte, error) {
	value, err := Get(appID, key)
	if err != nil {
		return nil, err
	}

	dataValue, ok := value.([]byte)
	if !ok {
		return nil, NewTypeMismatchError([]byte{}, value).WithKey(appID, key)
	}

	return dataValue, nil
}

// GetSlice retrieves a []any preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a []any.
func GetSlice(appID, key string) ([]any, error) {
	value, err := Get(appID, key)
	if err != nil {
		return nil, err
	}

	sliceValue, ok := value.([]any)
	if !ok {
		return nil, NewTypeMismatchError([]any{}, value).WithKey(appID, key)
	}

	return sliceValue, nil
}

// GetMap retrieves a map[string]any preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a map[string]any.
func GetMap(appID, key string) (map[string]any, error) {
	value, err := Get(appID, key)
	if err != nil {
		return nil, err
	}

	mapValue, ok := value.(map[string]any)
	if !ok {
		return nil, NewTypeMismatchError(map[string]any{}, value).WithKey(appID, key)
	}

	return mapValue, nil
}
