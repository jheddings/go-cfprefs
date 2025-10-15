package cfprefs

import (
	"fmt"
	"strings"

	"github.com/jheddings/go-cfprefs/internal"
)

// Removes a preference key for the given application ID
func Delete(appID, keypath string) error {
	segments := strings.Split(keypath, "/")

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
		return fmt.Errorf("keypath error: %s [%s] - segment '%s' is not a dictionary (type: %T)",
			keypath, appID, segments[0], rootValue)
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
			return fmt.Errorf("keypath error: %s [%s] - segment '%s' is not a dictionary (type: %T)",
				keypath, appID, segment, value)
		}
		currentDict = nextDict
	}

	// delete the final key from the parent dictionary
	finalKey := segments[len(segments)-1]
	delete(currentDict, finalKey)

	// write the modified root dictionary back
	return internal.Set(appID, segments[0], rootDict)
}

// Checks if a preference key exists for the given application ID.
func Exists(appID, keypath string) (bool, error) {
	segments := strings.Split(keypath, "/")

	// check if the root key exists
	exists, err := internal.Exists(appID, segments[0])
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	// if there's only one segment, we're done
	if len(segments) == 1 {
		return true, nil
	}

	// get the root value
	value, err := internal.Get(appID, segments[0])
	if err != nil {
		return false, nil
	}

	// traverse the remaining path segments
	for i := 1; i < len(segments); i++ {
		dict, ok := value.(map[string]any)
		if !ok {
			// segment is not a dictionary, path doesn't exist
			return false, nil
		}

		value, ok = dict[segments[i]]
		if !ok {
			// segment not found
			return false, nil
		}
	}

	return true, nil
}
