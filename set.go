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
	if kp.Path == "" {
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
	modified, err := setValueInNode(root, ptr.DecodedTokens(), value)
	if err != nil {
		return NewInternalError().Wrap(err).WithMsgF("failed to set: %s", kp.Path)
	}

	// write the modified root value back
	return internal.Set(appID, kp.Key, modified)
}

// setValueInNode recursively traverses the data structure using JSON Pointer tokens
// and sets the value at the final token, creating intermediate structures as needed.
func setValueInNode(node any, tokens []string, value any) (any, error) {
	// if there are no more tokens, exit early
	if len(tokens) == 0 {
		return value, nil
	}

	// work on the current token
	token := tokens[0]

	// handle array append operations
	if token == ArrayAppendOp {
		return setArrayAppend(node, tokens, value)
	}

	// handle array index tokens
	if idx, err := strconv.Atoi(token); err == nil {
		return setArrayIndex(node, idx, tokens, value)
	}

	// handle object key tokens
	return setObjectKey(node, token, tokens, value)
}

// setArrayAppend sets the value at the end of an array
func setArrayAppend(node any, tokens []string, value any) (any, error) {
	// ensure we have an array (create if needed for append operations)
	arr, ok := node.([]any)
	if !ok {
		arr = []any{}
	}

	// if this is the last token, append the value
	if len(tokens) == 1 {
		return append(arr, value), nil
	}

	// create a new element and continue setting
	newElement := newStructFor(tokens[1])
	modified, err := setValueInNode(newElement, tokens[1:], value)
	if err != nil {
		return nil, err
	}

	return append(arr, modified), nil
}

// setArrayIndex handles setting values in arrays
func setArrayIndex(node any, index int, tokens []string, value any) (any, error) {
	// ensure we have an array
	arr, ok := node.([]any)
	if !ok {
		return nil, NewKeyPathError().WithMsg("cannot index non-array value")
	}

	// validate array bounds
	if index < 0 || index >= len(arr) {
		return nil, NewKeyPathError().WithMsgF("array index out of bounds: %d (array length: %d)", index, len(arr))
	}

	// if this is the last token, set the value at the index
	if len(tokens) == 1 {
		arr[index] = value
		return arr, nil
	}

	// recursively set in the element
	modified, err := setValueInNode(arr[index], tokens[1:], value)
	if err != nil {
		return nil, err
	}

	arr[index] = modified
	return arr, nil
}

// setObjectKey handles setting values in objects
func setObjectKey(node any, key string, tokens []string, value any) (any, error) {
	// ensure we have an object
	obj, ok := node.(map[string]any)
	if !ok {
		// if data is not nil and not an object, we can't traverse through it
		if node != nil {
			return nil, NewKeyPathError().WithMsg("cannot traverse through non-object value")
		}
		// if data is nil, create a new object
		obj = make(map[string]any)
	}

	// if this is the last token, set the value at the key
	if len(tokens) == 1 {
		obj[key] = value
		return obj, nil
	}

	// get or create the child
	child, exists := obj[key]
	if !exists {
		child = newStructFor(tokens[1])
	}

	// recursively set in the child
	modified, err := setValueInNode(child, tokens[1:], value)
	if err != nil {
		return nil, err
	}

	obj[key] = modified
	return obj, nil
}

// newStructFor creates an empty array or map based on the token
func newStructFor(token string) any {
	// if the token is an array append operation...
	if token == ArrayAppendOp {
		return []any{}
	}

	// if the token is an array index...
	if _, err := strconv.Atoi(token); err == nil {
		return []any{}
	}

	// anything else is a map
	return make(map[string]any)
}
