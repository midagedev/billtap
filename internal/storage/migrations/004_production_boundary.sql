CREATE TABLE IF NOT EXISTS audit_log (
	id TEXT PRIMARY KEY,
	action TEXT NOT NULL,
	actor TEXT NOT NULL DEFAULT 'system',
	target_type TEXT NOT NULL,
	target_id TEXT NOT NULL,
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_audit_log_action ON audit_log(action);
CREATE INDEX IF NOT EXISTS idx_audit_log_target ON audit_log(target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON audit_log(created_at);
