package middlewares

import (
	"net/http"

	"github.com/codeshelldev/secured-signal-api/utils/config"
)

var Server Middleware = Middleware{
	Name: "Server",
	Use: serverHandler,
}

func serverHandler(next http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, config.ENV.FAVICON_PATH)
	})

	mux.Handle("/", next)

	return mux
}
