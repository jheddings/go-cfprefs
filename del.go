package cfprefs

import "github.com/jheddings/go-cfprefs/internal"

// Delete removes a preference value for the given key and application ID.
func Delete(appID, key string) error {
	return internal.Delete(appID, key)
}

// Exists checks if a preference key exists for the given application ID.
func Exists(appID, key string) (bool, error) {
	return internal.Exists(appID, key)
}
