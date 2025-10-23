package cfprefs

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

// pathSegment represents a single segment in a JSONPath
type pathSegment struct {
	field      string // field name for object access
	index      int    // index for array access
	isArrayIdx bool   // true if this is an array index access
}

// parseJSONPath parses a simple JSONPath expression into segments.
// Supports: $.field, $.field.nested, $.field[0], $.array[0].field
func parseJSONPath(path string) ([]pathSegment, error) {
	path = strings.TrimPrefix(path, "$")
	path = strings.TrimPrefix(path, ".")

	if path == "" {
		return []pathSegment{}, nil
	}

	// match field names and array indices
	re := regexp.MustCompile(`([^.\[\]]+)|\[(\d+)\]`)
	matches := re.FindAllStringSubmatch(path, -1)

	var segments []pathSegment
	for _, match := range matches {
		if match[1] != "" {
			// field access
			segments = append(segments, pathSegment{field: match[1], isArrayIdx: false})
		} else if match[2] != "" {
			// array index access
			idx, err := strconv.Atoi(match[2])
			if err != nil {
				return nil, fmt.Errorf("invalid array index: %s", match[2])
			}
			segments = append(segments, pathSegment{index: idx, isArrayIdx: true})
		}
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("invalid JSONPath: %s", path)
	}

	return segments, nil
}
