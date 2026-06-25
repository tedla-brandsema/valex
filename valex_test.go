package valex_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/tedla-brandsema/valex"
	_ "github.com/tedla-brandsema/valex/internal/valextest" // registers stub directives
)

// These exercise the engine itself (ValidateStruct, RegisterDirective,
// ValidatedValue, MustValidate) against the shared stub directives, with no
// dependency on the validator catalog.

func TestValidateStruct(t *testing.T) {
	type Person struct {
		Name string `val:"minlen,size=3"`
		Age  int    `val:"intrange,min=0,max=120"`
	}

	ok, err := valex.ValidateStruct(&Person{Name: "Alice", Age: 30})
	if !ok || err != nil {
		t.Fatalf("expected success, got ok=%v err=%v", ok, err)
	}

	ok, err = valex.ValidateStruct(&Person{Name: "Al", Age: 30})
	if ok || err == nil {
		t.Fatalf("expected failure on short name, got ok=%v err=%v", ok, err)
	}
	if !strings.Contains(err.Error(), "minlen") {
		t.Fatalf("expected error to reference directive, got %v", err)
	}
}

func TestValidateStructUnknownDirective(t *testing.T) {
	data := &struct {
		Field int `val:"nope"`
	}{Field: 1}

	ok, err := valex.ValidateStruct(data)
	if ok || err == nil {
		t.Fatalf("expected unknown-directive error, got ok=%v err=%v", ok, err)
	}
	if !strings.Contains(err.Error(), `unknown directive "nope"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidatedValue(t *testing.T) {
	nonNegative := valex.ValidatorFunc[int](func(val int) (bool, error) {
		if val < 0 {
			return false, errors.New("must be non-negative")
		}
		return true, nil
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
	positive := valex.ValidatorFunc[int](func(val int) (bool, error) {
		if val <= 0 {
			return false, errors.New("must be positive")
		}
		return true, nil
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
