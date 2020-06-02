package slice

import "context"

//go:generate moq -out dispatcher_test.go . Dispatcher

// Dispatcher controls application lifecycle.
type Dispatcher interface {
	// Run runs your application. Context will be canceled if application get
	// syscall.SIGINT or syscall.SIGTERM. Example implementation handles application
	// shutdown and worker errors:
	//
	// 	func(d *ExampleDispatcher) Run(ctx context.Context) (err error) {
	//		errChan := make(chan error)
	//		go func() {
	//			errChan <- d.Worker.Start()
	//		}
	// 		select {
	//		// application shutdown
	// 		case <-ctx.Done():
	//			return ctx.Err()
	//		case err = <-errChan:
	// 		}
	//		return err
	// 	}
	Run(ctx context.Context) (err error)
}
