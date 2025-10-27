package cfprefs

import (
	"errors"

	"github.com/jheddings/go-cfprefs/internal"
	"github.com/theory/jsonpath"
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
		if errors.Is(err, internal.ErrCFLookup) {
			return false, nil
		}
		return false, err
	}

	if query == "" || query == "$" {
		return true, nil
	}

	path, err := jsonpath.Parse(query)
	if err != nil {
		return false, NewKeyPathError().Wrap(err).WithMsgF("invalid query: %s", query)
	}

	results := path.Select(rootValue)
	return len(results) > 0, nil
}
