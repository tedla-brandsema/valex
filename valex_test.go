package valex_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
	_ "github.com/tedla-brandsema/valex/internal/stub" // registers stub directives
)

// These exercise the engine itself (ValidateStruct, RegisterDirective,
// ValidatedValue, MustValidate) against the shared stub directives, with no
// dependency on the validator catalog.

func TestValidateStruct(t *testing.T) {
	type Person struct {
		Name string `val:"minlen,size=3"`
		Age  int    `val:"intrange,min=0,max=120"`
	}

	if err := valex.ValidateStruct(&Person{Name: "Alice", Age: 30}); err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	err := valex.ValidateStruct(&Person{Name: "Al", Age: 30})
	if err == nil {
		t.Fatal("expected failure on short name")
	}
	if !strings.Contains(err.Error(), "minlen") {
		t.Fatalf("expected error to reference directive, got %v", err)
	}
}

func TestValidateStructUnknownDirective(t *testing.T) {
	data := &struct {
		Field int `val:"nope"`
	}{Field: 1}

	err := valex.ValidateStruct(data)
	if err == nil {
		t.Fatal("expected unknown-directive error")
	}
	if !strings.Contains(err.Error(), `unknown directive "nope"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidatedValue(t *testing.T) {
	nonNegative := valex.ValidatorFunc[int](func(val int) error {
		if val < 0 {
			return errors.New("must be non-negative")
		}
		return nil
	})

	v := valex.ValidatedValue[int]{Validator: nonNegative}
	if err := v.Set(5); err != nil {
		t.Fatalf("Set(5) error: %v", err)
	}
	if v.Get() != 5 {
		t.Fatalf("expected 5, got %d", v.Get())
	}

	if err := v.Set(-1); err == nil {
		t.Fatal("expected error on negative value")
	}
	if v.Get() != 5 {
		t.Fatalf("expected value unchanged after failed Set, got %d", v.Get())
	}
}

func TestValidatedValueNoValidator(t *testing.T) {
	var v valex.ValidatedValue[int]
	if err := v.Set(1); err == nil {
		t.Fatal("expected error when no validator is set")
	}
}

func TestMustValidate(t *testing.T) {
	positive := valex.ValidatorFunc[int](func(val int) error {
		if val <= 0 {
			return errors.New("must be positive")
		}
		return nil
	})

	if got := valex.MustValidate(5, positive); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on invalid value")
		}
	}()
	valex.MustValidate(-1, positive)
}

// regDupDirective is a uniquely-named throwaway directive for the registration
// error-path test. Its name is not used by any other test or example.
type regDupDirective struct{}

func (*regDupDirective) Name() string              { return "valex_test_regdup" }
func (*regDupDirective) Mode() tagex.DirectiveMode { return tagex.EvalMode }
func (*regDupDirective) Handle(s string) (string, error) {
	return s, nil
}

func TestRegisterDirectiveErrors(t *testing.T) {
	d := &regDupDirective{}

	if err := valex.RegisterDirective(d); err != nil {
		t.Fatalf("first RegisterDirective: %v", err)
	}

	// A second registration of the same name returns *DuplicateDirectiveError,
	// re-exported from tagex so callers need not import it.
	var dup *valex.DuplicateDirectiveError
	if err := valex.RegisterDirective(d); !errors.As(err, &dup) {
		t.Fatalf("want *DuplicateDirectiveError, got %v", err)
	}

	// MustRegisterDirective panics on that same duplicate.
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on duplicate MustRegisterDirective")
		}
	}()
	valex.MustRegisterDirective(d)
}

func TestValidateStructAllAndFieldErrors(t *testing.T) {
	// stub registers "minlen" (size) and "intrange" (min,max) on the default registry.
	type Person struct {
		Name string `val:"minlen,size=5"`
		Age  int    `val:"intrange,min=0,max=120"`
	}

	// First-fail returns only one error...
	if err := valex.ValidateStruct(&Person{Name: "Al", Age: 200}); err == nil {
		t.Fatal("expected failure")
	} else if got := len(valex.FieldErrors(err)); got != 1 {
		t.Fatalf("ValidateStruct should surface one field, got %d", got)
	}

	// ...ValidateStructAll accumulates both.
	err := valex.ValidateStructAll(&Person{Name: "Al", Age: 200})
	if err == nil {
		t.Fatal("expected failure")
	}
	fe := valex.FieldErrors(err)
	if len(fe) != 2 {
		t.Fatalf("want 2 field errors, got %d: %v", len(fe), fe)
	}
	var pe *valex.ProcessError
	if !errors.As(fe["Age"], &pe) || pe.FieldPath != "Age" {
		t.Fatalf("Age should be a *ProcessError with FieldPath Age, got %v", fe["Age"])
	}
	if fe["Name"] == nil {
		t.Fatalf("missing Name field error: %v", fe)
	}

	// Clean struct -> nil.
	if err := valex.ValidateStructAll(&Person{Name: "Alice", Age: 30}); err != nil {
		t.Fatalf("expected nil for valid struct, got %v", err)
	}
	if valex.FieldErrors(nil) != nil {
		t.Fatal("FieldErrors(nil) should be nil")
	}
}

func TestRegistryIsolation(t *testing.T) {
	type Box struct {
		S string `val:"valex_test_regdup"`
	}

	// Two independent registries hold the same directive name without colliding.
	a := valex.NewRegistry()
	b := valex.NewRegistry()
	if err := valex.RegisterDirectiveTo(a, &regDupDirective{}); err != nil {
		t.Fatalf("register on a: %v", err)
	}
	if err := valex.RegisterDirectiveTo(b, &regDupDirective{}); err != nil {
		t.Fatalf("register on b (must not collide with a): %v", err)
	}

	// a knows the directive and validates clean.
	if err := a.ValidateStruct(&Box{S: "x"}); err != nil {
		t.Fatalf("a should validate, got %v", err)
	}

	// A fresh registry does not see a's directives — that is the isolation.
	c := valex.NewRegistry()
	if err := c.ValidateStruct(&Box{S: "x"}); err == nil {
		t.Fatal("empty registry should not know the directive")
	}
}
