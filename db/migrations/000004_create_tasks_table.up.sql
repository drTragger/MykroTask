CREATE TABLE IF NOT EXISTS tasks
(
    id          UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      VARCHAR(50)  NOT NULL,
    priority    VARCHAR(50)  NOT NULL,
    assignee    UUID REFERENCES users (id),
    due_date    DATE,
    project_id  UUID REFERENCES projects (id) ON DELETE CASCADE,
    created_by  UUID REFERENCES users (id),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);