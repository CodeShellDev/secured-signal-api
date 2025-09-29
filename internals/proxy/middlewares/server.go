package middlewares

import (
	"net/http"

	"github.com/codeshelldev/secured-signal-api/utils/config"
)

type ServeMiddleware struct {
	Next http.Handler
}

func (data ServeMiddleware) Use() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, config.ENV.FAVICON_PATH)
	})

	mux.Handle("/", data.Next)

	return mux
}
