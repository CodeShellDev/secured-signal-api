package endpoints

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	request "github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/config"
	"github.com/codeshelldev/secured-signal-api/internals/config/structure"
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

		scheduling := conf.SETTINGS.MESSAGE.SCHEDULING.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.SCHEDULING)

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

			if !scheduling.Enabled {
				logger.Warn("Client tried scheduling message")
				WriteError(w, http.StatusForbidden, "scheduling is disabled")
				return
			}

			tm, err := parseTimestamp(sendAtStr)

			if err != nil {
				logger.Warn("Could not parse timestamp: ", err.Error())
				WriteError(w, http.StatusBadRequest, "invalid timestamp: " + err.Error())
				return
			}

			if scheduling.MaxHorizon.Set {
				if tm.After(time.Now().Add(scheduling.MaxHorizon.Value.Duration)) {
					logger.Warn("Request scheduled too far in the future: ", time.Until(tm).String())
					WriteError(w, http.StatusBadRequest, "invalid timestamp: " + "timestamp to far in the future")
					return
				}
			}

			err = handleScheduledMessage(tm, w, req)

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

func getSendCapabilities(conf *structure.CONFIG) []string {
	out := []string{}

	scheduling := conf.SETTINGS.MESSAGE.SCHEDULING.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.SCHEDULING)

	if scheduling.Enabled {
		out = append(out, "scheduling")
	}
	
	return out
}

func parseTimestamp(str string) (time.Time, error) {
	sendAt, err := strconv.Atoi(str)

	if err != nil {
		return time.Time{}, errors.New("invalid number string")
	}

	tm := time.Unix(int64(sendAt), 0)

	if tm.Before(time.Now()) {
		return time.Time{}, errors.New("timestamp expired")
	}

	return tm, nil
}

func handleScheduledMessage(tm time.Time, w http.ResponseWriter, req *http.Request) (error) {
	ChangeRequestDest(req, config.DEFAULT.API.URL.String() + req.URL.Path)

	reqID, err := scheduler.ScheduleRequest(tm, req)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	res := request.Body{
		Data: map[string]any{
			"id": reqID,
		},
	}

	w.WriteHeader(http.StatusAccepted)

	res.Write(w)

	return nil
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
