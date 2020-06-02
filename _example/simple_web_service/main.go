package main

import (
	"context"
	"net/http"

	"github.com/goava/di"

	"github.com/goava/slice"
)

// NewServer creates server.
func NewServer() *http.Server {
	return &http.Server{}
}

// NewMux creates new http mux. It require array of controllers that will be
// registered via mux.
func NewMux(controllers []Controller) *http.ServeMux {
	mux := http.NewServeMux()
	for _, ctrl := range controllers {
		ctrl.RegisterHandler(mux)
	}
	return mux
}

// Controller is a http controller. It can register own routes via mux.
type Controller interface {
	// RegisterHandler registers handler via mux.
	RegisterHandler(mux *http.ServeMux)
}

// Dispatcher dispatches Slice application lifecycle.
type Dispatcher struct {
	stop   chan<- error
	server *http.Server
}

// NewDispatcher creates new Slice application dispatcher.
func NewDispatcher(server *http.Server) *Dispatcher {
	return &Dispatcher{
		stop:   make(chan error),
		server: server,
	}
}

// Run runs Slice application.
func (d Dispatcher) Run(ctx context.Context) (err error) {
	errChan := make(chan error)
	go func() {
		errChan <- d.server.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errChan:
	}
	return err
}

func main() {
	slice.Run(
		slice.ConfigureContainer(
			di.Provide(NewDispatcher, di.As(new(slice.Dispatcher))),
			di.Provide(NewServer),
			di.Provide(NewMux),
		),
	)
}
