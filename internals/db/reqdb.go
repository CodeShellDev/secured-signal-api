package db

import (
	"database/sql"
	"time"
)

type ScheduledRequest struct {
	ID         	string
	Method     	string
	URL        	string
	Headers    	map[string][]string
	Body       	[]byte
	RunAt      	time.Time
	CreatedAt  	time.Time
}

type RequestResult struct {
	RunAt  		time.Time
	FinishedAt  *time.Time
	LastError	*string
	Status		*int
	Headers    	*map[string][]string
	Body       	*[]byte
}

type ScheduledRequestEntry struct {
	Status		ScheduledRequestStatus
	ID         	string
	Method     	string
	URL        	string
	Headers    	map[string][]string
	Body       	[]byte
	RunAt      	time.Time
	CreatedAt  	time.Time

	FinishedAt  *time.Time
	LastError	*string
	ResponseStatusCode *int
	ResponseBody *[]byte
	ResponseHeaders *map[string][]string
}

type ScheduledRequestStatus string

const (
	STATUS_PENDING ScheduledRequestStatus = "pending"
	STATUS_QUEUED ScheduledRequestStatus = "queued"
	STATUS_DONE ScheduledRequestStatus = "done"
	STATUS_FAILED ScheduledRequestStatus = "failed"
	STATUS_RUNNING ScheduledRequestStatus = "running"
)

type RequestSchedulerDB struct {
	db *sql.DB
}

func NewRequestSchedulerDB() *RequestSchedulerDB {
    return &RequestSchedulerDB{db: db}
}

func ShutdownRequestDB() {
	NewRequestSchedulerDB().RecoverStales(0)
}

func (s *RequestSchedulerDB) Insert(req *ScheduledRequest) error {
	_, err := s.db.Exec(`
		INSERT INTO scheduled_requests (id, status, method, url, created_at, run_at, request_headers, request_body)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.ID,
		STATUS_PENDING,
		req.Method,
		req.URL,
		time.Now().Unix(),
		req.RunAt.Unix(),
		Serialize(req.Headers),
		req.Body,
	)

	return err
}

func (s *RequestSchedulerDB) SetStatus(id string, status ScheduledRequestStatus) error {
	_, err := s.db.Exec(`
		UPDATE scheduled_requests
		SET status = ?
		WHERE id = ?`,
		string(status),
		id,
	)

	return err
}

func (s *RequestSchedulerDB) SetResponse(id string, e error, res RequestResult) error {
	var errMsg string

	if e != nil {
		errMsg = e.Error()
	}

	var finishedAt *int64

	if res.FinishedAt != nil {
		tm := res.FinishedAt.Unix()

		finishedAt = &tm
	}

	_, err := s.db.Exec(`
		UPDATE scheduled_requests
		SET finished_at = ?, last_error = ?, response_status_code = ?, response_headers = ?, response_body = ?
		WHERE id = ?`,
		finishedAt,
		errMsg,
		res.Status,
		Serialize(res.Headers),
		res.Body,
		id,
	)

	return err
}

func (s *RequestSchedulerDB) GetByID(id string) (*ScheduledRequestEntry, error) {
	row := s.db.QueryRow(`
		SELECT status, id, method, url, created_at, run_at, request_headers, request_body, finished_at, last_error, response_status_code, response_headers, response_body
		FROM scheduled_requests
		WHERE id = ?`,
		id,
	)

	res := &ScheduledRequestEntry{}

	var requestHeaderBytes []byte
	var responseHeaderBytes *[]byte

	var createdAt, runAt int64
	var finishedAt *int64

	err := row.Scan(&res.Status, &res.ID, &res.Method, &res.URL, &createdAt, &runAt, &requestHeaderBytes, &res.Body, &finishedAt, &res.LastError, &res.ResponseStatusCode, &responseHeaderBytes, &res.ResponseBody)

	if err != nil {
		return nil, err
	}

	res.Headers = Deserialize[map[string][]string](requestHeaderBytes)

	if responseHeaderBytes != nil {
		res.ResponseHeaders = Deserialize[*map[string][]string](*responseHeaderBytes)
	}

	res.CreatedAt = time.Unix(createdAt, 0)
	res.RunAt = time.Unix(runAt, 0)

	if finishedAt != nil {
		finished := time.Unix(*finishedAt, 0)
		res.FinishedAt = &finished
	}

	return res, nil
}

func (s *RequestSchedulerDB) DeleteByID(id string) error {
	_, err := s.db.Exec(`
		DELETE FROM scheduled_requests
		WHERE id = ?`,
		id,
	)

	return err
}

func (s *RequestSchedulerDB) Claim(id string) error {
	_, err := s.db.Exec(`
		UPDATE scheduled_requests
		SET claimed_at = ?
		WHERE id = ?`,
		time.Now().Unix(),
		id,
	)

	return err
}

func (s *RequestSchedulerDB) RecoverStales(threshold time.Duration) error {
	minClaimedAt := time.Now().Add(-threshold).Unix()

	_, err := s.db.Exec(`
		UPDATE scheduled_requests
		SET status = ?, claimed_at = null
		WHERE status != ? AND (claimed_at IS NULL OR claimed_at <= ?)`,
		STATUS_PENDING,
		STATUS_DONE,
		minClaimedAt,
	)

	return err
}

func (s *RequestSchedulerDB) CleanupDones(threshold time.Duration) error {
	minFinishedAt := time.Now().Add(-threshold).Unix()

	_, err := s.db.Exec(`
		DELETE FROM scheduled_requests
		WHERE status = ? AND finished_at <= ?`,
		STATUS_DONE,
		minFinishedAt,
	)

	return err
}

func (s *RequestSchedulerDB) FetchNext(amount int, within time.Duration) ([]*ScheduledRequest, error) {
	minRunAt := time.Now().Add(within).Unix()

	rows, err := s.db.Query(`
		SELECT id, method, url, created_at, run_at, request_headers, request_body
		FROM scheduled_requests
        WHERE status = ? AND run_at <= ?
		ORDER BY run_at ASC
		LIMIT ?`,
		STATUS_PENDING,
		minRunAt,
		amount,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*ScheduledRequest

	for rows.Next() {
		req := &ScheduledRequest{}

		var headerBytes []byte
		var createdAt, runAt int64

		err := rows.Scan(&req.ID, &req.Method, &req.URL, &createdAt, &runAt, &headerBytes, &req.Body)
		if err != nil {
			return nil, err
		}

		req.Headers = Deserialize[map[string][]string](headerBytes)
		req.CreatedAt = time.Unix(createdAt, 0)
		req.RunAt = time.Unix(runAt, 0)

		requests = append(requests, req)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return requests, nil
}