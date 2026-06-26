// Custom-directive extends the "val" tag with a directive of your own. A
// directive is any tagex.Directive[T]: implement Name, Mode, and Handle.
package main

import (
	"fmt"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
)

// EvenDirective accepts only even ints. Mode is EvalMode, so Handle's return
// value is used only to report success or failure, never written back.
type EvenDirective struct{}

func (*EvenDirective) Name() string              { return "even" }
func (*EvenDirective) Mode() tagex.DirectiveMode { return tagex.EvalMode }
func (*EvenDirective) Handle(n int) (int, error) {
	if n%2 != 0 {
		return n, fmt.Errorf("value %d is not even", n)
	}
	return n, nil
}

func main() {
	valex.RegisterDirective(&EvenDirective{})

	type Ticket struct {
		Seats int `val:"even"`
	}

	fmt.Println(valex.ValidateStruct(&Ticket{Seats: 4}))
	fmt.Println(valex.ValidateStruct(&Ticket{Seats: 3}))

	// Output:
	// <nil>
	// tag "val" error: directive processing field "Seats" directive "even": value 3 is not even
}
