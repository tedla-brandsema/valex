# chained

Apply several directives to one field by separating them with `;`. They run left
to right, and each `MutMode` result feeds the next: here `trim` then `lower`
(custom directives) normalize the username before the catalog `max` directive
checks its length. Order matters — `trim;lower;max` checks the cleaned value,
while `max;trim;lower` would check the raw one.

```bash
go run ./examples/chained
```

See [docs/struct-tags.md](../../docs/struct-tags.md#chaining-directives) for the
chaining rules (left-to-right, stop at first failure, quoting reserved
characters) and [docs/errors.md](../../docs/errors.md).
