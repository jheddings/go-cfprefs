package cfprefs

import (
	"github.com/go-openapi/jsonpointer"
	"github.com/jheddings/go-cfprefs/internal"
)

// Delete removes a preference value at the given keypath and application ID.
// If the keypath is a simple name, it deletes the entire key.
// If the keypath is a JSON Pointer path, it deletes the value at the specified path.
//
// Example usage:
//
//	// Delete a simple value
//	err := Delete("com.example.app", "username")
//
//	// Delete a nested value
//	err := Delete("com.example.app", "config/server/port")
//
// Returns an error if the keypath is invalid or the value cannot be deleted.
func Delete(appID, keypath string) error {
	kp, err := parseKeypath(keypath)
	if err != nil {
		return NewKeyPathError().Wrap(err).WithMsgF("invalid keypath: %s", keypath)
	}

	// if there is no pointer path, just delete the entire key
	if kp.IsRoot() {
		return internal.Delete(appID, kp.Key)
	}

	// check if the key exists
	exists, err := internal.Exists(appID, kp.Key)
	if err != nil {
		return NewInternalError().Wrap(err).WithMsgF("failed to check: %s", kp.Key)
	}

	// if key doesn't exist, return success (idempotent)
	if !exists {
		return nil
	}

	// get the current value
	root, err := internal.Get(appID, kp.Key)
	if err != nil {
		return NewInternalError().Wrap(err).WithMsgF("failed to get: %s", kp.Key)
	}

	ptr, err := jsonpointer.New(kp.Path)
	if err != nil {
		return NewKeyPathError().Wrap(err).WithMsgF("invalid path: %s", kp.Path)
	}

	// delete the value at the specified path
	modified, deleted, err := deleteValueAtPath(root, ptr.DecodedTokens())
	if err != nil {
		return err
	}

	// if nothing was deleted, return success (idempotent)
	if !deleted {
		return nil
	}

	// otherwise, write the modified data back
	return internal.Set(appID, kp.Key, modified)
}

// deleteValueAtPath uses a pointer walker to delete a value at the specified path.
func deleteValueAtPath(root any, tokens []string) (any, bool, error) {
	var modified bool
	var walker *pointerWalker

	handler := pathTokenHandler{
		onArrayIndex: func(arr []any, index int, remaining []string) (any, error) {
			// if this is the last token, drop the element at the index
			if len(remaining) == 0 {
				result := append(arr[:index], arr[index+1:]...)
				modified = true
				return result, nil
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
		onObjectKey: func(obj map[string]any, key string, remaining []string) (any, error) {
			// if the key doesn't exist, that's ok (idempotent)
			child, exists := obj[key]
			if !exists {
				return obj, nil
			}

			// if this is the last token, delete the key
			if len(remaining) == 0 {
				delete(obj, key)
				modified = true
				return obj, nil
			}

			// recursively delete in the child
			data, err := walker.walk(child, remaining)
			if err != nil {
				return nil, err
			}

			// update the object with any modifications
			obj[key] = data
			return obj, nil
		},
		onMissingElement: func(token string) (any, error) {
			return nil, NewKeyPathError().WithMsg("path not found")
		},
	}

	walker = newPointerWalker(&handler)
	data, err := walker.walk(root, tokens)
	return data, modified, err
}
