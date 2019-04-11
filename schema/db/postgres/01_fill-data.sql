
-- migrate up

-- sections
INSERT INTO "profile"."section" (id_section, "key", "name", description, position, "active", created_at, updated_at)
VALUES('1', 'home', 'Hello', 'Home Section', 1, true, NOW(), NOW());

INSERT INTO "profile"."section" (id_section, "key", "name", description, position, "active", created_at, updated_at)
VALUES('2', 'projects', 'Projects', 'Projects Section', 2, true, NOW(), NOW());

INSERT INTO "profile"."section" (id_section, "key", "name", description, position, "active", created_at, updated_at)
VALUES('3', 'about', 'Goodbye', 'About Section', 3, true, NOW(), NOW());


-- sections content
-- :: section home
INSERT INTO "profile"."content" (id_content, "key", fk_section, "content", position, "active", created_at, updated_at)
VALUES('1', 'hello', '1', '{"title": "I''m João Ribeiro", "description": "I like to code.", "url": "https://www.facebook.com/joaosoft"}', 1, true, NOW(), NOW());

-- :: section projects
INSERT INTO "profile"."content" (id_content, "key", fk_section, "content", position, "active", created_at, updated_at)
VALUES('2', 'dbr', '2', '{"title": "Dbr", "description": "A simple database client with support for master/slave databases.", "url": "https://github.com/joaosoft/dbr", "build": "https://travis-ci.org/joaosoft/dbr.svg?branch=master"}', 1, true, NOW(), NOW());

INSERT INTO "profile"."content" (id_content, "key", fk_section, "content", position, "active", created_at, updated_at)
VALUES('3', 'web', '2', '{"title": "Web", "description": "A simple and fast web server and client.", "url": "https://github.com/joaosoft/web", "build": "https://travis-ci.org/joaosoft/web.svg?branch=master"}', 2, true, NOW(), NOW());

INSERT INTO "profile"."content" (id_content, "key", fk_section, "content", position, "active", created_at, updated_at)
VALUES('4', 'validator', '2', '{"title": "Validator", "description": "A simple struct validator by tags.", "url": "https://github.com/joaosoft/validator", "build": "https://travis-ci.org/joaosoft/validator.svg?branch=master"}', 3, true, NOW(), NOW());

-- :: section about
INSERT INTO "profile"."content" (id_content, "key", fk_section, "content", position, "active", created_at, updated_at)
VALUES('5', 'goodbuye', '3', '{"title": "Thanks for reading", "description": "Find more about me...", "url": "https://www.facebook.com/joaosoft"}', 1, true, NOW(), NOW());


-- migrate down
DELETE FROM "profile"."section";
DELETE FROM "profile"."content";