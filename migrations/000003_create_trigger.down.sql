-- 000003_create_trigger.down.sql

-- Drop trigger and function
DROP TRIGGER IF EXISTS update_credits_updated_at ON credits;
DROP FUNCTION IF EXISTS update_updated_at_column();