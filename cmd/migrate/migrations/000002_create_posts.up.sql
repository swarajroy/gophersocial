CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    title text,
    content text,
    tags varchar(100) [],
    user_id bigint,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);