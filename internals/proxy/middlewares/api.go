package middlewares

import (
	"net/http"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	e "github.com/codeshelldev/secured-signal-api/internals/proxy/endpoints"
)

var InternalInsecureAPI Middleware = Middleware{
	Name: "_Internal_Insecure_API",
	Use: internalInsecureAPIHandler,
}

func internalInsecureAPIHandler(next http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, config.ENV.FAVICON_PATH)
	})

	mux.Handle("/", next)

	return mux
}

var InternalSecureAPI Middleware = Middleware{
	Name: "_Internal_Secure_API",
	Use: internalSecureAPIHandler,
}

func internalSecureAPIHandler(next http.Handler) http.Handler {
	mux := http.NewServeMux()

	e.AboutEndpoint.Use(mux)
	e.SendEnpoint.Use(mux)
	e.ScheduleEndpoint.Use(mux)

	mux.Handle("/", next)

	return mux
}