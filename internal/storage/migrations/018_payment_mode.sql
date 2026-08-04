ALTER TABLE checkout_sessions ADD COLUMN client_reference_id TEXT NOT NULL DEFAULT '';
ALTER TABLE checkout_sessions ADD COLUMN payment_intent_data TEXT NOT NULL DEFAULT '{}';
