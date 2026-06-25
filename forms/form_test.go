package forms

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/tedla-brandsema/valex"
	_ "github.com/tedla-brandsema/valex/internal/stub" // registers stub directives
)

// formRequest builds an x-www-form-urlencoded POST carrying values.
func formRequest(values url.Values) *http.Request {
	return rawFormRequest(values.Encode())
}

// rawFormRequest builds an x-www-form-urlencoded POST with a raw body, for
// exercising malformed input.
func rawFormRequest(body string) *http.Request {
	req := httptest.NewRequest("POST", "/submit", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

// validateForm builds a form request from values, then binds and validates dst
// with a Validator, returning the validation error (nil if valid).
func validateForm(t *testing.T, values url.Values, dst any) error {
	t.Helper()
	v, err := New(formRequest(values))
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	return v.Validate(dst)
}

// wantStatus asserts that err is an *Error carrying the given HTTP status code.
func wantStatus(t *testing.T, err error, status int) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with status %d, got nil", status)
	}
	var ferr *Error
	if !errors.As(err, &ferr) {
		t.Fatalf("expected *Error, got %v", err)
	}
	if ferr.StatusCode() != status {
		t.Fatalf("expected status %d, got %d", status, ferr.StatusCode())
	}
}

func TestFormValidatorBindAndValidate(t *testing.T) {
	type Input struct {
		Name   string   `field:"name, max=1, required=true, default=unused" val:"minlen,size=3"`
		Age    int      `field:"age, max=1, required=true, default=0" val:"intrange,min=0,max=120"`
		Active bool     `field:"active, max=1, required=false, default=false"`
		Score  float64  `field:"score, max=1, required=false, default=0"`
		Tags   []string `field:"tags, max=2, required=false, default=unused"`
		Count  *int     `field:"count, max=1, required=false, default=0"`
	}

	values := url.Values{}
	values.Set("name", "Alice")
	values.Set("age", "30")
	values.Set("active", "true")
	values.Set("score", "1.5")
	values.Add("tags", "a")
	values.Add("tags", "b")
	values.Set("count", "7")

	var input Input
	if err := validateForm(t, values, &input); err != nil {
		t.Fatalf("expected validation to succeed, got %v", err)
	}

	if input.Name != "Alice" || input.Age != 30 || !input.Active || input.Score != 1.5 {
		t.Fatalf("unexpected bound values: %+v", input)
	}
	if len(input.Tags) != 2 || input.Tags[0] != "a" || input.Tags[1] != "b" {
		t.Fatalf("unexpected tags: %+v", input.Tags)
	}
	if input.Count == nil || *input.Count != 7 {
		t.Fatalf("unexpected count: %v", input.Count)
	}
}

func TestFormValidatorRequiredMissing(t *testing.T) {
	type Nested struct {
		Name string `field:"Name, max=1, required=true, default=unused"`
	}
	type Input struct {
		Nested Nested
	}

	var input Input
	err := validateForm(t, url.Values{}, &input)
	if err == nil {
		t.Fatal("expected required error")
	}
	if !strings.Contains(err.Error(), `form field "Nested.Name": field is required`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFormValidatorDefaultValue(t *testing.T) {
	type Input struct {
		Mode string `field:"Mode, max=1, required=false, default=basic"`
	}

	var input Input
	if err := validateForm(t, url.Values{}, &input); err != nil {
		t.Fatalf("expected validation to succeed, got %v", err)
	}
	if input.Mode != "basic" {
		t.Fatalf("expected default value to be set, got %q", input.Mode)
	}
}

func TestFormValidatorConversionError(t *testing.T) {
	type Input struct {
		Count int `field:"count, max=1, required=false, default=0"`
	}

	values := url.Values{}
	values.Set("count", "nope")

	var input Input
	err := validateForm(t, values, &input)
	if err == nil {
		t.Fatal("expected conversion error")
	}
	if !strings.Contains(err.Error(), `form field "Count":`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFormValidatorDefaultsRequiredFalse(t *testing.T) {
	type Input struct {
		Note string `field:"note"`
	}

	var input Input
	if err := validateForm(t, url.Values{}, &input); err != nil {
		t.Fatalf("expected validation to succeed, got %v", err)
	}
	if input.Note != "" {
		t.Fatalf("expected empty value when required/default omitted, got %q", input.Note)
	}
}

func TestFormValidatorDefaultsMaxOne(t *testing.T) {
	type Input struct {
		Tags []string `field:"tags"`
	}

	values := url.Values{}
	values.Add("tags", "a")
	values.Add("tags", "b")

	var input Input
	err := validateForm(t, values, &input)
	if err == nil {
		t.Fatal("expected max default error")
	}
	if !strings.Contains(err.Error(), "too many values") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFormValidatorDefaultNonStringTypes(t *testing.T) {
	type Input struct {
		Count *int     `field:"count, default=9"`
		Rate  float64  `field:"rate, default=1.25"`
		Flag  bool     `field:"flag, default=true"`
		Note  string   `field:"note, default=hi"`
		Nums  []int    `field:"nums, max=3, default=5"`
		Flags []bool   `field:"flags, max=2, default=false"`
		Tags  []string `field:"tags, max=2, default=a"`
	}

	var input Input
	if err := validateForm(t, url.Values{}, &input); err != nil {
		t.Fatalf("expected validation to succeed, got %v", err)
	}
	if input.Count == nil || *input.Count != 9 {
		t.Fatalf("unexpected Count: %v", input.Count)
	}
	if input.Rate != 1.25 || !input.Flag || input.Note != "hi" {
		t.Fatalf("unexpected defaults: %+v", input)
	}
	if len(input.Nums) != 1 || input.Nums[0] != 5 {
		t.Fatalf("unexpected Nums: %+v", input.Nums)
	}
	if len(input.Flags) != 1 || input.Flags[0] != false {
		t.Fatalf("unexpected Flags: %+v", input.Flags)
	}
	if len(input.Tags) != 1 || input.Tags[0] != "a" {
		t.Fatalf("unexpected Tags: %+v", input.Tags)
	}
}

func TestFormValidatorRequiredEmptyValue(t *testing.T) {
	type Input struct {
		Name string `field:"name, required=true"`
	}

	values := url.Values{}
	values.Set("name", "")

	var input Input
	err := validateForm(t, values, &input)
	if err == nil {
		t.Fatal("expected required error")
	}
	if !strings.Contains(err.Error(), `form field "Name": field is required`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFormValidatorMaxInvalid(t *testing.T) {
	type Input struct {
		Tags []string `field:"tags, max=zero"`
	}

	values := url.Values{}
	values.Add("tags", "a")

	var input Input
	err := validateForm(t, values, &input)
	var convErr *valex.ConversionError
	if !errors.As(err, &convErr) {
		t.Fatalf("expected conversion error, got %v", err)
	}
}

func TestFormValidatorPointerMissingOptional(t *testing.T) {
	type Input struct {
		Count *int `field:"count"`
	}

	var input Input
	if err := validateForm(t, url.Values{}, &input); err != nil {
		t.Fatalf("expected validation to succeed, got %v", err)
	}
	if input.Count != nil {
		t.Fatalf("expected nil Count when missing optional, got %v", input.Count)
	}
}

func TestFormValidatorSliceConversionError(t *testing.T) {
	type Input struct {
		Nums []int `field:"nums, max=3"`
	}

	values := url.Values{}
	values.Add("nums", "1")
	values.Add("nums", "x")

	var input Input
	err := validateForm(t, values, &input)
	if err == nil {
		t.Fatal("expected conversion error")
	}
	if !strings.Contains(err.Error(), `form field "Nums":`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFormValidatorNestedErrorPath(t *testing.T) {
	type Inner struct {
		Code int `field:"code"`
	}
	type Outer struct {
		Inner Inner
	}

	values := url.Values{}
	values.Set("code", "bad")

	var input Outer
	err := validateForm(t, values, &input)
	if err == nil {
		t.Fatal("expected conversion error")
	}
	if !strings.Contains(err.Error(), `form field "Inner.Code":`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFormValidatorUnsupportedKind(t *testing.T) {
	type Input struct {
		Meta map[string]string `field:"meta"`
	}

	values := url.Values{}
	values.Set("meta", "x")

	var input Input
	err := validateForm(t, values, &input)
	if err == nil {
		t.Fatal("expected unsupported type error")
	}
	if !strings.Contains(err.Error(), "unsupported field type") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindFormValues(t *testing.T) {
	type Input struct {
		Name string `field:"name"`
		Age  int    `field:"age"`
	}

	values := url.Values{}
	values.Set("name", "Alice")
	values.Set("age", "29")

	var input Input
	if err := Bind(&input, values); err != nil {
		t.Fatalf("Bind error: %v", err)
	}
	if input.Name != "Alice" || input.Age != 29 {
		t.Fatalf("unexpected bound values: %+v", input)
	}
}

func TestValidateFormSuccess(t *testing.T) {
	type Input struct {
		Name string `field:"name" val:"minlen,size=3"`
	}

	values := url.Values{}
	values.Set("name", "Alice")

	var input Input
	if err := Validate(formRequest(values), &input); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestValidateFormStatusBadRequest(t *testing.T) {
	type Input struct {
		Tags []string `field:"tags, max=zero"`
	}

	values := url.Values{}
	values.Add("tags", "a")

	var input Input
	wantStatus(t, Validate(formRequest(values), &input), http.StatusBadRequest)
}

func TestValidateFormStatusMissingRequired(t *testing.T) {
	type Input struct {
		Name string `field:"name, required=true"`
	}

	var input Input
	wantStatus(t, Validate(formRequest(url.Values{}), &input), http.StatusUnprocessableEntity)
}

func TestValidateFormStatusUnprocessable(t *testing.T) {
	type Input struct {
		Name string `field:"name" val:"minlen,size=3"`
	}

	values := url.Values{}
	values.Set("name", "Al")

	var input Input
	wantStatus(t, Validate(formRequest(values), &input), http.StatusUnprocessableEntity)
}

func TestValidateFormParseError(t *testing.T) {
	type Input struct {
		Name string `field:"name"`
	}

	var input Input
	wantStatus(t, Validate(rawFormRequest("name=%zz"), &input), http.StatusBadRequest)
}

// The (*Validator).Validate method wraps failures in *Error too, not only the
// package-level Validate convenience wrapper.
func TestFormValidatorMethodWrapsError(t *testing.T) {
	type Input struct {
		Name string `field:"name" val:"minlen,size=3"`
	}

	values := url.Values{}
	values.Set("name", "Al") // shorter than min length 3

	var input Input
	wantStatus(t, validateForm(t, values, &input), http.StatusUnprocessableEntity)
}
