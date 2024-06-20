CREATE TABLE IF NOT EXISTS projects
(
    id          UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    start_date  DATE,
    end_date    DATE,
    owner_id    UUID         NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_owner
        FOREIGN KEY (owner_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);