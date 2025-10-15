package cfprefs

import (
	"fmt"

	"github.com/jheddings/go-cfprefs/internal"
)

// Set writes a preference value for the given keypath and application ID.
// If any intermediate dictionaries do not exist, they will be created.
// Returns an error if the operation fails or if any intermediate segment
// is not a dictionary.
func Set(appID, keypath string, value any) error {
	segments, err := splitKeypath(keypath)
	if err != nil {
		return err
	}

	// if there's only one segment, just set the value directly
	if len(segments) == 1 {
		return internal.Set(appID, keypath, value)
	}

	// get or create the root dictionary
	var rootDict map[string]any
	exists, err := Exists(appID, segments[0])
	if err != nil {
		return err
	} else if exists {
		if rootDict, err = GetMap(appID, segments[0]); err != nil {
			return err
		}
	} else {
		rootDict = make(map[string]any)
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
				return fmt.Errorf("keypath error: %s [%s] - segment '%s' exists but is not a dictionary (type: %T)", keypath, appID, segment, value)
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
