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
func SetName(name string) Option {
	return option(func(s *Application) {
		if s.info == nil {
			s.info = &Info{}
		}
		s.info.Name = name
	})
}

// RegisterBundles registers application bundles.
// todo: naming can be changed
func RegisterBundles(bundles ...Bundle) Option {
	return option(func(s *Application) {
		s.bundles = append(s.bundles, bundles...)
	})
}

// ConfigureContainer configures the dependency injection container. It saves container options for the compile stage.
// On compile stage they will be applied on container.
// todo: naming can be changed
func ConfigureContainer(options ...di.Option) Option {
	return option(func(s *Application) {
		s.di = append(s.di, options...)
	})
}

// BootTimeout sets application boot timeout.
func BootTimeout(timeout time.Duration) Option {
	return option(func(s *Application) {
		s.timeouts.boot = timeout
	})
}

// ShutdownTimeout sets application shutdown timeout.
func ShutdownTimeout(timeout time.Duration) Option {
	return option(func(s *Application) {
		s.timeouts.shutdown = timeout
	})
}

type option func(s *Application)

func (o option) apply(s *Application) { o(s) }
