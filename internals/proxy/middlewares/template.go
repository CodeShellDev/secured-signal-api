package middlewares

import (
	"net/http"

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
		
		variables := conf.SETTINGS.MESSAGE.VARIABLES.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.VARIABLES)

		body, err := request.GetReqBody(req)

		if err != nil {
			logger.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
			return
		}

		bodyData := map[string]any{}

		var modifiedBody bool

		if !body.Empty {
			var modified bool

			headerData := request.GetReqHeaders(req)

			bodyData, modified, err = TemplateBody(body.Data, headerData, variables)

			if err != nil {
				logger.Error("Error Templating JSON: ", err.Error())
			}

			if modified {
				modifiedBody = true
			}
		}

		if req.URL.RawQuery != "" {
			var modified bool

			req.URL.RawQuery, bodyData, modified, err = TemplateQuery(req.URL, bodyData, variables)

			if err != nil {
				logger.Error("Error Templating Query: ", err.Error())
			}

			if modified {
				modifiedBody = true
			}
		}

		if req.URL.Path != "" {
			var modified bool
			var templated bool

			req.URL.Path, bodyData, modified, templated, err = TemplatePath(req.URL, bodyData, variables)

			if err != nil {
				logger.Error("Error Templating Path: ", err.Error())
			}

			if modified {
				logger.Debug("Applied Path Templating: ", req.URL.Path)
			}

			if templated {
				modifiedBody = true
			}
		}

		if modifiedBody {
			body.Data = bodyData

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