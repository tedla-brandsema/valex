// Chained applies several directives to one field by separating them with ';'.
// They run left to right and each MutMode result feeds the next, so this trims a
// username, lowercases it, then enforces a maximum length — in that order. The
// length check comes from the valex/validators catalog; trim and lower are
// custom MutMode directives.
package main

import (
	"fmt"
	"strings"

	"github.com/tedla-brandsema/tagex"
	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/valex/validators"
)

// trimDirective removes surrounding whitespace (MutMode writes the result back).
type trimDirective struct{}

func (*trimDirective) Name() string              { return "trim" }
func (*trimDirective) Mode() tagex.DirectiveMode { return tagex.MutMode }
func (*trimDirective) Handle(s string) (string, error) {
	return strings.TrimSpace(s), nil
}

// lowerDirective lowercases the field (MutMode).
type lowerDirective struct{}

func (*lowerDirective) Name() string              { return "lower" }
func (*lowerDirective) Mode() tagex.DirectiveMode { return tagex.MutMode }
func (*lowerDirective) Handle(s string) (string, error) {
	return strings.ToLower(s), nil
}

func init() {
	valex.MustRegisterDirective(&trimDirective{})
	valex.MustRegisterDirective(&lowerDirective{})
	valex.MustRegisterDirective(&validators.MaxLengthValidator{})
}

type Account struct {
	Username string `val:"trim;lower;max,size=12"`
}

func main() {
	accounts := []Account{
		{Username: "  Ada  "},
		{Username: "  TheLongestUsername  "},
	}

	for i := range accounts {
		if err := valex.ValidateStruct(&accounts[i]); err != nil {
			fmt.Printf("rejected: %v\n", err)
			continue
		}
		fmt.Printf("normalized to %q\n", accounts[i].Username)
	}

	// Output:
	// normalized to "ada"
	// rejected: tag "val" error: directive processing field "Username" directive "max": value thelongestusername exceeds maximum length 12
}
