package internal

import (
	"errors"
	"fmt"
)

// Sentinel errors for CoreFoundation failures
var (
	// ErrCFOperation is the base error for CoreFoundation operations
	ErrCFOperation = errors.New("CoreFoundation operation failed")

	// ErrCFReference is returned when creating or accessing CF references fails
	ErrCFReference = errors.New("CoreFoundation reference error")

	// ErrCFType is returned when CF type conversions fail
	ErrCFType = errors.New("CoreFoundation type error")

	// ErrCFSync is returned when CF synchronization fails
	ErrCFSync = errors.New("CoreFoundation synchronization error")

	// ErrCFMemory is returned when CF memory operations fail
	ErrCFMemory = errors.New("CoreFoundation memory error")
)

// CFErr represents a generic CoreFoundation error with context
type CFErr struct {
	Op  string // Operation that failed
	Err error  // Underlying error
	Msg string // Additional context
}

// NewCFError creates a new CFErr for the given operation and error
func NewCFError(op string, err error) *CFErr {
	if err == nil {
		err = ErrCFOperation
	}
	return &CFErr{Op: op, Err: err}
}

// WithMsg adds context to the error
func (e *CFErr) WithMsg(msg string) *CFErr {
	e.Msg = msg
	return e
}

// WithMsgF adds formatted context to the error
func (e *CFErr) WithMsgF(format string, a ...any) *CFErr {
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

// Error returns the error message
func (e *CFErr) Error() string {
	var msg string
	if e.Op != "" {
		msg = e.Op
		if e.Msg != "" {
			msg += ": " + e.Msg
		}
	} else if e.Msg != "" {
		msg = e.Msg
	} else {
		msg = "CoreFoundation error"
	}

	if e.Err != nil && e.Err != ErrCFOperation {
		msg += ": " + e.Err.Error()
	}

	return msg
}

// Wrap wraps an error with the CFErr
func (e *CFErr) Wrap(err error) *CFErr {
	e.Err = errors.Join(e.Err, err)
	return e
}

// Unwrap returns the underlying error
func (e *CFErr) Unwrap() error {
	return e.Err
}

// Is implements support for errors.Is
func (e *CFErr) Is(target error) bool {
	switch e.Op {
	case "reference":
		if target == ErrCFReference {
			return true
		}
	case "type conversion":
		if target == ErrCFType {
			return true
		}
	case "synchronization":
		if target == ErrCFSync {
			return true
		}
	case "memory":
		if target == ErrCFMemory {
			return true
		}
	}

	return target == ErrCFOperation || target == e.Err
}

// CFRefError creates a new reference error
func CFRefError() *CFErr {
	return &CFErr{Op: "reference", Err: ErrCFReference}
}

// CFTypeError creates a new type conversion error
func CFTypeError() *CFErr {
	return &CFErr{Op: "type conversion", Err: ErrCFType}
}

// CFSyncError creates a new synchronization error
func CFSyncError() *CFErr {
	return &CFErr{Op: "synchronization", Err: ErrCFSync}
}

// CFMemoryError creates a new memory error
func CFMemoryError(msg string) *CFErr {
	return &CFErr{Op: "memory", Err: ErrCFMemory, Msg: msg}
}
