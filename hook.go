package slice

import "github.com/goava/di"

// Hook
type Hook struct {
	// BeforeStart invokes function before application start.
	BeforeStart di.Invocation
	// BeforeShutdown invokes function before application shutdown.
	BeforeShutdown di.Invocation
	// Deprecated: use BeforeStart
	Before di.Invocation
	// Deprecated: use BeforeShutdown
	After di.Invocation
}
