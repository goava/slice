package slice

import (
	"fmt"
)

// A Bundle  is a separate unit of application.
type Bundle struct {
	Name       string
	Parameters []Parameter
	Components []ComponentOption
	Hooks      []Hook
	Bundles    []Bundle
}

func (b Bundle) apply(app *Application) {
	for _, option := range b.Components {
		option.apply(app)
	}
}

type startErrors []error

// Error
func (e startErrors) Error() (r string) {
	for _, err := range e {
		r = fmt.Sprintf("%s- %s\n", r, err)
	}
	return r
}
