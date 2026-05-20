CREATE TABLE payment_intents_backup AS
SELECT
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
	metadata,
	created_at
FROM payment_intents;

DROP TABLE payment_intents;

CREATE TABLE invoices_new (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL REFERENCES customers(id),
	subscription_id TEXT REFERENCES subscriptions(id),
	status TEXT NOT NULL,
	currency TEXT NOT NULL,
	subtotal INTEGER NOT NULL,
	discount_amount INTEGER NOT NULL DEFAULT 0,
	discounts TEXT NOT NULL DEFAULT '[]',
	total INTEGER NOT NULL,
	amount_due INTEGER NOT NULL,
	amount_paid INTEGER NOT NULL,
	attempt_count INTEGER NOT NULL,
	next_payment_attempt TEXT,
	payment_intent_id TEXT,
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL
);

INSERT INTO invoices_new (
	id,
	customer_id,
	subscription_id,
	status,
	currency,
	subtotal,
	discount_amount,
	discounts,
	total,
	amount_due,
	amount_paid,
	attempt_count,
	next_payment_attempt,
	payment_intent_id,
	metadata,
	created_at
)
SELECT
	id,
	customer_id,
	subscription_id,
	status,
	currency,
	subtotal,
	discount_amount,
	discounts,
	total,
	amount_due,
	amount_paid,
	attempt_count,
	next_payment_attempt,
	payment_intent_id,
	'{}',
	created_at
FROM invoices;

DROP TABLE invoices;

ALTER TABLE invoices_new RENAME TO invoices;

CREATE TABLE payment_intents (
	id TEXT PRIMARY KEY,
	customer_id TEXT REFERENCES customers(id),
	invoice_id TEXT REFERENCES invoices(id),
	amount INTEGER NOT NULL,
	currency TEXT NOT NULL,
	status TEXT NOT NULL,
	capture_method TEXT NOT NULL DEFAULT 'automatic',
	failure_code TEXT NOT NULL DEFAULT '',
	failure_decline_code TEXT NOT NULL DEFAULT '',
	failure_message TEXT NOT NULL DEFAULT '',
	payment_method_id TEXT NOT NULL DEFAULT '',
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL
);

INSERT INTO payment_intents (
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
	metadata,
	created_at
)
SELECT
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
	metadata,
	created_at
FROM payment_intents_backup;

DROP TABLE payment_intents_backup;

CREATE TABLE IF NOT EXISTS invoice_items (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL REFERENCES customers(id),
	invoice_id TEXT NOT NULL REFERENCES invoices(id),
	amount INTEGER NOT NULL,
	currency TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_invoices_customer ON invoices(customer_id);
CREATE INDEX IF NOT EXISTS idx_invoices_subscription ON invoices(subscription_id);
CREATE INDEX IF NOT EXISTS idx_payment_intents_customer ON payment_intents(customer_id);
CREATE INDEX IF NOT EXISTS idx_payment_intents_invoice ON payment_intents(invoice_id);
CREATE INDEX IF NOT EXISTS idx_invoice_items_customer ON invoice_items(customer_id);
CREATE INDEX IF NOT EXISTS idx_invoice_items_invoice ON invoice_items(invoice_id);
