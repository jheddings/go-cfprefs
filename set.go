package cfprefs

import (
	"fmt"
	"strings"

	"github.com/jheddings/go-cfprefs/internal"
)

// Set writes a preference value for the given key and application ID.
// The key may be a keypath separated by forward slashes ("/") to traverse
// nested dictionaries. For example, "map-test/string" will set the "string"
// key within the "map-test" dictionary. If any intermediate dictionaries
// do not exist, they will be created.
func Set(appID, key string, value any) error {
	segments := strings.Split(key, "/")

	// if there's only one segment, just set the value directly
	if len(segments) == 1 {
		return internal.Set(appID, key, value)
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
				key, appID, segments[0], rootValue)
		}
	}

	// traverse remaining segments
	currentDict := rootDict
	for i := 1; i < len(segments)-1; i++ {
		segment := segments[i]

		// Check if the segment exists
		value, exists := currentDict[segment]
		if !exists {
			// Create a new dictionary for this segment
			newDict := make(map[string]any)
			currentDict[segment] = newDict
			currentDict = newDict
		} else {
			// Segment exists, verify it's a dictionary
			nextDict, ok := value.(map[string]any)
			if !ok {
				return fmt.Errorf("keypath error: %s [%s] - segment '%s' exists but is not a dictionary (type: %T)",
					key, appID, segment, value)
			}
			currentDict = nextDict
		}
	}

	// Set the final value
	finalKey := segments[len(segments)-1]
	currentDict[finalKey] = value

	// Write the modified root dictionary back
	return internal.Set(appID, segments[0], rootDict)
}
