package forms_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/valex/forms"
	"github.com/tedla-brandsema/valex/validators"
)

// Validate parses a request, binds its values into a struct, and validates the
// "val" tags, returning an *Error with an HTTP status code on failure.
func ExampleValidate() {
	valex.RegisterDirective(&validators.MinLengthValidator{})

	type Signup struct {
		Name string `field:"name" val:"min,size=3"`
	}

	submit := func(name string) {
		form := url.Values{"name": {name}}
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		var in Signup
		if err := forms.Validate(req, &in); err == nil {
			fmt.Printf("%q: ok\n", in.Name)
		} else {
			var ferr *forms.Error
			errors.As(err, &ferr)
			fmt.Printf("%q: rejected with status %d\n", name, ferr.StatusCode())
		}
	}

	submit("Gopher")
	submit("Al")
	// Output:
	// "Gopher": ok
	// "Al": rejected with status 422
}
