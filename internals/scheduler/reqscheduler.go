package scheduler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/codeshelldev/gotl/pkg/jsonutils"
	"github.com/codeshelldev/gotl/pkg/logger"
	"github.com/codeshelldev/gotl/pkg/request"
	scheduling "github.com/codeshelldev/gotl/pkg/scheduler"
	"github.com/codeshelldev/secured-signal-api/internals/db"
	"github.com/google/uuid"
)

var rsdb *db.RequestSchedulerDB

var reqscheduler = scheduling.New()
var cancel context.CancelFunc

const limit = 5
const withinTime = 5 * time.Minute

const recoveryThreshold = 10 * time.Minute

const doneStaleThreshold = 24 * time.Hour

func Stop() {
	cancel()
}

func StartRequestScheduler() {
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())

	go reqscheduler.Run(ctx)

	rsdb = db.NewRequestSchedulerDB()

	rsdb.CleanupDones(doneStaleThreshold)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			rsdb.RecoverStales(recoveryThreshold)

			if reqscheduler.Len() < limit {
				UpdateQueue()
			}
		}
	}()
}

func UpdateQueue() {
	requests, _ := rsdb.FetchNext(limit - reqscheduler.Len(), withinTime)

	for _, req := range requests {
		AddToQueue(req)
	}
}

func AddToQueue(req *db.ScheduledRequest) {
	rsdb.SetStatus(req.ID, db.STATUS_QUEUED)
	rsdb.Claim(req.ID)

	reqscheduler.AddAtWithID(req.ID, req.RunAt, func() {
		HandleScheduledRequest(req)
	})
}

func OnRequestScheduled(req *db.ScheduledRequest) {
	id, latest, exists := reqscheduler.Latest()

	if !exists {
		return
	}

	// get next request
	nexts, _ := rsdb.FetchNext(1, 0)

	if len(nexts) != 1 {
		return
	}

	// check if current request runs before the latest heap request
	// and if the next request in the db is also the current request
	if req.RunAt.Before(latest) && req.ID == nexts[0].ID {
		// add current job (earlier than latest)
		AddToQueue(req)

		// remove latest job
		reqscheduler.Cancel(id)

		// mark removed request as pending (not queued anymore)
		rsdb.SetStatus(id, db.STATUS_PENDING)
	}
}

func ScheduleRequest(tm time.Time, req *http.Request) (string, error) {
	if tm.Before(time.Now()) {
		return "", errors.New("time lies in the past")
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return "", err
	}

	id := uuid.NewString()

	scheduledReq := &db.ScheduledRequest{
		ID: id,
		Method: req.Method,
		URL: req.URL.String(),
		Headers: req.Header,
		Body: body,
		RunAt: tm,
		CreatedAt: time.Now(),
	}

	err = rsdb.Insert(scheduledReq)

	if err != nil {
		return "", err
	}

	OnRequestScheduled(scheduledReq)

	return id, nil
}

func HandleScheduledRequest(req *db.ScheduledRequest) {
	rsdb.SetStatus(req.ID, db.STATUS_RUNNING)

	res, err := fireScheduledRequest(req)
	result := db.RequestResult{}

	now := time.Now()
	result.FinishedAt = &now

	if err != nil {
		rsdb.SetStatus(req.ID, db.STATUS_FAILED)
		rsdb.SetResponse(req.ID, err, result)
		
		logger.Error("Could not send scheduled request: ", err.Error())
		return
	}

	body, err := request.GetResBody(res)

	if err != nil {
		body.Raw = nil
	}

	result.Status = &res.StatusCode

	headers := map[string][]string{}
	request.CopyHeaders(headers, res.Header)

	result.Headers = &headers

	bodyCopy := append([]byte(nil), body.Raw...)
	result.Body = &bodyCopy

	rsdb.SetStatus(req.ID, db.STATUS_DONE)
	rsdb.SetResponse(req.ID, nil, result)

	URL, _ := url.Parse(req.URL)

	if !logger.IsDev() {
		logger.Info("Fired request",
			" from ", req.CreatedAt.Local().Format("02.01.06 15:04:05"), ": ",
			req.Method, " ", URL.Path, " ", URL.RawQuery,
		)
	} else {
		if len(req.Body) != 0{
			logger.Dev("Fired request",
				" from ", req.CreatedAt.Local().Format("02.01.06 15:04:05"), ": ",
				req.Method, " ", URL.Path, " ", URL.RawQuery,
				jsonutils.GetJson[map[string]any](string(req.Body)),
			)
		} else {
			logger.Dev("Fired request",
				" from ", req.CreatedAt.Local().Format("02.01.06 15:04:05"), ": ",
				req.Method, " ", URL.Path, " ", URL.RawQuery,
			)
		}
	}
}

func fireScheduledRequest(req *db.ScheduledRequest) (*http.Response, error) {
    httpReq, _ := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Body))
	
	request.CopyHeaders(httpReq.Header, req.Headers)

	client := &http.Client{}

	return client.Do(httpReq)
}