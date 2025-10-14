package cfprefs

import "github.com/jheddings/go-cfprefs/internal"

// Set writes a preference value for the given key and application ID.
// The value can be any native Go type: string, bool, int, int64, float64,
// time.Time, []byte, []any, or map[string]any.
func Set(appID, key string, value any) error {
	return internal.Set(appID, key, value)
}
