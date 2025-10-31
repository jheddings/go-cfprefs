package cfprefs

import (
	"strconv"
)

// pathTokenHandler defines callbacks for handling path token operations.
type pathTokenHandler struct {
	// onArrayIndex is called when operating on an array element.
	// index is the array index (validated to be within bounds).
	// Returns the modified array or an error.
	onArrayIndex func(arr []any, index int, remaining []string) (any, error)

	// onArrayAppend is called for array append operations (ArrayAppendOp token).
	// Returns the modified array or an error.
	onArrayAppend func(arr []any, remaining []string) (any, error)

	// onObjectKey is called when operating on an object key.
	// Returns the modified object or an error.
	onObjectKey func(obj map[string]any, key string, remaining []string) (any, error)

	// onMissingElement is called when a required structure doesn't exist.
	// token is the current token being processed.
	// Returns the structure to use or an error.
	onMissingElement func(token string) (any, error)
}

// pointerWalker traverses JSON structures using JSON Pointer tokens.
type pointerWalker struct {
	handler *pathTokenHandler
}

// newPointerWalker creates a new pointer walker with the specified handler.
func newPointerWalker(handler *pathTokenHandler) *pointerWalker {
	return &pointerWalker{
		handler: handler,
	}
}

// walk traverses the data structure using JSON Pointer tokens and calls appropriate handlers.
func (w *pointerWalker) walk(node any, tokens []string) (any, error) {
	// base case: no more tokens
	if len(tokens) == 0 {
		return node, nil
	}

	token := tokens[0]
	remaining := tokens[1:]

	// handle array append operations
	if token == ArrayAppendOp {
		return w.walkArrayAppend(node, remaining)
	}

	// handle array index tokens
	if idx, err := strconv.Atoi(token); err == nil {
		return w.walkArrayIndex(node, idx, remaining)
	}

	// handle object key tokens
	return w.walkObjectKey(node, token, remaining)
}

// walkArrayAppend handles array append operations.
func (w *pointerWalker) walkArrayAppend(node any, remaining []string) (any, error) {
	// ensure we have an array
	arr, ok := node.([]any)
	if !ok {
		// if node is not an array, try to create one
		if w.handler.onMissingElement != nil {
			new, err := w.handler.onMissingElement(ArrayAppendOp)
			if err != nil {
				return nil, err
			}
			arr, ok = new.([]any)
			if !ok {
				return nil, NewInternalError().WithMsg("onMissing did not return an array for array append operation")
			}
		} else {
			return nil, NewKeyPathError().WithMsg("cannot append to non-array value")
		}
	}

	if w.handler.onArrayAppend != nil {
		return w.handler.onArrayAppend(arr, remaining)
	}

	return nil, NewInternalError().WithMsg("no handler for array append operation")
}

// walkArrayIndex handles array index operations.
func (w *pointerWalker) walkArrayIndex(node any, index int, remaining []string) (any, error) {
	// ensure we have an array
	arr, ok := node.([]any)
	if !ok {
		return nil, NewKeyPathError().WithMsg("cannot index non-array value")
	}

	// validate bounds
	if index < 0 || index >= len(arr) {
		return nil, NewKeyPathError().WithMsgF("array index out of bounds: %d (array length: %d)", index, len(arr))
	}

	if w.handler.onArrayIndex != nil {
		return w.handler.onArrayIndex(arr, index, remaining)
	}

	return nil, NewInternalError().WithMsg("no handler for array element operation")
}

// walkObjectKey handles object key operations.
func (w *pointerWalker) walkObjectKey(node any, key string, remaining []string) (any, error) {
	// ensure we have an object
	obj, ok := node.(map[string]any)
	if !ok {
		// if node is not nil and not an object, we can't traverse through it
		if node != nil {
			return nil, NewKeyPathError().WithMsg("cannot traverse through non-object value")
		}
		// node is nil, try to create an object
		if w.handler.onMissingElement != nil {
			new, err := w.handler.onMissingElement(key)
			if err != nil {
				return nil, err
			}
			obj, ok = new.(map[string]any)
			if !ok {
				return nil, NewInternalError().WithMsg("onMissing did not return an object for key operation")
			}
		} else {
			return nil, NewKeyPathError().WithMsg("cannot create object at path")
		}
	}

	if w.handler.onObjectKey != nil {
		return w.handler.onObjectKey(obj, key, remaining)
	}

	return nil, NewInternalError().WithMsg("no handler for object key operation")
}
