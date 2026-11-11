package errstack_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MahdiPezeshkian/onyx-errors"
)

func TestStackTraceStringExactFormat(t *testing.T) {
	st := errstack.StackTrace{
		{Function: "kanoon/usecase/usecase/authorization.(*authorizationClient).ValidateDchestCaptcha", File: "/app/usecase/usecase/authorization/impl.go", Line: 409},
		{Function: "kanoon/handler/handler.(*handlerClient).Otp", File: "/app/handler/handler/otp.go", Line: 12},
		{Function: "github.com/gin-gonic/gin.(*Context).Next", File: "/app/vendor/github.com/gin-gonic/gin/context.go", Line: 185},
	}
	want := "" +
		"kanoon/usecase/usecase/authorization.(*authorizationClient).ValidateDchestCaptcha\n" +
		"    /app/usecase/usecase/authorization/impl.go:409\n" +
		"\n" +
		"kanoon/handler/handler.(*handlerClient).Otp\n" +
		"    /app/handler/handler/otp.go:12\n" +
		"\n" +
		"github.com/gin-gonic/gin.(*Context).Next\n" +
		"    /app/vendor/github.com/gin-gonic/gin/context.go:185"
	if got := st.String(); got != want {
		t.Fatalf("format mismatch\nwant:\n%q\ngot:\n%q", want, got)
	}
}

func TestStackTraceNilAndEmpty(t *testing.T) {
	var nilST errstack.StackTrace
	if nilST.String() != "" {
		t.Fatal("nil stack must format empty")
	}
	if (errstack.StackTrace{}).String() != "" {
		t.Fatal("empty stack must format empty")
	}
}

func TestStackTraceIncompleteFrame(t *testing.T) {
	st := errstack.StackTrace{{Function: "pkg.Fn", File: "", Line: 0}}
	want := "pkg.Fn\n    :0"
	if got := st.String(); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestStackTraceUnicodeNames(t *testing.T) {
	st := errstack.StackTrace{{Function: "app.(*سرویس).اجرا", File: "/tmp/فایل.go", Line: 7}}
	got := st.String()
	if !strings.Contains(got, "app.(*سرویس).اجرا") || !strings.Contains(got, "/tmp/فایل.go:7") {
		t.Fatalf("unicode formatting failed: %q", got)
	}
}

func TestCapturePointsAtCaller(t *testing.T) {
	err := createSampleError()
	var e *errstack.Error
	if !errors.As(err, &e) {
		t.Fatal("expected *errstack.Error")
	}
	if len(e.StackTrace) == 0 {
		t.Fatal("expected captured frames")
	}
	top := e.StackTrace[0]
	if !strings.Contains(top.Function, "createSampleError") {
		t.Fatalf("top function=%q, want createSampleError", top.Function)
	}
	if filepath.Base(top.File) != "stack_test.go" {
		t.Fatalf("top file=%q", top.File)
	}
	if top.Line <= 0 {
		t.Fatal("line must be positive")
	}
	for _, f := range e.StackTrace {
		if strings.Contains(f.Function, "github.com/MahdiPezeshkian/onyx-errors.") {
			t.Fatalf("internal frame leaked: %+v", f)
		}
	}
}

func createSampleError() error {
	return errstack.New("sample") // line captured
}

func TestNestedCapture(t *testing.T) {
	err := outer()
	var e *errstack.Error
	if !errors.As(err, &e) {
		t.Fatal("as")
	}
	joined := e.FormatStackTrace()
	if !strings.Contains(joined, "inner") || !strings.Contains(joined, "outer") {
		// Stack is captured at New inside inner; outer is a caller frame.
		if !strings.Contains(joined, "inner") {
			t.Fatalf("stack missing inner: %s", joined)
		}
		if !strings.Contains(joined, "outer") {
			t.Fatalf("stack missing outer: %s", joined)
		}
	}
}

func outer() error { return inner() }
func inner() error { return errstack.New("nested") }

func TestJSONStructuredStack(t *testing.T) {
	err := errstack.NewCode("InvalidCaptcha", "Invalid captcha").(*errstack.Error)
	err.StackTrace = errstack.StackTrace{
		{Function: "kanoon/usecase/usecase/authorization.(*authorizationClient).ValidateDchestCaptcha", File: "/app/usecase/usecase/authorization/impl.go", Line: 409},
		{Function: "kanoon/handler/handler.(*handlerClient).Otp", File: "/app/handler/handler/otp.go", Line: 12},
	}
	b, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	var payload map[string]any
	if err := json.Unmarshal(b, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["code"] != "InvalidCaptcha" || payload["message"] != "Invalid captcha" {
		t.Fatalf("payload=%v", payload)
	}
	if _, ok := payload["cause"]; ok {
		t.Fatal("cause must not be serialized")
	}
	if _, ok := payload["status"]; ok {
		t.Fatal("status must not be serialized")
	}
	rawTrace, ok := payload["stack_trace"].([]any)
	if !ok || len(rawTrace) != 2 {
		t.Fatalf("stack_trace=%v", payload["stack_trace"])
	}
	frame0 := rawTrace[0].(map[string]any)
	if frame0["function"] == "" || frame0["file"] == "" || frame0["line"].(float64) != 409 {
		t.Fatalf("frame0=%v", frame0)
	}
}

func TestJSONEmptyAndNilStackOmitted(t *testing.T) {
	e := &errstack.Error{Message: "x", StackTrace: nil}
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "stack_trace") {
		t.Fatalf("nil stack should omit: %s", b)
	}
	e.StackTrace = errstack.StackTrace{}
	b, err = json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "stack_trace") {
		t.Fatalf("empty stack should omit: %s", b)
	}
}

func TestFormatMethodsAndFmt(t *testing.T) {
	e := &errstack.Error{
		Code:    "E1",
		Message: "boom",
		StackTrace: errstack.StackTrace{
			{Function: "pkg.Fn", File: "/a.go", Line: 1},
		},
	}
	if e.Error() != "boom" {
		t.Fatalf("Error()=%q", e.Error())
	}
	if e.FormatStackTrace() != "pkg.Fn\n    /a.go:1" {
		t.Fatalf("FormatStackTrace=%q", e.FormatStackTrace())
	}
	text := e.FormatText()
	if !strings.HasPrefix(text, "E1: boom\n") || !strings.Contains(text, "pkg.Fn") {
		t.Fatalf("FormatText=%q", text)
	}
	gotV := fmt.Sprintf("%v", e)
	if gotV != "boom" {
		t.Fatalf("%%v=%q", gotV)
	}
	plus := fmt.Sprintf("%+v", e)
	if !strings.Contains(plus, "E1: boom") || !strings.Contains(plus, "/a.go:1") {
		t.Fatalf("%%+v=%q", plus)
	}
	if strings.Contains(gotV, "/a.go") {
		t.Fatal("default format must not include stack")
	}
}

func TestWrapPreservesOriginalStack(t *testing.T) {
	root := errstack.New("root")
	var rootE *errstack.Error
	if !errors.As(root, &rootE) {
		t.Fatal("as root")
	}
	original := rootE.StackTrace.Clone()

	wrapped := errstack.Wrap(root, "layer")
	var we *errstack.Error
	if !errors.As(wrapped, &we) {
		t.Fatal("as wrapped")
	}
	if len(we.StackTrace) != len(original) {
		t.Fatalf("stack len changed: %d vs %d", len(we.StackTrace), len(original))
	}
	for i := range original {
		if we.StackTrace[i] != original[i] {
			t.Fatalf("frame %d changed: %+v vs %+v", i, we.StackTrace[i], original[i])
		}
	}
	if we.StackTrace[0].Function != original[0].Function {
		t.Fatal("top frame changed after wrap")
	}
	if !errstack.Is(wrapped, root) {
		t.Fatal("Is must match root error value identity via Unwrap chain")
	}
	// errors.Is compares by == for non-Is methods; root is *Error pointer.
	if !errors.Is(wrapped, root) {
		t.Fatal("stdlib Is")
	}
	if errstack.Cause(wrapped) != root && errstack.Cause(wrapped).Error() != "root" {
		// Cause should be the innermost.
		c := errstack.Cause(wrapped)
		if c.Error() != "root" {
			t.Fatalf("Cause=%v", c)
		}
	}
}

func TestWrapNonErrstackCapturesNewStack(t *testing.T) {
	root := errors.New("plain")
	wrapped := errstack.Wrap(root, "layer")
	var we *errstack.Error
	if !errors.As(wrapped, &we) {
		t.Fatal("as")
	}
	if len(we.StackTrace) == 0 {
		t.Fatal("expected new stack on plain wrap")
	}
	if !strings.Contains(we.StackTrace[0].Function, "TestWrapNonErrstackCapturesNewStack") {
		t.Fatalf("top=%q", we.StackTrace[0].Function)
	}
	if !errstack.Is(wrapped, root) {
		t.Fatal("Is plain root")
	}
}

func TestWithStackIdempotent(t *testing.T) {
	err := errstack.New("x")
	again := errstack.WithStack(err)
	if again != err {
		t.Fatal("WithStack should return same error when stack exists")
	}
	plain := errors.New("p")
	stacked := errstack.WithStack(plain)
	var e *errstack.Error
	if !errors.As(stacked, &e) || len(e.StackTrace) == 0 {
		t.Fatal("WithStack must capture for plain errors")
	}
}

func TestTraceHelper(t *testing.T) {
	err := errstack.New("x")
	st := errstack.Trace(err)
	if len(st) == 0 {
		t.Fatal("Trace")
	}
	st[0].Function = "mutated"
	st2 := errstack.Trace(err)
	if st2[0].Function == "mutated" {
		t.Fatal("Trace must clone")
	}
}
