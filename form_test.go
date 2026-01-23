package valex

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/tedla-brandsema/tagex"
)

func TestFormValidatorBindAndValidate(t *testing.T) {
	type Input struct {
		Name   string   `field:"name, max=1, required=true, default=unused" val:"min,size=3"`
		Age    int      `field:"age, max=1, required=true, default=0" val:"range,min=0,max=120"`
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

	req := httptest.NewRequest("POST", "/submit", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if !ok || err != nil {
		t.Fatalf("expected validation to succeed, got ok=%v err=%v", ok, err)
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

	req := httptest.NewRequest("POST", "/submit", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if ok || err == nil {
		t.Fatalf("expected required error, got ok=%v err=%v", ok, err)
	}
	if !strings.Contains(err.Error(), `form field "Nested.Name": field is required`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFormValidatorDefaultValue(t *testing.T) {
	type Input struct {
		Mode string `field:"Mode, max=1, required=false, default=basic"`
	}

	req := httptest.NewRequest("POST", "/submit", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if !ok || err != nil {
		t.Fatalf("expected validation to succeed, got ok=%v err=%v", ok, err)
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
	req := httptest.NewRequest("POST", "/submit", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if ok || err == nil {
		t.Fatalf("expected conversion error, got ok=%v err=%v", ok, err)
	}
	if !strings.Contains(err.Error(), `form field "Count":`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFormValidatorDefaultsRequiredFalse(t *testing.T) {
	type Input struct {
		Note string `field:"note"`
	}

	req := httptest.NewRequest("POST", "/submit", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if !ok || err != nil {
		t.Fatalf("expected validation to succeed, got ok=%v err=%v", ok, err)
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
	req := httptest.NewRequest("POST", "/submit", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if ok || err == nil {
		t.Fatalf("expected max default error, got ok=%v err=%v", ok, err)
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

	req := httptest.NewRequest("POST", "/submit", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if !ok || err != nil {
		t.Fatalf("expected validation to succeed, got ok=%v err=%v", ok, err)
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
	req := httptest.NewRequest("POST", "/submit", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if ok || err == nil {
		t.Fatalf("expected required error, got ok=%v err=%v", ok, err)
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
	req := httptest.NewRequest("POST", "/submit", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if ok || err == nil {
		t.Fatalf("expected max error, got ok=%v err=%v", ok, err)
	}
	var convErr *tagex.ConversionError
	if !errors.As(err, &convErr) {
		t.Fatalf("expected conversion error, got %v", err)
	}
}

func TestFormValidatorPointerMissingOptional(t *testing.T) {
	type Input struct {
		Count *int `field:"count"`
	}

	req := httptest.NewRequest("POST", "/submit", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if !ok || err != nil {
		t.Fatalf("expected validation to succeed, got ok=%v err=%v", ok, err)
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
	req := httptest.NewRequest("POST", "/submit", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if ok || err == nil {
		t.Fatalf("expected conversion error, got ok=%v err=%v", ok, err)
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
	req := httptest.NewRequest("POST", "/submit", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Outer
	ok, err := validator.Validate(&input)
	if ok || err == nil {
		t.Fatalf("expected conversion error, got ok=%v err=%v", ok, err)
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
	req := httptest.NewRequest("POST", "/submit", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	validator, err := NewFormValidator(req)
	if err != nil {
		t.Fatalf("NewFormValidator error: %v", err)
	}

	var input Input
	ok, err := validator.Validate(&input)
	if ok || err == nil {
		t.Fatalf("expected unsupported type error, got ok=%v err=%v", ok, err)
	}
	if !strings.Contains(err.Error(), "unsupported field type") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateFormStatusBadRequest(t *testing.T) {
	type Input struct {
		Tags []string `field:"tags, max=zero"`
	}

	values := url.Values{}
	values.Add("tags", "a")
	req := httptest.NewRequest("POST", "/submit", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var input Input
	ok, err := ValidateForm(req, &input)
	if ok || err == nil {
		t.Fatalf("expected validation error, got ok=%v err=%v", ok, err)
	}
	var formErr *FormError
	if !errors.As(err, &formErr) {
		t.Fatalf("expected FormError, got %v", err)
	}
	if formErr.StatusCode() != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, formErr.StatusCode())
	}
}

func TestValidateFormStatusUnprocessable(t *testing.T) {
	type Input struct {
		Name string `field:"name" val:"min,size=3"`
	}

	values := url.Values{}
	values.Set("name", "Al")
	req := httptest.NewRequest("POST", "/submit", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var input Input
	ok, err := ValidateForm(req, &input)
	if ok || err == nil {
		t.Fatalf("expected validation error, got ok=%v err=%v", ok, err)
	}
	var formErr *FormError
	if !errors.As(err, &formErr) {
		t.Fatalf("expected FormError, got %v", err)
	}
	if formErr.StatusCode() != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, formErr.StatusCode())
	}
}
