# Valex examples

Each example is a self-contained `main` program. Run one with:

```bash
go run ./examples/<name>
```

| Example | What it shows |
| --- | --- |
| [programmatic](programmatic/) | Validate values in code with `ValidatorFunc` and `ValidatedValue` — no tags. |
| [validate-struct](validate-struct/) | Register catalog directives and validate a struct with the `val` tag. |
| [custom-directive](custom-directive/) | Extend the `val` tag with your own directive via `RegisterDirective`. |
| [chained](chained/) | Apply several directives to one field with `;` — `trim;lower;max`, run left to right. |
| [forms](forms/) | Bind and validate an `net/http` request with `valex/forms`. |

For API reference and concept guides, see [docs/](../docs/index.md). Go testable
examples that render on pkg.go.dev live in `example_test.go` (and
`forms/example_*_test.go`).
