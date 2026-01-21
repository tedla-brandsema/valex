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

type CmpRangeValidator[T cmp.Ordered] struct {
	Min T
	Max T
}

func (v *CmpRangeValidator[T]) Validate(val T) (ok bool, err error) {
	if cmp.Less(val, v.Min) || cmp.Less(v.Max, val) {
		return false, fmt.Errorf("value %v is out of range [%v, %v]", val, v.Min, v.Max)
	}
	return true, nil
}

type IntRangeValidator struct {
	Min int `param:"min"`
	Max int `param:"max"`
}

func (v *IntRangeValidator) Validate(val int) (ok bool, err error) {
	if val < v.Min || val > v.Max {
		return false, fmt.Errorf("value %d is out of range [%d, %d]", val, v.Min, v.Max)
	}
	return true, nil
}

func (v *IntRangeValidator) Name() string {
	return "range"
}

func (v *IntRangeValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *IntRangeValidator) Handle(val int) (int, error) {
	_, err := v.Validate(val)
	return val, err
}

type NonNegativeIntValidator struct{}

func (v *NonNegativeIntValidator) Validate(val int) (ok bool, err error) {
	if val < 0 {
		return false, fmt.Errorf("value %d is a negative integer", val)
	}
	return true, nil
}

func (v *NonNegativeIntValidator) Name() string {
	return "pos"
}

func (v *NonNegativeIntValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *NonNegativeIntValidator) Handle(val int) (int, error) {
	_, err := v.Validate(val)
	return val, err
}

type NonPositiveIntValidator struct{}

func (v *NonPositiveIntValidator) Validate(val int) (ok bool, err error) {
	if val > 0 {
		return false, fmt.Errorf("value %d is a positive integer", val)
	}
	return true, nil
}

func (v *NonPositiveIntValidator) Name() string {
	return "neg"
}

func (v *NonPositiveIntValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *NonPositiveIntValidator) Handle(val int) (int, error) {
	_, err := v.Validate(val)
	return val, err
}

type UrlValidator struct{}

func (v *UrlValidator) Validate(val string) (ok bool, err error) {
	_, err = url.ParseRequestURI(val)
	if err == nil {
		ok = true
	}
	return
}

func (v *UrlValidator) Name() string {
	return "url"
}

func (v *UrlValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *UrlValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type EmailValidator struct{}

func (v *EmailValidator) Validate(val string) (ok bool, err error) {
	_, err = mail.ParseAddress(val)
	if err == nil {
		ok = true
	}
	return
}

func (v *EmailValidator) Name() string {
	return "email"
}

func (v *EmailValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *EmailValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type NonEmptyStringValidator struct{}

func (v *NonEmptyStringValidator) Validate(val string) (ok bool, err error) {
	if val == "" {
		return false, fmt.Errorf("string is empty")
	}
	return true, nil
}

func (v *NonEmptyStringValidator) Name() string {
	return "!empty"
}

func (v *NonEmptyStringValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *NonEmptyStringValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type MinLengthValidator struct {
	Size int `param:"size"`
}

func (v *MinLengthValidator) Validate(val string) (ok bool, err error) {
	if v.Size == 0 {
		return false, errors.New(`value of parameter "size" cannot be 0`)
	}
	if len(val) < v.Size {
		return false, fmt.Errorf("value %s exeeds minimum length %d", val, v.Size)
	}
	return true, nil
}

func (v *MinLengthValidator) Name() string {
	return "min"
}

func (v *MinLengthValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *MinLengthValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type MaxLengthValidator struct {
	Size int `param:"size"`
}

func (v *MaxLengthValidator) Validate(val string) (ok bool, err error) {
	if v.Size == 0 {
		return false, errors.New(`value of parameter "size" cannot be 0`)
	}
	if len(val) > v.Size {
		return false, fmt.Errorf("value %s exeeds maximum length %d", val, v.Size)
	}
	return true, nil
}

func (v *MaxLengthValidator) Name() string {
	return "max"
}

func (v *MaxLengthValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *MaxLengthValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type LengthRangeValidator struct {
	Min int `param:"min"`
	Max int `param:"max"`
}

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

func (v *LengthRangeValidator) Name() string {
	return "len"
}

func (v *LengthRangeValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *LengthRangeValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type RegexValidator struct {
	Pattern *regexp.Regexp `param:"pattern"`
}

func (v *RegexValidator) Validate(val string) (ok bool, err error) {
	if v.Pattern == nil {
		return false, errors.New("regex pattern not set")
	}
	if !v.Pattern.MatchString(val) {
		return false, fmt.Errorf("value %q does not match pattern %q", val, v.Pattern.String())
	}
	return true, nil
}

func (v *RegexValidator) Name() string {
	return "regex"
}

func (v *RegexValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

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

func (v *RegexValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type AlphaNumericValidator struct{}

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

func (v *AlphaNumericValidator) Name() string {
	return "alphanum"
}

func (v *AlphaNumericValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *AlphaNumericValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type MACAddressValidator struct{}

func (v *MACAddressValidator) Validate(val string) (ok bool, err error) {
	_, err = net.ParseMAC(val)
	if err != nil {
		return false, fmt.Errorf("invalid MAC address %q: %v", val, err)
	}
	return true, nil
}

func (v *MACAddressValidator) Name() string {
	return "mac"
}

func (v *MACAddressValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *MACAddressValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type IpValidator struct{}

func (v *IpValidator) Validate(val string) (ok bool, err error) {
	if ip := net.ParseIP(val); ip == nil {
		return false, fmt.Errorf("invalid IP address %q", val)
	}
	return true, nil
}

func (v *IpValidator) Name() string {
	return "ip"
}

func (v *IpValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *IpValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type IPv4Validator struct{}

func (v *IPv4Validator) Validate(val string) (ok bool, err error) {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() == nil {
		return false, fmt.Errorf("invalid IPv4 address %q", val)
	}
	return true, nil
}

func (v *IPv4Validator) Name() string {
	return "ipv4"
}

func (v *IPv4Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *IPv4Validator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type IPv6Validator struct{}

func (v *IPv6Validator) Validate(val string) (ok bool, err error) {
	ip := net.ParseIP(val)
	if ip == nil || ip.To4() != nil {
		return false, fmt.Errorf("invalid IPv6 address %q", val)
	}
	return true, nil
}

func (v *IPv6Validator) Name() string {
	return "ipv6"
}

func (v *IPv6Validator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *IPv6Validator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type XMLValidator struct{}

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

		if _, ok := tok.(xml.StartElement); ok { // atleast one tag
			hasElement = true
		}
	}

	if !hasElement {
		return false, fmt.Errorf("XML document must contain at least one element")
	}

	return true, nil
}

func (v *XMLValidator) Name() string {
	return "xml"
}

func (v *XMLValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *XMLValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type JSONValidator struct{}

func (v *JSONValidator) Validate(val string) (ok bool, err error) {
	if !json.Valid([]byte(val)) {
		return false, fmt.Errorf("invalid JSON")
	}
	return true, nil
}

func (v *JSONValidator) Name() string {
	return "json"
}

func (v *JSONValidator) Mode() tagex.DirectiveMode {
	return tagex.EvalMode
}

func (v *JSONValidator) Handle(val string) (string, error) {
	_, err := v.Validate(val)
	return val, err
}

type CompositeValidator[T cmp.Ordered] struct {
	Validators []Validator[T]
}

func (cv *CompositeValidator[T]) Validate(val T) (ok bool, err error) {
	for _, validator := range cv.Validators {
		if ok, err = validator.Validate(val); !ok {
			return false, err
		}
	}
	return true, nil
}
