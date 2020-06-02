package slice

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goava/di"
)

// ContainerBuilder provides container instance that squash bundle errors.
type ContainerBuilder interface {
	// Has checks that type exists in container, if not it is return false.
	Has(target interface{}, options ...di.ResolveOption) bool
	// Provide provides a reliable way of component building to the container.
	// The constructor will be invoked lazily on-demand. For more information about
	// constructors see Constructor interface. ProvideOption can add additional
	// behavior to the process of type resolving.
	Provide(constructor di.Constructor, options ...di.ProvideOption)
	// Invoke calls invocation. Note, that instances will be built on build stage.
	// It should be used only for conditional provides:
	//
	//	Invoke(func(env slice.Env, container *di.Container) error {
	//		if env.IsProduction() {
	//			return container.Provide(NewProductionType)
	//		}
	//		return nil
	//	})
	// todo: refactor this
	Invoke(invocation di.Invocation, options ...di.InvokeOption)
}

// Container is a dependency injection container.
type Container interface {
	// Has checks that type exists in container, if not it return false.
	Has(target interface{}, options ...di.ResolveOption) bool
	// Invoke calls provided function.
	Invoke(fn di.Invocation, options ...di.InvokeOption) error
	// Resolve builds instance of target type and fills target pointer.
	Resolve(into interface{}, options ...di.ResolveOption) error
}

// newContainerBuilder creates container builder for bundle.
func newContainerBuilder(container *di.Container) *containerBuilder {
	return &containerBuilder{
		container: container,
	}
}

// containerBuilder
type containerBuilder struct {
	container *di.Container
	errs      []error
}

// Has implements ContainerBuilder.
func (b *containerBuilder) Has(target interface{}, options ...di.ResolveOption) bool {
	return b.container.Has(target, options...)
}

// Provide implements ContainerBuilder.
func (b *containerBuilder) Provide(constructor di.Constructor, options ...di.ProvideOption) {
	if err := b.container.Provide(constructor, options...); err != nil {
		b.errs = append(b.errs, err) // append bundle provide error
	}
}

// Invoke invokes function and collect error.
func (b *containerBuilder) Invoke(invocation di.Invocation, options ...di.InvokeOption) {
	if err := b.container.Invoke(invocation, options...); err != nil {
		b.errs = append(b.errs, err)
	}
}

// Error returns collected build errors as one. If bundle build success returns nil.
// Error string representation will be multi line for easy reading. todo: use %+v
func (b *containerBuilder) Error() error {
	if len(b.errs) == 0 {
		return nil
	}
	sb := strings.Builder{}
	for _, err := range b.errs {
		sb.WriteString(fmt.Sprintf("\n\t- %s", err))
	}
	return errors.New(sb.String())
}
