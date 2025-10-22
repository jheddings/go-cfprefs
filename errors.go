package cfprefs

import "fmt"

type KeyNotFoundError struct {
	AppID string
	Key   string
	Msg   string
}

func (e *KeyNotFoundError) Error() string {
	if e.Msg == "" {
		return fmt.Sprintf("key not found: %s [%s]", e.Key, e.AppID)
	}

	return fmt.Sprintf("key not found: %s [%s] - %s", e.Key, e.AppID, e.Msg)
}

type TypeMismatchError struct {
	AppID string
	Key   string
	Type  any
	Value any
}

func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch: %s [%s] - expected %T, got %T", e.Key, e.AppID, e.Type, e.Value)
}
