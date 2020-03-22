package slice

import "github.com/goava/di"

// Option configure slice.
type Option interface {
	apply(s *lifecycle)
}

// Bundles registers reusable parts of application into registry. On compile stage each bundle will be processed.
func Bundles(bundles ...Bundle) Option {
	return option(func(s *lifecycle) {
		s.bundles = append(s.bundles, bundles...)
	})
}

// DependencyInjection configures the dependency injection container. It saves container options for the compile stage.
// On compile stage they will be applied on container.
func DependencyInjection(options ...di.Option) Option {
	return option(func(s *lifecycle) {
		s.di = append(s.di, options...)
	})
}

type option func(s *lifecycle)

func (o option) apply(s *lifecycle) { o(s) }
