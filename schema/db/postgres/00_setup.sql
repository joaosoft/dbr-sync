
-- migrate up

-- schema
CREATE SCHEMA IF NOT EXISTS "dbr-sync";

-- dbr-sync
CREATE TABLE "dbr-sync"."example" (
	id_example          TEXT NOT NULL,
	"name"    		    TEXT NOT NULL,
	description			TEXT NOT NULL,
	"active"			BOOLEAN DEFAULT TRUE NOT NULL,
	created_at			TIMESTAMP DEFAULT NOW(),
	updated_at			TIMESTAMP DEFAULT NOW()
);



-- migrate down

-- tables
DROP TABLE "dbr-sync"."example";

-- schema
DROP SCHEMA "dbr-sync";
