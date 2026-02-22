package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler interface {
	GetMetrics(w http.ResponseWriter, r *http.Request)
}

type handler struct{}

var _ Handler = (*handler)(nil)

func New() *handler {
	return &handler{}
}

func (h *handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
