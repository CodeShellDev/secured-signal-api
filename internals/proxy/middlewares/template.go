package middlewares

import (
	"net/http"
	"net/url"

	request "github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
)

var Template Middleware = Middleware{
	Name: "Template",
	Use: templateHandler,
}

func templateHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := GetLogger(req)

		conf := GetConfigByReq(req)
		
		templating := conf.SETTINGS.MESSAGE.TEMPLATING.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.TEMPLATING)
		injecting := conf.SETTINGS.MESSAGE.INJECTING.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.INJECTING)

		variables := conf.SETTINGS.MESSAGE.VARIABLES.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.VARIABLES)

		urlToBody := injecting.URLToBody.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.INJECTING.Value.URLToBody)

		body, err := request.GetReqBody(req)

		if err != nil {
			logger.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
			return
		}

		body.EnsureNotNil()

		var modifiedBody bool

		if !body.Empty && templating.Body {
			var modified bool

			headers := request.GetReqHeaders(req)

			body.Data, modified, err = GetTemplatedBody(body.Data, headers, variables)

			if err != nil {
				logger.Error("Error Templating Body: ", err.Error())
			}

			if modified {
				modifiedBody = true
			}
		}

		if req.URL.RawQuery != "" {
			oldRawQuery := req.URL.RawQuery

			if templating.Query {
				req.URL.RawQuery, err = TemplateQuery(req.URL.RawQuery, variables)

				if err != nil {
					logger.Error("Error Templating Query: ", err.Error())
				}
			}

			if urlToBody.Query {
				modified := InjectQueryIntoBody(req.URL.Query(), body.Data)

				if modified {
					modifiedBody = true
				}
			}

			if req.URL.RawQuery != oldRawQuery {
				decodedQuery, _ := url.QueryUnescape(req.URL.RawQuery)

				logger.Debug("Applied Query Templating: ", decodedQuery)
			}
		}

		if req.URL.Path != "" {
			if templating.Path {
				oldPath := req.URL.Path

				req.URL.Path, err = TemplatePath(req.URL.Path, variables)

				if err != nil {
					logger.Error("Error Templating Path: ", err.Error())
				}

				if req.URL.Path != oldPath {
					logger.Debug("Applied Path Templating: ", req.URL.Path)
				}
			}

			if urlToBody.Path {
				var modified bool
				req.URL.Path, modified = InjectPathIntoBody(req.URL.Path, body.Data)

				if modified {
					modifiedBody = true
				}
			}
		}

		if modifiedBody {
			err := body.UpdateReq(req)

			if err != nil {
				logger.Error("Could not write to Request Body: ", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			logger.Debug("Applied Body Templating: ", body.Data)
		}

		next.ServeHTTP(w, req)
	})
}