package errors_test

import (
	stdErr "errors"
	"testing"

	"github.com/mandacode-com/golib/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewPublicError validates creation of AppError with and without public message.
func TestNewPublicError(t *testing.T) {
	t.Run("creates AppError with public message", func(t *testing.T) {
		err := errors.NewPublicError("internal failure", "visible to user", "ERR_PUBLIC")
		require.Error(t, err)

		publicErr, ok := err.(errors.PublicError)
		require.True(t, ok, "should implement PublicError interface")

		assert.Equal(t, "visible to user", publicErr.Public())
		assert.Contains(t, publicErr.Location(), ".TestNewPublicError")
	})

	t.Run("creates AppError without public message", func(t *testing.T) {
		err := errors.NewPublicError("failure without public", "", "ERR_NO_PUBLIC")
		require.Error(t, err)

		publicErr, ok := err.(errors.PublicError)
		require.True(t, ok)

		assert.Equal(t, "internal error", publicErr.Public())
		assert.Contains(t, publicErr.Location(), "TestNewPublicError")
	})
}

// TestJoin validates the error chaining and message propagation.
func TestJoin(t *testing.T) {
	t.Run("joins with previous AppError", func(t *testing.T) {
		base := errors.NewPublicError("base error", "db failed", "ERR_DB")
		wrapped := errors.Join(base, "service failed")

		require.Error(t, wrapped)
		assert.Contains(t, wrapped.Error(), "service failed")
		assert.Contains(t, wrapped.Error(), "base error")
		assert.Equal(t, "db failed", wrapped.(errors.PublicError).Public())
		assert.Contains(t, wrapped.(errors.PublicError).Location(), "TestJoin")
	})

	t.Run("returns nil when base is nil", func(t *testing.T) {
		joined := errors.Join(nil, "should not wrap nil")
		assert.Nil(t, joined)
	})

	t.Run("fails when joining self referential error", func(t *testing.T) {
		base := errors.NewPublicError("self error", "self visible", "ERR_SELF")
		wrapped := errors.Join(base, "self join")

		require.Error(t, wrapped)
		assert.NotEqual(t, base, wrapped)
		assert.Contains(t, wrapped.Error(), "self join")
		assert.Contains(t, wrapped.Error(), "self error")
		assert.Equal(t, "self visible", wrapped.(errors.PublicError).Public())
		assert.Contains(t, wrapped.(errors.PublicError).Location(), "TestJoin")
	})
}

// TestIsHelpers verifies type check helpers for AppError and PublicError.
func TestIsHelpers(t *testing.T) {
	t.Run("identifies PublicError and AppError", func(t *testing.T) {
		err := errors.NewPublicError("fail", "visible", "ERR_CODE")
		assert.True(t, errors.IsAppError(err))
		assert.True(t, errors.IsPublicError(err))
	})

	t.Run("returns false on standard error", func(t *testing.T) {
		plain := stdErr.New("standard error")
		assert.False(t, errors.IsAppError(plain))
		assert.False(t, errors.IsPublicError(plain))
	})

	t.Run("returns false on nil", func(t *testing.T) {
		assert.False(t, errors.IsAppError(nil))
		assert.False(t, errors.IsPublicError(nil))
	})
}

// TestTrace ensures Trace prints full chain of error messages and locations.
func TestTrace(t *testing.T) {
	t.Run("prints error trace in correct format", func(t *testing.T) {
		base := errors.NewPublicError("db read failed", "try again later", "ERR_DB_READ")
		lvl2 := errors.Join(base, "repository error")
		lvl3 := errors.Join(lvl2, "usecase failed")

		trace := errors.Trace(lvl3)

		assert.Contains(t, trace, "usecase failed")
		assert.Contains(t, trace, "repository error")
		assert.Contains(t, trace, "db read failed")
		assert.Contains(t, trace, "caused by")

		t.Logf("\n%s", trace)
	})
}
