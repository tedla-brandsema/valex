package forms_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/valex/forms"
	"github.com/tedla-brandsema/valex/validators"
)

// formStore is a stand-in for a database; here it just appends records to a file.
type formStore struct{ path string }

func (s *formStore) save(record string) error {
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, record)
	return err
}

// Registration is a form that hooks into the processing lifecycle. The exported
// fields are bound and validated; the unexported fields carry per-request state
// and are ignored by binding.
type Registration struct {
	Name  string `field:"name" val:"min,size=2"`
	Email string `field:"email" val:"email"`

	store     *formStore // dependency, not a form field
	rejection string     // populated by the Failure hook
}

// Compile-time proof that Registration satisfies the tagex lifecycle hooks. valex
// runs validation through tagex, which invokes these if the target implements them.
var (
	_ tagex.PreProcessor         = (*Registration)(nil)
	_ tagex.SuccessPostProcessor = (*Registration)(nil)
	_ tagex.FailurePostProcessor = (*Registration)(nil)
)

// Before runs after binding but before the "val" directives, so it is the place
// to normalize input. (tagex.PreProcessor)
func (r *Registration) Before() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	return nil
}

// Success runs only after every directive passes — commit the record here.
// (tagex.SuccessPostProcessor)
func (r *Registration) Success() error {
	return r.store.save(fmt.Sprintf("%s <%s>", r.Name, r.Email))
}

// Failure runs when a directive fails. cause is the validation error; a real
// handler might log it and craft a response. Returning nil keeps cause as the
// error returned to the caller (so forms still maps it to an HTTP status).
// (tagex.FailurePostProcessor)
func (r *Registration) Failure(cause error) error {
	r.rejection = "please correct the highlighted fields"
	return nil
}

// This advanced example combines valex validation with tagex's processing
// lifecycle hooks to handle a form end to end: normalize on Before, persist on
// Success (to disk, standing in for a database), and prepare a response on
// Failure. Note that the form type is coupled to both valex (the "val" tag) and
// tagex (the hook interfaces).
func Example_lifecycleHooks() {
	valex.RegisterDirective(&validators.MinLengthValidator{})
	valex.RegisterDirective(&validators.EmailValidator{})

	dir, err := os.MkdirTemp("", "valex-hooks")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	store := &formStore{path: filepath.Join(dir, "registrations.txt")}

	submit := func(n int, name, email string) {
		form := url.Values{"name": {name}, "email": {email}}
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		fmt.Printf("submit #%d: name=%q email=%q\n", n, name, email)
		reg := Registration{store: store}
		if ok, err := forms.Validate(req, &reg); ok {
			fmt.Printf("  -> accepted; Success persisted: %s <%s>\n", reg.Name, reg.Email)
		} else {
			var ferr *forms.Error
			errors.As(err, &ferr)
			fmt.Printf("  -> rejected (HTTP %d); Failure said: %s\n", ferr.StatusCode(), reg.rejection)
		}
	}

	submit(1, "  Ada  ", "ADA@example.com") // normalized by Before, persisted by Success
	submit(2, "B", "not-an-email")          // fails validation; Failure handles it

	saved, _ := os.ReadFile(store.path)
	fmt.Printf("database contains %d record(s):\n", strings.Count(string(saved), "\n"))
	fmt.Print(string(saved))
	// Output:
	// submit #1: name="  Ada  " email="ADA@example.com"
	//   -> accepted; Success persisted: Ada <ada@example.com>
	// submit #2: name="B" email="not-an-email"
	//   -> rejected (HTTP 422); Failure said: please correct the highlighted fields
	// database contains 1 record(s):
	// Ada <ada@example.com>
}
