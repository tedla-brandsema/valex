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
}

// ValidateStruct validates struct fields using the "val" tag directives. It
// returns nil when the struct is valid. Additional tagex.Tag values can be
// provided to process more tags in the same pass.
func ValidateStruct(data any, tags ...*tagex.Tag) error {
	tags = append(tags, &tag)
	_, err := tagex.ProcessStruct(data, tags...)
	return err
}

// RegisterDirective registers a directive for use with the "val" struct tag.
func RegisterDirective[T any](d tagex.Directive[T]) {
	// Do not add mutex here; it is handled in tagex
	tagex.RegisterDirective(&tag, d)
}
