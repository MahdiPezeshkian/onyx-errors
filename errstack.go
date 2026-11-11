// Package errstack provides structured application/domain errors with
// independently formattable stack traces. It stays free of HTTP and transport
// concerns.
package errstack

import (
	"errors"
	"fmt"
)

// New creates an error with message and a stack trace captured at the call site.
func New(msg string) error {
	return &Error{
		Message:    msg,
		StackTrace: captureStack(3),
	}
}

// Newf creates a formatted error with a stack trace captured at the call site.
func Newf(format string, args ...interface{}) error {
	return &Error{
		Message:    fmt.Sprintf(format, args...),
		StackTrace: captureStack(3),
	}
}

// NewCode creates an error with a stable machine-readable code, message, and
// stack trace captured at the call site.
func NewCode(code, msg string) error {
	return &Error{
		Code:       code,
		Message:    msg,
		StackTrace: captureStack(3),
	}
}

// NewCodef creates a coded error with a formatted message and stack trace.
func NewCodef(code, format string, args ...interface{}) error {
	return &Error{
		Code:       code,
		Message:    fmt.Sprintf(format, args...),
		StackTrace: captureStack(3),
	}
}

// Wrap annotates err with a prefix message.
// If err (or a cause) already carries an errstack stack trace, that original
// stack is preserved and no new stack is captured. Otherwise a stack is
// captured at the Wrap call site.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return wrap(err, fmt.Sprintf("%s: %v", msg, err), 3)
}

// Wrapf annotates err with a formatted prefix message.
// Stack-trace preservation rules match Wrap.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	prefix := fmt.Sprintf(format, args...)
	return wrap(err, fmt.Sprintf("%s: %v", prefix, err), 3)
}

// WithStack ensures err has a stack trace. If err already has one (directly or
// via an errstack cause), it is returned unchanged. Otherwise a new Error is
// created that wraps err and captures a stack at this call site.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if len(stackOf(err)) > 0 {
		return err
	}
	return &Error{
		Message:    err.Error(),
		cause:      err,
		StackTrace: captureStack(3),
	}
}

func wrap(err error, message string, skip int) *Error {
	e := &Error{
		Message: message,
		cause:   err,
	}
	if st := stackOf(err); len(st) > 0 {
		e.StackTrace = st.Clone()
		if ce, ok := err.(*Error); ok {
			e.Code = ce.Code
		}
	} else {
		e.StackTrace = captureStack(skip + 1)
	}
	return e
}

// Cause returns the underlying cause by repeatedly unwrapping.
func Cause(err error) error {
	for err != nil {
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return err
		}
		next := u.Unwrap()
		if next == nil {
			return err
		}
		err = next
	}
	return nil
}

// Stack returns the human-readable stack trace for err when available.
// For non-errstack errors without a stack, it falls back to err.Error().
func Stack(err error) string {
	if err == nil {
		return ""
	}
	if st := stackOf(err); len(st) > 0 {
		return st.String()
	}
	return err.Error()
}

// ErrorWithStack returns the error message followed by its stack trace text.
func ErrorWithStack(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(*Error); ok {
		return e.FormatText()
	}
	st := stackOf(err)
	if len(st) == 0 {
		return err.Error()
	}
	return err.Error() + "\n" + st.String()
}

// Trace returns the structured StackTrace owned by err, if any.
func Trace(err error) StackTrace {
	return stackOf(err).Clone()
}

// Is reports whether err matches target using the standard library.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target using the
// standard library.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Must panics with a full text formatting when err is non-nil.
func Must(err error) {
	if err != nil {
		panic(ErrorWithStack(err))
	}
}

func stackOf(err error) StackTrace {
	for err != nil {
		if e, ok := err.(*Error); ok && len(e.StackTrace) > 0 {
			return e.StackTrace
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return nil
		}
		err = u.Unwrap()
	}
	return nil
}
