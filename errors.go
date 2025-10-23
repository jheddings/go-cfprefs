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
)

// KeyNotFoundErr represents an error when a preference key is not found
type KeyNotFoundErr struct {
	AppID string
	Key   string
	Msg   string
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

// Unwrap returns the underlying error
func (e *KeyNotFoundErr) Unwrap() error {
	return ErrKeyNotFound
}

// KeyPathErr represents an error with a key path
type KeyPathErr struct {
	AppID string
	Key   string
	Msg   string
	Err   error
}

// NewKeyPathError creates a new KeyPathErr
func NewKeyPathError(appID, key string) *KeyPathErr {
	return &KeyPathErr{AppID: appID, Key: key, Err: ErrInvalidKeyPath}
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
	var msg string
	if e.Msg != "" {
		msg = e.Msg
	} else if e.Err != nil {
		msg = e.Err.Error()
	} else {
		msg = "key path error"
	}
	return fmt.Sprintf("%s: %s [%s]", msg, e.Key, e.AppID)
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
}

// NewTypeMismatchError creates a new TypeMismatchErr
func NewTypeMismatchError(appID, key string, expected, actual any) *TypeMismatchErr {
	return &TypeMismatchErr{AppID: appID, Key: key, Expected: expected, Actual: actual}
}

// Error returns the error message
func (e *TypeMismatchErr) Error() string {
	return fmt.Sprintf("type mismatch: %s [%s] - expected %T, got %T", e.Key, e.AppID, e.Expected, e.Actual)
}

// Is implements support for errors.Is
func (e *TypeMismatchErr) Is(target error) bool {
	return target == ErrTypeMismatch
}

// Unwrap returns the underlying error
func (e *TypeMismatchErr) Unwrap() error {
	return ErrTypeMismatch
}
