PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS scheduled_requests (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL,

    method TEXT NOT NULL,
    url TEXT NOT NULL,

    created_at INTEGER NOT NULL,
    claimed_at INTEGER,
    run_at INTEGER NOT NULL,

    /* before run */
    request_headers BLOB NOT NULL,
    request_body BLOB NOT NULL,

    /* after run */
    finished_at INTEGER,
    last_error TEXT,
    response_status_code INTEGER,
    response_headers BLOB,
    response_body BLOB
);

CREATE INDEX IF NOT EXISTS idx_scheduled_requests_run_at
ON scheduled_requests (run_at)
WHERE status = 'pending';
