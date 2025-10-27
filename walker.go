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
	// ErrEmptyKeypath is returned when an empty keypath is provided
	ErrEmptyKeypath = errors.New("keypath cannot be empty")

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

// walkHandler defines callbacks for handling path walk operations.
type walkHandler struct {
	// onArrayLast is called when reaching the last segment of an array path.
	// Returns the modified array or an error.
	onArrayLast func(arr []any, index int) ([]any, error)

	// onArrayContinue is called when traversing an array path with more segments.
	// It receives the current array element (or nil if out of bounds) and remaining segments.
	// Returns the modified array or an error.
	onArrayContinue func(arr []any, index int, element any, remaining []*spec.Segment) ([]any, error)

	// onMapLast is called when reaching the last segment of a map path.
	// Returns the modified map or an error.
	onMapLast func(obj map[string]any, key string) (map[string]any, error)

	// onMapContinue is called when traversing a map path with more segments.
	// It receives the child element (or nil if key doesn't exist) and remaining segments.
	// Returns the modified map or an error.
	onMapContinue func(obj map[string]any, key string, child any, remaining []*spec.Segment) (map[string]any, error)
}

// walkPath recursively traverses segments of a path and applies the provided handlers.
// This generic walker can be used for both set and delete operations.
func walkPath(data any, segments []*spec.Segment, handler walkHandler) (any, error) {
	if len(segments) == 0 {
		return data, nil
	}

	segment := segments[0]
	selectors := segment.Selectors()
	if len(selectors) != 1 {
		return nil, NewInternalError().WithMsgF("expected single selector in segment, got %d", len(selectors))
	}

	selector := selectors[0]
	isLast := len(segments) == 1

	// handle array index selector
	if idx, ok := selector.(spec.Index); ok {
		arr, ok := data.([]any)
		if !ok {
			return nil, NewTypeMismatchError([]any{}, data)
		}

		index := int(idx)

		// handle last segment
		if isLast {
			return handler.onArrayLast(arr, index)
		}

		// for non-last segments, pass the element and remaining segments to the handler
		var element any
		if index >= 0 && index < len(arr) {
			element = arr[index]
		}
		// element is nil if index is out of bounds or special (like -1)

		return handler.onArrayContinue(arr, index, element, segments[1:])
	}

	// handle field name selector
	if field, ok := selector.(spec.Name); ok {
		obj, ok := data.(map[string]any)
		if !ok {
			return nil, NewTypeMismatchError(map[string]any{}, data)
		}

		name := string(field)

		// handle last segment
		if isLast {
			return handler.onMapLast(obj, name)
		}

		// pass the child (or nil if it doesn't exist) and remaining segments to the handler
		child, _ := obj[name]
		return handler.onMapContinue(obj, name, child, segments[1:])
	}

	return nil, NewInternalError().WithMsgF("unsupported selector type: %T", selector)
}
