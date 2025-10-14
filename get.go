package cfprefs

import (
	"fmt"

	"github.com/jheddings/go-cfprefs/internal"
)

// Get retrieves a preference value for the given key and application ID.
// The value is returned as a native Go type based on the stored CoreFoundation type.
func Get(appID, key string) (any, error) {
	return internal.Get(appID, key)
}

// GetStr retrieves a string preference value for the given key and application ID.
// Returns an error if the value is not a string.
func GetStr(appID, key string) (string, error) {
	value, err := internal.Get(appID, key)
	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("preference value is not a string: %T", value)
	}

	return strValue, nil
}

// GetInt retrieves an integer preference value for the given key and application ID.
// Returns an error if the value is not an integer.
func GetInt(appID, key string) (int64, error) {
	value, err := internal.Get(appID, key)
	if err != nil {
		return 0, err
	}

	intValue, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("preference value is not an integer: %T", value)
	}

	return intValue, nil
}

// GetBool retrieves a boolean preference value for the given key and application ID.
// Returns an error if the value is not a boolean.
func GetBool(appID, key string) (bool, error) {
	value, err := internal.Get(appID, key)
	if err != nil {
		return false, err
	}

	boolValue, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("preference value is not a boolean: %T", value)
	}

	return boolValue, nil
}

// GetFloat retrieves a float preference value for the given key and application ID.
// Returns an error if the value is not a float.
func GetFloat(appID, key string) (float64, error) {
	value, err := internal.Get(appID, key)
	if err != nil {
		return 0.0, err
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0.0, fmt.Errorf("preference value is not a float: %T", value)
	}

	return floatValue, nil
}
