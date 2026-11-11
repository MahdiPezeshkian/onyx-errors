package errstack_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/MahdiPezeshkian/onyx-errors"
)

func TestNewAndNewf(t *testing.T) {
	err := errstack.New("boom")
	if err == nil || err.Error() != "boom" {
		t.Fatalf("New: %v", err)
	}
	err = errstack.Newf("code=%d", 42)
	if err == nil || err.Error() != "code=42" {
		t.Fatalf("Newf: %v", err)
	}
}

func TestWrapNil(t *testing.T) {
	if errstack.Wrap(nil, "x") != nil {
		t.Fatal("Wrap(nil) must return nil")
	}
	if errstack.Wrapf(nil, "x %d", 1) != nil {
		t.Fatal("Wrapf(nil) must return nil")
	}
	if errstack.WithStack(nil) != nil {
		t.Fatal("WithStack(nil) must return nil")
	}
}

func TestWrapAndCause(t *testing.T) {
	root := errors.New("root")
	wrapped := errstack.Wrap(root, "layer")
	if errstack.Cause(wrapped) != root {
		t.Fatalf("Cause: got %v", errstack.Cause(wrapped))
	}
	if !errstack.Is(wrapped, root) {
		t.Fatal("Is must match root")
	}
}

func TestWrapf(t *testing.T) {
	root := errors.New("root")
	wrapped := errstack.Wrapf(root, "id=%s", "abc")
	if !strings.Contains(wrapped.Error(), "id=abc") {
		t.Fatalf("Wrapf message: %v", wrapped)
	}
}

func TestStackAndErrorWithStack(t *testing.T) {
	err := errstack.New("stacked")
	stack := errstack.Stack(err)
	if stack == "" || stack == "stacked" {
		t.Fatalf("expected stack frames, got %q", stack)
	}
	full := errstack.ErrorWithStack(err)
	if !strings.Contains(full, "stacked") || !strings.Contains(full, "\n") {
		t.Fatalf("ErrorWithStack: %q", full)
	}
	if errstack.Stack(nil) != "" || errstack.ErrorWithStack(nil) != "" {
		t.Fatal("nil stack helpers must return empty string")
	}
}

func TestAs(t *testing.T) {
	err := &ErrorString{S: "typed"}
	wrapped := errstack.Wrap(err, "outer")
	var target *ErrorString
	if !errstack.As(wrapped, &target) {
		t.Fatal("As must find typed error")
	}
	if target.S != "typed" {
		t.Fatalf("As target: %#v", target)
	}
}

type ErrorString struct{ S string }

func (e *ErrorString) Error() string { return e.S }

func TestMustPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Must must panic on error")
		}
	}()
	errstack.Must(errstack.New("panic-me"))
}

func TestMustNoPanic(t *testing.T) {
	errstack.Must(nil)
}
