# custom-directive

Extend the `val` tag with a directive of your own. A directive is any
`tagex.Directive[T]` — implement `Name`, `Mode`, and `Handle` for the field type
you handle, then register it with `valex.RegisterDirective`.

```bash
go run ./examples/custom-directive
```

The `even` directive accepts only even ints: `Seats: 4` passes, `Seats: 3` fails.
See [docs/struct-tags.md](../../docs/struct-tags.md#custom-directives).
