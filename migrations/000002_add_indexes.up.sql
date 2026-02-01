-- 000002_add_indexes.up.sql

-- Create indexes for all tables
CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_country ON clients(country);

CREATE INDEX idx_banks_type ON banks(type);

CREATE INDEX idx_credits_client_id ON credits(client_id);
CREATE INDEX idx_credits_bank_id ON credits(bank_id);
CREATE INDEX idx_credits_status ON credits(status);
CREATE INDEX idx_credits_created_at ON credits(created_at);