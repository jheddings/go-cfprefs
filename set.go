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
	arrayHandler := arraySegmentHandler{
		onLast: func(arr []any, index int) ([]any, error) {
			// handle append operation (empty index)
			if index == ArrayAppendIndex {
				return append(arr, value), nil
			}

			// validate array bounds
			if index < 0 || index >= len(arr) {
				return nil, NewKeyPathError().WithMsgF("array index out of bounds: %d (array length: %d)", index, len(arr))
			}

			// set the element at the index
			arr[index] = value
			return arr, nil
		},

		onContinue: func(arr []any, index int, element any, segments []*spec.Segment) ([]any, error) {
			// handle append operation (empty index)
			if index == ArrayAppendIndex {
				// create new element based on next segment type
				newElement := createStructureForSegment(segments)

				// recursively set value in the new element
				modified, err := walkOrSet(newElement, segments, value)
				if err != nil {
					return nil, NewInternalError().Wrap(err).WithMsg("failed to set in new array element")
				}

				// append the modified element to the array
				return append(arr, modified), nil
			}

			// validate bounds for normal indices
			if index < 0 || index >= len(arr) {
				return nil, NewKeyPathError().WithMsgF("array index out of bounds: %d (array length: %d)", index, len(arr))
			}

			// recursively set value in the existing element
			modified, err := walkOrSet(element, segments, value)
			if err != nil {
				return nil, NewInternalError().Wrap(err).WithMsgF("failed to set at array index %d", index)
			}

			// update the array with modified element
			arr[index] = modified
			return arr, nil
		},
	}

	mapHandler := mapSegmentHandler{
		onLast: func(obj map[string]any, key string) (map[string]any, error) {
			// set the value at the key
			obj[key] = value
			return obj, nil
		},

		onContinue: func(obj map[string]any, key string, child any, segments []*spec.Segment) (map[string]any, error) {
			// create child structure if it doesn't exist
			if child == nil {
				child = createStructureForSegment(segments)
			}

			// recursively set value in the child
			modified, err := walkOrSet(child, segments, value)
			if err != nil {
				return nil, NewInternalError().Wrap(err).WithMsgF("failed to set at key: %s", key)
			}

			// update the map with modified child
			obj[key] = modified
			return obj, nil
		},
	}

	return newPathWalker().WithHandler(&arrayHandler).WithHandler(&mapHandler).Walk(data, segments)
}

// createStructureForSegment creates an empty array or map based on the next segment type.
// Returns an array if the next segment is an array index, otherwise returns a map.
func createStructureForSegment(segments []*spec.Segment) any {
	if len(segments) > 0 {
		nextSelector := segments[0].Selectors()[0]
		if _, ok := nextSelector.(spec.Index); ok {
			return []any{}
		}
	}

	return make(map[string]any)
}
