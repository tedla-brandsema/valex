// Package valex provides generic validators and tag-based validation directives.
//
// It supports two primary workflows:
//
//  1. Programmatic validation using Validator or ValidatorFunc along with
//     ValidatedValue for guarded assignment.
//  2. Struct tag validation using the "val" tag and ValidateStruct.
//     You can pass additional tagex.Tag values to ValidateStruct to process
//     multiple tags in a single pass.
//
// For HTTP form binding, FormValidator parses requests and binds "field" tags
// before running validation, and ValidateForm provides a convenience wrapper
// with HTTP status mapping.
//
// The built-in directives cover common validations (ranges, lengths, URLs,
// emails, IPs, JSON/XML, and regex). You can extend tag validation by registering
// custom directives with RegisterDirective.
package valex
