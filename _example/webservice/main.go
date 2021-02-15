package main

import (
	"context"
	"net/http"

	"github.com/goava/di"

	"github.com/goava/slice"
	"github.com/goava/slice/_example/webservice/httpsrv"
)

func main() {
	slice.Run(
		slice.WithName("bundle-app"),
		slice.WithBundles(
			httpsrv.Bundle,
		),
		slice.WithComponents(
			di.Provide(NewDispatcher, di.As(new(slice.Dispatcher))),
		),
	)
}

type Dispatcher struct {
	server *http.Server
}

func NewDispatcher(server *http.Server) *Dispatcher {
	return &Dispatcher{server: server}
}

func (d Dispatcher) Run(ctx context.Context) (err error) {
	errch := make(chan error)
	go func() {
		errch <- d.server.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		if err := d.server.Close(); err != nil {
			return err
		}
		return ctx.Err()
	case <-errch:
		return err
	}
}
