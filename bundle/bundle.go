package bundle

import (
	"github.com/goava/slice"
)

// New creates bundle with provided options.
func New(options ...Option) slice.Bundle {
	b := &slice.Bundle{}
	for _, opt := range options {
		opt.apply(b)
	}
	return *b
}

// Option modifies bundle structure.
type Option interface {
	apply(bundle *slice.Bundle)
}

// WithName sets bundle name.
func WithName(name string) Option {
	return option(func(bundle *slice.Bundle) {
		bundle.Name = name
	})
}

// WithParameters add parameters to bundle.
func WithParameters(parameters ...slice.Parameter) Option {
	return option(func(bundle *slice.Bundle) {
		bundle.Parameters = append(bundle.Parameters, parameters...)
	})
}

// WithComponents configures bundle components.
func WithComponents(options ...slice.ComponentOption) Option {
	return option(func(bundle *slice.Bundle) {
		bundle.Components = append(bundle.Components, options...)
	})
}

// WithHooks adds application hooks.
func WithHooks(hooks ...slice.Hook) Option {
	return option(func(bundle *slice.Bundle) {
		bundle.Hooks = append(bundle.Hooks, hooks...)
	})
}

// WithBundles add dependency bundle.
func WithBundles(bundles ...slice.Bundle) Option {
	return option(func(bundle *slice.Bundle) {
		bundle.Bundles = append(bundle.Bundles, bundles...)
	})
}

type option func(bundle *slice.Bundle)

func (o option) apply(bundle *slice.Bundle) {
	o(bundle)
}
