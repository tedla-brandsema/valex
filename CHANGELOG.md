# Changelog

All notable changes to Valex are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and the project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Stability

Valex is **pre-1.0 (0.x)**. The public API is still settling: while in 0.x a
breaking change bumps the **minor** version (`0.x.0`) and is called out below
under *Changed* or *Removed* with a migration note; patch releases (`0.x.y`) are
additive or fixes only. Pin a version.

`1.0` will mean the public surface — `Validator`, `ValidatorFunc`,
`ValidatedValue`, `MustValidate`, `ValidateStruct`, `RegisterDirective`,
`MustRegisterDirective`, the re-exported error types, the `valex/validators`
catalog tag names, and the `valex/forms` API — is frozen, and breaking changes
thereafter require a major bump.

## [Unreleased]

### Added
- `Registry` and `NewRegistry` for an isolated directive set, with
  `RegisterDirectiveTo` / `MustRegisterDirectiveTo` and a `ValidateStruct`
  method. The package-level `RegisterDirective`, `MustRegisterDirective`, and
  `ValidateStruct` now wrap a shared default registry. Use a `Registry` for test
  isolation or to run two differently-configured validators in one process.
- `forms.NewWith` and `forms.ValidateWith`, which validate against a given
  `*valex.Registry` instead of the default — so isolated form validation stays a
  single call and keeps the `*forms.Error` status wrapping. `forms.New` /
  `forms.Validate` are unchanged and use the default registry.

## [0.1.0] - 2026-06-26

Initial release.

### Added
- Programmatic validation: the `Validator[T]` interface and `ValidatorFunc[T]`
  adapter, the `ValidatedValue[T]` guarded-assignment wrapper, and `MustValidate`.
- Struct-tag validation via the `val` tag: `ValidateStruct`, `RegisterDirective`
  (returns `*EmptyDirectiveNameError` for a blank name or `*DuplicateDirectiveError`
  for a name already registered), and `MustRegisterDirective`, which panics on
  those — the convenient choice for registering once at startup.
- `valex/validators`: an opt-in catalog of ready-made directives for `int`,
  `float64`, `string`, `time.Time`, `time.Duration`, `net.IP`, and `url.URL`,
  plus the generic `CmpRangeValidator`, `NonZeroValidator`, and
  `CompositeValidator` usable as `Validator` values directly.
- `valex/forms`: binds `net/http` request values into structs via the `field`
  tag and validates their `val` tags, returning a `*forms.Error` carrying an HTTP
  status code. Kept separate so the core engine never imports `net/http`.
- Error types re-exported from tagex (`ProcessError`, `TagError`, `HandleError`,
  `ConversionError`, `DuplicateDirectiveError`, …) plus the `ErrNoValidator`
  sentinel, so callers inspect failures with `errors.As` / `errors.Is` without
  importing tagex.
- Documentation set: README, a `docs/` guide, and runnable programs under
  `examples/`.

### Notes
- Built on [tagex](https://github.com/tedla-brandsema/tagex) v0.4.0, whose
  per-call directive cloning is what makes concurrent `ValidateStruct` on the
  shared `val` registry safe.
- Requires Go 1.22 or later.

[Unreleased]: https://github.com/tedla-brandsema/valex/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/tedla-brandsema/valex/releases/tag/v0.1.0
