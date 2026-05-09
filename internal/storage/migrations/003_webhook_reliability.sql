CREATE TABLE IF NOT EXISTS webhook_endpoints (
	id TEXT PRIMARY KEY,
	url TEXT NOT NULL,
	secret TEXT NOT NULL,
	enabled_events TEXT NOT NULL DEFAULT '[]',
	active INTEGER NOT NULL DEFAULT 1,
	retry_max_attempts INTEGER NOT NULL DEFAULT 5,
	retry_backoff TEXT NOT NULL DEFAULT '["10s","30s","2m","10m"]',
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS webhook_events (
	id TEXT PRIMARY KEY,
	type TEXT NOT NULL,
	created INTEGER NOT NULL,
	livemode INTEGER NOT NULL DEFAULT 0,
	api_version TEXT NOT NULL,
	pending_webhooks INTEGER NOT NULL DEFAULT 0,
	request_id TEXT NOT NULL DEFAULT '',
	idempotency_key TEXT NOT NULL DEFAULT '',
	object_payload TEXT NOT NULL,
	raw_payload TEXT NOT NULL,
	source TEXT NOT NULL,
	sequence INTEGER NOT NULL,
	scenario_run_id TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS delivery_attempts (
	id TEXT PRIMARY KEY,
	event_id TEXT NOT NULL REFERENCES webhook_events(id),
	endpoint_id TEXT NOT NULL REFERENCES webhook_endpoints(id),
	attempt_number INTEGER NOT NULL,
	status TEXT NOT NULL,
	scheduled_at TEXT NOT NULL,
	delivered_at TEXT,
	request_url TEXT NOT NULL,
	request_headers TEXT NOT NULL DEFAULT '{}',
	request_body TEXT NOT NULL DEFAULT '',
	response_status INTEGER NOT NULL DEFAULT 0,
	response_body TEXT NOT NULL DEFAULT '',
	error TEXT NOT NULL DEFAULT '',
	next_retry_at TEXT,
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_webhook_events_type ON webhook_events(type);
CREATE INDEX IF NOT EXISTS idx_webhook_events_sequence ON webhook_events(sequence);
CREATE INDEX IF NOT EXISTS idx_delivery_attempts_event ON delivery_attempts(event_id);
CREATE INDEX IF NOT EXISTS idx_delivery_attempts_endpoint ON delivery_attempts(endpoint_id);
CREATE INDEX IF NOT EXISTS idx_delivery_attempts_status ON delivery_attempts(status);
