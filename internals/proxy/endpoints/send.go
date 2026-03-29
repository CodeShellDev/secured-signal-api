package endpoints

import (
	"errors"
	"net/http"
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

func sendHandler(mux *http.ServeMux, next http.Handler) *http.ServeMux {
	mux.HandleFunc("POST /v2/send", func(w http.ResponseWriter, req *http.Request) {
		logger := GetLogger(req)

		conf := GetConfigByReq(req)

		variables := conf.SETTINGS.MESSAGE.VARIABLES.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.VARIABLES)
		templating := conf.SETTINGS.MESSAGE.TEMPLATING.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.TEMPLATING)

		scheduling := conf.SETTINGS.MESSAGE.SCHEDULING.OptOrEmpty(config.DEFAULT.SETTINGS.MESSAGE.SCHEDULING)

		body, err := request.GetReqBody(req)

		if err != nil {
			logger.Error("Could not get Request Body: ", err.Error())
			http.Error(w, "Bad Request: invalid body", http.StatusBadRequest)
			return
		}

		body.EnsureNotNil()

		var modifiedBody bool

		if !body.Empty {
			if templating.MessageTemplate != "" {
				headers := request.GetReqHeaders(req)

				templatedMessage, err := GetTemplatedMessage(templating.MessageTemplate, body.Data, headers, variables)

				if err != nil {
					logger.Error("Error Templating Message: ", err.Error())
				}

				if templatedMessage != body.Data[messageField] && templatedMessage != "" {
					body.Data[messageField] = templatedMessage

					logger.Debug("Applied Message Templating: \n", templatedMessage)

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
		}

		sendAt, ok := body.Data[sendAtField].(float64)

		if ok && body.Data[messageField] != "" && body.Data[messageField] != nil {
			delete(body.Data, sendAtField)

			body.UpdateReq(req)

			if !scheduling.Enabled {
				logger.Warn("Client tried scheduling message")
				WriteError(w, http.StatusForbidden, "scheduling is disabled")
				return
			}

			tm, err := parseTimestamp(int(sendAt))

			if err != nil {
				logger.Warn("Could not parse timestamp: ", err.Error())
				WriteError(w, http.StatusBadRequest, "invalid timestamp: " + err.Error())
				return
			}

			if scheduling.MaxHorizon.Set {
				if tm.After(time.Now().Add(time.Duration(*scheduling.MaxHorizon.Value))) {
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

		next.ServeHTTP(w, req)
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

func parseTimestamp(sendAt int) (time.Time, error) {
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

func GetTemplatedMessage(template string, body map[string]any, headers map[string][]string, VARIABLES map[string]any) (string, error) {
	const templatedSuffix = "_template"

	bodyCopy := map[string]any{
		messageField + templatedSuffix: template,
	}

	request.CopyMap(bodyCopy, body)

	data, _, err := GetTemplatedBody(bodyCopy, headers, VARIABLES)

	if err != nil || data == nil {
		return "", err
	}

	return data[messageField + templatedSuffix].(string), nil
}
