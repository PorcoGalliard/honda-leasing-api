-- Rollback: hapus constraint dan kolom user_id
ALTER TABLE dealer.customers
DROP CONSTRAINT IF EXISTS fk_customer_user;

ALTER TABLE dealer.customers
DROP COLUMN IF EXISTS user_id;
