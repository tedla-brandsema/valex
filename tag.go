package valex

import (
	"github.com/tedla-brandsema/tagex"
)

const tagKey = "val"

var (
	tag tagex.Tag
)

func init() {
	tag = tagex.NewTag(tagKey)

	// Int directives
	tagex.RegisterDirective(&tag, &IntRangeValidator{})
	tagex.RegisterDirective(&tag, &NonNegativeIntValidator{})
	tagex.RegisterDirective(&tag, &NonPositiveIntValidator{})

	// String directives
	tagex.RegisterDirective(&tag, &UrlValidator{})
	tagex.RegisterDirective(&tag, &EmailValidator{})
	tagex.RegisterDirective(&tag, &NonEmptyStringValidator{})
	tagex.RegisterDirective(&tag, &MinLengthValidator{})
	tagex.RegisterDirective(&tag, &MaxLengthValidator{})
	tagex.RegisterDirective(&tag, &LengthRangeValidator{})
	tagex.RegisterDirective(&tag, &RegexValidator{})
	tagex.RegisterDirective(&tag, &AlphaNumericValidator{})
	tagex.RegisterDirective(&tag, &MACAddressValidator{})
	tagex.RegisterDirective(&tag, &IpValidator{})
	tagex.RegisterDirective(&tag, &IPv4Validator{})
	tagex.RegisterDirective(&tag, &IPv6Validator{})
	tagex.RegisterDirective(&tag, &XMLValidator{})
	tagex.RegisterDirective(&tag, &JSONValidator{})
}

// ValidateStruct validates struct fields using the "val" tag directives.
// Additional tagex.Tag values can be provided to process more tags in the same pass.
func ValidateStruct(data interface{}, tags ...*tagex.Tag) (bool, error) {
	tags = append(tags, &tag)
	return tagex.ProcessStruct(data, tags...)
}

// RegisterDirective registers a directive for use with the "val" struct tag.
func RegisterDirective[T any](d tagex.Directive[T]) {
	// Do not add mutex here; it is handled in tagex
	tagex.RegisterDirective(&tag, d)
}
