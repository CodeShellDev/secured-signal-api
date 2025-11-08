package middlewares

import (
	"net/http"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	request "github.com/codeshelldev/secured-signal-api/utils/request"
)

var Message Middleware = Middleware{
	Name: "Message",
	Use: messageHandler,
}

func messageHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		settings := getSettingsByReq(req)

		variables := settings.MESSAGE.VARIABLES
		messageTemplate := settings.MESSAGE.TEMPLATE

		if variables == nil {
			variables = getSettings("*").MESSAGE.VARIABLES
		}

		if messageTemplate == "" {
			messageTemplate = getSettings("*").MESSAGE.TEMPLATE
		}

		body, err := request.GetReqBody(req)

		if err != nil {
			log.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
		}

		bodyData := map[string]any{}

		log.Dev("BodyData: ", body.Data)

		var modifiedBody bool

		if !body.Empty {
			bodyData = body.Data

			log.Dev("Template: ", messageTemplate)

			if messageTemplate != "" {
				headerData := request.GetReqHeaders(req)

				newData, err := TemplateMessage(messageTemplate, bodyData, headerData, variables)

				if err != nil {
					log.Error("Error Templating Message: ", err.Error())
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
				log.Error("Could not write to Request Body: ", err.Error())
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			log.Debug("Applied Message Templating: ", body.Data)
		}

		next.ServeHTTP(w, req)
	})
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
