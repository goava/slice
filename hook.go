package slice

import "github.com/goava/di"

// Hook
type Hook struct {
	Before di.Invocation
	After  di.Invocation
}
