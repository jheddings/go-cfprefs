package cfprefs

import (
	"errors"
	"fmt"
)

// Sentinel errors for common failures
var (
	// ErrKeyNotFound is returned when a requested key does not exist
	ErrKeyNotFound = errors.New("key not found")

	// ErrInvalidKeyPath is returned when a key path is malformed
	ErrInvalidKeyPath = errors.New("invalid key path")

	// ErrTypeMismatch is returned when a value is not of the expected type
	ErrTypeMismatch = errors.New("type mismatch")

	// ErrInternal is returned when an internal error occurs
	ErrInternal = errors.New("internal error")
)

// InternalErr represents an error that is internal to the library
type InternalErr struct {
	Err error
	Msg string
}

// NewInternalError creates a new InternalErr
func NewInternalError() *InternalErr {
	return &InternalErr{Err: ErrInternal}
}

// WithMsg adds a custom message to the error
func (e *InternalErr) WithMsg(msg string) *InternalErr {
	e.Msg = msg
	return e
}

// WithMsgF adds a formatted custom message to the error
func (e *InternalErr) WithMsgF(format string, a ...any) *InternalErr {
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

// Error returns the error message
func (e *InternalErr) Error() string {
	if e.Msg == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
}

// Is implements support for errors.Is
func (e *InternalErr) Is(target error) bool {
	return target == e.Err
}

// Wrap wraps an error with the InternalErr
func (e *InternalErr) Wrap(err error) *InternalErr {
	e.Err = errors.Join(e.Err, err)
	return e
}

// Unwrap returns the underlying error
func (e *InternalErr) Unwrap() error {
	return e.Err
}

// KeyNotFoundErr represents an error when a preference key is not found
type KeyNotFoundErr struct {
	AppID string
	Key   string
	Msg   string
	Err   error
}

// NewKeyNotFoundError creates a new KeyNotFoundErr
func NewKeyNotFoundError(appID, key string) *KeyNotFoundErr {
	return &KeyNotFoundErr{AppID: appID, Key: key}
}

// WithMsg adds a custom message to the error
func (e *KeyNotFoundErr) WithMsg(msg string) *KeyNotFoundErr {
	e.Msg = msg
	return e
}

// WithMsgF adds a formatted custom message to the error
func (e *KeyNotFoundErr) WithMsgF(format string, a ...any) *KeyNotFoundErr {
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

// Error returns the error message
func (e *KeyNotFoundErr) Error() string {
	if e.Msg == "" {
		return fmt.Sprintf("key not found: %s [%s]", e.Key, e.AppID)
	}
	return fmt.Sprintf("key not found: %s [%s] - %s", e.Key, e.AppID, e.Msg)
}

// Is implements support for errors.Is
func (e *KeyNotFoundErr) Is(target error) bool {
	return target == ErrKeyNotFound
}

// Wrap wraps an error with the KeyNotFoundErr
func (e *KeyNotFoundErr) Wrap(err error) *KeyNotFoundErr {
	e.Err = errors.Join(e.Err, err)
	return e
}

// Unwrap returns the underlying error
func (e *KeyNotFoundErr) Unwrap() error {
	return ErrKeyNotFound
}

// KeyPathErr represents an error with a key path
type KeyPathErr struct {
	Msg string
	Err error
}

// NewKeyPathError creates a new KeyPathErr
func NewKeyPathError() *KeyPathErr {
	return &KeyPathErr{Err: ErrInvalidKeyPath}
}

// WithMsg adds a custom message to the error
func (e *KeyPathErr) WithMsg(msg string) *KeyPathErr {
	e.Msg = msg
	return e
}

// WithMsgF adds a formatted custom message to the error
func (e *KeyPathErr) WithMsgF(format string, a ...any) *KeyPathErr {
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

// Error returns the error message
func (e *KeyPathErr) Error() string {
	if e.Msg != "" {
		return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
	}
	return e.Err.Error()
}

// Is implements support for errors.Is
func (e *KeyPathErr) Is(target error) bool {
	return target == ErrInvalidKeyPath
}

// Wrap wraps an error with the KeyPathErr
func (e *KeyPathErr) Wrap(err error) *KeyPathErr {
	e.Err = errors.Join(e.Err, err)
	return e
}

// Unwrap returns the underlying error
func (e *KeyPathErr) Unwrap() error {
	return ErrInvalidKeyPath
}

// TypeMismatchErr represents a type mismatch error
type TypeMismatchErr struct {
	AppID    string
	Key      string
	Expected any
	Actual   any
	Err      error
}

// NewTypeMismatchError creates a new TypeMismatchErr
func NewTypeMismatchError(expected, actual any) *TypeMismatchErr {
	return &TypeMismatchErr{Expected: expected, Actual: actual}
}

// Error returns the error message
func (e *TypeMismatchErr) Error() string {
	return fmt.Sprintf("type mismatch: %s [%s] - expected %T, got %T", e.Key, e.AppID, e.Expected, e.Actual)
}

// Is implements support for errors.Is
func (e *TypeMismatchErr) Is(target error) bool {
	return target == ErrTypeMismatch
}

// WithKey adds an appID and key to the error
func (e *TypeMismatchErr) WithKey(appID, key string) *TypeMismatchErr {
	e.AppID = appID
	e.Key = key
	return e
}

// Wrap wraps an error with the TypeMismatchErr
func (e *TypeMismatchErr) Wrap(err error) *TypeMismatchErr {
	e.Err = errors.Join(e.Err, err)
	return e
}

// Unwrap returns the underlying error
func (e *TypeMismatchErr) Unwrap() error {
	return ErrTypeMismatch
}
