CREATE TABLE IF NOT EXISTS quotes (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    last_modified timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    content text NOT NULL,
    author text NOT NULL,
    source_title text NOT NULL,
    source_type text NOT NULL, 
    tags text[] NOT NULL, 
    version integer NOT NULL DEFAULT 1
);