package valex

import (
	"cmp"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/tedla-brandsema/tagex"
)

// CmpRangeValidator validates that a value is within an inclusive range.
type CmpRangeValidator[T cmp.Ordered] struct {
	Min T
	Max T
}

// Validate checks whether the value is within the configured range.
func (v *CmpRangeValidator[T]) Validate(val T) (ok bool, err error) {
	return validateRange(val, v.Min, v.Max, "value %v is out of range [%v, %v]")
}

// NonZeroValidator validates that a value is not the zero value for its type.
type NonZeroValidator[T any] struct{}

// Validate checks whether the value is non-zero.
func (v *NonZeroValidator[T]) Validate(val T) (ok bool, err error) {
	rv := reflect.ValueOf(val)
	if !rv.IsValid() || rv.IsZero() {
		return false, fmt.Errorf("value is zero")
	}
	return true, nil
}

func validateNonZero[T any](val T) (ok bool, err error) {
	return (&NonZeroValidator[T]{}).Validate(val)
}

func validateRange[T cmp.Ordered](val, min, max T, format string) (ok bool, err error) {
	if cmp.Less(val, min) || cmp.Less(max, val) {
		return false, fmt.Errorf(format, val, min, max)
	}
	return true, nil
}

func validateMin[T cmp.Ordered](val, min T, format string) (ok bool, err error) {
	if cmp.Less(val, min) {
		return false, fmt.Errorf(format, val, min)
	}
	return true, nil
}

func validateMax[T cmp.Ordered](val, max T, format string) (ok bool, err error) {
	if cmp.Less(max, val) {
		return false, fmt.Errorf(format, val, max)
	}
	return true, nil
}

func validateNonNegative[T cmp.Ordered](val T, format string) (ok bool, err error) {
	var zero T
	if cmp.Less(val, zero) {
		return false, fmt.Errorf(format, val)
	}
	return true, nil
}

func validateNonPositive[T cmp.Ordered](val T, format string) (ok bool, err error) {
	var zero T
	if cmp.Less(zero, val) {
		return false, fmt.Errorf(format, val)
	}
	return true, nil
}

func validateOneOf[T comparable](val T, values []T, emptyErr error, format string) (ok bool, err error) {
	if len(values) == 0 {
		return false, emptyErr
	}
	if slices.Contains(values, val) {
		return true, nil
	}
	return false, fmt.Errorf(format, val)
}

// IntRangeValidator validates that an int is within an inclusive range.
type IntRangeValidator struct {
	Min int `param:"min"`
	Max int `param:"max"`
}

// Validate checks whether the value is within the configured range.
func (v *IntRangeValidator) Validate(val int) (ok bool, err error) {
	return validateRange(val, v.Min, v.Max, "value %d is out of range [%d, %d]")
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

// Float64RangeValidator validates that a float64 is within an inclusive range.
type Float64RangeValidator struct {
	Min float64 `param:"min"`
	Max float64 `param:"max"`
}

// Validate checks whether the value is within the configured range.
func (v *Float64RangeValidator) Validate(val float64) (ok bool, err error) {
	return validateRange(val, v.Min, v.Max, "value %g is out of range [%g, %g]")
}

// Name returns the directive identifier.
func (v *Float64RangeValidator) Name() string {
	return "rangef64"
}

// Mode returns the directive evaluation mode.
func (v *Float64RangeValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *Float64RangeValidator) Handle(val float64) (float64, error) {
	_, err := v.Validate(val)
	return val, err
}

// NonNegativeIntValidator validates that an int is not negative.
type NonNegativeIntValidator struct{}

// Validate checks whether the value is non-negative.
func (v *NonNegativeIntValidator) Validate(val int) (ok bool, err error) {
	return validateNonNegative(val, "value %d is a negative integer")
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

// NonNegativeFloat64Validator validates that a float64 is not negative.
type NonNegativeFloat64Validator struct{}

// Validate checks whether the value is non-negative.
func (v *NonNegativeFloat64Validator) Validate(val float64) (ok bool, err error) {
	return validateNonNegative(val, "value %g is a negative float")
}

// Name returns the directive identifier.
func (v *NonNegativeFloat64Validator) Name() string {
	return "posf64"
}

// Mode returns the directive evaluation mode.
func (v *NonNegativeFloat64Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *NonNegativeFloat64Validator) Handle(val float64) (float64, error) {
	_, err := v.Validate(val)
	return val, err
}

// NonPositiveIntValidator validates that an int is not positive.
type NonPositiveIntValidator struct{}

// Validate checks whether the value is non-positive.
func (v *NonPositiveIntValidator) Validate(val int) (ok bool, err error) {
	return validateNonPositive(val, "value %d is a positive integer")
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

// NonPositiveFloat64Validator validates that a float64 is not positive.
type NonPositiveFloat64Validator struct{}

// Validate checks whether the value is non-positive.
func (v *NonPositiveFloat64Validator) Validate(val float64) (ok bool, err error) {
	return validateNonPositive(val, "value %g is a positive float")
}

// Name returns the directive identifier.
func (v *NonPositiveFloat64Validator) Name() string {
	return "negf64"
}

// Mode returns the directive evaluation mode.
func (v *NonPositiveFloat64Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *NonPositiveFloat64Validator) Handle(val float64) (float64, error) {
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
	if ok, err := validateNonZero(val); !ok {
		return false, fmt.Errorf("string is empty")
	} else if err != nil {
		return false, err
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

// NonEmptyStringAliasValidator provides the legacy "nonempty" tag.
type NonEmptyStringAliasValidator struct {
	NonEmptyStringValidator
}

// Name returns the directive identifier.
func (v *NonEmptyStringAliasValidator) Name() string {
	return "nonempty"
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
	if v.Size < 0 {
		return false, errors.New(`value of parameter "size" cannot be negative`)
	}
	if len(val) < v.Size {
		return false, fmt.Errorf("value %s is shorter than minimum length %d", val, v.Size)
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
	if v.Size < 0 {
		return false, errors.New(`value of parameter "size" cannot be negative`)
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
	if v.Min < 0 || v.Max < 0 {
		return false, errors.New(`"min" and "max" cannot be negative`)
	}
	if v.Min > v.Max {
		return false, errors.New(`"min" cannot exceed "max"`)
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

// MinIntValidator validates that an int is greater than or equal to Min.
type MinIntValidator struct {
	Min int `param:"min"`
}

// Validate checks whether the value meets the minimum.
func (v *MinIntValidator) Validate(val int) (ok bool, err error) {
	return validateMin(val, v.Min, "value %d is less than minimum %d")
}

// Name returns the directive identifier.
func (v *MinIntValidator) Name() string {
	return "minint"
}

// Mode returns the directive evaluation mode.
func (v *MinIntValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *MinIntValidator) Handle(val int) (int, error) {
	_, err := v.Validate(val)
	return val, err
}

// MinFloat64Validator validates that a float64 is greater than or equal to Min.
type MinFloat64Validator struct {
	Min float64 `param:"min"`
}

// Validate checks whether the value meets the minimum.
func (v *MinFloat64Validator) Validate(val float64) (ok bool, err error) {
	return validateMin(val, v.Min, "value %g is less than minimum %g")
}

// Name returns the directive identifier.
func (v *MinFloat64Validator) Name() string {
	return "minf64"
}

// Mode returns the directive evaluation mode.
func (v *MinFloat64Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *MinFloat64Validator) Handle(val float64) (float64, error) {
	_, err := v.Validate(val)
	return val, err
}

// MaxIntValidator validates that an int is less than or equal to Max.
type MaxIntValidator struct {
	Max int `param:"max"`
}

// Validate checks whether the value meets the maximum.
func (v *MaxIntValidator) Validate(val int) (ok bool, err error) {
	return validateMax(val, v.Max, "value %d exceeds maximum %d")
}

// Name returns the directive identifier.
func (v *MaxIntValidator) Name() string {
	return "maxint"
}

// Mode returns the directive evaluation mode.
func (v *MaxIntValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *MaxIntValidator) Handle(val int) (int, error) {
	_, err := v.Validate(val)
	return val, err
}

// MaxFloat64Validator validates that a float64 is less than or equal to Max.
type MaxFloat64Validator struct {
	Max float64 `param:"max"`
}

// Validate checks whether the value meets the maximum.
func (v *MaxFloat64Validator) Validate(val float64) (ok bool, err error) {
	return validateMax(val, v.Max, "value %g exceeds maximum %g")
}

// Name returns the directive identifier.
func (v *MaxFloat64Validator) Name() string {
	return "maxf64"
}

// Mode returns the directive evaluation mode.
func (v *MaxFloat64Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *MaxFloat64Validator) Handle(val float64) (float64, error) {
	_, err := v.Validate(val)
	return val, err
}

// NonZeroIntValidator validates that an int is not zero.
type NonZeroIntValidator struct{}

// Validate checks whether the value is non-zero.
func (v *NonZeroIntValidator) Validate(val int) (ok bool, err error) {
	return validateNonZero(val)
}

// Name returns the directive identifier.
func (v *NonZeroIntValidator) Name() string {
	return "!zero"
}

// Mode returns the directive evaluation mode.
func (v *NonZeroIntValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *NonZeroIntValidator) Handle(val int) (int, error) {
	_, err := v.Validate(val)
	return val, err
}

// NonZeroFloat64Validator validates that a float64 is not zero.
type NonZeroFloat64Validator struct{}

// Validate checks whether the value is non-zero.
func (v *NonZeroFloat64Validator) Validate(val float64) (ok bool, err error) {
	return validateNonZero(val)
}

// Name returns the directive identifier.
func (v *NonZeroFloat64Validator) Name() string {
	return "!zerof64"
}

// Mode returns the directive evaluation mode.
func (v *NonZeroFloat64Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *NonZeroFloat64Validator) Handle(val float64) (float64, error) {
	_, err := v.Validate(val)
	return val, err
}

// NonZeroFloat64AliasValidator provides the legacy "nonzerof64" tag.
type NonZeroFloat64AliasValidator struct {
	NonZeroFloat64Validator
}

// Name returns the directive identifier.
func (v *NonZeroFloat64AliasValidator) Name() string {
	return "nonzerof64"
}

// NonZeroIntAliasValidator provides the legacy "nonzero" tag.
type NonZeroIntAliasValidator struct {
	NonZeroIntValidator
}

// Name returns the directive identifier.
func (v *NonZeroIntAliasValidator) Name() string {
	return "nonzero"
}

// NonZeroTimeValidator validates that a time.Time is not zero.
type NonZeroTimeValidator struct{}

// Validate checks whether the value is non-zero.
func (v *NonZeroTimeValidator) Validate(val time.Time) (ok bool, err error) {
	if val.IsZero() {
		return false, fmt.Errorf("time is zero")
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *NonZeroTimeValidator) Name() string {
	return "!zerotime"
}

// Mode returns the directive evaluation mode.
func (v *NonZeroTimeValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *NonZeroTimeValidator) Handle(val time.Time) (time.Time, error) {
	_, err := v.Validate(val)
	return val, err
}

// NonZeroTimeAliasValidator provides the legacy "nonzerotime" tag.
type NonZeroTimeAliasValidator struct {
	NonZeroTimeValidator
}

// Name returns the directive identifier.
func (v *NonZeroTimeAliasValidator) Name() string {
	return "nonzerotime"
}

// OneOfFloat64Validator validates that a float64 matches one of the configured values.
type OneOfFloat64Validator struct {
	Values []float64 `param:"values"`
}

// Validate checks whether the value is in the configured set.
func (v *OneOfFloat64Validator) Validate(val float64) (ok bool, err error) {
	return validateOneOf(val, v.Values, errors.New(`value of parameter "values" cannot be empty`), "value %g is not in allowed set")
}

// Name returns the directive identifier.
func (v *OneOfFloat64Validator) Name() string {
	return "oneoff64"
}

// Mode returns the directive evaluation mode.
func (v *OneOfFloat64Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// ConvertParam parses the values parameter.
func (v *OneOfFloat64Validator) ConvertParam(field reflect.StructField, fieldValue reflect.Value, raw string) error {
	if fieldValue.Type() != reflect.TypeOf([]float64{}) {
		return tagex.NewConversionError(field, raw, "[]float64")
	}
	items := splitList(raw)
	vals := make([]float64, 0, len(items))
	for _, item := range items {
		f, err := strconv.ParseFloat(item, 64)
		if err != nil {
			return fmt.Errorf("invalid float %q", item)
		}
		vals = append(vals, f)
	}
	fieldValue.Set(reflect.ValueOf(vals))
	return nil
}

// Handle validates the value and returns it unchanged.
func (v *OneOfFloat64Validator) Handle(val float64) (float64, error) {
	_, err := v.Validate(val)
	return val, err
}

// OneOfStringValidator validates that a string matches one of the configured values.
type OneOfStringValidator struct {
	Values []string `param:"values"`
}

// Validate checks whether the value is in the configured set.
func (v *OneOfStringValidator) Validate(val string) (ok bool, err error) {
	return validateOneOf(val, v.Values, errors.New(`value of parameter "values" cannot be empty`), "value %q is not in allowed set")
}

// Name returns the directive identifier.
func (v *OneOfStringValidator) Name() string {
	return "oneof"
}

// Mode returns the directive evaluation mode.
func (v *OneOfStringValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// ConvertParam parses the values parameter.
func (v *OneOfStringValidator) ConvertParam(field reflect.StructField, fieldValue reflect.Value, raw string) error {
	if fieldValue.Type() != reflect.TypeOf([]string{}) {
		return tagex.NewConversionError(field, raw, "[]string")
	}
	items := splitList(raw)
	fieldValue.Set(reflect.ValueOf(items))
	return nil
}

// Handle validates the value and returns it unchanged.
func (v *OneOfStringValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// OneOfIntValidator validates that an int matches one of the configured values.
type OneOfIntValidator struct {
	Values []int `param:"values"`
}

// Validate checks whether the value is in the configured set.
func (v *OneOfIntValidator) Validate(val int) (ok bool, err error) {
	return validateOneOf(val, v.Values, errors.New(`value of parameter "values" cannot be empty`), "value %d is not in allowed set")
}

// Name returns the directive identifier.
func (v *OneOfIntValidator) Name() string {
	return "oneofint"
}

// Mode returns the directive evaluation mode.
func (v *OneOfIntValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// ConvertParam parses the values parameter.
func (v *OneOfIntValidator) ConvertParam(field reflect.StructField, fieldValue reflect.Value, raw string) error {
	if fieldValue.Type() != reflect.TypeOf([]int{}) {
		return tagex.NewConversionError(field, raw, "[]int")
	}
	items := splitList(raw)
	vals := make([]int, 0, len(items))
	for _, item := range items {
		i, err := strconv.Atoi(item)
		if err != nil {
			return fmt.Errorf("invalid int %q", item)
		}
		vals = append(vals, i)
	}
	fieldValue.Set(reflect.ValueOf(vals))
	return nil
}

// Handle validates the value and returns it unchanged.
func (v *OneOfIntValidator) Handle(val int) (int, error) {
	_, err := v.Validate(val)
	return val, err
}

// PrefixValidator validates that a string has a given prefix.
type PrefixValidator struct {
	Value string `param:"value"`
}

// Validate checks whether the value has the configured prefix.
func (v *PrefixValidator) Validate(val string) (ok bool, err error) {
	if v.Value == "" {
		return false, errors.New(`value of parameter "value" cannot be empty`)
	}
	if !strings.HasPrefix(val, v.Value) {
		return false, fmt.Errorf("value %q does not have prefix %q", val, v.Value)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *PrefixValidator) Name() string {
	return "prefix"
}

// Mode returns the directive evaluation mode.
func (v *PrefixValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *PrefixValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// SuffixValidator validates that a string has a given suffix.
type SuffixValidator struct {
	Value string `param:"value"`
}

// Validate checks whether the value has the configured suffix.
func (v *SuffixValidator) Validate(val string) (ok bool, err error) {
	if v.Value == "" {
		return false, errors.New(`value of parameter "value" cannot be empty`)
	}
	if !strings.HasSuffix(val, v.Value) {
		return false, fmt.Errorf("value %q does not have suffix %q", val, v.Value)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *SuffixValidator) Name() string {
	return "suffix"
}

// Mode returns the directive evaluation mode.
func (v *SuffixValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *SuffixValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// ContainsValidator validates that a string contains a substring.
type ContainsValidator struct {
	Value string `param:"value"`
}

// Validate checks whether the value contains the configured substring.
func (v *ContainsValidator) Validate(val string) (ok bool, err error) {
	if v.Value == "" {
		return false, errors.New(`value of parameter "value" cannot be empty`)
	}
	if !strings.Contains(val, v.Value) {
		return false, fmt.Errorf("value %q does not contain %q", val, v.Value)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *ContainsValidator) Name() string {
	return "contains"
}

// Mode returns the directive evaluation mode.
func (v *ContainsValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *ContainsValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// UUIDValidator validates that a string is a RFC 4122 UUID.
// If Version is 0, version 4 is assumed.
type UUIDValidator struct {
	Version int `param:"version,required=false"`
}

// Validate checks whether the value is a UUID.
func (v *UUIDValidator) Validate(val string) (ok bool, err error) {
	re := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-([0-9a-fA-F])[0-9a-fA-F]{3}-([0-9a-fA-F])[0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	matches := re.FindStringSubmatch(val)
	if matches == nil {
		return false, fmt.Errorf("value %q is not a valid UUID", val)
	}
	versionChar := strings.ToLower(matches[1])
	variantChar := strings.ToLower(matches[2])
	version, err := strconv.ParseInt(versionChar, 16, 0)
	if err != nil {
		return false, fmt.Errorf("invalid UUID version %q", versionChar)
	}
	if variantChar != "8" && variantChar != "9" && variantChar != "a" && variantChar != "b" {
		return false, fmt.Errorf("value %q is not a valid UUID variant", val)
	}
	expected := v.Version
	if expected == 0 {
		expected = 4
	}
	if expected < 1 || expected > 8 {
		return false, fmt.Errorf("invalid UUID version %d", expected)
	}
	if int(version) != expected {
		return false, fmt.Errorf("value %q is not a UUIDv%d", val, expected)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *UUIDValidator) Name() string {
	return "uuid"
}

// Mode returns the directive evaluation mode.
func (v *UUIDValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// ConvertParam parses the version parameter.
func (v *UUIDValidator) ConvertParam(field reflect.StructField, fieldValue reflect.Value, raw string) error {
	if fieldValue.Kind() != reflect.Int {
		return tagex.NewConversionError(field, raw, "int")
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("version cannot be empty")
	}
	ver, err := strconv.Atoi(raw)
	if err != nil {
		return fmt.Errorf("invalid UUID version %q", raw)
	}
	fieldValue.SetInt(int64(ver))
	return nil
}

// Handle validates the value and returns it unchanged.
func (v *UUIDValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// HostnameValidator validates that a string is a valid hostname.
type HostnameValidator struct{}

var hostnamePattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

// Validate checks whether the value is a valid hostname.
func (v *HostnameValidator) Validate(val string) (ok bool, err error) {
	if val == "" {
		return false, fmt.Errorf("value %q is not a valid hostname", val)
	}
	if val == "localhost" {
		return true, nil
	}
	if !hostnamePattern.MatchString(val) {
		return false, fmt.Errorf("value %q is not a valid hostname", val)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *HostnameValidator) Name() string {
	return "hostname"
}

// Mode returns the directive evaluation mode.
func (v *HostnameValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *HostnameValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// IPCIDRValidator validates that a string is a valid CIDR notation.
type IPCIDRValidator struct{}

// Validate checks whether the value is a valid CIDR.
func (v *IPCIDRValidator) Validate(val string) (ok bool, err error) {
	if _, _, err := net.ParseCIDR(val); err != nil {
		return false, fmt.Errorf("invalid CIDR %q: %v", val, err)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *IPCIDRValidator) Name() string {
	return "cidr"
}

// Mode returns the directive evaluation mode.
func (v *IPCIDRValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *IPCIDRValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// Base64Validator validates that a string is valid base64.
type Base64Validator struct{}

// Validate checks whether the value is base64 encoded.
func (v *Base64Validator) Validate(val string) (ok bool, err error) {
	if val == "" {
		return false, fmt.Errorf("value is empty")
	}
	if _, err := base64.StdEncoding.DecodeString(val); err == nil {
		return true, nil
	}
	if _, err := base64.RawStdEncoding.DecodeString(val); err == nil {
		return true, nil
	}
	return false, fmt.Errorf("value %q is not valid base64", val)
}

// Name returns the directive identifier.
func (v *Base64Validator) Name() string {
	return "base64"
}

// Mode returns the directive evaluation mode.
func (v *Base64Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *Base64Validator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// HexValidator validates that a string is valid hex.
type HexValidator struct{}

// Validate checks whether the value is a hex string.
func (v *HexValidator) Validate(val string) (ok bool, err error) {
	if val == "" {
		return false, fmt.Errorf("value is empty")
	}
	clean := strings.TrimPrefix(val, "0x")
	clean = strings.TrimPrefix(clean, "0X")
	if _, err := hex.DecodeString(clean); err != nil {
		return false, fmt.Errorf("value %q is not valid hex", val)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *HexValidator) Name() string {
	return "hex"
}

// Mode returns the directive evaluation mode.
func (v *HexValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// Handle validates the value and returns it unchanged.
func (v *HexValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

// TimeValidator validates that a string matches a time layout.
// If Format is empty, time.RFC3339 is used.
type TimeValidator struct {
	Format string `param:"format,required=false"`
}

// Validate checks whether the value matches the configured layout.
func (v *TimeValidator) Validate(val string) (ok bool, err error) {
	layout := strings.TrimSpace(v.Format)
	if layout == "" {
		layout = time.RFC3339
	}
	if _, err := time.Parse(layout, val); err != nil {
		return false, fmt.Errorf("invalid time %q for layout %q: %v", val, layout, err)
	}
	return true, nil
}

// Name returns the directive identifier.
func (v *TimeValidator) Name() string {
	return "time"
}

// Mode returns the directive evaluation mode.
func (v *TimeValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

// ConvertParam maps well-known time layout names or accepts a raw layout string.
func (v *TimeValidator) ConvertParam(field reflect.StructField, fieldValue reflect.Value, raw string) error {
	if fieldValue.Type() != reflect.TypeOf("") {
		return tagex.NewConversionError(field, raw, "string")
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("format cannot be empty")
	}
	switch raw {
	case "ANSIC":
		raw = time.ANSIC
	case "UnixDate":
		raw = time.UnixDate
	case "RubyDate":
		raw = time.RubyDate
	case "RFC822":
		raw = time.RFC822
	case "RFC822Z":
		raw = time.RFC822Z
	case "RFC850":
		raw = time.RFC850
	case "RFC1123":
		raw = time.RFC1123
	case "RFC1123Z":
		raw = time.RFC1123Z
	case "RFC3339":
		raw = time.RFC3339
	case "RFC3339Nano":
		raw = time.RFC3339Nano
	case "Kitchen":
		raw = time.Kitchen
	case "Stamp":
		raw = time.Stamp
	case "StampMilli":
		raw = time.StampMilli
	case "StampMicro":
		raw = time.StampMicro
	case "StampNano":
		raw = time.StampNano
	}
	fieldValue.SetString(raw)
	return nil
}

// Handle validates the value and returns it unchanged.
func (v *TimeValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

func splitList(raw string) []string {
	parts := strings.Split(raw, "|")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		out = append(out, item)
	}
	return out
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
