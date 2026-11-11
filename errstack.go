// Package errstack provides error construction and wrapping helpers with
// stack traces, built on github.com/pkg/errors. It stays free of HTTP and
// transport concerns.
package errstack

import (
	"fmt"

	"github.com/pkg/errors"
)

func New(msg string) error {
	return errors.New(msg)
}

func Newf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, msg)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, format, args...)
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return errors.WithStack(err)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Stack(err error) string {
	if err == nil {
		return ""
	}
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	if st, ok := err.(stackTracer); ok {
		return fmt.Sprintf("%+v", st.StackTrace())
	}
	return err.Error()
}

func ErrorWithStack(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%v\n%s", err, Stack(err))
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Must(err error) {
	if err != nil {
		panic(ErrorWithStack(err))
	}
}
