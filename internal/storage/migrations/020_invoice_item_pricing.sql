-- Invoice items created via pricing[price] + quantity (newer Stripe API shape)
-- must survive reads: without these columns the create response carried the
-- fields but GET/LIST silently dropped them.
ALTER TABLE invoice_items ADD COLUMN price_id TEXT NOT NULL DEFAULT '';
ALTER TABLE invoice_items ADD COLUMN product_id TEXT NOT NULL DEFAULT '';
ALTER TABLE invoice_items ADD COLUMN quantity INTEGER NOT NULL DEFAULT 0;
