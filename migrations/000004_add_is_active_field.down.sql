-- 000004_add_is_active_field.down.sql

-- Remove is_active field from clients table
ALTER TABLE clients DROP COLUMN is_active;

-- Remove is_active field from banks table
ALTER TABLE banks DROP COLUMN is_active;

-- Remove is_active field from credits table
ALTER TABLE credits DROP COLUMN is_active;