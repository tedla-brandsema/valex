package valex

import (
	"errors"

	"github.com/tedla-brandsema/tagex"
)

// ErrNoValidator is returned by ValidatedValue.Set when no Validator is configured.
var ErrNoValidator = errors.New("valex: no validator set")

// The types below are re-exported from tagex. ValidateStruct — and the
// valex/forms helpers built on it — return these on failure, so callers can
// inspect them with errors.As / errors.Is without importing tagex directly:
//
//	var convErr *valex.ConversionError
//	if errors.As(err, &convErr) {
//		// ...
//	}
//
// They are type aliases, so a *valex.ConversionError and a *tagex.ConversionError
// are the same type; only the import path differs.
type (
	// Stage identifies the processing stage at which an error occurred.
	Stage = tagex.Stage

	// ProcessError describes a failure while processing a single struct field.
	ProcessError = tagex.ProcessError
	// TagError wraps all errors produced for a given struct-tag key.
	TagError = tagex.TagError
	// HandleError wraps an error returned by a directive's Handle method.
	HandleError = tagex.HandleError
	// HookError wraps an error returned by a pre- or post-processing hook.
	HookError = tagex.HookError

	// InvalidTargetError is returned when ValidateStruct gets a value that is not a pointer to a struct.
	InvalidTargetError = tagex.InvalidTargetError
	// NilTagError is returned when a nil tag is processed.
	NilTagError = tagex.NilTagError
	// MaxDepthError is returned when processing recurses past the nesting limit (usually cyclic data).
	MaxDepthError = tagex.MaxDepthError

	// UnknownDirectiveError is returned when a tag names an unregistered directive.
	UnknownDirectiveError = tagex.UnknownDirectiveError
	// EmptyDirectiveNameError is returned by RegisterDirective for a directive with a blank Name.
	EmptyDirectiveNameError = tagex.EmptyDirectiveNameError
	// DuplicateDirectiveError is returned by RegisterDirective when the directive name is already registered.
	DuplicateDirectiveError = tagex.DuplicateDirectiveError
	// DirectiveParseError is returned when a tag value omits the directive name.
	DirectiveParseError = tagex.DirectiveParseError
	// ParamParseError is returned for a malformed "key=value" parameter pair.
	ParamParseError = tagex.ParamParseError
	// MissingParamError is returned when a required parameter is absent.
	MissingParamError = tagex.MissingParamError
	// TypeMismatchError is returned when a field's type does not match the directive.
	TypeMismatchError = tagex.TypeMismatchError
	// ConversionError is returned when a raw parameter value cannot be converted.
	ConversionError = tagex.ConversionError
	// UnsupportedParamTypeError is returned for a parameter of an unsupported type.
	UnsupportedParamTypeError = tagex.UnsupportedParamTypeError
	// ParamConflictError is returned when a parameter sets both required and default.
	ParamConflictError = tagex.ParamConflictError
	// FieldAccessError is returned when a struct field cannot be read via reflection.
	FieldAccessError = tagex.FieldAccessError
	// FieldSetError is returned when a struct field cannot be set via reflection.
	FieldSetError = tagex.FieldSetError
)

// Processing stages, re-exported from tagex for use with Stage and ProcessError.
const (
	StageInput     = tagex.StageInput
	StagePre       = tagex.StagePre
	StageDirective = tagex.StageDirective
	StageParam     = tagex.StageParam
	StagePost      = tagex.StagePost
	StageStruct    = tagex.StageStruct
)

// FieldErrors flattens an error returned by ValidateStructAll into a map from
// field path to the error for that field, keeping the first error seen per
// field.
//
// The keys are struct field paths — the same values ProcessError.FieldPath
// carries, e.g. "Email" or "Items[2].SKU" — not request keys or display names;
// translate them yourself when rendering. It walks errors.Join trees, so it
// works on the accumulated result of ValidateStructAll. A nil error yields a nil
// map. Field-less errors (such as *InvalidTargetError) are omitted, so the
// original error stays authoritative — check err != nil first, then render the
// map on top.
func FieldErrors(err error) map[string]error {
	if err == nil {
		return nil
	}
	m := make(map[string]error)
	collectFieldErrors(err, m)
	if len(m) == 0 {
		return nil
	}
	return m
}

func collectFieldErrors(err error, m map[string]error) {
	if j, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range j.Unwrap() {
			collectFieldErrors(e, m)
		}
		return
	}
	var pe *ProcessError
	if errors.As(err, &pe) && pe.FieldPath != "" {
		if _, exists := m[pe.FieldPath]; !exists {
			m[pe.FieldPath] = err
		}
	}
}
