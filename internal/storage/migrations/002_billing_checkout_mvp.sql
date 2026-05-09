CREATE TABLE IF NOT EXISTS customers (
	id TEXT PRIMARY KEY,
	email TEXT NOT NULL DEFAULT '',
	name TEXT NOT NULL DEFAULT '',
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	active INTEGER NOT NULL DEFAULT 1,
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS prices (
	id TEXT PRIMARY KEY,
	product_id TEXT NOT NULL REFERENCES products(id),
	currency TEXT NOT NULL,
	unit_amount INTEGER NOT NULL,
	recurring_interval TEXT NOT NULL DEFAULT '',
	recurring_interval_count INTEGER NOT NULL DEFAULT 1,
	active INTEGER NOT NULL DEFAULT 1,
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS checkout_sessions (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL REFERENCES customers(id),
	mode TEXT NOT NULL,
	line_items TEXT NOT NULL,
	success_url TEXT NOT NULL DEFAULT '',
	cancel_url TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL,
	payment_status TEXT NOT NULL,
	subscription_id TEXT,
	invoice_id TEXT,
	payment_intent_id TEXT,
	created_at TEXT NOT NULL,
	completed_at TEXT
);

CREATE TABLE IF NOT EXISTS subscriptions (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL REFERENCES customers(id),
	status TEXT NOT NULL,
	items TEXT NOT NULL,
	current_period_start TEXT NOT NULL,
	current_period_end TEXT NOT NULL,
	cancel_at_period_end INTEGER NOT NULL DEFAULT 0,
	canceled_at TEXT,
	latest_invoice_id TEXT,
	metadata TEXT NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS invoices (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL REFERENCES customers(id),
	subscription_id TEXT NOT NULL REFERENCES subscriptions(id),
	status TEXT NOT NULL,
	currency TEXT NOT NULL,
	subtotal INTEGER NOT NULL,
	total INTEGER NOT NULL,
	amount_due INTEGER NOT NULL,
	amount_paid INTEGER NOT NULL,
	attempt_count INTEGER NOT NULL,
	next_payment_attempt TEXT,
	payment_intent_id TEXT,
	created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS payment_intents (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL REFERENCES customers(id),
	invoice_id TEXT NOT NULL REFERENCES invoices(id),
	amount INTEGER NOT NULL,
	currency TEXT NOT NULL,
	status TEXT NOT NULL,
	failure_code TEXT NOT NULL DEFAULT '',
	failure_message TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS timeline_entries (
	id TEXT PRIMARY KEY,
	action TEXT NOT NULL,
	message TEXT NOT NULL,
	object_type TEXT NOT NULL,
	object_id TEXT NOT NULL,
	customer_id TEXT NOT NULL DEFAULT '',
	checkout_session_id TEXT NOT NULL DEFAULT '',
	subscription_id TEXT NOT NULL DEFAULT '',
	invoice_id TEXT NOT NULL DEFAULT '',
	payment_intent_id TEXT NOT NULL DEFAULT '',
	data TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_checkout_sessions_customer ON checkout_sessions(customer_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_customer ON subscriptions(customer_id);
CREATE INDEX IF NOT EXISTS idx_invoices_customer ON invoices(customer_id);
CREATE INDEX IF NOT EXISTS idx_payment_intents_customer ON payment_intents(customer_id);
CREATE INDEX IF NOT EXISTS idx_timeline_customer ON timeline_entries(customer_id);
CREATE INDEX IF NOT EXISTS idx_timeline_checkout_session ON timeline_entries(checkout_session_id);
CREATE INDEX IF NOT EXISTS idx_timeline_subscription ON timeline_entries(subscription_id);
CREATE INDEX IF NOT EXISTS idx_timeline_invoice ON timeline_entries(invoice_id);
CREATE INDEX IF NOT EXISTS idx_timeline_payment_intent ON timeline_entries(payment_intent_id);

