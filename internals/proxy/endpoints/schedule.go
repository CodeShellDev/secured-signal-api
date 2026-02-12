package endpoints

import (
	"net/http"
	"strconv"

	"github.com/codeshelldev/gotl/pkg/jsonutils"
	"github.com/codeshelldev/gotl/pkg/request"
	"github.com/codeshelldev/secured-signal-api/internals/db"
	. "github.com/codeshelldev/secured-signal-api/internals/proxy/common"
)

var ScheduleEndpoint = Endpoint{
	Name: "Schedule",
	Handler: scheduleHandler,
}

func scheduleHandler(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("DELETE /v1/schedule/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")

		rsdb := db.NewRequestSchedulerDB()
		err := rsdb.DeleteByID(id)

		if err != nil {
			WriteError(w, http.StatusBadRequest, "request not found")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("GET /v1/schedule/{id}", func(w http.ResponseWriter, req *http.Request) {
		logger := GetLogger(req)

		id := req.PathValue("id")

		rsdb := db.NewRequestSchedulerDB()
		entry, err := rsdb.GetByID(id)

		if err != nil {
			WriteError(w, http.StatusBadRequest, "request not found")
			return
		}

		body, _ := request.CreateBody(map[string]any{
			"status": string(entry.Status),
			
			"method": entry.Method,
			"url": entry.URL,

			"created_at": strconv.Itoa(int(entry.CreatedAt.Unix())),
			"run_at": strconv.Itoa(int(entry.RunAt.Unix())),
		})

		if entry.Status != db.STATUS_DONE && entry.Status != db.STATUS_FAILED {
			body.Write(w)
			return
		}

		var finishedAt *string
		if entry.FinishedAt != nil {
			finished := entry.FinishedAt.Unix()

			tm := strconv.Itoa(int(finished))
			finishedAt = &tm
		}

		body.Data["finished_at"] = finishedAt
		body.Data["response_status_code"] = entry.ResponseStatusCode

		if entry.Status == db.STATUS_FAILED {
			body.Data["error"] = entry.LastError
			body.Write(w)
			return
		}

		var resBody *map[string]any
		if entry.ResponseBody != nil {
			resBody, err = jsonutils.GetJsonSafe[*map[string]any](string(*entry.ResponseBody))
		}

		if err != nil {
			logger.Error("Error parsing json string: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if resBody != nil {
			body.Data["response_body"] = resBody
		}

		if entry.Headers != nil {
			body.Data["response_headers"] = entry.ResponseHeaders
		}

		err = body.Write(w)

		if err != nil {
			logger.Error("Could not write to Response Body: ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		rsdb.DeleteByID(id)
	})

	return mux
}