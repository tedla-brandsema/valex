// Package valex provides generic validators and tag-based validation directives.
//
// It supports two primary workflows:
//
//  1. Programmatic validation using Validator or ValidatorFunc along with
//     ValidatedValue for guarded assignment.
//  2. Struct tag validation using the "val" tag and ValidateStruct.
//
// The built-in directives cover common validations (ranges, lengths, URLs,
// emails, IPs, JSON/XML, and regex). You can extend tag validation by registering
// custom directives with RegisterDirective.
package valex
