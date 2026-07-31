-- 000029: Revert - Add missing columns to expenses table
ALTER TABLE expenses DROP COLUMN IF EXISTS billable_client;
ALTER TABLE expenses DROP COLUMN IF EXISTS is_billable;
ALTER TABLE expenses DROP COLUMN IF EXISTS receipt_required;
ALTER TABLE expenses DROP COLUMN IF EXISTS total_amount;
ALTER TABLE expenses DROP COLUMN IF EXISTS tax_amount;
