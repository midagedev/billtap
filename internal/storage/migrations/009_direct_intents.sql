CREATE TABLE payment_intents_new (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL DEFAULT '',
	invoice_id TEXT NOT NULL DEFAULT '',
	amount INTEGER NOT NULL,
	currency TEXT NOT NULL,
	status TEXT NOT NULL,
	capture_method TEXT NOT NULL DEFAULT 'automatic',
	failure_code TEXT NOT NULL DEFAULT '',
	failure_decline_code TEXT NOT NULL DEFAULT '',
	failure_message TEXT NOT NULL DEFAULT '',
	payment_method_id TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL
);

INSERT INTO payment_intents_new (
	id,
	customer_id,
	invoice_id,
	amount,
	currency,
	status,
	capture_method,
	failure_code,
	failure_decline_code,
	failure_message,
	payment_method_id,
	created_at
)
SELECT
	id,
	customer_id,
	invoice_id,
	amount,
	currency,
	status,
	'automatic',
	failure_code,
	failure_decline_code,
	failure_message,
	payment_method_id,
	created_at
FROM payment_intents;

DROP TABLE payment_intents;

ALTER TABLE payment_intents_new RENAME TO payment_intents;

CREATE TABLE IF NOT EXISTS setup_intents (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL,
	usage TEXT NOT NULL DEFAULT 'off_session',
	failure_code TEXT NOT NULL DEFAULT '',
	failure_decline_code TEXT NOT NULL DEFAULT '',
	failure_message TEXT NOT NULL DEFAULT '',
	payment_method_id TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_payment_intents_customer ON payment_intents(customer_id);
CREATE INDEX IF NOT EXISTS idx_setup_intents_customer ON setup_intents(customer_id);
