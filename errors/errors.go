package errors

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// PublicError exposes error message and call location.
type PublicError interface {
	error
	Public() string
	Location() string
}

// AppError is a structured error with internal/public message and location.
type AppError struct {
	code      string
	msg       string
	publicMsg string
	location  string
	cause     error
}

func (e *AppError) Error() string {
	return e.msg
}

func (e *AppError) Public() string {
	return e.publicMsg
}

func (e *AppError) Location() string {
	return e.location
}

// NewPublicError creates a new AppError.
func New(msg string, publicMsg string, code string) error {
	return &AppError{
		code:      code,
		msg:       msg,
		publicMsg: publicMsg,
		location:  callerLocation(2),
	}
}

// Join wraps an existing error with a new message.
func Join(err error, msg string) error {
	if err == nil {
		return nil
	}

	// If depth(err) > 10 return err to avoid infinite wrapping
	if depth(err) > 10 {
		return err
	}

	publicMsg := ""
	if pub, ok := err.(PublicError); ok {
		publicMsg = pub.Public()
	}

	code := ""
	if ae, ok := err.(*AppError); ok {
		code = ae.code
	}

	return &AppError{
		code:      code,
		msg:       fmt.Sprintf("%s\n\t%s", msg, err.Error()),
		publicMsg: publicMsg,
		location:  callerLocation(2),
		cause:     err,
	}
}

func Upgrade(err error, code string, publicMsg string) error {
	if err == nil {
		return nil
	}

	if ae, ok := err.(*AppError); ok {
		// If it's already an AppError, just update the code and public message
		ae.code = code
		ae.publicMsg = publicMsg
		return ae
	}

	// If it's not an AppError, create a new one
	return &AppError{
		code:      code,
		msg:       err.Error(),
		publicMsg: publicMsg,
		location:  callerLocation(2),
		cause:     err,
	}
}

// Trace returns a multi-line trace of AppError chain.
func Trace(err error) string {
	var trace []string
	level := 0
	for err != nil {
		if len(trace) > 10 {
			trace = append(trace, fmt.Sprintf("... %d more errors", level))
			break
		}
		if ae, ok := err.(*AppError); ok {
			prefix := "!!"
			if level > 0 {
				prefix = "caused by:"
			}
			trace = append(trace, fmt.Sprintf("%s %s\n\tat %s", prefix, ae.msg, ae.location))

			// Point to the next error in the chain
			err = ae.cause
		} else {
			trace = append(trace, fmt.Sprintf("unknown error: %v", err))
			break
		}
		level++
	}
	return strings.Join(trace, "\n")
}

// IsPublicError checks if the error implements PublicError interface.
func IsPublicError(err error) bool {
	var e PublicError
	return errors.As(err, &e)
}

// IsAppError checks if the error is an AppError.
func IsAppError(err error) bool {
	var e *AppError
	return errors.As(err, &e)
}

// Is checks if the error matches a specific code.
func Is(err error, code string) bool {
	if err == nil {
		return false
	}

	var ae *AppError
	if errors.As(err, &ae) {
		return ae.code == code
	}

	return false
}

// Public extracts the public message from an error, if it implements PublicError.
func Public(err error) string {
	if err == nil {
		return ""
	}

	var e PublicError
	if errors.As(err, &e) {
		return e.Public()
	}

	return err.Error()
}

// Code extracts the error code from an AppError.
func Code(err error) string {
	if err == nil {
		return ""
	}

	var ae *AppError
	if errors.As(err, &ae) {
		// return ae.Code()
		return ae.code
	}

	return ""
}

// callerLocation returns the function name and file location of the caller.
func callerLocation(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s (%s:%d)", fn.Name(), filepath.Base(file), line)
}

// depth counts how many AppError layers are in the error chain.
func depth(err error) int {
	count := 0
	for {
		if ae, ok := err.(*AppError); ok {
			err = ae.cause
			count++
		} else {
			break
		}
	}
	return count
}
