-- Invoice items without an invoice are pending items (Stripe semantics): they attach
-- to the next invoice created for the customer unless that create passes
-- pending_invoice_items_behavior=exclude. The subscription column preserves the
-- caller's subscription reference (BO out-of-band subscription seeding).
CREATE TABLE invoice_items_backup (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL REFERENCES customers(id),
	invoice_id TEXT REFERENCES invoices(id),
	subscription_id TEXT NOT NULL DEFAULT '',
	amount INTEGER NOT NULL,
	currency TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	metadata TEXT NOT NULL DEFAULT '{}',
	created_at TEXT NOT NULL
);

INSERT INTO invoice_items_backup (id, customer_id, invoice_id, subscription_id, amount, currency, description, metadata, created_at)
SELECT id, customer_id, invoice_id, '', amount, currency, description, metadata, created_at FROM invoice_items;

DROP TABLE invoice_items;
ALTER TABLE invoice_items_backup RENAME TO invoice_items;

CREATE INDEX IF NOT EXISTS idx_invoice_items_customer ON invoice_items(customer_id);
CREATE INDEX IF NOT EXISTS idx_invoice_items_invoice ON invoice_items(invoice_id);

-- Out-of-band credit notes (paid outside Stripe) keep refunds empty and skip the
-- customer balance; ds2's out-of-band refund path requires both amounts verbatim.
ALTER TABLE credit_notes ADD COLUMN out_of_band_amount INTEGER NOT NULL DEFAULT 0;
ALTER TABLE credit_notes ADD COLUMN memo TEXT NOT NULL DEFAULT '';
