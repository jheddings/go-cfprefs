package cfprefs

import (
	"fmt"

	"github.com/theory/jsonpath/spec"
)

// arraySegmentHandler defines callbacks for handling array path walk operations.
type arraySegmentHandler struct {
	// onLast is called when reaching the last segment of an array path.
	// Returns the modified array or an error.
	onLast func(arr []any, index int) ([]any, error)

	// onContinue is called when traversing an array path with more segments.
	// It receives the current array element (or nil if out of bounds) and remaining segments.
	// Returns the modified array or an error.
	onContinue func(arr []any, index int, element any, remaining []*spec.Segment) ([]any, error)
}

// mapSegmentHandler defines callbacks for handling map path walk operations.
type mapSegmentHandler struct {
	// onLast is called when reaching the last segment of a map path.
	// Returns the modified map or an error.
	onLast func(obj map[string]any, key string) (map[string]any, error)

	// onContinue is called when traversing a map path with more segments.
	// It receives the child element (or nil if key doesn't exist) and remaining segments.
	// Returns the modified map or an error.
	onContinue func(obj map[string]any, key string, child any, remaining []*spec.Segment) (map[string]any, error)
}

type pathWalker struct {
	arrayHandler *arraySegmentHandler
	mapHandler   *mapSegmentHandler
}

func newPathWalker(handlers ...any) *pathWalker {
	w := &pathWalker{}
	for _, handler := range handlers {
		w.withHandler(handler)
	}
	return w
}

func (w *pathWalker) withHandler(handler any) *pathWalker {
	if arrayHandler, ok := handler.(*arraySegmentHandler); ok {
		w.arrayHandler = arrayHandler
	} else if arrayHandler, ok := handler.(arraySegmentHandler); ok {
		w.arrayHandler = &arrayHandler
	} else if mapHandler, ok := handler.(*mapSegmentHandler); ok {
		w.mapHandler = mapHandler
	} else if mapHandler, ok := handler.(mapSegmentHandler); ok {
		w.mapHandler = &mapHandler
	} else {
		panic(fmt.Sprintf("invalid handler type: %T", handler))
	}

	return w
}

// walk recursively traverses segments of a path and calls the appropriate handler.
func (w *pathWalker) walk(data any, segments []*spec.Segment) (any, error) {

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
	if idx, ok := selector.(spec.Index); ok && w.arrayHandler != nil {
		arr, ok := data.([]any)
		if !ok {
			return nil, NewTypeMismatchError([]any{}, data)
		}

		index := int(idx)

		// handle last segment
		if isLast {
			return w.arrayHandler.onLast(arr, index)
		}

		// for non-last segments, pass the element and remaining segments to the handler
		var element any
		if index >= 0 && index < len(arr) {
			element = arr[index]
		}

		// element is nil if index is out of bounds or special (like -1)
		return w.arrayHandler.onContinue(arr, index, element, segments[1:])
	}

	// handle field name selector
	if field, ok := selector.(spec.Name); ok && w.mapHandler != nil {
		obj, ok := data.(map[string]any)
		if !ok {
			return nil, NewTypeMismatchError(map[string]any{}, data)
		}

		name := string(field)

		// handle last segment
		if isLast {
			return w.mapHandler.onLast(obj, name)
		}

		// pass the child (or nil if it doesn't exist) and remaining segments to the handler
		child, _ := obj[name]
		return w.mapHandler.onContinue(obj, name, child, segments[1:])
	}

	return nil, NewInternalError().WithMsgF("unsupported selector type: %T", selector)
}
