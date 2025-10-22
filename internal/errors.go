package internal

import "fmt"

// CFErr represents a generic CoreFoundation error.
type CFErr struct {
	Err error
	Msg string
}

// CFError creates a new CFErr from the given error.
func CFError() *CFErr {
	return &CFErr{}
}

// WithErr sets the error for the CFErr.
func (e *CFErr) WithErr(err error) *CFErr {
	e.Err = err
	return e
}

// WithMsg sets the message for the CFErr.
func (e *CFErr) WithMsg(msg string) *CFErr {
	e.Msg = msg
	return e
}

// WithMsgF sets the message for the CFErr using a format string.
func (e *CFErr) WithMsgF(format string, a ...any) *CFErr {
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

// Error returns the error message for the CFErr.
func (e *CFErr) Error() string {
	if e.Msg == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
}

// Unwrap returns the underlying error for the CFErr.
func (e *CFErr) Unwrap() error {
	return e.Err
}

// CFRefErr represents a CoreFoundation reference error.
type CFRefErr struct {
	Err error
	Msg string
}

// CFRefError creates a new CFRefErr from the given error.
func CFRefError(err error) *CFRefErr {
	return &CFRefErr{Err: err}
}

// WithMsg sets the message for the CFRefErr.
func (e *CFRefErr) WithMsg(msg string) *CFRefErr {
	e.Msg = msg
	return e
}

// WithMsgF sets the message for the CFRefErr using a format string.
func (e *CFRefErr) WithMsgF(format string, a ...any) *CFRefErr {
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

// Error returns the error message for the CFRefErr.
func (e *CFRefErr) Error() string {
	if e.Msg == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
}

// Unwrap returns the underlying error for the CFRefErr.
func (e *CFRefErr) Unwrap() error {
	return e.Err
}

// CFTypeErr represents a CoreFoundation type error.
type CFTypeErr struct {
	Err error
	Msg string
}

// CFTypeError creates a new CFTypeErr from the given error.
func CFTypeError(err error) *CFTypeErr {
	return &CFTypeErr{Err: err}
}

// WithMsg sets the message for the CFTypeErr.
func (e *CFTypeErr) WithMsg(msg string) *CFTypeErr {
	e.Msg = msg
	return e
}

// WithMsgF sets the message for the CFTypeErr using a format string.
func (e *CFTypeErr) WithMsgF(format string, a ...any) *CFTypeErr {
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

// Error returns the error message for the CFTypeErr.
func (e *CFTypeErr) Error() string {
	if e.Msg == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
}

// Unwrap returns the underlying error for the CFTypeErr.
func (e *CFTypeErr) Unwrap() error {
	return e.Err
}

// CFSynvErr represents a CoreFoundation synchronization error.
type CFSyncErr struct {
	Msg string
}

// CFSyncError creates a new CFSyncErr from the given error.
func CFSyncError() *CFSyncErr {
	return &CFSyncErr{}
}

// WithMsg sets the message for the CFSyncErr.
func (e *CFSyncErr) WithMsg(msg string) *CFSyncErr {
	e.Msg = msg
	return e
}

// WithMsgF sets the message for the CFSyncErr using a format string.
func (e *CFSyncErr) WithMsgF(format string, a ...any) *CFSyncErr {
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

// Error returns the error message for the CFSyncErr.
func (e *CFSyncErr) Error() string {
	if e.Msg == "" {
		return "synchronization error"
	}
	return fmt.Sprintf("synchronization error: %s", e.Msg)
}
