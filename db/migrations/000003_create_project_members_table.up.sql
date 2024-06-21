CREATE TABLE IF NOT EXISTS project_members
(
    project_id UUID REFERENCES projects (id) ON DELETE CASCADE,
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    role       VARCHAR(50) NOT NULL,
    joined_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, user_id)
);