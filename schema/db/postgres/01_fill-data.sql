
-- migrate up

-- sections
INSERT INTO "profile"."section" (id_section, "key", "name", description, position, "active", created_at, updated_at)
VALUES('1', 'home', 'Home', 'Home Section', 1, true, NOW(), NOW());

INSERT INTO "profile"."section" (id_section, "key", "name", description, position, "active", created_at, updated_at)
VALUES('2', 'projects', 'Projects', 'Projects Section', 2, true, NOW(), NOW());

INSERT INTO "profile"."section" (id_section, "key", "name", description, position, "active", created_at, updated_at)
VALUES('3', 'about', 'About', 'About Section', 3, true, NOW(), NOW());


-- sections content
INSERT INTO "profile"."content" (id_content, "key", fk_section, "content", position, "active", created_at, updated_at)
VALUES('1', 'dbr', '2', '{"title": "Dbr", "description": "A simple database client with support for master/slave databases.", "url": "https://github.com/joaosoft/dbr", "build": "https://travis-ci.org/joaosoft/dbr.svg?branch=master"}', 1, true, NOW(), NOW());

INSERT INTO "profile"."content" (id_content, "key", fk_section, "content", position, "active", created_at, updated_at)
VALUES('2', 'web', '2', '{"title": "Web", "description": "A simple and fast web server and client.", "url": "https://github.com/joaosoft/web", "build": "https://travis-ci.org/joaosoft/web.svg?branch=master"}', 2, true, NOW(), NOW());

INSERT INTO "profile"."content" (id_content, "key", fk_section, "content", position, "active", created_at, updated_at)
VALUES('3', 'validator', '2', '{"title": "Validator", "description": "A simple struct validator by tags.", "url": "https://github.com/joaosoft/validator", "build": "https://travis-ci.org/joaosoft/validator.svg?branch=master"}', 3, true, NOW(), NOW());



-- migrate down
DELETE FROM "profile"."section";
DELETE FROM "profile"."content";