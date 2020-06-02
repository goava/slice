package slice

import (
	"context"
	"reflect"

	"github.com/kelseyhightower/envconfig"
)

//go:generate moq -out bundle_mock_test.go . Bundle BootShutdown

// Bundle registers reusable set of components.
type Bundle interface {
	// Build provides bundle components to container builder.
	Build(builder ContainerBuilder)
}

// BootShutdown is a bundle that have boot and shutdown stages.
// todo: naming?
type BootShutdown interface {
	Bundle
	// Boot provides way to interact with dependency injection container on the
	// boot stage. On boot stage main dependencies already provided to container.
	// And on this stage bundle can interact with them.
	// Boot can return error if process failed. It will be handled by Slice.
	Boot(ctx context.Context, container Container) error
	// Shutdown provides way to interact with dependency injection container
	// on shutdown stage. It can compensate things that was be made on boot stage.
	// Shutdown can return error if process failed. It will be handled by application.
	Shutdown(ctx context.Context, container Container) error
}

// A DependOn describe that bundle depends on another bundle.
type DependOn interface {
	Bundle
	// DependOn returns dependent bundle.
	DependOn() []Bundle
}

// bundleConfigurator is a function that loads bundle configuration by some way
type bundleConfigurator func(bundle Bundle) error

// todo: replace envconfig
func defaultBundleConfigurator() bundleConfigurator {
	return func(bundle Bundle) error {
		return envconfig.Process("", bundle)
	}
}

// inspectBundles prepares bundles.
func inspectBundles(bundles ...Bundle) (result []bundle) {
	for _, b := range bundles {
		result = append(result, bundle{
			name:   bundleName(b),
			Bundle: b,
		})
	}
	return result
}

// bundle
type bundle struct {
	Bundle
	name string
}

// bundleName gets bundle string representation.
func bundleName(bundle Bundle) string {
	return reflect.TypeOf(bundle).String()
}
