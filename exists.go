package cfprefs

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/jheddings/go-cfprefs/internal"
)

// Exists checks if a preference key exists for the given application ID.
// Returns true if the key exists, false otherwise.
func Exists(appID, key string) (bool, error) {
	return internal.Exists(appID, key)
}

// ExistsQ checks if a value exists using JSONPath syntax within a specific root key.
// The JSONPath query is applied to the value stored under the rootKey.
// Returns true if the query resolves to a valid value, false otherwise.
//
// Example usage:
//
//	// Check if a nested value exists: $.user.name
//	exists, err := ExistsQ("com.example.app", "userData", "$.user.name")
//
//	// Check if an array has items: $.items[0]
//	exists, err := ExistsQ("com.example.app", "data", "$.items[0]")
//
//	// Check if filtered array has results: $.items[?(@.active == true)]
//	exists, err := ExistsQ("com.example.app", "data", "$.items[?(@.active == true)]")
func ExistsQ(appID, rootKey, query string) (bool, error) {
	rootValue, err := internal.Get(appID, rootKey)
	if err != nil {
		return false, nil
	}

	if query == "" || query == "$" {
		return true, nil
	}

	result, err := jsonpath.Get(query, rootValue)
	if err != nil {
		return false, NewKeyPathError(appID, rootKey).WithMsgF("JSONPath query failed for path '%s'", query)
	}

	return result != nil, nil
}
