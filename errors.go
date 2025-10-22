package cfprefs

import "fmt"

type KeyNotFoundErr struct {
	AppID string
	Key   string
	Msg   string
}

func KeyNotFoundError(appID, key string) *KeyNotFoundErr {
	return &KeyNotFoundErr{AppID: appID, Key: key}
}

func (e *KeyNotFoundErr) WithMsg(msg string) *KeyNotFoundErr {
	e.Msg = msg
	return e
}

func (e *KeyNotFoundErr) WithMsgF(format string, a ...any) *KeyNotFoundErr {
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

func (e *KeyNotFoundErr) Error() string {
	if e.Msg == "" {
		return fmt.Sprintf("key not found: %s [%s]", e.Key, e.AppID)
	}
	return fmt.Sprintf("key not found: %s [%s] - %s", e.Key, e.AppID, e.Msg)
}

type KeyPathErr struct {
	AppID string
	Key   string
	Msg   string
}

func KeyPathError(appID, key string) *KeyPathErr {
	return &KeyPathErr{AppID: appID, Key: key}
}

func (e *KeyPathErr) WithMsg(msg string) *KeyPathErr {
	e.Msg = msg
	return e
}

func (e *KeyPathErr) WithMsgF(format string, a ...any) *KeyPathErr {
	e.Msg = fmt.Sprintf(format, a...)
	return e
}

func (e *KeyPathErr) Error() string {
	if e.Msg == "" {
		return fmt.Sprintf("key path error: %s [%s]", e.Key, e.AppID)
	}
	return fmt.Sprintf("key path error: %s [%s] - %s", e.Key, e.AppID, e.Msg)
}

type TypeMismatchErr struct {
	AppID string
	Key   string
	Type  any
	Value any
}

func TypeMismatchError(appID, key string, expected, actual any) *TypeMismatchErr {
	return &TypeMismatchErr{AppID: appID, Key: key, Type: expected, Value: actual}
}

func (e *TypeMismatchErr) Error() string {
	return fmt.Sprintf("type mismatch: %s [%s] - expected %T, got %T", e.Key, e.AppID, e.Type, e.Value)
}
