# onyx-errors

[![Go Version](https://img.shields.io/badge/go-1.27-00ADD8.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Structured application/domain errors with independently formattable stack
traces. Package name: `errstack`.

JSON is the structured representation. Human-readable formatting is a separate
presentation layer over the same `StackTrace` data. The library does not depend
on any logging provider.

## Installation

```bash
go get github.com/MahdiPezeshkian/onyx-errors@latest
```

```go
import "github.com/MahdiPezeshkian/onyx-errors"
```

## Quick start

### 1. Create an error with a stack trace

```go
err := errstack.New("invalid captcha")
// or with a stable code:
err = errstack.NewCode("InvalidCaptcha", "invalid captcha")
```

The stack is captured at the call site via `runtime.Callers` /
`runtime.CallersFrames`, not during formatting or JSON marshaling.

### 2. Serialize to structured JSON

```go
b, _ := json.Marshal(err)
// {
//   "code": "InvalidCaptcha",
//   "message": "invalid captcha",
//   "stack_trace": [
//     {"function": "...Create...", "file": "/path/file.go", "line": 12}
//   ]
// }
```

`stack_trace` is always an array of objects (`function`, `file`, `line`).
It is never a preformatted string. The wrapped cause is not serialized.

### 3. Print only the stack trace

```go
var e *errstack.Error
if errors.As(err, &e) {
    fmt.Print(e.FormatStackTrace())
}
// or:
fmt.Print(errstack.Stack(err))
```

Human-readable layout (exact):

```text
package.(*Type).Method
    /full/path/file.go:12

caller.Function
    /full/path/other.go:34
```

### 4. Print the complete error

```go
fmt.Print(e.FormatText())
// InvalidCaptcha: invalid captcha
// package.(*Type).Method
//     /full/path/file.go:12
```

### 5. Use `%+v`

```go
fmt.Printf("%v\n", err)  // message only
fmt.Printf("%+v\n", err) // FormatText: metadata + stack
```

`Error()` always returns the message only and never a huge stack dump.

### 6. Wrap without losing the original stack

```go
root := errstack.New("root failure")
wrapped := errstack.Wrap(root, "while saving user")

// wrapped keeps root's original StackTrace
fmt.Print(errstack.Stack(wrapped))
```

Wrapping a plain (non-errstack) error captures a new stack at the wrap site.
`Unwrap`, `errors.Is`, and `errors.As` remain supported.

## Public API

| API | Purpose |
|---|---|
| `New` / `Newf` | Create stacked errors |
| `NewCode` / `NewCodef` | Create coded stacked errors |
| `Wrap` / `Wrapf` / `WithStack` | Annotate / attach stack (nil-safe) |
| `Cause` | Walk to the innermost cause |
| `Trace` | Structured `StackTrace` clone |
| `Stack` / `ErrorWithStack` | Package-level text helpers |
| `Is` / `As` | Standard comparison helpers |
| `Must` | Panic with full text when err != nil |
| `Error` | Structured error (`Code`, `Message`, `StackTrace`) |
| `Frame` / `StackTrace` | Structured stack model |
| `(*Error).FormatStackTrace` | Stack text only |
| `(*Error).FormatText` | Full human-readable text |
| `StackTrace.String` | Independent stack formatting |

## Design decisions

- **No HTTP semantics**: there is no status field. Handlers map domain errors to
  HTTP statuses.
- Stack traces are owned by `Error` and captured once at creation.
- Formatting and JSON serialization never re-capture stacks.
- Internal `onyx-errors` frames are filtered from captured traces.
- Zero external dependencies (standard library only).

## Dependencies

None beyond the Go standard library.

## License

MIT — see [LICENSE](LICENSE).
