
-- migrate up

-- sections
INSERT INTO "profile"."section" VALUES('1', 'home', 'Home', 'Home Section', true, NOW(), NOW());
INSERT INTO "profile"."section" VALUES('2', 'projects', 'Projects', 'Projects Section', true, NOW(), NOW());
INSERT INTO "profile"."section" VALUES('3', 'about', 'About', 'About Section', true, NOW(), NOW());

-- sections content
INSERT INTO "profile"."content" VALUES('1', 'dbr', '2', '{"title": "dbr"}', true, NOW(), NOW());
INSERT INTO "profile"."content" VALUES('2', 'web', '2', '{"title": "web"}', true, NOW(), NOW());
INSERT INTO "profile"."content" VALUES('3', 'validator', '2', '{"title": "validator"}', true, NOW(), NOW());



-- migrate down
DELETE FROM "profile"."section";
DELETE FROM "profile"."content";