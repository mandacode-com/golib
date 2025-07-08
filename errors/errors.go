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
	msg       string
	publicMsg string
	location  string
	cause     error
}

func (e *AppError) Error() string {
	return e.msg
}

func (e *AppError) Public() string {
	if e.publicMsg != "" {
		return e.publicMsg
	}
	return "internal error"
}

func (e *AppError) Location() string {
	return e.location
}

func (e *AppError) Unwrap() error {
	return e.cause
}

// NewPublicError creates a new AppError.
func NewPublicError(msg string, publicMsg ...string) error {
	return &AppError{
		msg:       msg,
		publicMsg: first(publicMsg),
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

	// prevent nesting if message already contains same content
	if ae, ok := err.(*AppError); ok && ae.msg == msg {
		return err // skip wrapping
	}

	publicMsg := "internal error"
	if pub, ok := err.(PublicError); ok {
		publicMsg = pub.Public()
	}

	return &AppError{
		msg:       fmt.Sprintf("%s\n\t%s", msg, err.Error()),
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
			trace = append(trace, fmt.Sprintf("↳ unknown error: %v", err))
			break
		}
		level++
	}
	return strings.Join(trace, "\n")
}

func IsPublicError(err error) bool {
	var e PublicError
	return errors.As(err, &e)
}

func IsAppError(err error) bool {
	var e *AppError
	return errors.As(err, &e)
}

func callerLocation(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s (%s:%d)", fn.Name(), filepath.Base(file), line)
}

func first(msgs []string) string {
	if len(msgs) > 0 {
		return msgs[0]
	}
	return ""
}

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
