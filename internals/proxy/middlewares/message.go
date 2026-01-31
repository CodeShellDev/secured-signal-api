package middlewares

import (
	"net/http"
	"path"
	"strings"

	request "github.com/codeshelldev/gotl/pkg/request"
)

var Message Middleware = Middleware{
	Name: "Message",
	Use: messageHandler,
}

const templateMessageEndpoint = "/v2/send"

func messageHandler(next http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(templateMessageEndpoint, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			next.ServeHTTP(w, req)
		}

		logger := getLogger(req)

		conf := getConfigByReq(req)

		variables := conf.SETTINGS.MESSAGE.VARIABLES
		messageTemplate := conf.SETTINGS.MESSAGE.TEMPLATE

		if variables == nil {
			variables = getConfig("").SETTINGS.MESSAGE.VARIABLES
		}

		if strings.TrimSpace(messageTemplate) == "" {
			messageTemplate = getConfig("").SETTINGS.MESSAGE.TEMPLATE
		}

		body, err := request.GetReqBody(req)

		if err != nil {
			logger.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
			return
		}

		bodyData := map[string]any{}

		var modifiedBody bool

		if !body.Empty && path.Clean(req.URL.Path) == templateMessageEndpoint {
			bodyData = body.Data

			if messageTemplate != "" {
				headerData := request.GetReqHeaders(req)

				newData, err := TemplateMessage(messageTemplate, bodyData, headerData, variables)

				if err != nil {
					logger.Error("Error Templating Message: ", err.Error())
				}

				if newData["message"] != bodyData["message"] && newData["message"] != "" && newData["message"] != nil {
					bodyData = newData
					modifiedBody = true
				}
			}
		}

		if modifiedBody {
			body.Data = bodyData

			err := body.Write(req)

			if err != nil {
				logger.Error("Could not write to Request Body: ", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			logger.Debug("Applied Message Templating: ", body.Data)
		}

		next.ServeHTTP(w, req)
	})

	mux.Handle("/", next)

	return mux
}

func TemplateMessage(template string, bodyData map[string]any, headerData map[string][]string, variables map[string]any) (map[string]any, error) {
	bodyData["message_template"] = template

	data, _, err := TemplateBody(bodyData, headerData, variables)

	if err != nil || data == nil {
		return bodyData, err
	}

	data["message"] = data["message_template"]

	delete(data, "message_template")

	return data, nil
}
