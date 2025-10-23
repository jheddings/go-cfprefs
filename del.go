package cfprefs

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/jheddings/go-cfprefs/internal"
)

// Removes a preference value for the given keypath and application ID.
// Returns an error if any intermediate segment exists but is not a dictionary.
func Delete(appID, keypath string) error {
	segments, err := splitKeypath(keypath)
	if err != nil {
		return err
	}

	// if there's only one segment, delete directly
	if len(segments) == 1 {
		return internal.Delete(appID, keypath)
	}

	rootValue, err := internal.Get(appID, segments[0])
	if err != nil {
		// root doesn't exist, nothing to delete
		return nil
	}

	// verify root is a dictionary
	rootDict, ok := rootValue.(map[string]any)
	if !ok {
		return NewKeyPathError(appID, keypath).WithMsgF("root is not a dictionary")
	}

	// traverse to the parent of the final key
	currentDict := rootDict
	for i := 1; i < len(segments)-1; i++ {
		segment := segments[i]

		// get the next dictionary
		value, exists := currentDict[segment]
		if !exists {
			// path doesn't exist, nothing to delete
			return nil
		}

		nextDict, ok := value.(map[string]any)
		if !ok {
			return NewKeyPathError(appID, keypath).WithMsgF("segment '%s' is not a dictionary", segment)
		}
		currentDict = nextDict
	}

	// delete the final key from the parent dictionary
	finalKey := segments[len(segments)-1]
	delete(currentDict, finalKey)

	// write the modified root dictionary back
	return internal.Set(appID, segments[0], rootDict)
}

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
