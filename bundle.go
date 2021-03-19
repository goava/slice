package slice

import (
	"fmt"

	"github.com/goava/di"
)

// Bundle
type Bundle struct {
	Name       string
	Parameters []Parameter
	Components []di.Option
	Hooks      []Hook
	Bundles    []Bundle
}

type startErrors []error

// Error
func (e startErrors) Error() (r string) {
	for _, err := range e {
		r = fmt.Sprintf("%s- %s\n", r, err)
	}
	return r
}
