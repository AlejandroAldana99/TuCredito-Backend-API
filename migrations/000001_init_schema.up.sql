-- 000001_init_schema.up.sql

-- Tu Credito: Credit Decision & Management Service - Initial Schema

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Clients:
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    birth_date DATE NOT NULL,
    country VARCHAR(2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Banks:
CREATE TABLE banks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('PRIVATE', 'GOVERNMENT')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Credits:
CREATE TABLE credits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    bank_id UUID NOT NULL REFERENCES banks(id) ON DELETE RESTRICT,
    min_payment DECIMAL(15, 2) NOT NULL,
    max_payment DECIMAL(15, 2) NOT NULL,
    term_months INT NOT NULL CHECK (term_months > 0),
    credit_type VARCHAR(20) NOT NULL CHECK (credit_type IN ('AUTO', 'MORTGAGE', 'COMMERCIAL')),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);