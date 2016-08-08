DROP TABLE IF EXISTS users CASCADE;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS feeds CASCADE;
DROP INDEX IF EXISTS idx_feeds_created_by_user_uuid;
DROP INDEX IF EXISTS idx_feeds_lat_long;
DROP TABLE IF EXISTS pokemon CASCADE;
DROP INDEX IF EXISTS idx_pokemon_name;
