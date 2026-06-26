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

`Validate` wraps every failure in `*forms.Error`:

| Method | Returns |
| --- | --- |
| `Error()` | the wrapped message |
| `Unwrap()` | the underlying error (for `errors.As` / `errors.Is`) |
| `StatusCode()` | the HTTP status |

`forms.Status(err)` exposes the status mapping for any error: **422** for
validation failures (a `*valex.TagError`) and missing required fields
(`ErrFieldRequired`), **400** for binding and parse problems. Because the wrapped
validation errors are the types [valex re-exports](errors.md), you can inspect
them with `errors.As` without importing `tagex`.

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
