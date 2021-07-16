package slice

import "github.com/goava/di"

// ComponentOption modifies application components.
type ComponentOption interface {
	apply(s *Application)
}

// Provide returns container option that provides to container reliable way to build type. The constructor will
// be invoked lazily on-demand. For more information about constructors see di.Constructor interface. di.ProvideOption can
// add additional behavior to the process of type resolving.
func Provide(constructor di.Constructor, options ...di.ProvideOption) ComponentOption {
	return option(func(s *Application) {
		s.providers = append(s.providers, di.Provide(constructor, options...))
	})
}

// Supply provides value as is.
func Supply(value di.Value, options ...di.ProvideOption) ComponentOption {
	return option(func(s *Application) {
		s.providers = append(s.providers, di.ProvideValue(value, options...))
	})
}
