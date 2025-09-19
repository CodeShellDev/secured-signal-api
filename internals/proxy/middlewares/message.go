package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/codeshelldev/secured-signal-api/utils/jsonutils"
	log "github.com/codeshelldev/secured-signal-api/utils/logger"
	request "github.com/codeshelldev/secured-signal-api/utils/request"
)

type MessageMiddleware struct {
	Next      http.Handler
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
			templatedMessage, err := TemplateMessage(messageTemplate, bodyData, variables)

			if err != nil {
				log.Error("Error Templating Message: ", err.Error())
			}

			if templatedMessage != bodyData["message"] && templatedMessage != "" {
				bodyData["message"] = templatedMessage
				modifiedBody = true
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

			log.Debug("Applied Body Templating: ", strData)

			req.ContentLength = int64(len(strData))
			req.Header.Set("Content-Length", strconv.Itoa(len(strData)))
		}

		req.Body = io.NopCloser(bytes.NewReader(body.Raw))

		next.ServeHTTP(w, req)
	})
}

func TemplateMessage(template string, data map[string]any, VARIABLES any) (string, error) {
	data, ok, err := TemplateBody(data, VARIABLES)

	if err != nil || !ok || data == nil {
		return template, err
	}

	jsonStr, err := jsonutils.ToJsonSafe(data)

	if err != nil {
		return template, err
	}

	return jsonStr, nil
}