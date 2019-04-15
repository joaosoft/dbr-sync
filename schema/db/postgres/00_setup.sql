
-- migrate up

-- schema
CREATE SCHEMA IF NOT EXISTS "profile";

-- functions
CREATE OR REPLACE FUNCTION "profile".function_updated_at()
  RETURNS TRIGGER AS $$
  BEGIN
   NEW.updated_at = now();
   RETURN NEW;
  END;
  $$ LANGUAGE 'plpgsql';

-- sections
CREATE TABLE "profile"."section" (
	id_section 		    TEXT PRIMARY KEY,
	"key"               TEXT NOT NULL UNIQUE,
	"name"    		    TEXT NOT NULL,
	description			TEXT NOT NULL,
	position            INTEGER NOT NULL UNIQUE,
	"active"			BOOLEAN DEFAULT TRUE NOT NULL,
	created_at			TIMESTAMP DEFAULT NOW(),
	updated_at			TIMESTAMP DEFAULT NOW()
);


-- content type
CREATE TABLE "profile"."content_type" (
    id_content_type 		    TEXT PRIMARY KEY,
    "key"                       TEXT NOT NULL UNIQUE,
	"name"                      TEXT NOT NULL,
	"active"			        BOOLEAN DEFAULT TRUE NOT NULL,
	created_at			        TIMESTAMP DEFAULT NOW(),
	updated_at			        TIMESTAMP DEFAULT NOW()
);

-- section contents
CREATE TABLE "profile"."content" (
    id_content 		            TEXT PRIMARY KEY,
    "key"                       TEXT NOT NULL UNIQUE,
	fk_section                  TEXT NOT NULL REFERENCES "profile"."section" (id_section),
	fk_content_type             TEXT NOT NULL REFERENCES "profile"."content_type" (id_content_type),
	"content"                   JSONB NOT NULL,
    position                    INTEGER NOT NULL,
	"active"			        BOOLEAN DEFAULT TRUE NOT NULL,
	created_at			        TIMESTAMP DEFAULT NOW(),
	updated_at			        TIMESTAMP DEFAULT NOW(),
	UNIQUE(fk_section, "position")
);

-- triggers
CREATE TRIGGER trigger_section_updated_at BEFORE UPDATE
  ON "profile"."section" FOR EACH ROW EXECUTE PROCEDURE "profile".function_updated_at();

CREATE TRIGGER trigger_content_updated_at BEFORE UPDATE
  ON "profile"."content" FOR EACH ROW EXECUTE PROCEDURE "profile".function_updated_at();






-- migrate down

-- triggers
DROP TRIGGER trigger_section_updated_at ON profile."section";
DROP TRIGGER trigger_content_updated_at ON profile."content";

-- tables
DROP TABLE "profile"."section";
DROP TABLE "profile"."content";

-- functions
DROP FUNCTION "profile".function_updated_at;

-- schema
DROP SCHEMA "profile";
