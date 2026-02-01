-- 000003_create_trigger.up.sql

-- Trigger to update updated_at on credits
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_credits_updated_at
    BEFORE UPDATE ON credits
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();
