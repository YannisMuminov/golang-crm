
CREATE TYPE task_status AS ENUM (
    'new',
    'in_progress',
    'done'
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status task_status NOT NULL DEFAULT 'new',
    deal_id INT NOT NULL REFERENCES deals(id) ON DELETE RESTRICT,
    assigned_to INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_by INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    due_data TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_deal_id ON tasks(deal_id);
CREATE INDEX idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX idx_tasks_status ON tasks(status);

INSERT INTO permissions (name, description) VALUES
    ('tasks:read', 'Просмотр задач'),
    ('tasks:write', 'Создание и редактирование задач'),
    ('tasks:delete', 'Удаление задач')
ON CONFLICT (name) DO NOTHING;


INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin'
    AND p.name IN ('tasks:read', 'tasks:write', 'tasks:delete')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'manager'
    AND p.name IN ('tasks:read', 'tasks:write')
ON CONFLICT DO NOTHING;

 
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'users'
  AND p.name IN ('tasks:read')
ON CONFLICT DO NOTHING;