package cfprefs

import (
	"fmt"
	"strings"

	"github.com/jheddings/go-cfprefs/internal"
)

// Get retrieves a preference value for the given key and application ID.
func Get(appID, keypath string) (any, error) {
	segments := strings.Split(keypath, "/")

	// start with the root value
	value, err := internal.Get(appID, segments[0])
	if err != nil {
		return nil, err
	}

	// traverse the remaining path segments
	for i := 1; i < len(segments); i++ {
		// Current value must be a map to traverse further
		dict, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("key not found: %s [%s] - segment '%s' is not a dictionary",
				keypath, appID, segments[i-1])
		}

		// Get the next value from the dictionary
		value, ok = dict[segments[i]]
		if !ok {
			return nil, fmt.Errorf("key not found: %s [%s] - segment '%s' not found in dictionary",
				keypath, appID, segments[i])
		}
	}

	return value, nil
}

// GetStr retrieves a string preference value for the given key and application ID.
func GetStr(appID, keypath string) (string, error) {
	value, err := Get(appID, keypath)
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
func GetInt(appID, keypath string) (int64, error) {
	value, err := Get(appID, keypath)
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
func GetBool(appID, keypath string) (bool, error) {
	value, err := Get(appID, keypath)
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
func GetFloat(appID, keypath string) (float64, error) {
	value, err := Get(appID, keypath)
	if err != nil {
		return 0.0, err
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0.0, fmt.Errorf("preference value is not a float: %T", value)
	}

	return floatValue, nil
}
