package cfprefs

import (
	"strconv"

	"github.com/go-openapi/jsonpointer"
	"github.com/jheddings/go-cfprefs/internal"
)

const (
	// ArrayAppendOp is a special operator used to indicate an append operation.
	// When used in a keypath (e.g., "array-test/items/~]"), it signals that a
	// new element should be appended to end of the array.
	ArrayAppendOp = "~]"

	// ArrayPrependOp is a special operator used to indicate a prepend operation.
	// When used in a keypath (e.g., "array-test/items/~["), it signals that a
	// new element should be prepended to the beginning of the array.
	ArrayPrependOp = "~["
)

// Set writes a preference value for the given key and application ID.
//
// The keypath can be a simple name or include a JSON Pointer path (e.g.,
// "config/server/port") to access nested values within the preference.
//
// Example usage:
//
//	// Set a simple value
//	err := Set("com.example.app", "username", "John Doe")
//
//	// Set a nested value
//	err := Set("com.example.app", "config/server/port", 8080)
func Set(appID, keypath string, value any) error {
	kp, err := parseKeypath(keypath)
	if err != nil {
		return NewKeyPathError().Wrap(err).WithMsgF("invalid keypath: %s", keypath)
	}

	// if there is no pointer, just set the value
	if kp.IsRoot() {
		return internal.Set(appID, kp.Key, value)
	}

	// get or create the root value
	var root any
	exists, err := internal.Exists(appID, kp.Key)
	if err != nil {
		return NewInternalError().Wrap(err).WithMsgF("failed to check: %s", kp.Key)
	}

	if exists {
		root, err = internal.Get(appID, kp.Key)
		if err != nil {
			return NewInternalError().Wrap(err).WithMsgF("failed to get: %s", kp.Key)
		}
	} else {
		root = make(map[string]any)
	}

	ptr, err := jsonpointer.New(kp.Path)
	if err != nil {
		return NewKeyPathError().Wrap(err).WithMsgF("invalid path: %s", kp.Path)
	}

	// set the value at the specified path
	modified, err := setValueAtPath(root, ptr.DecodedTokens(), value)
	if err != nil {
		return NewInternalError().Wrap(err).WithMsgF("failed to set: %s", kp.Path)
	}

	// write the modified root value back
	return internal.Set(appID, kp.Key, modified)
}

// setValueAtPath uses a pointer walker to set a value at the specified path.
func setValueAtPath(root any, tokens []string, value any) (any, error) {
	var walker *pointerWalker

	handler := pathTokenHandler{
		onArrayIndex: func(arr []any, index int, remaining []string) (any, error) {
			// if this is the last token, set the value at the index
			if len(remaining) == 0 {
				arr[index] = value
				return arr, nil
			}

			// continue to walk the path
			data, err := walker.walk(arr[index], remaining)
			if err != nil {
				return nil, err
			}

			// update the array with any modifications
			arr[index] = data
			return arr, nil
		},
		onArrayAppend: func(arr []any, remaining []string) (any, error) {
			// if this is the last token, append the value to the array
			if len(remaining) == 0 {
				return append(arr, value), nil
			}

			// construct the remaining path elements
			new := createStructureFor(remaining[0])
			data, err := walker.walk(new, remaining)
			if err != nil {
				return nil, err
			}

			// update the array with any modifications
			return append(arr, data), nil
		},
		onObjectKey: func(obj map[string]any, key string, remaining []string) (any, error) {
			// if this is the last token, set the value at the key
			if len(remaining) == 0 {
				obj[key] = value
				return obj, nil
			}

			// get or create the child
			child, exists := obj[key]
			if !exists {
				child = createStructureFor(remaining[0])
			}

			// construct the remaining path elements
			data, err := walker.walk(child, remaining)
			if err != nil {
				return nil, err
			}

			// update the object with any modifications
			obj[key] = data
			return obj, nil
		},
		onMissingElement: func(token string) (any, error) {
			return createStructureFor(token), nil
		},
	}

	walker = newPointerWalker(&handler)
	return walker.walk(root, tokens)
}

// createStructureFor creates an empty array or map based on the next token.
// This is a helper for handlers that need to create intermediate structures.
func createStructureFor(token string) any {
	// array append operation
	if token == ArrayAppendOp {
		return []any{}
	}

	// array index
	if _, err := strconv.Atoi(token); err == nil {
		return []any{}
	}

	// object key
	return make(map[string]any)
}
