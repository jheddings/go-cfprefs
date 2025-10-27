package cfprefs

import (
	"sort"

	"github.com/jheddings/go-cfprefs/internal"
	"github.com/theory/jsonpath"
	"github.com/theory/jsonpath/spec"
)

// Delete removes a preference value for the given key and application ID.
// This deletes the entire key from the preferences.
func Delete(appID, key string) error {
	return internal.Delete(appID, key)
}

// DeleteQ removes values using JSONPath syntax within a specific root key.
// If the query matches multiple nodes in the root key, all matched nodes
// will be deleted.
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
	data, err := internal.Get(appID, rootKey)
	if err != nil {
		return nil
	}

	// if the query is empty or "$", delete the entire root key
	if query == "" || query == "$" {
		return internal.Delete(appID, rootKey)
	}

	// parse the query
	path, err := jsonpath.Parse(query)
	if err != nil {
		return NewKeyPathError().Wrap(err).WithMsgF("invalid query: %s", query)
	}

	// select the nodes to delete
	located := path.SelectLocated(data)

	// if no nodes were found, return nil
	if len(located) == 0 {
		return nil
	}

	// when deleting array elements, delete from highest index to lowest to avoid shifting issues
	sort.Slice(located, func(i, j int) bool {
		return located[i].Path.String() > located[j].Path.String()
	})

	// apply deletions to the data structure - rebuild as we go
	for _, node := range located {
		data, err = deleteValueAtPath(data, node.Path)
		if err != nil {
			return NewInternalError().Wrap(err).WithMsgF("failed to delete at path: %s", path)
		}
	}

	// write the modified data back to the root
	return internal.Set(appID, rootKey, data)
}

// deleteValueAtPath removes a value at the given path from the data structure.
// Returns the modified data structure and an error if the path cannot be deleted.
func deleteValueAtPath(data any, norm spec.NormalizedPath) (any, error) {
	path, err := jsonpath.Parse(norm.String())
	if err != nil {
		return nil, NewKeyPathError().Wrap(err).WithMsg("failed to parse normalized path")
	}

	query := path.Query()
	segments := query.Segments()

	// if no segments (shouldn't happen after validation), return as-is
	if len(segments) == 0 {
		return data, nil
	}

	// recursively delete using the segments
	return deleteFinalSegment(data, segments)
}

// deleteFinalSegment recursively traverses the data structure and deletes the value
// at the final segment of the path, returning the modified structure.
func deleteFinalSegment(data any, segments []*spec.Segment) (any, error) {
	if len(segments) == 0 {
		return data, nil
	}

	segment := segments[0]
	selectors := segment.Selectors()
	if len(selectors) != 1 {
		return nil, NewInternalError().WithMsgF("expected single selector in segment, got %d", len(selectors))
	}

	selector := selectors[0]
	isLast := len(segments) == 1

	// check if this is an index selector
	if idx, ok := selector.(spec.Index); ok {
		arr, ok := data.([]any)
		if !ok {
			return nil, NewTypeMismatchError([]any{}, data)
		}

		index := int(idx)
		if index < 0 || index >= len(arr) {
			return nil, NewKeyPathError().WithMsgF("array index out of bounds: %d (array length: %d)", index, len(arr))
		}

		if isLast {
			// delete the element at this index
			arr = append(arr[:index], arr[index+1:]...)

		} else {
			// continue traversing the remaining segments
			modified, err := deleteFinalSegment(arr[index], segments[1:])
			if err != nil {
				return nil, NewInternalError().Wrap(err).WithMsgF("failed to delete at array index %d", index)
			}

			// update the array in place
			arr[index] = modified
		}

		return arr, nil
	}

	// check if this is a name selector
	if name, ok := selector.(spec.Name); ok {
		obj, ok := data.(map[string]any)
		if !ok {
			return nil, NewTypeMismatchError(map[string]any{}, data)
		}

		key := string(name)
		if _, exists := obj[key]; !exists {
			return nil, NewKeyNotFoundError("", key)
		}

		if isLast {
			// delete the key from the map
			delete(obj, key)

		} else {
			// continue traversing
			modified, err := deleteFinalSegment(obj[key], segments[1:])
			if err != nil {
				return nil, NewInternalError().Wrap(err).WithMsgF("failed to delete at key: %s", key)
			}

			// update the map in place
			obj[key] = modified
		}

		return obj, nil
	}

	return nil, NewInternalError().WithMsgF("unsupported selector type for deletion: %T", selector)
}
