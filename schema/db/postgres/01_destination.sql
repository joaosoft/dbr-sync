
-- migrate up

-- schema
CREATE SCHEMA IF NOT EXISTS "dbr-sync-destination";

-- dbr-sync
CREATE TABLE "dbr-sync-destination"."example" (
	id_example          TEXT NOT NULL,
	"name"    		    TEXT NOT NULL,
	description			TEXT NOT NULL,
	"active"			BOOLEAN DEFAULT TRUE NOT NULL,
	created_at			TIMESTAMP DEFAULT NOW(),
	updated_at			TIMESTAMP DEFAULT NOW()
);



-- migrate down

-- tables
DROP TABLE "dbr-sync-destination"."example";

-- schema
DROP SCHEMA "dbr-sync-destination";
