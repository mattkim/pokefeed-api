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

CREATE TABLE feed_items (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    message TEXT NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL,
    formatted_address TEXT NOT NULL,
    created_by_user_uuid VARCHAR(255) NOT NULL REFERENCES users(uuid),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_feed_items_created_by_user_uuid on feed_items (created_by_user_uuid);
CREATE INDEX idx_feed_items_lat_long on feed_items (lat, long);

CREATE TABLE feed_tags (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    type VARCHAR(255) NOT NULL,
    name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    image_url TEXT NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

CREATE INDEX idx_feed_tags_type_name on feed_tags (type, name);

CREATE TABLE feed_items_feed_tags (
    feed_item_uuid VARCHAR(255) NOT NULL REFERENCES feed_items(uuid),
    feed_tag_uuid VARCHAR(255) NOT NULL REFERENCES feed_tags(uuid),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

CREATE TABLE comments (
    uuid VARCHAR(255) PRIMARY KEY NOT NULL,
    feed_item_uuid VARCHAR(255) NOT NULL REFERENCES feed_items(uuid),
    message TEXT NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL,
    formatted_address TEXT NOT NULL,
    created_by_user_uuid VARCHAR(255) NOT NULL REFERENCES users(uuid),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);
