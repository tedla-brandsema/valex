// Validate-struct registers catalog directives and validates a struct with the
// "val" tag.
package main

import (
	"fmt"

	"github.com/tedla-brandsema/valex"
	"github.com/tedla-brandsema/valex/validators"
)

// Register the directives the struct uses. The engine ships none of its own.
func init() {
	valex.MustRegisterDirective(&validators.MinLengthValidator{})
	valex.MustRegisterDirective(&validators.EmailValidator{})
	valex.MustRegisterDirective(&validators.IntRangeValidator{})
}

type User struct {
	Name  string `val:"min,size=3"`
	Email string `val:"email"`
	Age   int    `val:"rangeint,min=0,max=120"`
}

func main() {
	users := []User{
		{Name: "Gopher", Email: "gopher@example.com", Age: 30},
		{Name: "Al", Email: "gopher@example.com", Age: 30},
		{Name: "Gopher", Email: "gopher@example.com", Age: 200},
	}

	for _, u := range users {
		if err := valex.ValidateStruct(&u); err != nil {
			fmt.Printf("%q invalid: %v\n", u.Name, err)
			continue
		}
		fmt.Printf("%q valid\n", u.Name)
	}

	// Output:
	// "Gopher" valid
	// "Al" invalid: tag "val" error: directive processing field "Name" directive "min": value Al is shorter than minimum length 3
	// "Gopher" invalid: tag "val" error: directive processing field "Age" directive "rangeint": value 200 is out of range [0, 120]
}
