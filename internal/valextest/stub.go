// Package valextest provides shared test fixtures for the valex module.
//
// It lives under internal/ and is imported only from _test.go files, so it is
// never part of a production build, and internal/ keeps it out of reach of other
// modules. Importing it (a blank import is enough) registers a couple of stub
// directives against valex's "val" tag, letting tests exercise the validation
// pipeline without depending on the valex/validators catalog.
package valextest

import (
	"fmt"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
)

// MinLen is a stub directive ("minlen") enforcing a string minimum length.
type MinLen struct {
	Size int `param:"size"`
}

func (d *MinLen) Name() string              { return "minlen" }
func (d *MinLen) Mode() tagex.DirectiveMode { return tagex.EvalMode }
func (d *MinLen) Handle(val string) (string, error) {
	if len(val) < d.Size {
		return val, fmt.Errorf("value %q is shorter than minimum length %d", val, d.Size)
	}
	return val, nil
}

// IntRange is a stub directive ("intrange") enforcing an inclusive int range.
type IntRange struct {
	Min int `param:"min"`
	Max int `param:"max"`
}

func (d *IntRange) Name() string              { return "intrange" }
func (d *IntRange) Mode() tagex.DirectiveMode { return tagex.EvalMode }
func (d *IntRange) Handle(val int) (int, error) {
	if val < d.Min || val > d.Max {
		return val, fmt.Errorf("value %d is out of range [%d, %d]", val, d.Min, d.Max)
	}
	return val, nil
}

func init() {
	valex.RegisterDirective[string](&MinLen{})
	valex.RegisterDirective[int](&IntRange{})
}
