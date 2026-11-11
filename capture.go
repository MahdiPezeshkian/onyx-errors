package errstack

import (
	"runtime"
	"strings"
)

const (
	defaultStackDepth = 64
	modulePathPrefix  = "github.com/MahdiPezeshkian/onyx-errors"
)

// captureStack records the current call stack.
// skip is passed to runtime.Callers: 0 identifies Callers itself, so callers
// typically pass 2+ to omit captureStack and their immediate constructor.
func captureStack(skip int) StackTrace {
	pcs := make([]uintptr, defaultStackDepth)
	n := runtime.Callers(skip, pcs)
	if n == 0 {
		return nil
	}
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)
	out := make(StackTrace, 0, n)
	for {
		f, more := frames.Next()
		if f.Function != "" || f.File != "" || f.Line != 0 {
			if !isInternalFrame(f.Function) {
				out = append(out, Frame{
					Function: f.Function,
					File:     f.File,
					Line:     f.Line,
				})
			}
		}
		if !more {
			break
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func isInternalFrame(function string) bool {
	if function == "" {
		return false
	}
	// Skip frames belonging to this module's implementation, including
	// methods on types defined here (e.g. (*Error).Format).
	if strings.HasPrefix(function, modulePathPrefix+".") ||
		strings.HasPrefix(function, modulePathPrefix+"/") {
		return true
	}
	return false
}
