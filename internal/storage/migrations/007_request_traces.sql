CREATE TABLE IF NOT EXISTS request_traces (
	id TEXT PRIMARY KEY,
	method TEXT NOT NULL,
	path TEXT NOT NULL,
	query TEXT NOT NULL DEFAULT '',
	status INTEGER NOT NULL,
	duration_ms INTEGER NOT NULL DEFAULT 0,
	request_id TEXT NOT NULL DEFAULT '',
	idempotency_key TEXT NOT NULL DEFAULT '',
	request_headers TEXT NOT NULL DEFAULT '{}',
	request_body TEXT NOT NULL DEFAULT '',
	response_body TEXT NOT NULL DEFAULT '',
	response_object TEXT NOT NULL DEFAULT '',
	response_object_id TEXT NOT NULL DEFAULT '',
	error_type TEXT NOT NULL DEFAULT '',
	error_code TEXT NOT NULL DEFAULT '',
	error_param TEXT NOT NULL DEFAULT '',
	related_ids TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_request_traces_created_at ON request_traces(created_at);
CREATE INDEX IF NOT EXISTS idx_request_traces_path ON request_traces(path);
CREATE INDEX IF NOT EXISTS idx_request_traces_status ON request_traces(status);
CREATE INDEX IF NOT EXISTS idx_request_traces_request_id ON request_traces(request_id);
CREATE INDEX IF NOT EXISTS idx_request_traces_response_object_id ON request_traces(response_object_id);
