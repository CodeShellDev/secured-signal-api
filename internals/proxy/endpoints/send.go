package endpoints

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	request "github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
	"github.com/codeshelldev/secured-signal-api/internals/scheduler"
)

var SendEnpoint = Endpoint{
	Name: "Send",
	Handler: sendHandler,
}

const messageField = "message"
const sendAtField = "send_at"

func sendHandler(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /v2/send", func(w http.ResponseWriter, req *http.Request) {
		logger := GetLogger(req)

		conf := GetConfigByReq(req)

		variables := conf.SETTINGS.MESSAGE.VARIABLES.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.VARIABLES)
		messageTemplate := conf.SETTINGS.MESSAGE.TEMPLATE.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.TEMPLATE)

		body, err := request.GetReqBody(req)

		if err != nil {
			logger.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
			return
		}

		bodyData := map[string]any{}

		var modifiedBody bool

		if !body.Empty {
			bodyData = body.Data

			if messageTemplate != "" {
				headerData := request.GetReqHeaders(req)

				newData, err := TemplateMessage(messageTemplate, bodyData, headerData, variables)

				if err != nil {
					logger.Error("Error Templating Message: ", err.Error())
				}

				if newData[messageField] != bodyData[messageField] && newData[messageField] != "" && newData[messageField] != nil {
					bodyData = newData
					modifiedBody = true
				}
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

			logger.Debug("Applied Message Templating: ", body.Data)
		}

		sendAtStr, ok := bodyData[sendAtField].(string)

		if ok && bodyData[messageField] != "" && bodyData[messageField] != nil {
			delete(bodyData, sendAtField)

			body.Data = bodyData

			body.UpdateReq(req)

			tm, err := handleScheduledMessage(sendAtStr, w, req)

			if err != nil {
				logger.Warn("Could not schedule request: ", err.Error())
				return
			}

			logger.Debug("Scheduled message for ", tm.Local().Format("02.01.06 15:04:05"))

			return
		}
	})

	return mux
}

func handleScheduledMessage(sendAtStr string, w http.ResponseWriter, req *http.Request) (time.Time, error) {
	sendAt, err := strconv.Atoi(sendAtStr)

	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid timestamp: invalid unix number string")
		return time.Time{}, errors.New("invalid timestamp")
	}

	tm := time.Unix(int64(sendAt), 0)

	if tm.Before(time.Now()) {
		WriteError(w, http.StatusBadRequest, "invalid timestamp: time lies in the past")
		return time.Time{}, errors.New("timestamp expired")
	}

	ChangeRequestDest(req, config.DEFAULT.API.URL.String() + req.URL.Path)

	reqID, err := scheduler.ScheduleRequest(tm, req)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return time.Time{}, err
	}

	res := request.Body{
		Data: map[string]any{
			"id": reqID,
		},
	}

	w.WriteHeader(http.StatusAccepted)

	res.Write(w)

	return tm, nil
}

func TemplateMessage(template string, bodyData map[string]any, headerData map[string][]string, variables map[string]any) (map[string]any, error) {
	bodyData["message_template"] = template

	data, _, err := TemplateBody(bodyData, headerData, variables)

	if err != nil || data == nil {
		return bodyData, err
	}

	data[messageField] = data["message_template"]

	delete(data, "message_template")

	return data, nil
}
