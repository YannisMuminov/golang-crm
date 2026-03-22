CREATE TYPE deal_status AS ENUM (
    'new',
    'negotiation',
    'won',
    'lost'
    );

CREATE TABLE IF NOT EXISTS deals (
id SERIAL PRIMARY KEY,
title VARCHAR(255) NOT NULL,
amount NUMERIC(12, 2) NOT NULL,
status deal_status NOT NULL DEFAULT 'new',
client_id INT NOT NULL REFERENCES clients(id) ON DELETE RESTRICT, assigned_to INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
created_by INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
closed_at TIMESTAMPTZ,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_deals_client_id ON deals(client_id);
CREATE INDEX idx_deals_assigned_to ON deals(assigned_to);
CREATE INDEX idx_deals_status ON deals(status);