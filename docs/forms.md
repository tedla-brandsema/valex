# HTTP forms

`valex/forms` binds `net/http` request values into a struct, then validates the
struct's `val` tags. It is a separate package so the core engine never imports
`net/http` — programs that only validate in code or by tag don't pull it in.

Two tags are involved:

- `field` maps a struct field to a request key and controls binding.
- `val` (the engine's tag) validates the bound value.

```go
type Signup struct {
	Name  string `field:"name"  val:"min,size=3"`
	Email string `field:"email" val:"email"`
}
```

## Validate

`forms.Validate` parses the request, binds it, validates it, and returns `nil` or
an `*Error` carrying an HTTP status code:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	var in Signup
	if err := forms.Validate(r, &in); err != nil {
		var ferr *forms.Error
		errors.As(err, &ferr)
		http.Error(w, err.Error(), ferr.StatusCode())
		return
	}
	// ... use in
}
```

Binding reads from `request.ParseForm`, which merges the POST body **and** the
URL query string, so GET requests with query parameters work too. The directives
referenced by your `val` tags must be registered first (see
[struct-tags.md](struct-tags.md#registering-directives)).

To validate against an isolated [registry](struct-tags.md#registries) instead of
the global default — for test isolation, or two differently-configured form
validators in one process — use `ValidateWith` (or `NewWith` for the reusable
`Validator`); both still return a `*forms.Error` with its status code:

```go
err := forms.ValidateWith(r, &in, reg) // reg is a *valex.Registry
```

## The field tag

The first token is the request key; the rest are `key=value` options:

| Option | Default | Description |
| --- | --- | --- |
| *(key)* | field name | request key to read |
| `max` | `1` | maximum number of values accepted (for slice fields) |
| `required` | `false` | report `ErrFieldRequired` when the value is missing or empty |
| `default` | — | value to bind when the field is missing or empty |

```go
type Search struct {
	Q      string   `field:"q,required=true"`
	Tags   []string `field:"tag,max=5"`
	Page   int      `field:"page,default=1"`
}
```

## Binding without an HTTP handler

`forms.Bind` binds a `url.Values` into a struct using `field` tags only — no
parsing, no validation — for use outside a request:

```go
err := forms.Bind(&in, url.Values{"name": {"Gopher"}})
```

`forms.New(r)` returns a `*Validator` you can call `Validate` on repeatedly if
you want to separate parsing from validation.

## Errors

`Validate` and `ValidateAll` wrap every failure in `*forms.Error`:

| Method | Returns |
| --- | --- |
| `Error()` | the wrapped message |
| `Unwrap()` | the underlying error (for `errors.As` / `errors.Is`) |
| `StatusCode()` | the HTTP status |

`forms.Status(err)` maps errors to a status: **422** for field-level problems —
a validation failure (`*valex.TagError`) or a binding failure (a value that can't
convert to the field type, too many values, a missing required field) — and
**400** for a request that couldn't be parsed at all, or a malformed `field` tag
(a developer error). A field-level error is 422 whether or not a neighbor also
failed. Because the wrapped errors are the types [valex re-exports](errors.md),
inspect them with `errors.As` without importing `tagex`.

### Every error at once

`Validate` stops at the first failure. To report everything wrong with a
submission, use `ValidateAll` (or `ValidateAllWith`) with `FieldErrors`:

```go
if err := forms.ValidateAll(r, &in); err != nil {
	for field, fe := range forms.FieldErrors(err) {
		fmt.Printf("%s: %v\n", field, fe)
	}
}
```

`FieldErrors` merges binding *and* validation failures into a map keyed by struct
field path (`Name`, `Address.Zip` — not the request key; translate when
rendering). When a field fails both — e.g. `age=abc` on an `int` with a range
rule — the binding error wins, since "not a number" is the actionable message.
Non-field errors (an unparseable request) are omitted, so keep `err` itself
authoritative and render the map on top.

## Lifecycle hooks

Because validation runs through tagex, a form can opt into the processing
lifecycle by implementing tagex's hook interfaces — useful for normalizing input
before validation and acting on the result:

```go
func (r *Registration) Before() error          { /* normalize before validation */ return nil }
func (r *Registration) Success() error          { /* persist on success */ return nil }
func (r *Registration) Failure(cause error) error { /* handle rejection */ return nil }
```

`Before` runs after binding but before the `val` directives; `Success` runs only
when every directive passes; `Failure` runs on the first failure with the
validation error as `cause`. See the runnable
[forms example](../examples/forms/) and the `Example_lifecycleHooks` testable
example in `forms/example_hooks_test.go`.
