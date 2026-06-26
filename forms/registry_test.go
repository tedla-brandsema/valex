package forms_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/valex/forms"
)

// onlyRegDirective is registered only on a custom Registry, never on the global
// default, so a test can prove which registry validation actually used.
type onlyRegDirective struct{}

func (*onlyRegDirective) Name() string              { return "forms_test_onlyreg" }
func (*onlyRegDirective) Mode() tagex.DirectiveMode { return tagex.EvalMode }
func (*onlyRegDirective) Handle(s string) (string, error) {
	if s == "bad" {
		return s, errors.New("value is bad")
	}
	return s, nil
}

func newReq(value string) *http.Request {
	form := url.Values{"f": {value}}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestValidateWithIsolatedRegistry(t *testing.T) {
	type Form struct {
		F string `field:"f" val:"forms_test_onlyreg"`
	}

	reg := valex.NewRegistry()
	valex.MustRegisterDirectiveTo(reg, &onlyRegDirective{})

	// ValidateWith uses reg, which knows the directive: "bad" fails as a
	// validation error (HTTP 422), and the *forms.Error status wrapping is kept.
	var bad Form
	err := forms.ValidateWith(newReq("bad"), &bad, reg)
	if err == nil {
		t.Fatal("expected validation failure from reg's directive")
	}
	var ferr *forms.Error
	if !errors.As(err, &ferr) || ferr.StatusCode() != http.StatusUnprocessableEntity {
		t.Fatalf("want *forms.Error with 422, got %v", err)
	}

	// "ok" passes on reg.
	var ok Form
	if err := forms.ValidateWith(newReq("ok"), &ok, reg); err != nil {
		t.Fatalf("reg should accept \"ok\", got %v", err)
	}
	if ok.F != "ok" {
		t.Fatalf("binding failed: F=%q", ok.F)
	}

	// The default-registry path does NOT know the directive — proving ValidateWith
	// used reg, not the global.
	var viaDefault Form
	if err := forms.Validate(newReq("ok"), &viaDefault); err == nil {
		t.Fatal("default registry should not know forms_test_onlyreg")
	}
}
