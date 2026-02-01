-- 000002_add_indexes.down.sql

-- Drop indexes for all tables
DROP INDEX IF EXISTS idx_clients_email;
DROP INDEX IF EXISTS idx_clients_country;

DROP INDEX IF EXISTS idx_banks_type;

DROP INDEX IF EXISTS idx_credits_client_id;
DROP INDEX IF EXISTS idx_credits_bank_id;
DROP INDEX IF EXISTS idx_credits_status;
DROP INDEX IF EXISTS idx_credits_created_at;