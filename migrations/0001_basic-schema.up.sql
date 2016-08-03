CREATE TABLE users (
    uuid TEXT PRIMARY KEY NOT NULL,
    email TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_email on users (email);

CREATE TABLE feeds (
    uuid TEXT PRIMARY KEY NOT NULL,
    message TEXT NOT NULL,
    username TEXT NOT NULL,
    pokemon TEXT NOT NULL,
    lat DECIMAL NOT NULL,
    long DECIMAL NOT NULL,
    geocodes TEXT NOT NULL,
    display_type TEXT,
    created_by_user_uuid TEXT NOT NULL REFERENCES users(uuid),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_feeds_created_by_user_uuid on feeds (created_by_user_uuid);
CREATE INDEX idx_feeds_lat_long on feeds (lat, long);
