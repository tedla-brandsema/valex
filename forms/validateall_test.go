package forms_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/valex/forms"
	"github.com/tedla-brandsema/valex/validators"
)

func postForm(values url.Values) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestValidateAllAndFieldErrors(t *testing.T) {
	type Signup struct {
		Name  string `field:"name" val:"min,size=3"`
		Email string `field:"email" val:"email"`
		Age   int    `field:"age" val:"rangeint,min=1,max=150"`
	}

	reg := valex.NewRegistry()
	valex.MustRegisterDirectiveTo(reg, &validators.MinLengthValidator{})
	valex.MustRegisterDirectiveTo(reg, &validators.EmailValidator{})
	valex.MustRegisterDirectiveTo(reg, &validators.IntRangeValidator{})

	form := url.Values{
		"name":  {"Al"},   // validation fails: too short
		"email": {"nope"}, // validation fails: not an email
		"age":   {"abc"},  // binds AND validates wrong: not an int, and 0 < min
	}

	err := forms.ValidateAllWith(postForm(form), &Signup{}, reg)
	if err == nil {
		t.Fatal("expected failures")
	}

	// Field-level problems => 422.
	var ferr *forms.Error
	if !errors.As(err, &ferr) || ferr.StatusCode() != http.StatusUnprocessableEntity {
		t.Fatalf("want *forms.Error with 422, got %v", err)
	}

	fe := forms.FieldErrors(err)
	if len(fe) != 3 {
		t.Fatalf("want 3 field errors, got %d: %v", len(fe), fe)
	}

	// Name and Email are validation errors (reachable as *ProcessError).
	var pe *tagex.ProcessError
	if !errors.As(fe["Name"], &pe) {
		t.Fatalf("Name should be a validation error, got %v", fe["Name"])
	}
	if !errors.As(fe["Email"], &pe) {
		t.Fatalf("Email should be a validation error, got %v", fe["Email"])
	}

	// Age failed BOTH bind and validation -> the bind error wins (bind-precedence):
	// it is NOT a *ProcessError, and its message is the actionable parse failure,
	// not the range complaint on the unset zero value.
	if errors.As(fe["Age"], &pe) {
		t.Fatalf("Age should hold the bind error, got validation error %v", fe["Age"])
	}
	if !strings.Contains(fe["Age"].Error(), "abc") {
		t.Fatalf("Age should report the parse failure, got %v", fe["Age"])
	}
	if strings.Contains(fe["Age"].Error(), "range") {
		t.Fatalf("Age should NOT report the range complaint, got %v", fe["Age"])
	}
}
