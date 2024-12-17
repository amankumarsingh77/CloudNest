CREATE TABLE IF NOT EXISTS user_quotas (
    user_id uuid PRIMARY KEY REFERENCES users(id),
    storage_used BIGINT DEFAULT 0,
    storage_limit BIGINT NOT NULL DEFAULT 10737418240,
    last_calculated_at TIMESTAMP With Time Zone DEFAULT current_timestamp,
    CONSTRAINT positive_storage  CHECK ( storage_used >=0 AND storage_used>=0 )
)