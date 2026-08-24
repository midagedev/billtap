-- Credit notes persist Stripe's memo / out_of_band_amount / refund_amount so
-- create and retrieve echo the same allocation. Amount still stores the total.
ALTER TABLE credit_notes ADD COLUMN memo TEXT NOT NULL DEFAULT '';
ALTER TABLE credit_notes ADD COLUMN out_of_band_amount INTEGER NOT NULL DEFAULT 0;
ALTER TABLE credit_notes ADD COLUMN refund_amount INTEGER NOT NULL DEFAULT 0;
