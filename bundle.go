package slice

import (
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
