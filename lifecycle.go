package slice

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/goava/di"
	"github.com/oklog/run"
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

// before is a step of application bootstrap. It iterates over all registered bundles and call their Boot()
// method. If bundle boot are success shutdown function will be returned in shutdowns. In case, that boot
// failed process of booting application will be stopped.
func beforeStart(ctx context.Context, container *di.Container, bundles ...Bundle) (after []hook, _ error) {
	var errs startErrors
	for _, bundle := range bundles {
		if err := ctx.Err(); err != nil {
			return after, fmt.Errorf("boot %s bundle failed: %w", bundle.Name, err)
		}
		for _, h := range bundle.Hooks {
			if h.BeforeStart != nil {
				if err := container.Invoke(h.BeforeStart); err != nil {
					errs = append(errs, fmt.Errorf("boot %s bundle failed: %w", bundle.Name, err))
				}
				if h.BeforeShutdown != nil {
					after = append(after, hook{
						name: bundle.Name,
						hook: h.BeforeShutdown,
					})
				}
			}
		}
	}
	if len(errs) != 0 {
		return nil, errs
	}
	return after, nil
}

// dispatch is a part of application lifecycle. It resolves application dispatcher via container and call Run() method.
func dispatch(ctx context.Context, logger Logger, stop func(), dispatchers []Dispatcher) error {
	var once sync.Once
	// start all dispatchers
	var workers run.Group
	for _, d := range dispatchers {
		dispatcher := d
		dt := reflect.TypeOf(dispatcher)
		execute := func() error {
			logger.Printf("Start %s", dt)
			if err := dispatcher.Run(ctx); err != nil {
				return fmt.Errorf("%s: %w", dt, err)
			}
			once.Do(func() {
				logger.Printf("Terminate signal from %s", dt)
				stop()
			})
			logger.Printf("Stopped %s", dt)
			return nil
		}
		interrupt := func(err error) {
			once.Do(func() {
				stop()
			})
		}
		workers.Add(execute, interrupt)
	}
	if err := workers.Run(); err != nil {
		return fmt.Errorf("failure: %w", err)
	}
	return nil
}

// beforeShutdown invoke hooks in reverse order.
func beforeShutdown(ctx context.Context, container *di.Container, hooks []hook) error {
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
