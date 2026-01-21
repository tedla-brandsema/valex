package valex

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/tedla-brandsema/tagex"
	"io"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

// CmpRangeValidator validates that a value is within an inclusive range.
type CmpRangeValidator[T cmp.Ordered] struct {
	Min T
	Max T
}

// Validate checks whether the value is within the configured range.
func (v *CmpRangeValidator[T]) Validate(val T) (ok bool, err error) {
	if cmp.Less(val, v.Min) || cmp.Less(v.Max, val) {
		return false, fmt.Errorf("value %v is out of range [%v, %v]", val, v.Min, v.Max)
	}
	return true, nil
}

// IntRangeValidator validates that an int is within an inclusive range.
type IntRangeValidator struct {
	Min int `param:"min"`
	Max int `param:"max"`
}

// Validate checks whether the value is within the configured range.
func (v *IntRangeValidator) Validate(val int) (ok bool, err error) {
	if val < v.Min || val > v.Max {
		return false, fmt.Errorf("value %d is out of range [%d, %d]", val, v.Min, v.Max)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *IntRangeValidator) Name() string {
	return "range"
}

// Mode returns the directive evaluation mode.
func (v *IntRangeValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *IntRangeValidator) Handle(val int) (int, error) {
	_, err := v.Validate(val)
	return val, err
}

// NonNegativeIntValidator validates that an int is not negative.
type NonNegativeIntValidator struct{}

// Validate checks whether the value is non-negative.
func (v *NonNegativeIntValidator) Validate(val int) (ok bool, err error) {
	if val < 0 {
		return false, fmt.Errorf("value %d is a negative integer", val)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *NonNegativeIntValidator) Name() string {
	return "pos"
}

// Mode returns the directive evaluation mode.
func (v *NonNegativeIntValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *NonNegativeIntValidator) Handle(val int) (int, error) {
	_, err := v.Validate(val)
	return val, err
}

// NonPositiveIntValidator validates that an int is not positive.
type NonPositiveIntValidator struct{}

// Validate checks whether the value is non-positive.
func (v *NonPositiveIntValidator) Validate(val int) (ok bool, err error) {
	if val > 0 {
		return false, fmt.Errorf("value %d is a positive integer", val)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *NonPositiveIntValidator) Name() string {
	return "neg"
}

// Mode returns the directive evaluation mode.
func (v *NonPositiveIntValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *NonPositiveIntValidator) Handle(val int) (int, error) {
	_, err := v.Validate(val)
	return val, err
}

// UrlValidator validates that a string is a valid URL.
type UrlValidator struct{}

// Validate checks whether the value is a valid URL.
func (v *UrlValidator) Validate(val string) (ok bool, err error) {
	_, err = url.ParseRequestURI(val)
	if err == nil {
		ok = true
	}
	return
}

// Name returns the directive identifier.
func (v *UrlValidator) Name() string {
	return "url"
}

// Mode returns the directive evaluation mode.
func (v *UrlValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *UrlValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// EmailValidator validates that a string is a valid email address.
type EmailValidator struct{}

// Validate checks whether the value is a valid email address.
func (v *EmailValidator) Validate(val string) (ok bool, err error) {
	_, err = mail.ParseAddress(val)
	if err == nil {
		ok = true
	}
	return
}

// Name returns the directive identifier.
func (v *EmailValidator) Name() string {
	return "email"
}

// Mode returns the directive evaluation mode.
func (v *EmailValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *EmailValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// NonEmptyStringValidator validates that a string is not empty.
type NonEmptyStringValidator struct{}

// Validate checks whether the value is non-empty.
func (v *NonEmptyStringValidator) Validate(val string) (ok bool, err error) {
	if val == "" {
		return false, fmt.Errorf("string is empty")
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *NonEmptyStringValidator) Name() string {
	return "!empty"
}

// Mode returns the directive evaluation mode.
func (v *NonEmptyStringValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *NonEmptyStringValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// MinLengthValidator validates that a string meets a minimum length.
type MinLengthValidator struct {
	Size int `param:"size"`
}

// Validate checks whether the value meets the minimum length.
func (v *MinLengthValidator) Validate(val string) (ok bool, err error) {
	if v.Size == 0 {
		return false, errors.New(`value of parameter "size" cannot be 0`)
	}
	if len(val) < v.Size {
		return false, fmt.Errorf("value %s exceeds minimum length %d", val, v.Size)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *MinLengthValidator) Name() string {
	return "min"
}

// Mode returns the directive evaluation mode.
func (v *MinLengthValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *MinLengthValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// MaxLengthValidator validates that a string does not exceed a maximum length.
type MaxLengthValidator struct {
	Size int `param:"size"`
}

// Validate checks whether the value does not exceed the maximum length.
func (v *MaxLengthValidator) Validate(val string) (ok bool, err error) {
	if v.Size == 0 {
		return false, errors.New(`value of parameter "size" cannot be 0`)
	}
	if len(val) > v.Size {
		return false, fmt.Errorf("value %s exceeds maximum length %d", val, v.Size)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *MaxLengthValidator) Name() string {
	return "max"
}

// Mode returns the directive evaluation mode.
func (v *MaxLengthValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *MaxLengthValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// LengthRangeValidator validates that a string length is within an inclusive range.
type LengthRangeValidator struct {
	Min int `param:"min"`
	Max int `param:"max"`
}

// Validate checks whether the value length is within the configured range.
func (v *LengthRangeValidator) Validate(val string) (ok bool, err error) {
	l := len(val)
	if v.Min == 0 {
		return false, errors.New(`"min" value cannot be 0`)
	}
	if v.Max == 0 {
		return false, errors.New(`"max" value cannot be 0`)
	}
	if l < v.Min || l > v.Max {
		return false, fmt.Errorf("value %q with length %d is not in range [%d, %d]", val, l, v.Min, v.Max)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *LengthRangeValidator) Name() string {
	return "len"
}

// Mode returns the directive evaluation mode.
func (v *LengthRangeValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *LengthRangeValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// RegexValidator validates that a string matches a regular expression.
type RegexValidator struct {
	Pattern *regexp.Regexp `param:"pattern"`
}

// Validate checks whether the value matches the configured pattern.
func (v *RegexValidator) Validate(val string) (ok bool, err error) {
	if v.Pattern == nil {
		return false, errors.New("regex pattern not set")
	}
	if !v.Pattern.MatchString(val) {
		return false, fmt.Errorf("value %q does not match pattern %q", val, v.Pattern.String())
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *RegexValidator) Name() string {
	return "regex"
}

// Mode returns the directive evaluation mode.
func (v *RegexValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// ConvertParam compiles the regex pattern parameter.
func (v *RegexValidator) ConvertParam(field reflect.StructField, fieldValue reflect.Value, raw string) error {
	if fieldValue.Type() != reflect.TypeOf((*regexp.Regexp)(nil)) {
		return tagex.NewConversionError(field, raw, "*regexp.Regexp")
	}
	r, err := regexp.Compile(raw)
	if err != nil {
		return fmt.Errorf("invalid regex pattern %q: %v", raw, err)
	}
	fieldValue.Set(reflect.ValueOf(r))
	return nil
}

// Handle validates the value and returns it unchanged.
func (v *RegexValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// AlphaNumericValidator validates that a string contains only alphanumeric characters.
type AlphaNumericValidator struct{}

// Validate checks whether the value is alphanumeric.
func (v *AlphaNumericValidator) Validate(val string) (ok bool, err error) {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, val)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, fmt.Errorf("value %q is not alphanumeric", val)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *AlphaNumericValidator) Name() string {
	return "alphanum"
}

// Mode returns the directive evaluation mode.
func (v *AlphaNumericValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *AlphaNumericValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// MACAddressValidator validates that a string is a valid MAC address.
type MACAddressValidator struct{}

// Validate checks whether the value is a valid MAC address.
func (v *MACAddressValidator) Validate(val string) (ok bool, err error) {
	_, err = net.ParseMAC(val)
	if err != nil {
		return false, fmt.Errorf("invalid MAC address %q: %v", val, err)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *MACAddressValidator) Name() string {
	return "mac"
}

// Mode returns the directive evaluation mode.
func (v *MACAddressValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *MACAddressValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// IpValidator validates that a string is a valid IP address.
type IpValidator struct{}

// Validate checks whether the value is a valid IP address.
func (v *IpValidator) Validate(val string) (ok bool, err error) {
	if ip := net.ParseIP(val); ip == nil {
		return false, fmt.Errorf("invalid IP address %q", val)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *IpValidator) Name() string {
	return "ip"
}

// Mode returns the directive evaluation mode.
func (v *IpValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *IpValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// IPv4Validator validates that a string is a valid IPv4 address.
type IPv4Validator struct{}

// Validate checks whether the value is a valid IPv4 address.
func (v *IPv4Validator) Validate(val string) (ok bool, err error) {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() == nil {
		return false, fmt.Errorf("invalid IPv4 address %q", val)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *IPv4Validator) Name() string {
	return "ipv4"
}

// Mode returns the directive evaluation mode.
func (v *IPv4Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *IPv4Validator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// IPv6Validator validates that a string is a valid IPv6 address.
type IPv6Validator struct{}

// Validate checks whether the value is a valid IPv6 address.
func (v *IPv6Validator) Validate(val string) (ok bool, err error) {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() != nil {
		return false, fmt.Errorf("invalid IPv6 address %q", val)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *IPv6Validator) Name() string {
	return "ipv6"
}

// Mode returns the directive evaluation mode.
func (v *IPv6Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *IPv6Validator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// XMLValidator validates that a string is well-formed XML with at least one element.
type XMLValidator struct{}

// Validate checks whether the value is valid XML with at least one element.
func (v *XMLValidator) Validate(val string) (ok bool, err error) {
	decoder := xml.NewDecoder(strings.NewReader(val))
	var hasElement bool

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return false, fmt.Errorf("XML parsing error: %w", err)
		}

		if _, ok := tok.(xml.StartElement); ok { // at least one tag
			hasElement = true
		}
	}

	if !hasElement {
		return false, fmt.Errorf("XML document must contain at least one element")
	}

	return true, nil
}

// Name returns the directive identifier.
func (v *XMLValidator) Name() string {
	return "xml"
}

// Mode returns the directive evaluation mode.
func (v *XMLValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *XMLValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// JSONValidator validates that a string is valid JSON.
type JSONValidator struct{}

// Validate checks whether the value is valid JSON.
func (v *JSONValidator) Validate(val string) (ok bool, err error) {
	if !json.Valid([]byte(val)) {
		return false, fmt.Errorf("invalid JSON")
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *JSONValidator) Name() string {
	return "json"
}

// Mode returns the directive evaluation mode.
func (v *JSONValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *JSONValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// CompositeValidator validates a value by running multiple validators in order.
type CompositeValidator[T cmp.Ordered] struct {
	Validators []Validator[T]
}

// Validate checks the value against each validator in order.
func (cv *CompositeValidator[T]) Validate(val T) (ok bool, err error) {
	for _, validator := range cv.Validators {
		if ok, err = validator.Validate(val); !ok {
			return false, err
		}
	}
	return true, nil
}
