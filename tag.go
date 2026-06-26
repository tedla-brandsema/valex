package valex

import (
	"github.com/tedla-brandsema/tagex"
)

const tagKey = "val"

// Registry is an isolated set of "val" directives. Most programs use the
// package-level RegisterDirective / MustRegisterDirective / ValidateStruct, which
// share one process-global default registry. Create your own with NewRegistry
// when you need an independent one — for test isolation, or to run two
// differently-configured validators in the same process.
type Registry struct {
	tag *tagex.Tag
}

// NewRegistry returns a new, empty Registry with its own directive set.
func NewRegistry() *Registry {
	return &Registry{tag: tagex.NewTag(tagKey)}
}

// defaultRegistry backs the package-level functions.
var defaultRegistry = NewRegistry()

// ValidateStruct validates struct fields against the registry's "val" directives.
// It returns nil when the struct is valid. Additional tagex.Tag values can be
// provided to process more tags in the same pass.
func (r *Registry) ValidateStruct(data any, tags ...*tagex.Tag) error {
	return tagex.ProcessStruct(data, append(tags, r.tag)...)
}

// ValidateStructAll is like ValidateStruct but does not stop at the first
// failure: it validates every field and returns errors.Join of the per-field
// errors (nil when all pass). Use FieldErrors to turn the result into a map
// keyed by field path.
func (r *Registry) ValidateStructAll(data any, tags ...*tagex.Tag) error {
	return tagex.ProcessStructAll(data, append(tags, r.tag)...)
}

// RegisterDirectiveTo registers a directive on r. It is a free function rather
// than a method because Go methods cannot have type parameters. It returns
// *EmptyDirectiveNameError if the directive's Name is blank, or
// *DuplicateDirectiveError if that name is already registered on r; use
// MustRegisterDirectiveTo to panic on these instead.
func RegisterDirectiveTo[T any](r *Registry, d tagex.Directive[T]) error {
	return tagex.RegisterDirective(r.tag, d)
}

// MustRegisterDirectiveTo is like RegisterDirectiveTo but panics if registration
// fails — the convenient choice for registering directives once at startup.
func MustRegisterDirectiveTo[T any](r *Registry, d tagex.Directive[T]) {
	tagex.MustRegisterDirective(r.tag, d)
}

// ValidateStruct validates struct fields using the default registry's "val"
// directives. It returns nil when the struct is valid. Additional tagex.Tag
// values can be provided to process more tags in the same pass.
func ValidateStruct(data any, tags ...*tagex.Tag) error {
	return defaultRegistry.ValidateStruct(data, tags...)
}

// ValidateStructAll is like ValidateStruct but does not stop at the first
// failure: it validates every field against the default registry and returns
// errors.Join of the per-field errors (nil when all pass). Use FieldErrors to
// turn the result into a map keyed by field path.
func ValidateStructAll(data any, tags ...*tagex.Tag) error {
	return defaultRegistry.ValidateStructAll(data, tags...)
}

// RegisterDirective registers a directive on the default registry for use with
// the "val" struct tag. It returns *EmptyDirectiveNameError if the directive's
// Name is blank, or *DuplicateDirectiveError if that name is already registered
// (it does not overwrite). Use MustRegisterDirective to panic on these instead,
// which is usually what you want when registering at startup.
func RegisterDirective[T any](d tagex.Directive[T]) error {
	return RegisterDirectiveTo(defaultRegistry, d)
}

// MustRegisterDirective is like RegisterDirective but panics if registration
// fails. It is the convenient choice for registering directives once at startup
// (typically in an init function), where a duplicate or empty directive name is
// a programming error that should fail fast.
func MustRegisterDirective[T any](d tagex.Directive[T]) {
	MustRegisterDirectiveTo(defaultRegistry, d)
}
