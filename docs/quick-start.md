# Quick start

This guide validates a struct with the `val` tag using directives from the
`valex/validators` catalog.

## Before you begin

- Go 1.22 or later.

```bash
go get github.com/tedla-brandsema/valex@latest
```

## 1. Register the directives you use

The engine ships **no** directives of its own — you register the ones you need,
typically once at startup in an `init` function:

```go
import (
	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/valex/validators"
)

func init() {
	valex.RegisterDirective(&validators.MinLengthValidator{})
	valex.RegisterDirective(&validators.EmailValidator{})
	valex.RegisterDirective(&validators.IntRangeValidator{})
}
```

Each directive maps to a `val` tag name — `MinLengthValidator` is `min`,
`EmailValidator` is `email`, `IntRangeValidator` is `rangeint`. The full catalog
is in [struct-tags.md](struct-tags.md#catalog).

## 2. Annotate a struct

```go
type User struct {
	Name  string `val:"min,size=3"`
	Email string `val:"email"`
	Age   int    `val:"rangeint,min=0,max=120"`
}
```

The tag is `directive,arg=value,…`. `min,size=3` selects the `min` directive and
passes it `size=3`.

## 3. Validate

```go
err := valex.ValidateStruct(&User{Name: "Al", Email: "invalid", Age: 200})
// err: tag "val" error: directive processing field "Name" directive "min": ...
```

`ValidateStruct` takes a pointer to a struct, walks its fields, and returns `nil`
when every directive passes — or a typed [error](errors.md) describing the first
failure.

## Next steps

- Validate values in code, without tags — [programmatic validation](programmatic.md).
- Add your own rule — [custom directives](struct-tags.md#custom-directives).
- Bind and validate HTTP requests — [forms](forms.md).
- Inspect failures by type — [errors](errors.md).
