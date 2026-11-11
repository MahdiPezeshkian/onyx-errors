package errstack

import (
	"strconv"
	"strings"
)

// Frame is a single program counter resolved to a function, file, and line.
type Frame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// StackTrace is an ordered list of stack frames, captured at error creation.
// It is independently usable and does not require an Error value.
type StackTrace []Frame

// String formats the stack trace in the library's human-readable layout:
//
//	function.name
//	    /path/to/file.go:LINE
//
// Consecutive frames are separated by one blank line. Nil and empty traces
// format as an empty string. Formatting never mutates the receiver.
func (s StackTrace) String() string {
	if len(s) == 0 {
		return ""
	}
	var b strings.Builder
	for i, f := range s {
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(f.Function)
		b.WriteString("\n    ")
		b.WriteString(f.File)
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(f.Line))
	}
	return b.String()
}

// Clone returns a shallow copy of the frame slice so callers can share stacks
// without risk of accidental append aliasing.
func (s StackTrace) Clone() StackTrace {
	if s == nil {
		return nil
	}
	out := make(StackTrace, len(s))
	copy(out, s)
	return out
}
