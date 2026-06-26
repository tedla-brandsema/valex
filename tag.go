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

// RegisterDirective registers a directive for use with the "val" struct tag. It
// returns *EmptyDirectiveNameError if the directive's Name is blank, or
// *DuplicateDirectiveError if that name is already registered (it does not
// overwrite). Use MustRegisterDirective to panic on these instead, which is
// usually what you want when registering at startup.
func RegisterDirective[T any](d tagex.Directive[T]) error {
	return tagex.RegisterDirective(tag, d)
}

// MustRegisterDirective is like RegisterDirective but panics if registration
// fails. It is the convenient choice for registering directives once at startup
// (typically in an init function), where a duplicate or empty directive name is
// a programming error that should fail fast.
func MustRegisterDirective[T any](d tagex.Directive[T]) {
	tagex.MustRegisterDirective(tag, d)
}
