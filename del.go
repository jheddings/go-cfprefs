package cfprefs

import (
	"fmt"
	"strings"

	"github.com/jheddings/go-cfprefs/internal"
)

// Delete removes a preference value for the given key and application ID.
// The key may be a keypath separated by forward slashes ("/") to traverse
// nested dictionaries. For example, "map-test/string" will delete the "string"
// key from the "map-test" dictionary while preserving the parent dictionary
// and any other keys within it.
func Delete(appID, key string) error {
	segments := strings.Split(key, "/")

	// if there's only one segment, delete directly
	if len(segments) == 1 {
		return internal.Delete(appID, key)
	}

	// get the root dictionary
	rootValue, err := internal.Get(appID, segments[0])
	if err != nil {
		// root doesn't exist, nothing to delete
		return nil
	}

	// verify root is a dictionary
	rootDict, ok := rootValue.(map[string]any)
	if !ok {
		return fmt.Errorf("keypath error: %s [%s] - segment '%s' is not a dictionary (type: %T)",
			key, appID, segments[0], rootValue)
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
				key, appID, segment, value)
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
// The key may be a keypath separated by forward slashes ("/") to traverse
// nested dictionaries. Returns true only if all elements in the keypath exist.
// For example, "map-test/string" will return true only if both the "map-test"
// dictionary exists and it contains a "string" key.
func Exists(appID, key string) (bool, error) {
	segments := strings.Split(key, "/")

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
		// Current value must be a map to traverse further
		dict, ok := value.(map[string]any)
		if !ok {
			// segment is not a dictionary, path doesn't exist
			return false, nil
		}

		// Check if the next segment exists in the dictionary
		value, ok = dict[segments[i]]
		if !ok {
			// segment not found
			return false, nil
		}
	}

	return true, nil
}
