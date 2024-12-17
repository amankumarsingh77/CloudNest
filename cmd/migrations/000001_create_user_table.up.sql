CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(100) NOT NULL ,
    username varchar(50) NOT NULL ,
    email citext UNIQUE NOT NULL ,
    password bytea NOT NULL ,
    status varchar(20) DEFAULT 'active' CHECK ( status IN ('active','inactive','suspended')),
    email_verified boolean DEFAULT false,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
)