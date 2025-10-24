package cfprefs

import (
	"time"

	"github.com/jheddings/go-cfprefs/internal"
	"github.com/theory/jsonpath"
)

// GetKeys retrieves all keys for the given appID.
// Returns an error if the appID is not found.
func GetKeys(appID string) ([]string, error) {
	return internal.GetKeys(appID)
}

// Get retrieves a preference value for the given key and application ID.
// Returns the value at the specified path or an error if not found.
func Get(appID, key string) (any, error) {
	value, err := internal.Get(appID, key)
	if err != nil {
		return nil, NewKeyNotFoundError(appID, key).Wrap(err)
	}

	return value, nil
}

// GetQ retrieves a value using JSONPath syntax within a specific root key.
// The JSONPath query is applied to the value stored under the rootKey.
// Returns the result of the JSONPath query or an error if not found.
//
// Example usage:
//
//	// Get a nested value: $.user.name
//	value, err := GetQ("com.example.app", "userData", "$.user.name")
//
//	// Get all items in an array: $.items[*]
//	items, err := GetQ("com.example.app", "data", "$.items[*]")
//
//	// Filter array items: $.items[?(@.active == true)]
//	activeItems, err := GetQ("com.example.app", "data", "$.items[?(@.active == true)]")
func GetQ(appID, rootKey, query string) (any, error) {
	rootValue, err := internal.Get(appID, rootKey)
	if err != nil {
		return nil, NewKeyNotFoundError(appID, rootKey).Wrap(err)
	}

	if query == "" || query == "$" {
		return rootValue, nil
	}

	path, err := jsonpath.Parse(query)
	if err != nil {
		return nil, NewKeyPathError().Wrap(err).WithMsgF("invalid query: %s", query)
	}

	results := path.Select(rootValue)
	if len(results) == 0 {
		return nil, NewKeyNotFoundError(appID, rootKey).WithMsgF("no results for query: %s", query)
	}

	return results[0], nil
}

// GetStr retrieves a string preference value for the given key and application ID.
// Returns an error if the key doesn't exist or if the value is not a string.
func GetStr(appID, key string) (string, error) {
	return GetStrQ(appID, key, "$")
}

// GetStrQ retrieves a string preference value for the given JSONPath query using an application ID and root key.
// Returns an error if the query is invalid or if the value is not a string.
func GetStrQ(appID, key, query string) (string, error) {
	value, err := GetQ(appID, key, query)

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
	return GetBoolQ(appID, key, "$")
}

// GetBoolQ retrieves a boolean preference value for the given JSONPath query using an application ID and root key.
// Returns an error if the query is invalid or if the value is not a boolean.
func GetBoolQ(appID, key, query string) (bool, error) {
	value, err := GetQ(appID, key, query)
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
	return GetIntQ(appID, key, "$")
}

// GetIntQ retrieves an integer preference value for the given JSONPath query using an application ID and root key.
// Returns an error if the query is invalid or if the value is not an integer.
func GetIntQ(appID, key, query string) (int64, error) {
	value, err := GetQ(appID, key, query)
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
	return GetFloatQ(appID, key, "$")
}

// GetFloatQ retrieves a float preference value for the given JSONPath query using an application ID and root key.
// Returns an error if the query is invalid or if the value is not a float.
func GetFloatQ(appID, key, query string) (float64, error) {
	value, err := GetQ(appID, key, query)
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
	return GetDateQ(appID, key, "$")
}

// GetDateQ retrieves a time.Time preference value for the given JSONPath query using an application ID and root key.
// Returns an error if the query is invalid or if the value is not a time.Time.
func GetDateQ(appID, key, query string) (time.Time, error) {
	value, err := GetQ(appID, key, query)
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
	return GetDataQ(appID, key, "$")
}

// GetDataQ retrieves a []byte preference value for the given JSONPath query using an application ID and root key.
// Returns an error if the query is invalid or if the value is not a []byte.
func GetDataQ(appID, key, query string) ([]byte, error) {
	value, err := GetQ(appID, key, query)
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
	return GetSliceQ(appID, key, "$")
}

// GetSliceQ retrieves a []any preference value for the given JSONPath query using an application ID and root key.
// Returns an error if the query is invalid or if the value is not a []any.
func GetSliceQ(appID, key, query string) ([]any, error) {
	value, err := GetQ(appID, key, query)
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
	return GetMapQ(appID, key, "$")
}

// GetMapQ retrieves a map[string]any preference value for the given JSONPath query using an application ID and root key.
// Returns an error if the query is invalid or if the value is not a map[string]any.
func GetMapQ(appID, key, query string) (map[string]any, error) {
	value, err := GetQ(appID, key, query)
	if err != nil {
		return nil, err
	}

	mapValue, ok := value.(map[string]any)
	if !ok {
		return nil, NewTypeMismatchError(map[string]any{}, value).WithKey(appID, key)
	}

	return mapValue, nil
}
