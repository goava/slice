package bundle

import (
	"github.com/goava/di"

	"github.com/goava/slice"
)

// New
func New(options ...Option) slice.Bundle {
	b := &slice.Bundle{}
	for _, opt := range options {
		opt.apply(b)
	}
	return *b
}

type Option interface {
	apply(bundle *slice.Bundle)
}

func WithName(name string) Option {
	return option(func(bundle *slice.Bundle) {
		bundle.Name = name
	})
}

func WithParameters(parameters ...slice.Parameter) Option {
	return option(func(bundle *slice.Bundle) {
		bundle.Parameters = append(bundle.Parameters, parameters...)
	})
}

func WithComponents(components ...di.Option) Option {
	return option(func(bundle *slice.Bundle) {
		bundle.Components = append(bundle.Components, components...)
	})
}

// WithHooks
func WithHooks(hooks ...slice.Hook) Option {
	return option(func(bundle *slice.Bundle) {
		bundle.Hooks = append(bundle.Hooks, hooks...)
	})
}

// WithBundles
func WithBundles(bundles ...slice.Bundle) Option {
	return option(func(bundle *slice.Bundle) {
		bundle.Bundles = append(bundle.Bundles, bundles...)
	})
}

type option func(bundle *slice.Bundle)

func (o option) apply(bundle *slice.Bundle) {
	o(bundle)
}
