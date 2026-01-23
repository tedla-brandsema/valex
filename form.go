package valex

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/tedla-brandsema/tagex"
)

type FieldDirective struct {
	Key          string `param:"key,required=false"`
	Max          int    `param:"max,default=1"`
	Required     bool   `param:"required,required=false"`
	DefaultValue string `param:"default,required=false"`
}

var ErrFieldRequired = errors.New("field is required")

type FormValidator struct {
	tags      []*tagex.Tag
	rawValues url.Values
}

// NewFormValidator parses the request and prepares a validator.
// ParseForm handles both POST bodies and URL query parameters, so GET requests
// with query values are supported.
func NewFormValidator(r *http.Request) (*FormValidator, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	return &FormValidator{
			tags: []*tagex.Tag{
				&tag,
			},
			rawValues: r.Form,
		},
		nil
}

// Validate binds form values into dst and validates "val" tags.
func (v *FormValidator) Validate(dst any) (bool, error) {
	if err := bindFormValues(dst, v.rawValues); err != nil {
		return false, err
	}
	return tagex.ProcessStruct(dst, v.tags...)
}

func bindFormValues(dst any, values url.Values) error {
	val, err := pointerStruct(dst)
	if err != nil {
		return err
	}
	return bindStructFields(val, values, "")
}

// BindFormValues binds url.Values into a struct pointer using "field" tags.
func BindFormValues(dst any, values url.Values) error {
	return bindFormValues(dst, values)
}

func pointerStruct(v any) (reflect.Value, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return val, fmt.Errorf("expected a pointer to a struct but got %T", v)
	}
	return val.Elem(), nil
}

func bindStructFields(val reflect.Value, values url.Values, path string) error {
	for n := 0; n < val.NumField(); n++ {
		field := val.Type().Field(n)
		if field.PkgPath != "" {
			continue
		}

		fieldValue := val.FieldByName(field.Name)
		if tagValue, ok := field.Tag.Lookup("field"); ok {
			args, err := splitFormTag(tagValue)
			if err != nil {
				return wrapFormFieldError(path, field.Name, err)
			}
			var directive FieldDirective
			if err := tagex.ProcessParams(&directive, args); err != nil {
				return wrapFormFieldError(path, field.Name, err)
			}

			key := strings.TrimSpace(directive.Key)
			if key == "" {
				key = field.Name
			}

			raw, ok := values[key]
			if !ok || len(raw) == 0 || raw[0] == "" {
				if err := applyDefaultOrRequired(fieldValue, directive, path, field.Name); err != nil {
					return err
				}
			} else {
				if err := enforceMax(raw, directive.Max); err != nil {
					return wrapFormFieldError(path, field.Name, err)
				}
				if err := setValueFromRaw(fieldValue, raw); err != nil {
					return wrapFormFieldError(path, field.Name, err)
				}
			}
		}

		switch fieldValue.Kind() {
		case reflect.Struct:
			nextPath := field.Name
			if path != "" {
				nextPath = path + "." + field.Name
			}
			if err := bindStructFields(fieldValue, values, nextPath); err != nil {
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
			nextPath := field.Name
			if path != "" {
				nextPath = path + "." + field.Name
			}
			if err := bindStructFields(elem, values, nextPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func applyDefaultOrRequired(fieldValue reflect.Value, directive FieldDirective, path, fieldName string) error {
	if directive.Required {
		return wrapFormFieldError(path, fieldName, ErrFieldRequired)
	}
	if strings.TrimSpace(directive.DefaultValue) != "" {
		if err := setValueFromRaw(fieldValue, []string{directive.DefaultValue}); err != nil {
			return wrapFormFieldError(path, fieldName, err)
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

func wrapFormFieldError(path, fieldName string, err error) error {
	if path == "" {
		return fmt.Errorf("form field %q: %w", fieldName, err)
	}
	return fmt.Errorf("form field %q: %w", path+"."+fieldName, err)
}
