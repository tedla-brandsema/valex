// Forms binds an net/http request into a struct and validates its "val" tags.
// It uses httptest so the program is self-contained and runnable.
package main

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

// The "field" tag binds a request value; the "val" tag validates it.
type Signup struct {
	Name  string `field:"name"  val:"min,size=3"`
	Email string `field:"email" val:"email"`
}

func init() {
	valex.RegisterDirective(&validators.MinLengthValidator{})
	valex.RegisterDirective(&validators.EmailValidator{})
}

func submit(name, email string) {
	form := url.Values{"name": {name}, "email": {email}}
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var in Signup
	if err := forms.Validate(req, &in); err != nil {
		var ferr *forms.Error
		errors.As(err, &ferr)
		fmt.Printf("name=%q rejected (HTTP %d)\n", name, ferr.StatusCode())
		return
	}
	fmt.Printf("name=%q accepted\n", in.Name)
}

func main() {
	submit("Gopher", "gopher@example.com")
	submit("Al", "gopher@example.com")
	submit("Gopher", "not-an-email")

	// Output:
	// name="Gopher" accepted
	// name="Al" rejected (HTTP 422)
	// name="Gopher" rejected (HTTP 422)
}
