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

	// ErrInvalidJSONPath is returned when an invalid JSONPath is provided
	ErrInvalidJSONPath = errors.New("invalid JSONPath expression")
)

// pathSegment represents a single segment in a JSONPath
type pathSegment struct {
	field      string // field name for object access
	index      int    // index for array access
	isArrayIdx bool   // true if this is an array index access
}

// parseJSONPath parses a simple JSONPath expression into segments.
// Supports: $.field, $.field.nested, $.field[0], $.array[0].field, $.array[]
func parseJSONPath(path string) ([]pathSegment, error) {
	path = strings.TrimPrefix(path, "$")
	path = strings.TrimPrefix(path, ".")

	if path == "" {
		return []pathSegment{}, nil
	}

	// match field names, array indices, and empty array brackets
	re := regexp.MustCompile(`([^.\[\]]+)|\[(\d*)\]`)
	matches := re.FindAllStringSubmatch(path, -1)

	var segments []pathSegment
	for _, match := range matches {
		if match[1] != "" {
			// field access
			segments = append(segments, pathSegment{field: match[1], isArrayIdx: false})
		} else if match[0] == "[]" {
			// empty array brackets - append operation
			segments = append(segments, pathSegment{index: -1, isArrayIdx: true})
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
