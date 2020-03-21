package slice

import (
	"context"
)

//go:generate moq -out bundle_test.go . BootShutdown

// Bundle registers reusable set of components.
type Bundle interface {
	// DependencyInjection provides bundle components to container builder.
	DependencyInjection(builder ContainerBuilder)
}

// BootShutdown is a bundle that have boot and shutdown stages.
type BootShutdown interface {
	Bundle
	// Boot provides way to interact with dependency injection container on the
	// boot stage. On boot stage main dependencies already provided to container.
	// And on this stage bundle can interact with them.
	// Boot can return error if process failed. It will be handled by Slice.
	Boot(ctx context.Context, container Container) error
	// Shutdown provides way to interact with dependency injection container
	// on shutdown stage. It can compensate things that was be made on boot stage.
	// Shutdown can return error if process failed. It will be handled by Slice.
	Shutdown(ctx context.Context, container Container) error
}
