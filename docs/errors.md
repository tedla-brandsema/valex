# Errors

`ValidateStruct` and the `valex/forms` helpers return structured, typed errors so
callers can inspect *where* and *why* validation failed rather than matching
strings. Use `errors.As` to reach a concrete type and `errors.Is` / `Unwrap` to
follow the chain.

## Re-exported types

The typed errors come from tagex, but valex **re-exports** them as type aliases,
so you inspect failures without importing `tagex`:

```go
err := valex.ValidateStruct(&User{ /* ... */ })

var conv *valex.ConversionError
switch {
case errors.As(err, &conv):
	// a parameter value could not be converted to the field type
case errors.Is(err, valex.ErrNoValidator):
	// a ValidatedValue had no Validator set
}
```

A `*valex.ConversionError` and a `*tagex.ConversionError` are the *same* type;
only the import path differs.

| Type | Returned when |
| --- | --- |
| `ProcessError` | any failure while processing a field (carries `Stage`, field path, directive, param) |
| `TagError` | a tag's processing failed (carries the tag key) |
| `HandleError` | a directive's `Handle` rejected the value |
| `HookError` | a `Before`/`Success`/`Failure` hook returned an error |
| `UnknownDirectiveError` | a tag named a directive that wasn't registered |
| `DirectiveParseError` | a tag value omitted the directive name |
| `ParamParseError` | a tag arg isn't a `key=value` pair |
| `MissingParamError` | a required parameter was absent |
| `ParamConflictError` | a `param` sets both `required` and `default` |
| `ConversionError` | a parameter value couldn't be converted to the field type |
| `UnsupportedParamTypeError` | a `param` field has an unsupported type |
| `TypeMismatchError` | a directive was applied to a field of the wrong type |
| `FieldAccessError` | a field value couldn't be read |
| `FieldSetError` | a `MutMode` result couldn't be written back |

`ErrNoValidator` is valex's own sentinel — returned by `ValidatedValue.Set` when
no `Validator` is configured. The rest are the re-exported tagex types.

## The shape of a validation failure

`ValidateStruct` wraps a failed `val` directive in a `*TagError` (carrying the
tag key) around a `*ProcessError` (carrying the field path and directive):

```
tag "val" error: directive processing field "Age" directive "rangeint": value 200 is out of range [0, 120]
```

Reach into it for the parts you need:

```go
var pe *valex.ProcessError
if errors.As(err, &pe) {
	log.Printf("field=%s directive=%s: %v", pe.FieldPath, pe.Directive, pe.Cause)
}
```

`Stage` reports which phase failed. valex re-exports the stage constants
`StagePre`, `StageDirective`, `StageParam`, and `StagePost`.

## Distinguishing a rejected value from a wiring bug

A directive's `Handle` rejecting a value is wrapped in `*HandleError` — a *domain*
failure (the input is invalid), as opposed to a `*TypeMismatchError` or
`*FieldSetError`, which are programming mistakes. Branch on it to decide whether
to surface the error to a user or treat it as an internal bug:

```go
var he *valex.HandleError
if errors.As(err, &he) {
	return badRequest(he)   // a rule fired: tell the user
}
return internalError(err)   // wrong field type / unsettable field: fix the code
```

## Forms

`valex/forms` wraps validation and binding failures in `*forms.Error`, which adds
an HTTP status code. It unwraps to the same re-exported types, so all the
`errors.As` patterns above still work — see [forms.md](forms.md#errors).
