package slice

import (
	"time"

	"github.com/goava/di"
)

// Option configure slice.
type Option interface {
	apply(s *Application)
}

// SetName sets application name.
// In case you need to change the name, you can use the APP_NAME environment variable.
// Deprecated: use WithName().
func SetName(name string) Option {
	return option(func(s *Application) {
		s.Name = name
	})
}

// RegisterBundles registers application bundles.
// deprecated: use WithBundles()
func RegisterBundles(bundles ...Bundle) Option {
	return option(func(s *Application) {
		s.Bundles = append(s.Bundles, bundles...)
	})
}

// WithName sets application name.
// In case you need to change the name, you can use the APP_NAME environment variable.
func WithName(name string) Option {
	return option(func(s *Application) {
		s.Name = name
	})
}

// RegisterBundles registers application bundles.
func WithBundles(bundles ...Bundle) Option {
	return option(func(s *Application) {
		s.Bundles = append(s.Bundles, bundles...)
	})
}

// ConfigureContainer configures the dependency injection container. It saves container options for the compile stage.
// On compile stage they will be applied on container.
// Deprecated: use WithComponents()
func ConfigureContainer(options ...di.Option) Option {
	return option(func(s *Application) {
		s.Components = append(s.Components, options...)
	})
}

// ConfigureContainer configures the dependency injection container. It saves container options for the compile stage.
// On compile stage they will be applied on container.
func WithComponents(components ...di.Option) Option {
	return option(func(s *Application) {
		s.Components = append(s.Components, components...)
	})
}

// BootTimeout sets application boot timeout.
// Deprecated: use StartTimeout()
func BootTimeout(timeout time.Duration) Option {
	return option(func(s *Application) {
		s.StartTimeout = timeout
	})
}

// ShutdownTimeout sets application shutdown timeout.
// Deprecated: use StopTimeout()
func ShutdownTimeout(timeout time.Duration) Option {
	return option(func(s *Application) {
		s.StopTimeout = timeout
	})
}

// BootTimeout sets application boot timeout.
func StartTimeout(timeout time.Duration) Option {
	return option(func(s *Application) {
		s.StartTimeout = timeout
	})
}

// ShutdownTimeout sets application shutdown timeout.
func StopTimeout(timeout time.Duration) Option {
	return option(func(s *Application) {
		s.StopTimeout = timeout
	})
}

// WithParameters adds parameters to container. On boot stage all parameter structures
// will be processed via ParameterParser.
func WithParameters(parameters ...Parameter) Option {
	return option(func(s *Application) {
		s.Parameters = append(s.Parameters, parameters...)
	})
}

type option func(s *Application)

func (o option) apply(s *Application) { o(s) }
