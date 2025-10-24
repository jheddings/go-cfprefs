package cfprefs

import (
	"errors"
	"testing"

	"github.com/jheddings/go-cfprefs/testutil"
)

func TestSentinelErrors(t *testing.T) {
	t.Run("KeyNotFoundErr", func(t *testing.T) {
		err := NewKeyNotFoundError("com.test.app", "missing-key")

		if !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("expected errors.Is(err, ErrKeyNotFound) to be true")
		}

		if unwrapped := errors.Unwrap(err); unwrapped != ErrKeyNotFound {
			t.Errorf("expected Unwrap to return ErrKeyNotFound, got %v", unwrapped)
		}

		var knfErr *KeyNotFoundErr
		if !errors.As(err, &knfErr) {
			t.Errorf("expected errors.As to work with *KeyNotFoundErr")
		}

		expected := "key not found: missing-key [com.test.app]"
		if err.Error() != expected {
			t.Errorf("expected error message %q, got %q", expected, err.Error())
		}

		err = err.WithMsg("custom context")
		expected = "key not found: missing-key [com.test.app] - custom context"
		if err.Error() != expected {
			t.Errorf("expected error message %q, got %q", expected, err.Error())
		}
	})

	t.Run("KeyPathErr", func(t *testing.T) {
		err := NewKeyPathError().Wrap(errors.New("invalid/path"))

		if !errors.Is(err, ErrInvalidKeyPath) {
			t.Errorf("expected errors.Is(err, ErrInvalidKeyPath) to be true")
		}

		if unwrapped := errors.Unwrap(err); unwrapped != ErrInvalidKeyPath {
			t.Errorf("expected Unwrap to return ErrInvalidKeyPath, got %v", unwrapped)
		}

		var kpErr *KeyPathErr
		if !errors.As(err, &kpErr) {
			t.Errorf("expected errors.As to work with *KeyPathErr")
		}

		err = err.WithMsgF("error at segment %d", 2)
		expected := "error at segment 2: invalid key path\ninvalid/path"
		if err.Error() != expected {
			t.Errorf("expected error message %q, got %q", expected, err.Error())
		}
	})

	t.Run("TypeMismatchErr", func(t *testing.T) {
		err := NewTypeMismatchError(int64(0), "string value").WithKey("com.test.app", "type-key")

		if !errors.Is(err, ErrTypeMismatch) {
			t.Errorf("expected errors.Is(err, ErrTypeMismatch) to be true")
		}

		if unwrapped := errors.Unwrap(err); unwrapped != ErrTypeMismatch {
			t.Errorf("expected Unwrap to return ErrTypeMismatch, got %v", unwrapped)
		}

		var tmErr *TypeMismatchErr
		if !errors.As(err, &tmErr) {
			t.Errorf("expected errors.As to work with *TypeMismatchErr")
		}

		expected := "type mismatch: type-key [com.test.app] - expected int64, got string"
		if err.Error() != expected {
			t.Errorf("expected error message %q, got %q", expected, err.Error())
		}

		if tmErr.AppID != "com.test.app" || tmErr.Key != "type-key" {
			t.Errorf("expected fields to be accessible")
		}
	})
}

func TestErrorChaining(t *testing.T) {
	baseErr := NewKeyNotFoundError("com.test.app", "key")
	wrappedErr := errors.Join(errors.New("operation failed"), baseErr)

	if !errors.Is(wrappedErr, ErrKeyNotFound) {
		t.Errorf("expected wrapped error to still match ErrKeyNotFound")
	}

	var knfErr *KeyNotFoundErr
	if !errors.As(wrappedErr, &knfErr) {
		t.Errorf("expected errors.As to work through wrapping")
	}
}

func TestRealWorldErrors(t *testing.T) {
	appID := "com.jheddings.cfprefs.testing.errors"

	t.Run("GetMissingKey", func(t *testing.T) {
		// Ensure key doesn't exist
		Delete(appID, "missing-key")

		_, err := Get(appID, "missing-key")
		testutil.AssertError(t, err, "getting missing key")

		// Should be able to check with errors.Is
		if !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("expected error to match ErrKeyNotFound sentinel")
		}

		// Should be able to extract details
		var knfErr *KeyNotFoundErr
		if errors.As(err, &knfErr) {
			if knfErr.AppID != appID || knfErr.Key != "missing-key" {
				t.Errorf("expected error details to match")
			}
		} else {
			t.Errorf("expected to extract KeyNotFoundErr")
		}
	})

	t.Run("TypeMismatch", func(t *testing.T) {
		// Set a string value
		err := Set(appID, "string-value", "hello")
		testutil.AssertNoError(t, err, "setting string value")
		defer Delete(appID, "string-value")

		// Try to get it as int
		_, err = GetInt(appID, "string-value")
		testutil.AssertError(t, err, "getting string as int")

		// Should match type mismatch
		if !errors.Is(err, ErrTypeMismatch) {
			t.Errorf("expected error to match ErrTypeMismatch sentinel")
		}

		// Should be able to extract details
		var tmErr *TypeMismatchErr
		if errors.As(err, &tmErr) {
			if tmErr.AppID != appID || tmErr.Key != "string-value" {
				t.Errorf("expected error details to match")
			}
			// Check that Expected and Actual are set
			if tmErr.Expected == nil || tmErr.Actual == nil {
				t.Errorf("expected both Expected and Actual to be set")
			}
		} else {
			t.Errorf("expected to extract TypeMismatchErr")
		}
	})

	t.Run("InvalidKeyPath", func(t *testing.T) {
		// Set a non-dict value
		err := Set(appID, "not-a-dict", "string value")
		testutil.AssertNoError(t, err, "setting string value")
		defer Delete(appID, "not-a-dict")

		// Try to access through it
		_, err = Get(appID, "not-a-dict/nested")
		testutil.AssertError(t, err, "accessing through non-dict")

		// Should be able to detect it's a key not found error
		if !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("expected error to match ErrKeyNotFound sentinel")
		}
	})

	t.Run("SetInvalidPath", func(t *testing.T) {
		// Try to set through a non-object
		err := Set(appID, "scalar-value", "not an object")
		testutil.AssertNoError(t, err, "setting scalar value")
		defer Delete(appID, "scalar-value")

		// Try to set a nested field on a scalar
		err = SetQ(appID, "scalar-value", "$.nested.field", "value")
		testutil.AssertError(t, err, "setting through scalar")

		// Should be a key path error
		if !errors.Is(err, ErrInvalidKeyPath) {
			t.Errorf("expected error to match ErrInvalidKeyPath sentinel")
		}
	})
}
