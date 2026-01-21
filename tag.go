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

func ValidateStruct(data interface{}) (bool, error) {
	return tag.ProcessStruct(data)
}
