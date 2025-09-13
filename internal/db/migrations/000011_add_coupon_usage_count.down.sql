-- Remove usage_count column from coupons table
ALTER TABLE coupons
DROP COLUMN IF EXISTS usage_count;
