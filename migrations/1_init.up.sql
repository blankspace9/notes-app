CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pass_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    note TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes (user_id);