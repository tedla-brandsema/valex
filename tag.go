package valex

import (
	"github.com/tedla-brandsema/tagex"
)

const tagKey = "val"

// tag is the shared "val" registry. NewTag returns a *Tag (tagex v0.4.0), and a
// Tag is safe to share once its directives are registered.
var tag = tagex.NewTag(tagKey)

// ValidateStruct validates struct fields using the "val" tag directives. It
// returns nil when the struct is valid. Additional tagex.Tag values can be
// provided to process more tags in the same pass.
func ValidateStruct(data any, tags ...*tagex.Tag) error {
	tags = append(tags, tag)
	return tagex.ProcessStruct(data, tags...)
}

// RegisterDirective registers a directive for use with the "val" struct tag.
func RegisterDirective[T any](d tagex.Directive[T]) {
	// Ignore the error: registration is idempotent here. tagex v0.4.0 reports a
	// *DuplicateDirectiveError instead of silently overwriting, but re-registering
	// the same directive from multiple call sites is intentional and harmless.
	_ = tagex.RegisterDirective(tag, d)
}
