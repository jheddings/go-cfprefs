package cfprefs

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/theory/jsonpath/spec"
)

const (
	// ArrayAppendIndex is a special index value used to indicate an append operation.
	// When used in a JSONPath expression with empty brackets (e.g., $.items[]),
	// it signals that a new element should be appended to the array.
	ArrayAppendIndex = -1
)

var (
	// ErrInvalidJSONPath is returned when an invalid JSONPath is provided
	ErrInvalidJSONPath = errors.New("invalid JSONPath expression")
)

type keypath struct {
	Key  string
	Path string
}

func newKeypath(pref, path string) *keypath {
	return &keypath{Key: pref, Path: path}
}

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

// parseJSONPath parses a simple JSONPath expression into segments.
// Supports: $.field, $.field.nested, $.field[0], $.array[0].field, $.array[]
//
// This function is a simplified version of the jsonpath.Parse function, which
// adds support for append operations (empty array brackets []).
func parseJSONPath(path string) ([]*spec.Segment, error) {
	if path == "" {
		return []*spec.Segment{}, nil
	}

	// match field names, array indices, and empty array brackets
	re := regexp.MustCompile(`([^.\[\]]+)|\[(\d*)\]`)
	matches := re.FindAllStringSubmatch(path, -1)

	var segments []*spec.Segment
	for _, match := range matches {
		var selector spec.Selector
		if match[0] == "$" {
			// root element - no selector
			continue
		} else if match[0] == "[]" {
			// empty array brackets - append operation
			selector = spec.Index(ArrayAppendIndex)
		} else if match[1] != "" {
			// field access
			selector = spec.Name(match[1])
		} else if match[2] != "" {
			// array index access
			idx, err := strconv.Atoi(match[2])
			if err != nil {
				return nil, NewKeyPathError().WithMsgF("invalid array index: %s", match[2])
			}
			selector = spec.Index(idx)
		}
		segments = append(segments, spec.Child(selector))
	}

	if len(segments) == 0 {
		return nil, NewKeyPathError().WithMsgF("invalid JSONPath: %s", path)
	}

	return segments, nil
}
