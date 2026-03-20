package middlewares

import (
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
	"github.com/codeshelldev/secured-signal-api/utils/urlutils"
)

var CORS Middleware = Middleware{
	Name: "CORS",
	Use: corsHandler,
}

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conf := GetConfigByReq(req)

		cors := conf.SETTINGS.ACCESS.CORS.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.CORS)

		defaultMethods := cors.Methods.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.CORS.Value.Methods)
		defaultHeaders := cors.Headers.OptOrEmpty(config.DEFAULT.SETTINGS.ACCESS.CORS.Value.Headers)

		if len(cors.Origins) == 0 {
			next.ServeHTTP(w, req)
			return
		}

		origin := req.Header.Get("Origin")

		if origin == "" {
			next.ServeHTTP(w, req)
			return
		}

		originURL, err := url.Parse(origin)

		var matchingOrigin *structure.Origin

		if err == nil {
			for _, o := range cors.Origins {
				if urlutils.NormalizeURL(originURL) == urlutils.NormalizeURL((*url.URL)(&o.URL)) {
					matchingOrigin = &o
				}
			}
		}

		if matchingOrigin == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)

		// CORS preflight request
		if req.Method == "OPTIONS" {
			requestedMethod := req.Header.Get("Access-Control-Request-Method")

			if requestedMethod != "" {
				allowedMethods := matchingOrigin.Methods.ValueOrFallback(defaultMethods)

				if len(allowedMethods) != 0 {
					// only set if any (matching) methods
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
				}
			}

			requestedHeaders := req.Header.Get("Access-Control-Request-Headers")

			if requestedHeaders != "" {
				allowedHeaders := matchingOrigin.Headers.ValueOrFallback(defaultHeaders)

				matchingHeaders := []string{}

				// echo back allowed and requested headers
				for header := range strings.SplitSeq(requestedHeaders, ",") {
					header = strings.TrimSpace(header)

					var match string

					if slices.ContainsFunc(allowedHeaders, func(allowed string) bool {
						if strings.EqualFold(header, allowed) {
							match = allowed
							return true
						}

						return false
					}) {
						matchingHeaders = append(matchingHeaders, match)
					}
				}

				if len(matchingHeaders) != 0 {
					// only set if any (matching) headers
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(matchingHeaders, ","))
				}
			}

			w.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(w, req)
	})
}