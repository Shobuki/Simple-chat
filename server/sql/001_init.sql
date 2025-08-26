-- Jalankan sekali untuk membuat schema
CREATE TABLE IF NOT EXISTS users (
id BIGSERIAL PRIMARY KEY,
email TEXT UNIQUE NOT NULL,
display_name TEXT NOT NULL,
password_hash TEXT NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS messages (
id BIGSERIAL PRIMARY KEY,
user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
content TEXT NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at DESC);