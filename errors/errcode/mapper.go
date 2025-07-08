package errcode

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// MapCodeToHTTP maps internal error codes to HTTP status codes.
func MapCodeToHTTP(code string) int {
	switch code {
	case
		ErrInvalidInput,
		ErrMissingRequiredField,
		ErrInvalidFormat,
		ErrTooLarge:
		return http.StatusBadRequest

	case ErrTooManyRequests:
		return http.StatusTooManyRequests

	case ErrUnauthorized, ErrInvalidToken, ErrTokenExpired:
		return http.StatusUnauthorized

	case ErrForbidden:
		return http.StatusForbidden

	case ErrNotFound:
		return http.StatusNotFound

	case ErrAlreadyExists, ErrConflict:
		return http.StatusConflict

	case ErrUserNotVerified,
		ErrAccountDisabled,
		ErrInsufficientBalance:
		return http.StatusForbidden // or 422 Unprocessable Entity

	case ErrDependencyFailure, ErrTimeout:
		return http.StatusFailedDependency

	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable

	case ErrInternalFailure:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}

// MapCodeToGRPC maps internal error codes to gRPC status codes.
func MapCodeToGRPC(code string) codes.Code {
	switch code {
	case
		ErrInvalidInput,
		ErrMissingRequiredField,
		ErrInvalidFormat,
		ErrTooLarge:
		return codes.InvalidArgument

	case ErrTooManyRequests:
		return codes.ResourceExhausted

	case ErrUnauthorized, ErrInvalidToken, ErrTokenExpired:
		return codes.Unauthenticated

	case ErrForbidden,
		ErrUserNotVerified,
		ErrAccountDisabled,
		ErrInsufficientBalance:
		return codes.PermissionDenied

	case ErrNotFound:
		return codes.NotFound

	case ErrAlreadyExists, ErrConflict:
		return codes.AlreadyExists

	case ErrDependencyFailure:
		return codes.FailedPrecondition

	case ErrTimeout:
		return codes.DeadlineExceeded

	case ErrServiceUnavailable:
		return codes.Unavailable

	case ErrInternalFailure:
		fallthrough
	default:
		return codes.Internal
	}
}
