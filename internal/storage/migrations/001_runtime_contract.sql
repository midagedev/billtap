CREATE TABLE IF NOT EXISTS runtime_metadata (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO runtime_metadata (key, value)
VALUES ('runtime_contract', 'g1')
ON CONFLICT(key) DO UPDATE SET
	value = excluded.value,
	updated_at = CURRENT_TIMESTAMP;

