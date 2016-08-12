DROP TABLE IF EXISTS users CASCADE;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS pokemon CASCADE;
DROP INDEX IF EXISTS idx_pokemon_name;

DROP TABLE IF EXISTS feed_items CASCADE;
DROP INDEX IF EXISTS idx_feed_items_created_by_user_uuid;
DROP INDEX IF EXISTS idx_feed_items_lat_long;
DROP TABLE IF EXISTS feed_tags;
DROP INDEX IF EXISTS idx_feed_tags_type_name;
DROP TABLE IF EXISTS feed_items_feed_tags;
