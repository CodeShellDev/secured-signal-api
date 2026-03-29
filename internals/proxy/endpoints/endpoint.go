package endpoints

import (
	"net/http"
)

type Endpoint struct {
	Name string
	Handler func(mux *http.ServeMux, next http.Handler) *http.ServeMux
}

func (endpoint Endpoint) Use(mux *http.ServeMux, next http.Handler) *http.ServeMux {
	return endpoint.Handler(mux, next)
}