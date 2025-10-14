package cfprefs

import (
	"fmt"

	"github.com/jheddings/go-cfprefs/internal"
)

// retrieves a value for the given key and appID.
func Get(appID, key string) (any, error) {
	return internal.Get(appID, key)
}

// retrieves a string value for the given key and appID.
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

// retrieves an integer value for the given key and appID.
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

// retrieves a boolean value for the given key and appID.
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

// retrieves a float value for the given key and appID.
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
