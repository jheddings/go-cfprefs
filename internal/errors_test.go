package internal

import (
	"errors"
	"testing"
)

func TestCFErrors(t *testing.T) {
	t.Run("NewCFError", func(t *testing.T) {
		err := NewCFError("test operation", nil)

		// Should default to ErrCFOperation
		if !errors.Is(err, ErrCFOperation) {
			t.Errorf("expected NewCFError with nil to use ErrCFOperation")
		}

		expected := "test operation"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}

		customErr := errors.New("custom error")
		err = NewCFError("operation", customErr)
		expected = "operation: custom error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}

		if !errors.Is(err, ErrCFOperation) {
			t.Errorf("expected Is(ErrCFOperation) to be true")
		}
		if !errors.Is(err, customErr) {
			t.Errorf("expected Is(customErr) to be true")
		}
	})

	t.Run("CFError with message", func(t *testing.T) {
		err := NewCFError("read", nil).WithMsg("failed to read preference")
		expected := "read: failed to read preference"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}

		err = NewCFError("write", nil).WithMsgF("key %q not found", "test-key")
		expected = "write: key \"test-key\" not found"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("CFRefError", func(t *testing.T) {
		baseErr := errors.New("allocation failed")
		err := CFRefError().Wrap(baseErr)

		if err.Op != "reference" {
			t.Errorf("expected Op to be 'reference', got %q", err.Op)
		}

		if !errors.Is(err, baseErr) {
			t.Errorf("expected to wrap base error")
		}

		err = CFRefError().Wrap(nil)
		if !errors.Is(err, ErrCFReference) {
			t.Errorf("expected nil to use ErrCFReference")
		}
	})

	t.Run("CFTypeError", func(t *testing.T) {
		err := CFTypeError().Wrap(nil).WithMsg("unsupported type")

		// Should have correct operation
		if err.Op != "type conversion" {
			t.Errorf("expected Op to be 'type conversion', got %q", err.Op)
		}

		// Should match ErrCFType
		if !errors.Is(err, ErrCFType) {
			t.Errorf("expected to match ErrCFType")
		}

		expected := "type conversion: unsupported type: CoreFoundation type error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("CFSyncError", func(t *testing.T) {
		err := CFSyncError()

		// Should match ErrCFSync
		if !errors.Is(err, ErrCFSync) {
			t.Errorf("expected to match ErrCFSync")
		}

		// Test with message
		err = err.WithMsg("failed to sync preferences")
		expected := "synchronization: failed to sync preferences: CoreFoundation synchronization error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("CFMemoryError", func(t *testing.T) {
		err := CFMemoryError("allocation failed")

		// Should match ErrCFMemory
		if !errors.Is(err, ErrCFMemory) {
			t.Errorf("expected to match ErrCFMemory")
		}

		expected := "memory: allocation failed: CoreFoundation memory error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})
}

func TestErrorWrapping(t *testing.T) {
	t.Run("Unwrap chain", func(t *testing.T) {
		innerErr := errors.New("inner error")
		cfErr := NewCFError("operation", innerErr)

		// Direct unwrap
		if errors.Unwrap(cfErr) != innerErr {
			t.Errorf("expected Unwrap to return inner error")
		}

		// errors.Is should work through the chain
		if !errors.Is(cfErr, innerErr) {
			t.Errorf("expected Is to find inner error")
		}
	})

	t.Run("Multiple sentinel matches", func(t *testing.T) {
		// Create error that matches multiple sentinels
		err := CFTypeError().Wrap(ErrCFMemory).WithMsg("type conversion failed due to memory")

		// Should match both sentinels
		if !errors.Is(err, ErrCFType) {
			t.Errorf("expected to match ErrCFType")
		}
		if !errors.Is(err, ErrCFMemory) {
			t.Errorf("expected to match ErrCFMemory")
		}
	})
}
