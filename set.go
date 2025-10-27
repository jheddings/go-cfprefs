package cfprefs

import (
	"fmt"

	"github.com/jheddings/go-cfprefs/internal"
)

// Set writes a preference value for the given key and application ID.
// This replaces the entire value at the specified key.
func Set(appID, key string, value any) error {
	return internal.Set(appID, key, value)
}

// SetQ sets a value using JSONPath syntax within a specific root key.
// The JSONPath query is applied to the value stored under the rootKey.
// Missing field segments will create new maps. Array indices must be valid
// or use [] to append to an existing array.
//
// Example usage:
//
//	// Set a nested value: $.user.name
//	err := SetQ("com.example.app", "userData", "$.user.name", "John Doe")
//
//	// Set array element: $.items[0]
//	err := SetQ("com.example.app", "data", "$.items[0]", item)
//
//	// Append to array: $.items[]
//	err := SetQ("com.example.app", "data", "$.items[]", newItem)
func SetQ(appID, rootKey, query string, value any) error {
	// if the query is empty or "$", replace the entire root key
	if query == "" || query == "$" {
		return internal.Set(appID, rootKey, value)
	}

	// get or create the root value
	var rootValue any
	exists, err := internal.Exists(appID, rootKey)
	if err != nil {
		return err
	}

	if exists {
		rootValue, err = internal.Get(appID, rootKey)
		if err != nil {
			return err
		}
	} else {
		// create a new map as the root value
		rootValue = make(map[string]any)
	}

	// set the value at the specified path
	modified, err := setValueAtPath(rootValue, query, value)
	if err != nil {
		return NewKeyPathError().Wrap(err).WithMsgF("failed to set at path '%s'", query)
	}

	// write the modified root value back
	return internal.Set(appID, rootKey, modified)
}

// setValueAtPath sets a value at the given JSONPath in the data structure.
// Returns the modified structure or an error if the path cannot be set.
func setValueAtPath(data any, path string, value any) (any, error) {
	segments, err := parseJSONPath(path)
	if err != nil {
		return nil, err
	}

	// if no segments (shouldn't happen after validation), return the value
	if len(segments) == 0 {
		return value, nil
	}

	return walkOrSet(data, segments, value)
}

// walkOrSet navigates to the target and sets the value, creating
// intermediate structures as needed.
func walkOrSet(data any, segments []pathSegment, value any) (any, error) {
	if len(segments) == 0 {
		return value, nil
	}

	segment := segments[0]
	isLast := len(segments) == 1

	if segment.isArrayIdx {
		// handle array index
		arr, ok := data.([]any)
		if !ok {
			return nil, fmt.Errorf("expected array but got %T", data)
		}

		// check for append operation (empty index)
		if segment.index == -1 {
			if isLast {
				// append to array
				return append(arr, value), nil
			}
			// append in the middle of path - create new element
			var newElement any
			if len(segments) > 1 && segments[1].isArrayIdx {
				newElement = []any{}
			} else {
				newElement = make(map[string]any)
			}

			// continue setting in the new element
			modified, err := walkOrSet(newElement, segments[1:], value)
			if err != nil {
				return nil, err
			}

			// append the modified element to the array
			return append(arr, modified), nil
		}

		// validate array bounds
		if segment.index < 0 || segment.index >= len(arr) {
			return nil, fmt.Errorf("array index out of bounds: %d", segment.index)
		}

		if isLast {
			// set the element at the index
			arr[segment.index] = value
			return arr, nil
		}

		// continue traversing
		modified, err := walkOrSet(arr[segment.index], segments[1:], value)
		if err != nil {
			return nil, err
		}
		arr[segment.index] = modified
		return arr, nil
	} else {
		// handle field access
		obj, ok := data.(map[string]any)
		if !ok {
			// if not a map, we can't set a field on it
			return nil, fmt.Errorf("expected object but got %T", data)
		}

		if isLast {
			// set the field
			obj[segment.field] = value
			return obj, nil
		}

		// get or create the child
		child, exists := obj[segment.field]
		if !exists {
			// create new structure based on next segment
			if len(segments) > 1 && segments[1].isArrayIdx {
				child = []any{}
			} else {
				child = make(map[string]any)
			}
		}

		modified, err := walkOrSet(child, segments[1:], value)
		if err != nil {
			return nil, err
		}
		obj[segment.field] = modified
		return obj, nil
	}
}
