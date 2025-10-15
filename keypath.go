package cfprefs

import (
	"errors"
	"strings"
)

var (
	// ErrEmptyKeypath is returned when an empty keypath is provided
	ErrEmptyKeypath = errors.New("keypath cannot be empty")
)

// splitKeypath splits a keypath into segments and validates input.
// Returns the segments or an error if the input is invalid.
func splitKeypath(keypath string) ([]string, error) {
	if keypath == "" {
		return nil, ErrEmptyKeypath
	}

	// Split and filter out empty segments (in case of double slashes)
	segments := strings.Split(keypath, "/")
	filtered := make([]string, 0, len(segments))

	for _, segment := range segments {
		if segment != "" {
			filtered = append(filtered, segment)
		}
	}

	if len(filtered) == 0 {
		return nil, ErrEmptyKeypath
	}

	return filtered, nil
}
