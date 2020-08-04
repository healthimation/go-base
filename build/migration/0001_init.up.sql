CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE goal (
    id           BIGSERIAL PRIMARY KEY,
    user_id      UUID NOT NULL,    
    type         TEXT NOT NULL,
    category     TEXT NOT NULL,
    value        BIGINT NOT NULL,
    name         TEXT,
    description  TEXT,
    
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP,
    
    starts_at    TIMESTAMP NOT NULL,
    expires_at   TIMESTAMP NOT NULL,

    completed_at TIMESTAMP,
    skipped_at   TIMESTAMP
);

CREATE INDEX user_type_category_idx ON goal(user_id, type, category);