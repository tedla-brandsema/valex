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
