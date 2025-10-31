package cfprefs

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidJSONPath is returned when an invalid JSONPath is provided
	ErrInvalidJSONPath = errors.New("invalid JSONPath expression")
)

// keypath represents a key and an optional JSON Pointer path.
type keypath struct {
	Key  string
	Path string
}

// newKeypath creates a new keypath with the given key and path.
func newKeypath(pref, path string) *keypath {
	return &keypath{Key: pref, Path: path}
}

// parseKeypath parses a keypath string into a keypath struct.
func parseKeypath(keypath string) (*keypath, error) {
	if keypath == "" {
		return nil, NewKeyPathError().WithMsg("invalid keypath; empty string")
	}

	if strings.HasPrefix(keypath, "/") {
		return nil, NewKeyPathError().WithMsg("invalid keypath; empty pref")
	}

	parts := strings.SplitN(keypath, "/", 2)

	if len(parts) < 2 {
		return newKeypath(parts[0], ""), nil
	}

	path := parts[1]
	if path != "" {
		path = "/" + path
	}

	return newKeypath(parts[0], path), nil
}

// String returns the keypath as a string.
func (k *keypath) String() string {
	return k.Key + k.Path
}

// IsRoot returns true if the keypath is a root keypath (empty pointer).
func (k *keypath) IsRoot() bool {
	return k.Path == "" || k.Path == "/"
}
