package cfprefs

import "fmt"

type KeyNotFoundError struct {
	AppID string
	Key   string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key not found: %s [%s]", e.Key, e.AppID)
}

type QueryTypeMismatchError struct {
	AppID string
	Key   string
	Query string
	Type  any
	Value any
}

func (e *QueryTypeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch: %s [%s] - expected %T, got %T", e.Key, e.AppID, e.Type, e.Value)
}
