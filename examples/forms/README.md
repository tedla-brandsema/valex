# forms

Bind an `net/http` request into a struct with `field` tags, then validate its
`val` tags. `forms.Validate` returns a `*forms.Error` carrying an HTTP status
code — 422 for validation failures. The example uses `httptest` so it runs
without a server.

```bash
go run ./examples/forms
```

The first submission passes; the others fail on a too-short name and an invalid
email, each rejected with HTTP 422. See [docs/forms.md](../../docs/forms.md).

For an advanced flow that normalizes input and persists on success using the
lifecycle hooks, see `Example_lifecycleHooks` in `forms/example_hooks_test.go`.
