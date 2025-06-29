-- users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE users IS 'This table stores user information';

CREATE TRIGGER update_updated_at_trigger_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- item types (moved before items to resolve foreign key dependency)
CREATE TABLE item_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_item_type UNIQUE (name)
);

COMMENT ON TABLE item_types IS 'This table is used to store item types';

CREATE TRIGGER update_updated_at_trigger_item_types
BEFORE UPDATE ON item_types
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- items (now references item_types which exists)
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type_id INTEGER REFERENCES item_types(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE items IS 'This table is used to store master items';

CREATE TRIGGER update_updated_at_trigger_items
BEFORE UPDATE ON items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
