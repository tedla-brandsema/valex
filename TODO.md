# Valex backlog

Working notes for valex. Completed work lives in git history, not here. Checkbox
= actionable; the prose under each item is the *why*, so the reasoning survives
even when the original conversation doesn't.

## In Progress

v0.2.0, driven by third-party review feedback (the three things most likely to
bite an adopter):

- [x] **Instance API for isolation (`NewRegistry`), core *and* forms.** Done.
  `valex.NewRegistry` gives an independent directive set with `RegisterDirectiveTo`
  / `MustRegisterDirectiveTo` and a `ValidateStruct` method; the package-level
  functions now wrap a default registry. `forms.NewWith` / `forms.ValidateWith`
  take a `*valex.Registry`, so isolated form validation is one call and keeps the
  `*forms.Error` status wrapping (the original deferral warned forms would need
  this — it's now covered, not just core). Fixes the test-isolation footgun and
  lets two differently-configured validators coexist in one process.

- [x] **Field-keyed multiple errors (the headline).** Done, tagex-first. tagex
  v0.4.1 added `ProcessStructAll` (accumulate, returns `errors.Join` of the typed
  per-field errors). valex surfaces it as `ValidateStructAll` (+ `Registry`
  method) and `FieldErrors(err) map[string]error`; `forms.ValidateAll` /
  `ValidateAllWith` accumulate binding *and* validation, and `forms.FieldErrors`
  merges them, bind error winning on a same-field collision. Keys are struct
  field paths (the only unambiguous key — request keys collide on nested
  `field:"id"`). Field-level binding failures now map to 422 (input problem),
  while a malformed `field` tag stays 400 (developer error).

- [x] **Fuzz the `forms` binding path.** Done — `FuzzBind` over adversarial
  `url.Values`, ~900k execs, no crashers.

## Backlog

- [ ] **Request-key keying for `forms.FieldErrors` (follow-up to field-keyed errors).**
  `FieldErrors` keys by struct field path (`Email`, `Sub.N`), not the request key
  (`email`, `n`) a frontend posts. Request-key keying is what you'd render next to
  inputs, but it needs a struct-path → request-key translation, and today's binder
  uses *flat* request keys, so two nested structs each with `field:"id"` collide —
  the keying would need hierarchical request keys in the binder first. Real work,
  not vapor; do it when an adopter needs input-keyed errors. Until then, callers
  map Go paths → input names themselves (documented on `FieldErrors`).

- [ ] **(Deferred, on demand) Allow `=` inside `field` tag param values via `SplitN`.**
  `splitFormTag` in `forms/form.go` uses `strings.Split(pair, "=")` and requires
  exactly two parts, so a `default=` value can't itself contain `=` (a query
  string, base64 padding, `default=a=b` all fail). The fix is one word —
  `strings.SplitN(pair, "=", 2)` — and is backward-compatible: every currently-
  valid tag parses identically, only previously-rejected `2+ =` input becomes
  accepted. tagex shipped the matching relaxation in v0.5.0 — its `kv` now splits
  on the first `=`, and single-quoted values embed `=`/`,`/`;` literally — so the
  `val` tag already accepts `=` in a value; `forms`' own `splitFormTag` is the
  lone holdout now, which sharpens the asymmetry. Still deferred: no third-party
  user demands it yet, and it slightly weakens loud-fail typo detection, so it
  waits until a real adopter needs it.

- [ ] **Bind into slices and maps of structs in `forms`.**
  A *binding* gap, not a validation one — `forms` does its own reflection in
  `bindStructFields` (forms/form.go) to copy request values into the struct before
  tagex runs; tagex never sees the request. That recursion handles nested structs
  and non-nil pointers-to-struct, but not slices or maps *of* structs (the
  `reflect.Slice` case only binds slices of scalars like `[]string`). So a
  `[]Address` field is *validated* — tagex's walk descends into collections — but
  never *populated*. That asymmetry is the bug.

  **Why we can't copy tagex.** tagex's collection support is a validation walk over
  an *already-built* value (`val.Index(i)` / `val.MapIndex`); `forms` has the
  inverse job — *construct* the collection from the wire. The recursion shape
  transfers; the hard part — turning flat request data into nested elements — does
  not. Different problem, not a shortcut we skipped.

  **Root cause is stdlib.** `url.Values` is `map[string][]string`, and
  `ParseForm`/`ParseQuery` do zero nesting: `lines[0].sku` is one opaque flat key
  (Go has no `a[b][c]` convention — that's PHP/Rails). Only repeated-key scalar
  slices work today. `encoding/json`, by contrast, is hierarchical natively.

  **Two routes when this is picked up:**
  - *Stay on `url.Values`* — invent/adopt a flat-key convention (`a[0][b]`), parse
    it ourselves, plus string coercion and a sparse-index/DoS bound (cap via the
    `field` tag's `max`). Stdlib gives nothing here.
  - *JSON-body source* — `json.Unmarshal` → `map[string]any` (hierarchy for free),
    then bind through the existing `field`/`val` tags by path. The library's model
    stays intact; only the binder's *source* changes (`map[string]any` instead of
    `map[string][]string`, path resolution instead of a flat lookup). Note this is
    a new content type, not a swap — urlencoded stays the flat case.

  Worth doing — a go-to form-validation lib will be asked for nested/repeated
  groups — but it needs deliberate design, not a hasty convention. Deferred; flat
  and singly-nested forms (the common case) are unaffected.

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

- **No rename of `forms.Validator` to avoid the `valex.Validator` name overlap.**
  `valex.Validator[T]` (the interface) and `forms.Validator` (the request
  binder) share a name but live in different packages, so the import path
  disambiguates them — exactly like `text/template.Template` vs
  `html/template.Template`. Go treats the package name as part of the identifier.
  The in-package pair `forms.Validate` (one-shot func) and `(*Validator).Validate`
  (method) is the same idiom as `http.ListenAndServe`. A rename would churn the
  API for no clarity gain. Not changing.

- **The per-validator `Name`/`Mode`/`Handle` boilerplate (and the int/float64
  twins) stays — it is deliberate, not debt.**
  Every catalog directive repeats an identical `Mode()` returning `EvalMode` and
  a `Handle()` that delegates to `Validate`, and the ordered-type families
  (Range, Min, Max, NonNegative, NonPositive, NonZero, OneOf) exist once for
  `int` and once for `float64`. A line-count audit flags this as ~450 lines to
  cut with a generic eval-adapter and `cmp.Ordered` generics. We are not doing
  it, on purpose:
  - **Boring beats clever.** Each validator is a flat, self-contained, greppable
    unit you can read top to bottom. A generic adapter trades that for
    indirection you decode at 3am.
  - **It would churn the public API.** Directives are registered as themselves
    (`RegisterDirective(&EmailValidator{})`); they implement `tagex.Directive`
    directly. An adapter that supplies `Mode`/`Handle` means the types stop being
    directives, changing registration. Go embedding doesn't dispatch to the
    concrete `Validate`, so there's no free lunch via embedding either.
  - **The generic collapse barely pays.** Per-`T` `Name()` and tag-arg
    conversion still need concrete types, so collapsing the int/float64 twins
    saves only the comparison bodies (~20 lines), not the structs or the
    `Name`/`Mode`/`Handle` trio. `CmpRangeValidator[T cmp.Ordered]` already
    exists for callers who want the generic form programmatically.

  In a ~50-entry catalog, uniform, dumb repetition is a feature: predictable,
  diff-friendly, and impossible to get subtly wrong. Revisit only if a validator
  ever needs `MutMode` or genuinely divergent dispatch.
