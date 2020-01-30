package services

const (
	DefaultURL = "http://localhost:8001"

	CREATE_MIGRATION_TABLES = `
CREATE SCHEMA IF NOT EXISTS "%s";


DO $$ 
BEGIN
    PERFORM ('"'||current_schema()||'"'||'.mode')::regtype;
EXCEPTION
    WHEN undefined_object THEN
        CREATE TYPE mode AS ENUM ('database', 'rabbitmq');
END $$;


CREATE TABLE IF NOT EXISTS migration (
  id_migration      TEXT NOT NULL,
  mode           	mode,
  "user"            TEXT DEFAULT user,
  executed_at       TIMESTAMP DEFAULT NOW(),
  CONSTRAINT migration_id_pkey PRIMARY KEY (id_migration, mode)
);
`
)
