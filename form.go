package valex

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/tedla-brandsema/tagex"
)

type FieldDirective struct {
	Key          string   `param:"key"`
	Max          int      `param:"max"`
	Values       []string `param:"values"`
	Required     bool     `param:"required"`
	DefaultValue string   `param:"default"`
}

var ErrFieldRequired = errors.New("field is required")

func (d *FieldDirective) Name() string {
	return "field"
}

func (d *FieldDirective) Mode() tagex.DirectiveMode {
	return tagex.MutMode
}

func (d *FieldDirective) Handle(val string) (string, error) {
	if val == "" && d.Required {
	 	// Do not set DefaultValue here
		return val, &tagex.HandleError{Nested: ErrFieldRequired}
	}
	if val == "" {
		val = d.DefaultValue
	}

	return val, nil
}

type FormValidator struct {
	tags      []*tagex.Tag
	rawValues url.Values
}

func NewFormValidator(r *http.Request) (*FormValidator, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	formTag := tagex.NewTag("form")
	tagex.RegisterDirective(&formTag, &FieldDirective{})

	tagex.RegisterDirective(&formTag, &NonEmptyStringValidator{})
	return &FormValidator{
			tags: []*tagex.Tag{
				&formTag,
				&tag,
			},
			rawValues: r.Form,
		},
		nil
}

func (v *FormValidator) Validate(dst any) (bool, error) {
	return tagex.ProcessStruct(dst, v.tags...)
}
