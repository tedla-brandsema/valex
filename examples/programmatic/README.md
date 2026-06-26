# programmatic

Validate values in code, with no struct tags. A `ValidatorFunc` wraps a plain
function as a `Validator`, and `ValidatedValue` stores a value only when it
passes — leaving the previous value in place on failure.

```bash
go run ./examples/programmatic
```

`Set(42)` is stored; `Set(-1)` is rejected and the value stays `42`. See
[docs/programmatic.md](../../docs/programmatic.md).
