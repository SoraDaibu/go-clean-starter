-- Initialize schema
CREATE TABLE IF NOT EXISTS schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL,
    PRIMARY KEY (version)
);

-- Function to automatically update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
