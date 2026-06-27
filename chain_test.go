package valex_test

import (
	"strings"
	"testing"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/valex/validators"
)

// trimDirective is a MutMode string directive used to prove chain ordering: it
// shortens the value, so it interacts with a following length check. The catalog
// is all EvalMode, so a custom MutMode directive is the only way to exercise the
// "each segment sees the previous one's output" rule.
type trimDirective struct{}

func (*trimDirective) Name() string                    { return "trim" }
func (*trimDirective) Mode() tagex.DirectiveMode       { return tagex.MutMode }
func (*trimDirective) Handle(s string) (string, error) { return strings.TrimSpace(s), nil }

func chainRegistry(t *testing.T) *valex.Registry {
	t.Helper()
	reg := valex.NewRegistry()
	valex.MustRegisterDirectiveTo(reg, &trimDirective{})
	valex.MustRegisterDirectiveTo(reg, &validators.MinLengthValidator{})
	return reg
}

// Distinct struct types because Go struct tags are static — the chain order has
// to be written into the tag, not passed at runtime.
type (
	trimThenMin struct {
		Name string `val:"trim;min,size=3"`
	}
	minThenTrim struct {
		Name string `val:"min,size=3;trim"`
	}
	minBigThenTrim struct {
		Name string `val:"min,size=10;trim"`
	}
	trimThenMinBig struct {
		Name string `val:"trim;min,size=10"`
	}
	stray struct {
		Name string `val:";trim;;min,size=1;"`
	}
)

// Order is semantic: trim then min sees the shortened value; min then trim sees
// the original. Same input, opposite outcomes.
func TestChainedDirectives_OrderMatters(t *testing.T) {
	// trim first: "  ab  " -> "ab" (len 2) -> min size=3 fails.
	if err := chainRegistry(t).ValidateStruct(&trimThenMin{Name: "  ab  "}); err == nil {
		t.Error("trim;min: expected length failure after trim, got nil")
	}

	// min first: "  ab  " (len 6) passes -> trim -> "ab".
	v := &minThenTrim{Name: "  ab  "}
	if err := chainRegistry(t).ValidateStruct(v); err != nil {
		t.Fatalf("min;trim: unexpected error: %v", err)
	}
	if v.Name != "ab" {
		t.Errorf("min;trim: expected trimmed %q, got %q", "ab", v.Name)
	}
}

// A failed segment stops the chain: a later MutMode segment must not run.
func TestChainedDirectives_StopsAtFirstFailure(t *testing.T) {
	v := &minBigThenTrim{Name: "  ab  "} // len 6, below min=10
	if err := chainRegistry(t).ValidateStruct(v); err == nil {
		t.Fatal("expected min failure, got nil")
	}
	if v.Name != "  ab  " {
		t.Errorf("trim after a failed min must not run; value changed to %q", v.Name)
	}
}

// Under ValidateStructAll a MutMode segment that already ran leaves the field
// mutated even though a later segment in the same chain fails: the partial
// (trimmed) value persists alongside the recorded error.
func TestChainedDirectives_PartialMutationPersistsUnderAll(t *testing.T) {
	v := &trimThenMinBig{Name: "  ab  "}
	if err := chainRegistry(t).ValidateStructAll(v); err == nil {
		t.Fatal("expected min failure, got nil")
	}
	if v.Name != "ab" {
		t.Errorf("expected partial mutation %q to persist, got %q", "ab", v.Name)
	}
}

// Empty, doubled, leading, and trailing ';' segments are skipped, not errors.
func TestChainedDirectives_SkipsEmptySegments(t *testing.T) {
	v := &stray{Name: "  ab  "}
	if err := chainRegistry(t).ValidateStruct(v); err != nil {
		t.Fatalf("stray ';' separators should be skipped, got: %v", err)
	}
	if v.Name != "ab" {
		t.Errorf("expected %q, got %q", "ab", v.Name)
	}
}

// A single-quoted parameter value lets a regex pattern carry a comma (`{1,3}`),
// which would otherwise split parameters. The backslash is doubled because the
// struct-tag layer unquotes the value before valex sees it, so the tag source
// `'^\\d{1,3}$'` reaches the directive as `^\d{1,3}$`. Mirrors tagex's
// TestProcessDirective_QuotedReservedChars — the motivating `\d{1,3}` case, not
// a backslash-free stand-in.
type quotedPattern struct {
	Code string `val:"regex,pattern='^\\d{1,3}$'"`
}

func TestQuotedParamValue(t *testing.T) {
	reg := valex.NewRegistry()
	valex.MustRegisterDirectiveTo(reg, &validators.RegexValidator{})

	if err := reg.ValidateStruct(&quotedPattern{Code: "42"}); err != nil {
		t.Errorf("quoted pattern should parse (comma inside quotes) and match %q: %v", "42", err)
	}
	if err := reg.ValidateStruct(&quotedPattern{Code: "4242"}); err == nil {
		t.Errorf(`expected "4242" to fail ^\d{1,3}$`)
	}
}

// FuzzChainedDirectives throws arbitrary field values through a chained tag and
// asserts the chain path never panics. The chain splitter itself is fuzzed
// upstream in tagex (FuzzSplitChain); this covers valex's value path through it.
type fuzzChainTarget struct {
	Name string `val:"trim;min,size=3"`
}

func FuzzChainedDirectives(f *testing.F) {
	reg := valex.NewRegistry()
	valex.MustRegisterDirectiveTo(reg, &trimDirective{})
	valex.MustRegisterDirectiveTo(reg, &validators.MinLengthValidator{})

	for _, s := range []string{"  ab  ", "abc", "", "   ", "a", strings.Repeat("x", 100)} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		_ = reg.ValidateStruct(&fuzzChainTarget{Name: s})    // must not panic
		_ = reg.ValidateStructAll(&fuzzChainTarget{Name: s}) // must not panic
	})
}
