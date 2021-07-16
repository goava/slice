package slice

import (
	"time"
)

// Option configure slice.
type Option interface {
	apply(s *Application)
}

// WithName sets application name.
// In case you need to change the name, you can use the APP_NAME environment variable.
func WithName(name string) Option {
	return option(func(s *Application) {
		s.Name = name
	})
}

// WithParameters adds parameters to container. On boot stage all parameter structures
// will be processed via ParameterParser.
func WithParameters(parameters ...Parameter) Option {
	return option(func(s *Application) {
		s.Parameters = append(s.Parameters, parameters...)
	})
}

// WithBundles registers application bundles.
func WithBundles(bundles ...Bundle) Option {
	return option(func(s *Application) {
		s.Bundles = append(s.Bundles, bundles...)
	})
}

// WithComponents contains component options.
func WithComponents(components ...ComponentOption) Option {
	return option(func(s *Application) {
		for _, c := range components {
			c.apply(s)
		}
	})
}

// WithParameterParser sets parser for application.
func WithParameterParser(parser ParameterParser) Option {
	return option(func(s *Application) {
		s.ParameterParser = parser
	})
}

// WithLogger sets application logger.
func WithLogger(logger Logger) Option {
	return option(func(s *Application) {
		s.Logger = logger
	})
}

// StartTimeout sets application boot timeout.
func StartTimeout(timeout time.Duration) Option {
	return option(func(s *Application) {
		s.StartTimeout = timeout
	})
}

// ShutdownTimeout sets application shutdown timeout.
func ShutdownTimeout(timeout time.Duration) Option {
	return option(func(s *Application) {
		s.ShutdownTimeout = timeout
	})
}

type option func(s *Application)

func (o option) apply(s *Application) { o(s) }
