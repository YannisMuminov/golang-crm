CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(100),
    company VARCHAR(50),
    notes TEXT,
    created_by INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT  NULL DEFAULT NOW()
);

CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_created_by ON clients(created_by);