package cfprefs

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/theory/jsonpath/spec"
)

var (
	// ErrInvalidJSONPath is returned when an invalid JSONPath is provided
	ErrInvalidJSONPath = errors.New("invalid JSONPath expression")
)

// parseJSONPath parses a simple JSONPath expression into segments.
// Supports: $.field, $.field.nested, $.field[0], $.array[0].field, $.array[]
//
// This function is a simplified version of the jsonpath.Parse function, which
// adds support for append operations (empty array brackets []).
func parseJSONPath(path string) ([]*spec.Segment, error) {
	path = strings.TrimPrefix(path, "$")
	path = strings.TrimPrefix(path, ".")

	if path == "" {
		return []*spec.Segment{}, nil
	}

	// match field names, array indices, and empty array brackets
	re := regexp.MustCompile(`([^.\[\]]+)|\[(\d*)\]`)
	matches := re.FindAllStringSubmatch(path, -1)

	var segments []*spec.Segment
	for _, match := range matches {
		var selector spec.Selector
		if match[1] != "" {
			// field access
			selector = spec.Name(match[1])
		} else if match[0] == "[]" {
			// empty array brackets - append operation
			selector = spec.Index(-1)
		} else if match[2] != "" {
			// array index access
			idx, err := strconv.Atoi(match[2])
			if err != nil {
				return nil, fmt.Errorf("invalid array index: %s", match[2])
			}
			selector = spec.Index(idx)
		}
		segments = append(segments, spec.Child(selector))
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("invalid JSONPath: %s", path)
	}

	return segments, nil
}
