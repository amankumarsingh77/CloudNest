CREATE TABLE IF NOT EXISTS files(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) NOT NULL ,
    folder_id uuid REFERENCES folders(id),
    remote_file_name varchar(255) NOT NULL ,
    mime_type varchar(100) NOT NULL ,
    size bigint NOT NULL ,
    checksum varchar(64),
    encryption_key TEXT,
    url TEXT,
    created_by uuid NOT NULL REFERENCES users(id),
    path Text,
    is_deleted boolean DEFAULT FALSE,
    deleted_at timestamp WITH TIME ZONE,
    permanent_deletion_at timestamp WITH TIME ZONE,
    last_accessed_at timestamp WITH TIME ZONE,
    created_at timestamp WITH TIME ZONE DEFAULT current_timestamp,
    updated_at timestamp With Time Zone DEFAULT current_timestamp
)