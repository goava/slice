package slice

import "context"

//go:generate moq -out kernel_test.go . Kernel

// Kernel runs application.
type Kernel interface {
	Run(ctx context.Context) error
}
