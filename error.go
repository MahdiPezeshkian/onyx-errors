package errstack

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Error is a structured application/domain error that owns an optional
// StackTrace captured at creation time. It does not carry HTTP status codes;
// transport layers remain responsible for status mapping.
//
// JSON serialization exposes code, message, and stack_trace only. The wrapped
// cause is available through Unwrap / errors.Is / errors.As and is not
// serialized.
type Error struct {
	Code       string     `json:"code,omitempty"`
	Message    string     `json:"message"`
	StackTrace StackTrace `json:"stack_trace,omitempty"`

	cause error
}

// Error implements the error interface. It returns the message only and never
// includes the stack trace.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// Unwrap returns the wrapped cause, if any.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// FormatStackTrace returns only the human-readable stack trace text.
// It does not mutate the error.
func (e *Error) FormatStackTrace() string {
	if e == nil {
		return ""
	}
	return e.StackTrace.String()
}

// FormatText returns a human-readable representation of the error metadata
// followed by its stack trace (when present).
func (e *Error) FormatText() string {
	if e == nil {
		return ""
	}
	var b strings.Builder
	if e.Code != "" {
		b.WriteString(e.Code)
		b.WriteString(": ")
	}
	b.WriteString(e.Message)
	if st := e.FormatStackTrace(); st != "" {
		b.WriteByte('\n')
		b.WriteString(st)
	}
	return b.String()
}

// Format implements fmt.Formatter.
//
//	%v / %s  → Error() message only
//	%+v      → FormatText() including stack trace
func (e *Error) Format(s fmt.State, verb rune) {
	if e == nil {
		return
	}
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.FormatText())
			return
		}
		fallthrough
	case 's', 'q':
		if verb == 'q' {
			_, _ = fmt.Fprintf(s, "%q", e.Error())
			return
		}
		_, _ = io.WriteString(s, e.Error())
	}
}

// MarshalJSON ensures a nil receiver marshals as null and that an empty stack
// is omitted. Exported fields already use appropriate json tags.
func (e *Error) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	type wire struct {
		Code       string     `json:"code,omitempty"`
		Message    string     `json:"message"`
		StackTrace StackTrace `json:"stack_trace,omitempty"`
	}
	return json.Marshal(wire{
		Code:       e.Code,
		Message:    e.Message,
		StackTrace: e.StackTrace,
	})
}
