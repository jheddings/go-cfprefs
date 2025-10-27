package cfprefs

import (
	"github.com/jheddings/go-cfprefs/internal"
	"github.com/theory/jsonpath/spec"
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

// walkOrSet navigates to the target and sets the value, creating intermediate structures as needed.
func walkOrSet(data any, segments []*spec.Segment, value any) (any, error) {
	if len(segments) == 0 {
		return value, nil
	}

	handler := walkHandler{
		onArrayLast: func(arr []any, index int) ([]any, error) {
			// check for append operation (empty index)
			if index == -1 {
				return append(arr, value), nil
			}

			// validate array bounds
			if index < 0 || index >= len(arr) {
				return nil, NewKeyPathError().WithMsgF("array index out of bounds: %d", index)
			}

			// set the element at the index
			arr[index] = value
			return arr, nil
		},

		onArrayContinue: func(arr []any, index int, element any, remaining []*spec.Segment) ([]any, error) {
			// check for append operation (empty index)
			if index == -1 {
				// append in the middle of path - create new element based on next segment
				var newElement any
				if len(remaining) > 0 {
					nextSelector := remaining[0].Selectors()[0]
					if nextIdx, ok := nextSelector.(spec.Index); ok && int(nextIdx) == -1 {
						newElement = []any{}
					} else {
						newElement = make(map[string]any)
					}
				} else {
					newElement = make(map[string]any)
				}

				// continue setting in the new element
				modified, err := walkOrSet(newElement, remaining, value)
				if err != nil {
					return nil, err
				}

				// append the modified element to the array
				return append(arr, modified), nil
			}

			// validate bounds for normal indices
			if index < 0 || index >= len(arr) {
				return nil, NewKeyPathError().WithMsgF("array index out of bounds: %d", index)
			}

			// continue traversing the existing element
			modified, err := walkOrSet(element, remaining, value)
			if err != nil {
				return nil, err
			}

			arr[index] = modified
			return arr, nil
		},

		onMapLast: func(obj map[string]any, key string) (map[string]any, error) {
			obj[key] = value
			return obj, nil
		},

		onMapContinue: func(obj map[string]any, key string, child any, remaining []*spec.Segment) (map[string]any, error) {
			// if child is nil, create it based on the next segment
			if child == nil {
				// create new structure based on next segment
				if len(remaining) > 0 {
					nextSelector := remaining[0].Selectors()[0]
					if nextIdx, ok := nextSelector.(spec.Index); ok && int(nextIdx) == -1 {
						child = []any{}
					} else {
						child = make(map[string]any)
					}
				} else {
					child = make(map[string]any)
				}
			}

			// continue traversing with the child
			modified, err := walkOrSet(child, remaining, value)
			if err != nil {
				return nil, err
			}

			obj[key] = modified
			return obj, nil
		},
	}

	return walkPath(data, segments, handler)
}
