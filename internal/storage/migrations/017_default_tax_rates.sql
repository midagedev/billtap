ALTER TABLE checkout_sessions ADD COLUMN default_tax_rates TEXT NOT NULL DEFAULT '[]';
ALTER TABLE invoices ADD COLUMN default_tax_rates TEXT NOT NULL DEFAULT '[]';
