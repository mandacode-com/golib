package errors_test

import (
	stdErr "errors"
	"testing"

	"github.com/mandacode-com/golib/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewPublicError validates creation of AppError with and without public message.
func TestNew(t *testing.T) {
	t.Run("creates AppError with public message", func(t *testing.T) {
		err := errors.New("internal failure", "visible to user", "ERR_PUBLIC")
		require.Error(t, err)

		publicErr, ok := err.(errors.PublicError)
		require.True(t, ok, "should implement PublicError interface")

		assert.Equal(t, "visible to user", publicErr.Public())
		assert.Contains(t, publicErr.Location(), ".TestNew")
	})

	t.Run("creates AppError without public message", func(t *testing.T) {
		err := errors.New("failure without public", "", "ERR_NO_PUBLIC")
		require.Error(t, err)

		publicErr, ok := err.(errors.PublicError)
		require.True(t, ok)

		assert.Equal(t, "", publicErr.Public())
		assert.Contains(t, publicErr.Location(), "TestNew")
	})
}

// TestJoin validates the error chaining and message propagation.
func TestJoin(t *testing.T) {
	t.Run("joins with previous AppError", func(t *testing.T) {
		base := errors.New("base error", "db failed", "ERR_DB")
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

	t.Run("joins with standard error", func(t *testing.T) {
		stdErr := stdErr.New("standard error")
		wrapped := errors.Join(stdErr, "additional context")

		require.Error(t, wrapped)
		assert.Contains(t, wrapped.Error(), "additional context")
		assert.Contains(t, wrapped.Error(), "standard error")
		assert.Equal(t, "", wrapped.(errors.PublicError).Public())
		assert.Contains(t, wrapped.(errors.PublicError).Location(), "TestJoin")
	})
}

// TestIsHelpers verifies type check helpers for AppError and PublicError.
func TestIsHelpers(t *testing.T) {
	t.Run("identifies PublicError and AppError", func(t *testing.T) {
		err := errors.New("fail", "visible", "ERR_CODE")
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
		base := errors.New("db read failed", "try again later", "ERR_DB_READ")
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

// TestUpgrade checks that Upgrade modifies existing AppError or creates a new one.
func TestUpgrade(t *testing.T) {
	t.Run("upgrades existing AppError", func(t *testing.T) {
		base := errors.New("initial error", "initial visible", "ERR_INIT")
		updated := errors.Upgrade(base, "ERR_UPDATED", "updated visible")

		require.Error(t, updated)
		assert.Equal(t, "ERR_UPDATED", errors.Code(updated))
		assert.Equal(t, "updated visible", updated.(errors.PublicError).Public())
		assert.Contains(t, updated.(errors.PublicError).Location(), "TestUpgrade")
	})

	t.Run("creates new AppError from standard error", func(t *testing.T) {
		stdErr := stdErr.New("standard error")
		upgraded := errors.Upgrade(stdErr, "ERR_STD_UPGRADED", "visible from std")

		require.Error(t, upgraded)
		assert.Equal(t, "ERR_STD_UPGRADED", errors.Code(upgraded))
		assert.Equal(t, "visible from std", upgraded.(errors.PublicError).Public())
		assert.Contains(t, upgraded.(errors.PublicError).Location(), "TestUpgrade")
	})
}

// TestIs checks if the error matches a specific code.
func TestIs(t *testing.T) {
	t.Run("matches AppError by code", func(t *testing.T) {
		err := errors.New("test error", "visible", "ERR_TEST")
		assert.True(t, errors.Is(err, "ERR_TEST"))
		assert.False(t, errors.Is(err, "ERR_NOT_FOUND"))
	})

	t.Run("returns false for nil error", func(t *testing.T) {
		assert.False(t, errors.Is(nil, "ERR_ANY"))
	})

	t.Run("returns false for standard error", func(t *testing.T) {
		stdErr := stdErr.New("standard error")
		assert.False(t, errors.Is(stdErr, "ERR_ANY"))
	})
}

// TestPublic extracts the public message from an error.
func TestPublic(t *testing.T) {
	t.Run("returns public message from AppError", func(t *testing.T) {
		err := errors.New("test error", "visible to user", "ERR_TEST")
		assert.Equal(t, "visible to user", errors.Public(err))
	})

	t.Run("returns error() string for standard error", func(t *testing.T) {
		stdErr := stdErr.New("standard error")
		assert.Equal(t, "standard error", errors.Public(stdErr))
	})

	t.Run("returns empty string for nil error", func(t *testing.T) {
		assert.Equal(t, "", errors.Public(nil))
	})
}

// TestCode extracts the error code from an AppError.
func TestCode(t *testing.T) {
	t.Run("returns code from AppError", func(t *testing.T) {
		err := errors.New("test error", "visible to user", "ERR_TEST")
		assert.Equal(t, "ERR_TEST", errors.Code(err))
	})

	t.Run("returns empty string for standard error", func(t *testing.T) {
		stdErr := stdErr.New("standard error")
		assert.Equal(t, "", errors.Code(stdErr))
	})

	t.Run("returns empty string for nil error", func(t *testing.T) {
		assert.Equal(t, "", errors.Code(nil))
	})
}
