ALTER TABLE checkout_sessions ADD COLUMN allow_promotion_codes INTEGER NOT NULL DEFAULT 0;
ALTER TABLE checkout_sessions ADD COLUMN trial_period_days INTEGER NOT NULL DEFAULT 0;
