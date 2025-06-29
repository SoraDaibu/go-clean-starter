-- Drop users table
DROP TRIGGER IF EXISTS update_updated_at_trigger_users ON users;
DROP TABLE IF EXISTS users;

-- Drop items table first (since it references item_types)
DROP TRIGGER IF EXISTS update_updated_at_trigger_items ON items;
DROP TABLE IF EXISTS items;

-- Drop item_types table last
DROP TRIGGER IF EXISTS update_updated_at_trigger_item_types ON item_types;
DROP TABLE IF EXISTS item_types;
