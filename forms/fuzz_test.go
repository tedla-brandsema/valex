package forms_test

import (
	"net/url"
	"testing"

	"github.com/tedla-brandsema/valex/forms"
)

type fuzzNested struct {
	N string `field:"n"`
}

// fuzzTarget exercises every branch of the binding type switch: scalars of each
// kind, slices (with a max), a pointer, required/default options, and a nested
// struct (recursion).
type fuzzTarget struct {
	S   string   `field:"s"`
	I   int      `field:"i"`
	I8  int8     `field:"i8"`
	U   uint     `field:"u"`
	F   float64  `field:"f"`
	B   bool     `field:"b"`
	Sl  []string `field:"sl,max=3"`
	ISl []int    `field:"isl,max=3"`
	Ptr *int     `field:"ptr"`
	Req string   `field:"req,required=true"`
	Def string   `field:"def,default=x"`
	Sub fuzzNested
}

// FuzzBind throws arbitrary request values at the binding path — the only place
// untrusted input enters valex/forms — and asserts it never panics. Binding
// errors (type mismatches, too-many-values, missing-required) are expected and
// ignored; only a panic fails the test.
func FuzzBind(f *testing.F) {
	f.Add("s=hi&i=3&b=true&req=x")
	f.Add("i=notanint&u=-1&f=abc&sl=a&sl=b&sl=c&sl=d")
	f.Add("ptr=5&isl=1&isl=2&def=&n=deep")
	f.Add("")
	f.Fuzz(func(t *testing.T, raw string) {
		values, err := url.ParseQuery(raw)
		if err != nil {
			return // ParseForm would reject the same input upstream
		}
		var dst fuzzTarget
		_ = forms.Bind(&dst, values) // must not panic
	})
}
