package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	request "github.com/codeshelldev/secured-signal-api/utils/request"
)

type MessageMiddleware struct {
	Next http.Handler
}

func (data MessageMiddleware) Use() http.Handler {
	next := data.Next

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		settings := getSettingsByReq(req)

		variables := settings.VARIABLES
		messageTemplate := settings.MESSAGE_TEMPLATE

		if variables == nil {
			variables = getSettings("*").VARIABLES
		}

		if messageTemplate == "" {
			messageTemplate = getSettings("*").MESSAGE_TEMPLATE
		}

		body, err := request.GetReqBody(w, req)

		if err != nil {
			log.Error("Could not get Request Body: ", err.Error())
		}

		bodyData := map[string]any{}

		var modifiedBody bool

		if !body.Empty {
			bodyData = body.Data

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
			modifiedBody, err := request.CreateBody(bodyData)

			if err != nil {
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			body = modifiedBody

			strData := body.ToString()

			log.Debug("Applied Message Templating: ", strData)

			req.ContentLength = int64(len(strData))
			req.Header.Set("Content-Length", strconv.Itoa(len(strData)))
		}

		req.Body = io.NopCloser(bytes.NewReader(body.Raw))

		next.ServeHTTP(w, req)
	})
}

func TemplateMessage(template string, bodyData map[string]any, headerData map[string]any, variables map[string]any) (map[string]any, error) {
	bodyData["message_template"] = template

	data, _, err := TemplateBody(bodyData, headerData, variables)

	if err != nil || data == nil {
		return bodyData, err
	}

	data["message"] = data["message_template"]

	delete(data, "message_template")

	return data, nil
}
