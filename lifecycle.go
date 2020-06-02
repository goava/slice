package slice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/goava/di"
)

// configureBundles is a step of application bootstrap. It iterates over all registered bundles and
// loads configuration via bundleConfigurator
func configureBundles(configure bundleConfigurator, bundles ...bundle) error {
	for _, bundle := range bundles {
		if err := configure(bundle); err != nil {
			return fmt.Errorf("configure %s bundle failed: %s", bundle.name, err)
		}
	}
	return nil
}

// createContainer is a step of application bootstrap. It collects user dependency injection
// options and creates container with them. Invalid dependency injection option will cause error.
func createContainer(diopts ...di.Option) (*di.Container, error) {
	// create container and validate user dependency injection options
	container, err := di.New(diopts...)
	if err != nil {
		return nil, fmt.Errorf("create container failed: %s", err)
	}
	return container, nil
}

// buildBundles is a step of application bootstrap. It iterates over all registered bundles and
// builds their dependencies. Build errors will be combined into one by containerBuilder.
func buildBundles(container *di.Container, bundles ...bundle) error {
	for _, bundle := range bundles {
		builder := newContainerBuilder(container)
		bundle.Build(builder)
		if err := builder.Error(); err != nil {
			return fmt.Errorf("build %s bundle failed: %s", bundle.name, err)
		}
	}
	return nil
}

// boot is a step of application bootstrap. It iterates over all registered bundles and call their Boot()
// method. If bundle boot are success shutdown function will be returned in shutdowns. In case, that boot
// failed process of booting application will be stopped.
func boot(ctx context.Context, container *di.Container, bundles ...bundle) (shutdowns shutdowns, _ error) {
	for _, bundle := range bundles {
		if err := ctx.Err(); err != nil {
			return shutdowns, fmt.Errorf("boot %s bundle failed: %s", bundle.name, err)
		}
		if boot, ok := bundle.Bundle.(BootShutdown); ok {
			// boot bundle
			if err := boot.Boot(ctx, container); err != nil {
				return shutdowns, fmt.Errorf("boot %s bundle failed: %s", bundle.name, err)
			}
			// append successfully booted bundle shutdown
			shutdowns = append(shutdowns, bundleShutdown{
				name:     bundle.name,
				shutdown: boot.Shutdown,
			})
		}
	}
	return shutdowns, nil
}

// run is a part of application lifecycle. It resolves application dispatcher via container and call Run() method.
func run(ctx context.Context, container *di.Container) error {
	// resolve dispatcher
	var dispatcher Dispatcher
	if err := container.Resolve(&dispatcher); err != nil {
		return fmt.Errorf("resolve dispatcher failed: %s", err)
	}
	// dispatcher run
	if err := dispatcher.Run(ctx); err != nil {
		return fmt.Errorf("failure: %s", err)
	}
	return nil
}

// reverseShutdown shutdowns in reverse order.
func reverseShutdown(timeout time.Duration, container *di.Container, shutdowns shutdowns) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// shutdown bundles in reverse order
	var errs errShutdown
	for i := len(shutdowns) - 1; i >= 0; i-- {
		// bundle shutdown
		bs := shutdowns[i]
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("shutdown failed: %s", err)
		}
		if err := bs.shutdown(ctx, container); err != nil {
			errs = append(errs, fmt.Errorf("shutdown %s failed: %s", bs.name, err))
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("shutdown failed: %s", errs)
	}
	return nil
}

type bundleShutdown struct {
	name     string
	shutdown func(ctx context.Context, container Container) error
}

type shutdowns []bundleShutdown

type errShutdown []error

func (e errShutdown) Error() string {
	var s []string
	for _, err := range e {
		s = append(s, err.Error())
	}
	return strings.Join(s, "; ")
}
