-- Pending invoice items have no invoice yet. invoice_id was NOT NULL with a
-- foreign key, so an empty string could not store a customer-attached pending
-- item. Recreate the table with a nullable invoice_id and persist the optional
-- subscription id the create surface already accepts.
CREATE TABLE invoice_items_new (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL REFERENCES customers(id),
	invoice_id TEXT REFERENCES invoices(id),
	amount INTEGER NOT NULL,
	currency TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL,
	price_id TEXT NOT NULL DEFAULT '',
	product_id TEXT NOT NULL DEFAULT '',
	quantity INTEGER NOT NULL DEFAULT 0,
	subscription_id TEXT NOT NULL DEFAULT ''
);

INSERT INTO invoice_items_new (
	id,
	customer_id,
	invoice_id,
	amount,
	currency,
	description,
	metadata,
	created_at,
	price_id,
	product_id,
	quantity,
	subscription_id
)
SELECT
	id,
	customer_id,
	invoice_id,
	amount,
	currency,
	description,
	metadata,
	created_at,
	price_id,
	product_id,
	quantity,
	''
FROM invoice_items;

DROP TABLE invoice_items;

ALTER TABLE invoice_items_new RENAME TO invoice_items;

CREATE INDEX IF NOT EXISTS idx_invoice_items_customer ON invoice_items(customer_id);
CREATE INDEX IF NOT EXISTS idx_invoice_items_invoice ON invoice_items(invoice_id);
