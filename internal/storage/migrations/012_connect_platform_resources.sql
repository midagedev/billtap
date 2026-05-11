CREATE TABLE IF NOT EXISTS connect_resources (
	id TEXT PRIMARY KEY,
	object TEXT NOT NULL,
	account_id TEXT NOT NULL DEFAULT '',
	parent_id TEXT NOT NULL DEFAULT '',
	amount INTEGER NOT NULL DEFAULT 0,
	currency TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT '',
	description TEXT NOT NULL DEFAULT '',
	destination TEXT NOT NULL DEFAULT '',
	source_transaction TEXT NOT NULL DEFAULT '',
	country TEXT NOT NULL DEFAULT '',
	bank_name TEXT NOT NULL DEFAULT '',
	last4 TEXT NOT NULL DEFAULT '',
	routing_number TEXT NOT NULL DEFAULT '',
	arrival_date TEXT NOT NULL DEFAULT '',
	metadata TEXT NOT NULL DEFAULT '{}',
	data TEXT NOT NULL DEFAULT '{}',
	deleted INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_connect_resources_object ON connect_resources(object);
CREATE INDEX IF NOT EXISTS idx_connect_resources_account ON connect_resources(account_id);
CREATE INDEX IF NOT EXISTS idx_connect_resources_parent ON connect_resources(parent_id);
CREATE INDEX IF NOT EXISTS idx_connect_resources_destination ON connect_resources(destination);
