package cfprefs

import (
	"fmt"

	"github.com/jheddings/go-cfprefs/internal"
)

// Delete removes a preference value for the given key and application ID.
// This deletes the entire key from the preferences.
func Delete(appID, key string) error {
	return internal.Delete(appID, key)
}

// DeleteQ removes a value using JSONPath syntax within a specific root key.
// The JSONPath query is applied to the value stored under the rootKey.
// Currently supports simple queries that resolve to a single item.
//
// Example usage:
//
//	// Delete a nested field: $.user.name
//	err := DeleteQ("com.example.app", "userData", "$.user.name")
//
//	// Delete an array element: $.items[0]
//	err := DeleteQ("com.example.app", "data", "$.items[0]")
//
//	// Delete the entire root (empty query or "$")
//	err := DeleteQ("com.example.app", "data", "$")
func DeleteQ(appID, rootKey, query string) error {
	rootValue, err := internal.Get(appID, rootKey)
	if err != nil {
		return nil
	}

	// if the query is empty or "$", delete the entire root key
	if query == "" || query == "$" {
		return internal.Delete(appID, rootKey)
	}

	// verify the path exists before attempting to delete
	exists, err := ExistsQ(appID, rootKey, query)
	if err != nil {
		return NewKeyPathError(appID, rootKey).Wrap(err).WithMsg("JSONPath query failed")
	}
	if !exists {
		return nil
	}

	// proceed to delete the specified path
	modified, err := deleteAtPath(rootValue, query)
	if err != nil {
		return NewKeyPathError(appID, rootKey).Wrap(err).WithMsgF("failed to delete at path '%s'", query)
	}

	// write the modified root dictionary back
	return internal.Set(appID, rootKey, modified)
}

// deleteAtPath removes a value at the given JSONPath from the data structure.
// Returns the modified structure or an error if the path cannot be deleted.
func deleteAtPath(data any, path string) (any, error) {
	segments, err := parseJSONPath(path)
	if err != nil {
		return nil, err
	}

	// if no segments (shouldn't happen after validation), return as-is
	if len(segments) == 0 {
		return data, nil
	}

	return deleteSegments(data, segments)
}

// deleteSegments navigates to the target and deletes it.
func deleteSegments(data any, segments []pathSegment) (any, error) {
	if len(segments) == 0 {
		return data, nil
	}

	segment := segments[0]
	isLast := len(segments) == 1

	if segment.isArrayIdx {
		// if the segment is an array index, delete the element at the index
		arr, ok := data.([]any)
		if !ok {
			return nil, fmt.Errorf("expected array but got %T", data)
		}

		if segment.index < 0 || segment.index >= len(arr) {
			return nil, fmt.Errorf("array index out of bounds: %d", segment.index)
		}

		if isLast {
			// if this is the last segment, delete the element at the index
			return append(arr[:segment.index], arr[segment.index+1:]...), nil
		}

		// delete the element at the index and return the modified array
		modified, err := deleteSegments(arr[segment.index], segments[1:])
		if err != nil {
			return nil, err
		}
		arr[segment.index] = modified
		return arr, nil
	} else {
		// if the segment is a field, delete the field
		obj, ok := data.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected object but got %T", data)
		}

		if isLast {
			// if this is the last segment, delete the field
			delete(obj, segment.field)
			return obj, nil
		}

		// delete the field and return the modified object
		child, exists := obj[segment.field]
		if !exists {
			return nil, fmt.Errorf("field not found: %s", segment.field)
		}

		modified, err := deleteSegments(child, segments[1:])
		if err != nil {
			return nil, err
		}
		obj[segment.field] = modified
		return obj, nil
	}
}
