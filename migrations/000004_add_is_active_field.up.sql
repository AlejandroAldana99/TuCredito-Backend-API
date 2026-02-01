-- 000004_add_is_active_field.up.sql

-- Add is_active field to clients table
ALTER TABLE clients ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- Add is_active field to banks table
ALTER TABLE banks ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- Add is_active field to credits table
ALTER TABLE credits ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;