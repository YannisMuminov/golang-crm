 CREATE TABLE IF NOT EXISTS roles (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50)  NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permissions (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       INT NOT NULL REFERENCES roles(id)       ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);
 

CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    role_id       INT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
 

CREATE INDEX idx_users_email   ON users(email);
CREATE INDEX idx_users_role_id ON users(role_id);


INSERT INTO roles (name, description) VALUES
    ('admin', 'Полный доступ'),
    ('manager', 'Управление клиентами и сделками'),
    ('users', 'Только чтение')
ON CONFLICT (name) DO NOTHING;


INSERT INTO permissions (name, description) VALUES
    ('users:read',    'Просмотр пользователей'),
    ('users:write',   'Создание и редактирование'),
    ('users:delete',  'Удаление пользователей'),
    ('clients:read',  'Просмотр клиентов'),
    ('clients:write', 'Создание и редактирование клиентов'),
    ('clients:delete','Удаление клиентов'),
    ('deals:read',    'Просмотр сделок'),
    ('deals:write',   'Создание и редактирование сделок'),
    ('deals:delete',  'Удаление сделок')
ON CONFLICT (name) DO NOTHING;


INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'manager'
    AND p.name IN (
        'users:read',
        'clients:read', 'clients:write',
        'deals:read', 'deals:write'
    )
ON CONFLICT DO NOTHING;


INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'users'
    AND p.name IN('users:read', 'clients:read', 'deals:read')
 ON CONFLICT DO NOTHING;