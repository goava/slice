package slice

import "github.com/goava/di"

// Hook
type Hook struct {
	// Deprecated: use BeforeStart
	Before di.Invocation
	// Deprecated: use BeforeShutdown
	After          di.Invocation
	BeforeStart    di.Invocation
	BeforeShutdown di.Invocation
}
