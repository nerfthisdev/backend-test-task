CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_guid UUID NOT NULL REFERENCES users(guid) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
