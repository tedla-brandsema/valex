// Package forms binds HTTP request values into structs and validates them with
// the valex engine's "val" tag.
//
// It is a separate package from the core valex engine so that programs which
// only need programmatic or struct-tag validation do not pull in net/http.
//
// # Binding and validation
//
// Two struct tags are involved. The "field" tag maps a struct field to a request
// value and controls binding; the "val" tag (from the valex engine) validates the
// bound value. Validation directives are opt-in and must be registered with
// valex.MustRegisterDirective — for example from the valex/validators subpackage.
//
//	type Signup struct {
//		Name  string `field:"name" val:"min,size=3"`
//		Email string `field:"email" val:"email"`
//	}
//
// New calls request.ParseForm, which reads both POST bodies and URL query
// parameters, so GET requests are supported. Validate is a convenience wrapper
// that parses, binds, validates, and returns a *Error carrying an HTTP status
// code. Bind binds url.Values without validating, for use outside HTTP handlers.
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//		var in Signup
//		if err := forms.Validate(r, &in); err != nil {
//			var ferr *forms.Error
//			errors.As(err, &ferr)
//			http.Error(w, err.Error(), ferr.StatusCode())
//			return
//		}
//		// ... use in
//	}
//
// # The field tag
//
// The first (positional) value is the request key; the remaining comma-separated
// options are key=value pairs:
//
//	Option    Default  Description
//	--------- -------- ---------------------------------------------------------
//	(key)     field    request key to read; defaults to the struct field name
//	max       1        maximum number of values accepted (for slice fields)
//	required  false    report ErrFieldRequired when the value is missing or empty
//	default   -        value to bind when the field is missing or empty
//
// # Errors
//
// Validate wraps failures in *Error, whose StatusCode reports an HTTP status:
// 400 for binding and parse problems, 422 for validation failures and missing
// required fields. Status exposes the same mapping for an arbitrary error, and
// ErrFieldRequired is returned for missing required fields. Validation failures
// are the error types re-exported by the valex package, so they can be inspected
// with errors.As / errors.Is without importing tagex.
package forms
