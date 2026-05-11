CREATE TABLE IF NOT EXISTS test_clocks (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL,
	frozen_time TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS refunds (
	id TEXT PRIMARY KEY,
	charge_id TEXT NOT NULL DEFAULT '',
	payment_intent_id TEXT NOT NULL DEFAULT '',
	invoice_id TEXT NOT NULL DEFAULT '',
	customer_id TEXT NOT NULL DEFAULT '',
	amount INTEGER NOT NULL,
	currency TEXT NOT NULL,
	reason TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL,
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS credit_notes (
	id TEXT PRIMARY KEY,
	invoice_id TEXT NOT NULL,
	customer_id TEXT NOT NULL DEFAULT '',
	amount INTEGER NOT NULL,
	currency TEXT NOT NULL,
	reason TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL,
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_test_clocks_status ON test_clocks(status);
CREATE INDEX IF NOT EXISTS idx_refunds_charge ON refunds(charge_id);
CREATE INDEX IF NOT EXISTS idx_refunds_payment_intent ON refunds(payment_intent_id);
CREATE INDEX IF NOT EXISTS idx_refunds_invoice ON refunds(invoice_id);
CREATE INDEX IF NOT EXISTS idx_refunds_customer ON refunds(customer_id);
CREATE INDEX IF NOT EXISTS idx_credit_notes_invoice ON credit_notes(invoice_id);
CREATE INDEX IF NOT EXISTS idx_credit_notes_customer ON credit_notes(customer_id);
