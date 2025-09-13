-- Add usage_count column to coupons table for tracking coupon usage
ALTER TABLE coupons
ADD COLUMN IF NOT EXISTS usage_count INT NOT NULL DEFAULT 0;
