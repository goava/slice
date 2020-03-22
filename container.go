package slice

import "github.com/goava/di"

// ContainerBuilder builds the container. It is used to providing components in bundles.
type ContainerBuilder interface {
	// Has checks that type exists in container, if not it is return false.
	Has(target interface{}, options ...di.ResolveOption) bool
	// Provide provides a reliable way of component building to the container.
	// The constructor will be invoked lazily on-demand. For more information about
	// constructors see Constructor interface. ProvideOption can add additional
	// behavior to the process of type resolving.
	Provide(constructor di.Constructor, options ...di.ProvideOption)
}

// Container is a compiled container.
type Container interface {
	// Has checks that type exists in container, if not it return false.
	Has(target interface{}, options ...di.ResolveOption) bool
	// Invoke calls provided function.
	Invoke(fn di.Invocation, options ...di.InvokeOption) error
	// Resolve builds instance of target type and fills target pointer.
	Resolve(into interface{}, options ...di.ResolveOption) error
}

type bundleContainerBuilder struct {
	container *di.Container
	bundleErr bundleDIErrors
}

// Has implements ContainerBuilder.
func (b *bundleContainerBuilder) Has(target interface{}, options ...di.ResolveOption) bool {
	return b.container.Has(target, options...)
}

// Provide implements ContainerBuilder.
func (b *bundleContainerBuilder) Provide(constructor di.Constructor, options ...di.ProvideOption) {
	if err := b.container.Provide(constructor, options...); err != nil {
		b.bundleErr.list = append(b.bundleErr.list, bundleDIError{err})
	}
}
