// Package valex is a small value-validation engine built on the tagex
// struct-tag processor.
//
// Requires Go 1.22 or later.
//
// It supports two primary workflows:
//
//  1. Programmatic validation using Validator or ValidatorFunc, with
//     ValidatedValue for guarded assignment and MustValidate for fail-fast use.
//  2. Struct-tag validation using the "val" tag and ValidateStruct. Register
//     directives with MustRegisterDirective (or RegisterDirective, which returns
//     an error instead of panicking); pass additional tagex.Tag values to
//     ValidateStruct to process multiple tags in a single pass.
//
// The engine ships no directives of its own. Ready-made validators live in the
// github.com/tedla-brandsema/valex/validators subpackage; register the ones you
// need with MustRegisterDirective. HTTP request binding and validation live in the
// github.com/tedla-brandsema/valex/forms subpackage, which keeps net/http out of
// the core engine.
//
// # Concurrency
//
// Registering a directive and ValidateStruct are safe for concurrent use: the "val"
// tag's directive registry is guarded by a mutex, and ValidateStruct only reads
// it. The intended pattern is to register directives once at startup (typically
// in an init function) and validate from many goroutines thereafter. Registering
// while other goroutines validate is safe but unusual.
//
// # Registries
//
// The package-level RegisterDirective, MustRegisterDirective, and ValidateStruct
// share one process-global default registry — like flag.CommandLine or
// http.DefaultServeMux, it belongs to the application, which registers once at
// startup and validates anywhere.
//
// Libraries, and tests that need isolation, should create their own with
// NewRegistry instead of touching the global. Each Registry has an independent
// directive set, so two can hold the same directive name without colliding.
// Register on one with the free functions RegisterDirectiveTo /
// MustRegisterDirectiveTo (free functions because Go methods can't be generic),
// and validate with its ValidateStruct method.
package valex
