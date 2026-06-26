package forms

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
)

type fieldDirective struct {
	Key          string `param:"key,required=false"`
	Max          int    `param:"max,default=1"`
	Required     bool   `param:"required,required=false"`
	DefaultValue string `param:"default,required=false"`
}

// ErrFieldRequired is returned when a required form field is missing or empty.
var ErrFieldRequired = errors.New("field is required")

// Validator parses an HTTP request and validates bound structs using the
// valex "val" tag.
type Validator struct {
	rawValues url.Values
	reg       *valex.Registry // nil uses valex's default registry
}

// New parses the request and prepares a Validator that validates against valex's
// default registry.
// ParseForm handles both POST bodies and URL query parameters, so GET requests
// with query values are supported.
func New(r *http.Request) (*Validator, error) {
	return NewWith(r, nil)
}

// NewWith is like New but validates against reg instead of the default registry.
// A nil reg uses the default. Use it for an isolated directive set — test
// isolation, or two differently-configured form validators in one process.
func NewWith(r *http.Request, reg *valex.Registry) (*Validator, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	return &Validator{rawValues: r.Form, reg: reg}, nil
}

// Validate binds form values into dst and validates its "val" tags. It returns
// nil on success, or an *Error carrying an HTTP status code on failure (see Status).
func (v *Validator) Validate(dst any) error {
	if err := bindFormValues(dst, v.rawValues); err != nil {
		return &Error{status: Status(err), Err: err}
	}
	if err := v.validate(dst); err != nil {
		return &Error{status: Status(err), Err: err}
	}
	return nil
}

// validate runs the "val" directives against the Validator's registry, or the
// default registry when none was set.
func (v *Validator) validate(dst any) error {
	if v.reg != nil {
		return v.reg.ValidateStruct(dst)
	}
	return valex.ValidateStruct(dst)
}

// ValidateAll binds form values into dst and validates its "val" tags, collecting
// every binding and validation failure instead of stopping at the first. It
// returns nil on success, or an *Error whose underlying error is an errors.Join
// of the per-field failures — pass it to FieldErrors for a field-keyed map.
func (v *Validator) ValidateAll(dst any) error {
	if _, err := pointerStruct(dst); err != nil {
		return &Error{status: Status(err), Err: err}
	}
	var parts []error
	if bindErr := bindFormValuesAll(dst, v.rawValues); bindErr != nil {
		parts = append(parts, bindErr)
	}
	if valErr := v.validateAll(dst); valErr != nil {
		parts = append(parts, valErr)
	}
	if len(parts) == 0 {
		return nil
	}
	joined := errors.Join(parts...)
	return &Error{status: Status(joined), Err: joined}
}

// validateAll runs the "val" directives in accumulate mode against the
// Validator's registry, or the default registry when none was set.
func (v *Validator) validateAll(dst any) error {
	if v.reg != nil {
		return v.reg.ValidateStructAll(dst)
	}
	return valex.ValidateStructAll(dst)
}

// Bind binds url.Values into a struct pointer using "field" tags, stopping at
// the first error.
func Bind(dst any, values url.Values) error {
	return bindFormValues(dst, values)
}

// bindError is a field-scoped binding failure (type mismatch, too many values,
// or a missing required value). It carries the struct field path so FieldErrors
// can key it, and Unwraps to the underlying cause.
type bindError struct {
	Field string
	Err   error
}

func (e *bindError) Error() string { return fmt.Sprintf("form field %q: %v", e.Field, e.Err) }
func (e *bindError) Unwrap() error { return e.Err }

// bindFormValues binds into dst, stopping at the first field error.
func bindFormValues(dst any, values url.Values) error {
	val, err := pointerStruct(dst)
	if err != nil {
		return err
	}
	return bindStructFields(val, values, "", nil)
}

// bindFormValuesAll binds into dst, accumulating every field error and returning
// them as errors.Join (nil when all fields bind).
func bindFormValuesAll(dst any, values url.Values) error {
	val, err := pointerStruct(dst)
	if err != nil {
		return err
	}
	errs := make([]error, 0)
	_ = bindStructFields(val, values, "", &errs)
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func pointerStruct(v any) (reflect.Value, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return val, fmt.Errorf("expected a pointer to a struct but got %T", v)
	}
	return val.Elem(), nil
}

// bindStructFields binds val's "field"-tagged fields. When errs is nil it stops
// at the first error; when non-nil, each field error accumulates into it (as a
// *bindError) and binding continues.
func bindStructFields(val reflect.Value, values url.Values, path string, errs *[]error) error {
	for n := 0; n < val.NumField(); n++ {
		field := val.Type().Field(n)
		if field.PkgPath != "" {
			continue
		}

		fieldValue := val.FieldByName(field.Name)
		fieldPath := field.Name
		if path != "" {
			fieldPath = path + "." + field.Name
		}

		if _, ok := field.Tag.Lookup("field"); ok {
			if err := bindField(field, fieldValue, values, fieldPath); err != nil {
				if errs == nil {
					return err
				}
				*errs = append(*errs, err)
			}
		}

		switch fieldValue.Kind() {
		case reflect.Struct:
			if err := bindStructFields(fieldValue, values, fieldPath, errs); err != nil {
				return err
			}
		case reflect.Ptr:
			if fieldValue.IsNil() {
				continue
			}
			elem := fieldValue.Elem()
			if elem.Kind() != reflect.Struct {
				continue
			}
			if err := bindStructFields(elem, values, fieldPath, errs); err != nil {
				return err
			}
		}
	}

	return nil
}

// bindField binds one "field"-tagged struct field from values. Input-driven
// failures (type mismatch, too many values, a missing required value) are
// returned as a *bindError keyed by fieldPath — Status maps those to 422 and
// FieldErrors surfaces them. Failures from a malformed field tag itself
// (splitFormTag, ProcessParams) are developer errors, returned unwrapped, so
// Status maps them to 400 and FieldErrors omits them.
func bindField(field reflect.StructField, fieldValue reflect.Value, values url.Values, fieldPath string) error {
	args, err := splitFormTag(field.Tag.Get("field"))
	if err != nil {
		return err
	}
	var directive fieldDirective
	if err := tagex.ProcessParams(&directive, args); err != nil {
		return err
	}

	key := strings.TrimSpace(directive.Key)
	if key == "" {
		key = field.Name
	}

	raw, ok := values[key]
	if !ok || len(raw) == 0 || raw[0] == "" {
		if err := applyDefaultOrRequired(fieldValue, directive); err != nil {
			return &bindError{Field: fieldPath, Err: err}
		}
		return nil
	}
	if err := enforceMax(raw, directive.Max); err != nil {
		return &bindError{Field: fieldPath, Err: err}
	}
	if err := setValueFromRaw(fieldValue, raw); err != nil {
		return &bindError{Field: fieldPath, Err: err}
	}
	return nil
}

func applyDefaultOrRequired(fieldValue reflect.Value, directive fieldDirective) error {
	if directive.Required {
		return ErrFieldRequired
	}
	if strings.TrimSpace(directive.DefaultValue) != "" {
		if err := setValueFromRaw(fieldValue, []string{directive.DefaultValue}); err != nil {
			return err
		}
	}
	return nil
}

func enforceMax(raw []string, max int) error {
	if max <= 0 {
		return fmt.Errorf("invalid max %d", max)
	}
	if len(raw) > max {
		return fmt.Errorf("too many values (%d), max %d", len(raw), max)
	}
	return nil
}

func splitFormTag(tagVal string) (map[string]string, error) {
	parts := strings.Split(tagVal, ",")
	if len(parts) == 0 || strings.TrimSpace(parts[0]) == "" {
		return nil, fmt.Errorf("field tag value is required")
	}

	args := make(map[string]string)
	args["key"] = strings.TrimSpace(parts[0])
	for _, pair := range parts[1:] {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("malformed key value pair %q, expected format is \"key=value\"", pair)
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		if key == "" || val == "" {
			return nil, fmt.Errorf("malformed key value pair %q, expected format is \"key=value\"", pair)
		}
		args[key] = val
	}
	return args, nil
}

func setValueFromRaw(fieldValue reflect.Value, raw []string) error {
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		}
		return setValueFromRaw(fieldValue.Elem(), raw)
	}

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(raw[0])
		return nil
	case reflect.Bool:
		b, err := strconv.ParseBool(raw[0])
		if err != nil {
			return err
		}
		fieldValue.SetBool(b)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(raw[0], 10, fieldValue.Type().Bits())
		if err != nil {
			return err
		}
		fieldValue.SetInt(i)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u, err := strconv.ParseUint(raw[0], 10, fieldValue.Type().Bits())
		if err != nil {
			return err
		}
		fieldValue.SetUint(u)
		return nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(raw[0], fieldValue.Type().Bits())
		if err != nil {
			return err
		}
		fieldValue.SetFloat(f)
		return nil
	case reflect.Slice:
		return setSliceFromRaw(fieldValue, raw)
	default:
		return fmt.Errorf("unsupported field type %s", fieldValue.Type())
	}
}

func setSliceFromRaw(fieldValue reflect.Value, raw []string) error {
	elemType := fieldValue.Type().Elem()
	slice := reflect.MakeSlice(fieldValue.Type(), 0, len(raw))
	for _, item := range raw {
		elem := reflect.New(elemType).Elem()
		if err := setValueFromRaw(elem, []string{item}); err != nil {
			return err
		}
		slice = reflect.Append(slice, elem)
	}
	fieldValue.Set(slice)
	return nil
}
