package cfprefs

import (
	"fmt"
	"strings"

	"github.com/jheddings/go-cfprefs/internal"
)

// Set writes a preference value for the given key and application ID.
// If any intermediate dictionaries do not exist, they will be created.
func Set(appID, keypath string, value any) error {
	segments := strings.Split(keypath, "/")

	// if there's only one segment, just set the value directly
	if len(segments) == 1 {
		return internal.Set(appID, keypath, value)
	}

	// get or create the root dictionary
	var rootDict map[string]any
	rootValue, err := internal.Get(appID, segments[0])
	if err != nil {
		// root doesn't exist, create a new empty dictionary
		rootDict = make(map[string]any)
	} else {
		// root exists, verify it's a dictionary
		var ok bool
		rootDict, ok = rootValue.(map[string]any)
		if !ok {
			return fmt.Errorf("keypath error: %s [%s] - segment '%s' exists but is not a dictionary (type: %T)",
				keypath, appID, segments[0], rootValue)
		}
	}

	// traverse remaining segments
	currentDict := rootDict
	for i := 1; i < len(segments)-1; i++ {
		segment := segments[i]

		value, exists := currentDict[segment]

		if !exists {
			// create a new dictionary for this segment
			newDict := make(map[string]any)
			currentDict[segment] = newDict
			currentDict = newDict

		} else {
			// segments must be dictionaries to traverse further
			nextDict, ok := value.(map[string]any)
			if !ok {
				return fmt.Errorf("keypath error: %s [%s] - segment '%s' exists but is not a dictionary (type: %T)",
					keypath, appID, segment, value)
			}
			currentDict = nextDict
		}
	}

	// set the final value
	finalKey := segments[len(segments)-1]
	currentDict[finalKey] = value

	// write the modified root dictionary back
	return internal.Set(appID, segments[0], rootDict)
}
