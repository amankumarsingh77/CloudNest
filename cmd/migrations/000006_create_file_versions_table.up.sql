CREATE TABLE IF NOT EXISTS file_versions(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id uuid NOT NULL REFERENCES files(id),
    version_number INTEGER NOT NULL ,
    remote_file_name VARCHAR(255) NOT NULL ,
    size BIGINT NOT NULL ,
    checksum VARCHAR(64),
    created_at timestamp WITH TIME ZONE DEFAULT current_timestamp,
    CONSTRAINT positive_version CHECK ( version_number > 0 ),
    UNIQUE(file_id, version_number)
)