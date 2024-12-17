CREATE TABLE IF NOT EXISTS folders (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) NOT NULL ,
    parent_folder_id uuid REFERENCES folders(id),
    created_by uuid NOT NULL REFERENCES users(id),
    path text NOT NULL ,
    color varchar(20),
    description TEXT,
    is_deleted boolean DEFAULT FALSE,
    deleted_at timestamp WITH TIME ZONE,
    permanent_deletion_at timestamp WITH TIME ZONE,
    created_at timestamp WITH TIME ZONE DEFAULT current_timestamp,
    updated_at timestamp WITH TIME ZONE DEFAULT current_timestamp
)