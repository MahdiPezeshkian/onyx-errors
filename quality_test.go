package errstack_test

import (
	"errors"
	"testing"

	"github.com/MahdiPezeshkian/onyx-errors"
)

func TestCodePreservedAcrossNestedWraps(t *testing.T) {
	root := errstack.NewCode("E_ROOT", "root")
	mid := errstack.Wrap(root, "mid")
	outer := errstack.Wrapf(mid, "outer %d", 1)

	if got := errstack.CodeOf(outer); got != "E_ROOT" {
		t.Fatalf("CodeOf(outer)=%q", got)
	}
	var e *errstack.Error
	if !errors.As(outer, &e) || e.Code != "E_ROOT" {
		t.Fatalf("outer code=%q", e.Code)
	}
	// Intermediate wrap of a coded *Error also keeps the code.
	if !errors.As(mid, &e) || e.Code != "E_ROOT" {
		t.Fatalf("mid code=%q", e.Code)
	}
}

func TestCodeOfWalksNonErrstackWrappers(t *testing.T) {
	root := errstack.NewCode("X", "msg")
	wrapped := &wrapper{err: root}
	if got := errstack.CodeOf(wrapped); got != "X" {
		t.Fatalf("CodeOf=%q", got)
	}
}

type wrapper struct{ err error }

func (w *wrapper) Error() string { return "wrapper: " + w.err.Error() }
func (w *wrapper) Unwrap() error { return w.err }

func TestMustPreservesErrorIdentity(t *testing.T) {
	root := errstack.NewCode("PANIC", "boom")
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		err, ok := r.(error)
		if !ok {
			t.Fatalf("panic value type %T", r)
		}
		if !errors.Is(err, root) {
			t.Fatalf("recovered error is not original: %v", err)
		}
		var e *errstack.Error
		if !errors.As(err, &e) || e.Code != "PANIC" {
			t.Fatalf("As failed: %#v", err)
		}
	}()
	errstack.Must(root)
}

func TestErrorsJoinCompatibility(t *testing.T) {
	a := errstack.NewCode("A", "one")
	b := errors.New("two")
	joined := errors.Join(a, b)
	if !errors.Is(joined, a) {
		t.Fatal("Join must match errstack error")
	}
	if !errors.Is(joined, b) {
		t.Fatal("Join must match plain error")
	}
	var e *errstack.Error
	if !errors.As(joined, &e) || e.Code != "A" {
		t.Fatalf("As through Join: %#v", e)
	}
}

func TestNilAPIs(t *testing.T) {
	if errstack.Wrap(nil, "x") != nil || errstack.Wrapf(nil, "x") != nil || errstack.WithStack(nil) != nil {
		t.Fatal("nil wrap helpers")
	}
	if errstack.Stack(nil) != "" || errstack.ErrorWithStack(nil) != "" || errstack.CodeOf(nil) != "" {
		t.Fatal("nil inspectors")
	}
	if errstack.Trace(nil) != nil {
		t.Fatal("Trace(nil)")
	}
	errstack.Must(nil)
}
