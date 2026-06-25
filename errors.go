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

	// UnknownDirectiveError is returned when a tag names an unregistered directive.
	UnknownDirectiveError = tagex.UnknownDirectiveError
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
	StagePre       = tagex.StagePre
	StageDirective = tagex.StageDirective
	StageParam     = tagex.StageParam
	StagePost      = tagex.StagePost
)
