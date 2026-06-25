// Package valex is a small value-validation engine built on the tagex
// struct-tag processor.
//
// It supports two primary workflows:
//
//  1. Programmatic validation using Validator or ValidatorFunc, with
//     ValidatedValue for guarded assignment and MustValidate for fail-fast use.
//  2. Struct-tag validation using the "val" tag and ValidateStruct. Register
//     directives with RegisterDirective; pass additional tagex.Tag values to
//     ValidateStruct to process multiple tags in a single pass.
//
// The engine ships no directives of its own. Ready-made validators live in the
// github.com/tedla-brandsema/valex/validators subpackage; register the ones you
// need with RegisterDirective. HTTP request binding and validation live in the
// github.com/tedla-brandsema/valex/forms subpackage, which keeps net/http out of
// the core engine.
//
// # Concurrency
//
// RegisterDirective and ValidateStruct are safe for concurrent use: the "val"
// tag's directive registry is guarded by a mutex, and ValidateStruct only reads
// it. The intended pattern is to register directives once at startup (typically
// in an init function) and validate from many goroutines thereafter. Registering
// while other goroutines validate is safe but unusual.
package valex
