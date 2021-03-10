package httpctrl

import (
	"net/http"
)

// NopController
type NopController struct {
}

// NewNopController constructs nop controller.
func NewNopController() *NopController {
	return &NopController{}
}

func (n NopController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", n.Nop)
}

// Nop writes "nop" to response writer.
func (n NopController) Nop(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("nop"))
}
