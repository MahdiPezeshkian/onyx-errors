# onyx-errors

[![Go Version](https://img.shields.io/badge/go-1.27-00ADD8.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Reusable application/domain error primitives with stack traces. Extracted from
`Onyx.Common/errstack`.

## Installation

```bash
go get github.com/MahdiPezeshkian/onyx-errors@v0.1.0
```

## Import path

```go
import "github.com/MahdiPezeshkian/onyx-errors"
```

Package name: `errstack` (avoids colliding with the standard `errors` package).

## Basic usage

```go
err := fetch()
return errstack.Wrapf(err, "failed to load resource %s", id)
```

## Public API overview

| API | Purpose |
|---|---|
| `New` / `Newf` | Create stacked errors |
| `Wrap` / `Wrapf` / `WithStack` | Wrap existing errors (nil-safe) |
| `Cause` | Unwrap to the root cause |
| `Stack` / `ErrorWithStack` | Format stack traces |
| `Is` / `As` | Standard comparison helpers |
| `Must` | Panic with stack text when err != nil |

## Design decisions

- Built on `github.com/pkg/errors` for stack-trace compatibility with the
  original Onyx.Common implementation.
- **No HTTP semantics**: this library does not choose status codes and does not
  import `net/http` or web frameworks.
- `Wrap(nil, ...)` returns `nil`.

## Dependencies

- `github.com/pkg/errors`

## Limitations

- Does not define application-specific error codes or catalogs.
- Stack formatting depends on `pkg/errors` stack tracer behavior.

## License

MIT — see [LICENSE](LICENSE).
