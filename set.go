package cfprefs

import "github.com/jheddings/go-cfprefs/internal"

// Set writes a preference value for the given key and application ID.
func Set(appID, key string, value any) error {
	return internal.Set(appID, key, value)
}
