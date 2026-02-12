package endpoints

import (
	"net/http"
)

type Endpoint struct {
	Name string
	Handler func(mux *http.ServeMux) *http.ServeMux
}

func (endpoint Endpoint) Use(mux *http.ServeMux) *http.ServeMux {
	return endpoint.Handler(mux)
}