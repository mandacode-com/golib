package errors_test

import (
	"testing"

	"github.com/mandacode-com/golib/errors"
)

func TestNewPublicError(t *testing.T) {
	t.Run("with public message", func(t *testing.T) {
		err := errors.NewPublicError("test error", "public message")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		publicErr, ok := err.(errors.PublicError)
		if !ok {
			t.Fatal("expected PublicError type")
		}
		if publicErr.Public() != "public message" {
			t.Errorf("expected public message 'public message', got '%s'", publicErr.Public())
		}
	})

	t.Run("without public message", func(t *testing.T) {
		err := errors.NewPublicError("test error")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		publicErr, ok := err.(errors.PublicError)
		if !ok {
			t.Fatal("expected PublicError type")
		}

		if publicErr.Public() != "internal error" {
			t.Errorf("expected public message 'internal error', got '%s'", publicErr.Public())
		}
	})

	t.Run("with empty message", func(t *testing.T) {
		err := errors.NewPublicError("")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		publicErr, ok := err.(errors.PublicError)
		if !ok {
			t.Fatal("expected PublicError type")
		}

		if publicErr.Public() != "internal error" {
			t.Errorf("expected public message 'internal error', got '%s'", publicErr.Public())
		}
	})

	t.Run("with nil message", func(t *testing.T) {
		err := errors.NewPublicError("")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		publicErr, ok := err.(errors.PublicError)
		if !ok {
			t.Fatal("expected PublicError type")
		}

		if publicErr.Public() != "internal error" {
			t.Errorf("expected public message 'internal error', got '%s'", publicErr.Public())
		}
	})
}

func TestJoin(t *testing.T) {
	t.Run("join with PublicError", func(t *testing.T) {
		err1 := errors.NewPublicError("first error", "public first error")
		err2 := errors.Join(err1, "second error")

		if err2 == nil {
			t.Fatal("expected error, got nil")
		}

		publicErr, ok := err2.(errors.PublicError)
		if !ok {
			t.Fatal("expected PublicError type")
		}

		if publicErr.Public() != "public first error" {
			t.Errorf("expected public message 'public first error', got '%s'", publicErr.Public())
		}
	})

	t.Run("join with non-PublicError", func(t *testing.T) {
		err1 := errors.NewPublicError("first error")
		err2 := errors.Join(err1, "second error")

		if err2 == nil {
			t.Fatal("expected error, got nil")
		}

		if err2.Error() != "second error\n\tfirst error" {
			t.Errorf("expected 'second error: first error', got '%s'", err2.Error())
		}
	})

	t.Run("join with nil", func(t *testing.T) {
		err2 := errors.Join(nil, "second error")
		if err2 != nil {
			t.Fatal("expected nil, got an error")
		}
	})
}

func TestIsPublicError(t *testing.T) {
	t.Run("is PublicError", func(t *testing.T) {
		err := errors.NewPublicError("test error", "public message")
		if !errors.IsPublicError(err) {
			t.Fatal("expected true, got false")
		}
	})

	t.Run("is not PublicError", func(t *testing.T) {
		err := errors.NewPublicError("test error")
		if !errors.IsPublicError(err) {
			t.Fatal("expected false, got true")
		}
	})

	t.Run("nil error", func(t *testing.T) {
		if errors.IsPublicError(nil) {
			t.Fatal("expected false, got true")
		}
	})
}

func TestIsAppError(t *testing.T) {
	t.Run("is AppError", func(t *testing.T) {
		err := errors.NewPublicError("test error", "public message")
		if !errors.IsAppError(err) {
			t.Fatal("expected true, got false")
		}
	})

	t.Run("is not AppError", func(t *testing.T) {
		err := errors.NewPublicError("test error")
		if !errors.IsAppError(err) {
			t.Fatal("expected true, got false")
		}
	})

	t.Run("nil error", func(t *testing.T) {
		if errors.IsAppError(nil) {
			t.Fatal("expected false, got true")
		}
	})
}
