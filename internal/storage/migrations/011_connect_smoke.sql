CREATE TABLE IF NOT EXISTS connect_accounts (
	id TEXT PRIMARY KEY,
	type TEXT NOT NULL,
	country TEXT NOT NULL,
	email TEXT NOT NULL DEFAULT '',
	business_type TEXT NOT NULL DEFAULT '',
	default_currency TEXT NOT NULL,
	charges_enabled INTEGER NOT NULL DEFAULT 1,
	payouts_enabled INTEGER NOT NULL DEFAULT 1,
	details_submitted INTEGER NOT NULL DEFAULT 1,
	capabilities TEXT NOT NULL DEFAULT '{}',
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_connect_accounts_country ON connect_accounts(country);
CREATE INDEX IF NOT EXISTS idx_connect_accounts_type ON connect_accounts(type);
