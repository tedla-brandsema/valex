# validate-struct

Register directives from the `valex/validators` catalog, then validate a struct
with the `val` tag. The first user passes; the others fail on a too-short name
and an out-of-range age. The error carries the tag key, field path, and directive
name.

```bash
go run ./examples/validate-struct
```

See [docs/quick-start.md](../../docs/quick-start.md),
[docs/struct-tags.md](../../docs/struct-tags.md), and
[docs/errors.md](../../docs/errors.md).
