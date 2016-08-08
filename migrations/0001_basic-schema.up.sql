CREATE TABLE users (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

CREATE UNIQUE INDEX idx_users_email on users (email);

CREATE TABLE pokemon (
    id INT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    image_url TEXT NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

CREATE UNIQUE INDEX idx_pokemon_name on pokemon (name);

CREATE TABLE feeds (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    message TEXT NOT NULL,
    pokemon_name TEXT NOT NULL REFERENCES pokemon(name),
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL,
    formatted_address TEXT NOT NULL,
    created_by_user_uuid VARCHAR(255) NOT NULL REFERENCES users(uuid),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_feeds_created_by_user_uuid on feeds (created_by_user_uuid);
CREATE INDEX idx_feeds_lat_long on feeds (lat, long);
