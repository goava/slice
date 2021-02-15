package slice

import (
	"context"
	"fmt"
	"strings"

	"github.com/goava/di"
)

// createContainer is a step of application bootstrap. It collects user dependency injection
// options and creates container with them. Invalid dependency injection option will cause error.
func createContainer(diopts ...di.Option) (*di.Container, error) {
	// create container and validate user dependency injection options
	container, err := di.New(diopts...)
	if err != nil {
		return nil, fmt.Errorf("create container failed: %w", err)
	}
	return container, nil
}

// buildBundles is a step of application bootstrap. It iterates over all registered bundles and
// builds their dependencies. Build errors will be combined into one by containerBuilder.
func buildBundles(container *di.Container, bundles ...Bundle) error {
	for _, bundle := range bundles {
		if err := container.Apply(bundle.Components...); err != nil {
			return fmt.Errorf("build %s bundle failed: %w", bundle.Name, err)
		}
	}
	return nil
}

// before is a step of application bootstrap. It iterates over all registered bundles and call their Boot()
// method. If bundle boot are success shutdown function will be returned in shutdowns. In case, that boot
// failed process of booting application will be stopped.
func before(ctx context.Context, container *di.Container, bundles ...Bundle) (after []hook, _ error) {
	for _, bundle := range bundles {
		if err := ctx.Err(); err != nil {
			return after, fmt.Errorf("boot %s bundle failed: %w", bundle.Name, err)
		}
		for _, h := range bundle.Hooks {
			if err := container.Invoke(h.Before); err != nil {
				return nil, fmt.Errorf("boot %s bundle failed: %w", bundle.Name, err)
			}
			if h.After != nil {
				after = append(after, hook{
					name: bundle.Name,
					hook: h.After,
				})
			}
		}
	}
	return after, nil
}

// run is a part of application lifecycle. It resolves application dispatcher via container and call Run() method.
func run(ctx context.Context, container *di.Container) error {
	// resolve dispatcher
	var dispatcher Dispatcher
	if err := container.Resolve(&dispatcher); err != nil {
		return fmt.Errorf("resolve Dispatcher failed: %w", err)
	}
	// dispatcher run
	if err := dispatcher.Run(ctx); err != nil {
		return fmt.Errorf("failure: %w", err)
	}
	return nil
}

// after invoke hooks in reverse order.
func after(ctx context.Context, container *di.Container, hooks []hook) error {
	done := make(chan struct{})
	var errs errShutdown
	go func() {
		// shutdown bundles in reverse order
		for i := len(hooks) - 1; i >= 0; i-- {
			// bundle shutdown
			h := hooks[i]
			if err := container.Invoke(h.hook); err != nil {
				errs = append(errs, fmt.Errorf("shutdown %s failed: %w", h.name, err))
			}
		}
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return fmt.Errorf("shutdown failed: %w", ctx.Err())
	case <-done:
		if len(errs) != 0 {
			return fmt.Errorf("shutdown failed: %w", errs)
		}
		return nil
	}
}

type hook struct {
	name string
	hook di.Invocation
}

type errShutdown []error

func (e errShutdown) Error() string {
	var s []string
	for _, err := range e {
		s = append(s, err.Error())
	}
	return strings.Join(s, "; ")
}
