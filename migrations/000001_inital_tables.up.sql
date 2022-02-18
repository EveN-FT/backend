CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    address TEXT NOT NULL UNIQUE,
    owner_address TEXT NOT NULL
);
