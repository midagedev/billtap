ALTER TABLE payment_intents ADD COLUMN failure_decline_code TEXT NOT NULL DEFAULT '';
ALTER TABLE payment_intents ADD COLUMN payment_method_id TEXT NOT NULL DEFAULT '';
