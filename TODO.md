# Valex backlog

Working notes for valex. Completed work lives in git history, not here. Checkbox
= actionable; the prose under each item is the *why*, so the reasoning survives
even when the original conversation doesn't.

## In Progress

_(nothing in flight)_

## Backlog

- [ ] **(Deferred, on demand) Allow `=` inside `field` tag param values via `SplitN`.**
  `splitFormTag` in `forms/form.go` uses `strings.Split(pair, "=")` and requires
  exactly two parts, so a `default=` value can't itself contain `=` (a query
  string, base64 padding, `default=a=b` all fail). The fix is one word —
  `strings.SplitN(pair, "=", 2)` — and is backward-compatible: every currently-
  valid tag parses identically, only previously-rejected `2+ =` input becomes
  accepted. Mirrors the same deferred relaxation in tagex's `kv`. No third-party
  users demand it yet, and it slightly weakens loud-fail typo detection, so it
  waits until a real adopter needs it.

- [ ] **Bind into slices and maps of structs in `forms`.**
  This is a *binding* gap, not a validation one — `forms` does its own reflection
  in `bindStructFields` (forms/form.go) to copy request values into the struct
  before tagex runs; tagex never sees the request. That recursion handles nested
  structs and non-nil pointers-to-struct, but not slices or maps *of* structs
  (the `reflect.Slice` case only binds slices of scalars like `[]string`). So a
  `[]Address` field is never populated from the request. tagex v0.4.0 closed the
  matching gap on its side (its validation walk now descends into collections),
  which makes the asymmetry sharper: such a field gets validated but never bound.
  Either extend binding to cover repeated nested groups or document the boundary
  explicitly. Deferred until a real form needs repeated nested groups; flat and
  singly-nested forms (the common case) are unaffected.

## Decided / Won't do

- **The engine ships no directives; the catalog is opt-in.**
  Importing `valex` pulls in only tagex — not the catalog's `net`, `regexp`,
  `encoding/json`, `time`, and friends. Adopters register exactly the directives
  they use. The cost — a few `RegisterDirective` lines per program — is the
  point: the dependency surface and the active directive set stay explicit. Not
  changing.

- **`net/http` stays out of the core; `forms` is a separate package.**
  `valex/forms` is the only package that imports `net/http`. Programs that
  validate in code or by tag never pull it in. Folding forms into the core would
  force that dependency on everyone for a feature many won't use. Not changing.
