package httpsrv

import (
	"net/http"
	"time"

	"github.com/goava/di"

	"github.com/goava/slice"
	"github.com/goava/slice/bundle"
)

// Parameters contains application configuration.
type Parameters struct {
	Addr         string        `envconfig:"addr" required:"true" desc:"Server address"`
	ReadTimeout  time.Duration `envconfig:"read_timeout" required:"true" desc:"Server read timeout"`
	WriteTimeout time.Duration `envconfig:"write_timeout" required:"true" desc:"Server write timeout"`
}

// DefaultParameters returns default application parameters.
func DefaultParameters() *Parameters {
	return &Parameters{
		// default values
		Addr:         ":http",
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
}

// Bundle
var Bundle = bundle.New(
	bundle.WithName("http"),
	bundle.WithParameters(
		DefaultParameters(),
	),
	bundle.WithComponents(
		di.Provide(NewHTTPServer),
		di.Provide(http.NewServeMux, di.As(new(http.Handler))),
	),
	bundle.WithHooks(
		slice.Hook{
			BeforeStart: RegisterHTTPControllers,
		},
	),
)

// NewHTTPServer
func NewHTTPServer(logger slice.Logger, handler http.Handler, params *Parameters) *http.Server {
	logger.Printf("Server Addr %s", params.Addr)
	logger.Printf("Server ReadTimeout %s", params.ReadTimeout)
	logger.Printf("Server WriteTimeout %s", params.WriteTimeout)
	return &http.Server{
		Addr:         params.Addr,
		Handler:      handler,
		ReadTimeout:  params.ReadTimeout,
		WriteTimeout: params.WriteTimeout,
	}
}

// Controller
type Controller interface {
	RegisterRoutes(mux *http.ServeMux)
}

// RegisterHTTPControllers
func RegisterHTTPControllers(logger slice.Logger, container *di.Container, mux *http.ServeMux) error {
	var controllers []Controller
	has, err := container.Has(&controllers)
	if err != nil {
		return err
	}
	if !has {
		logger.Printf("Controllers not found")
		return nil
	}
	if err := container.Resolve(&controllers); err != nil {
		return err
	}
	for _, ctrl := range controllers {
		ctrl.RegisterRoutes(mux)
	}
	return err
}
