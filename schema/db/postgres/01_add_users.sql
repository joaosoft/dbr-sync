
-- migrate up

-- the password is 'teste'
INSERT INTO "session"."user" VALUES(1, 'joao', 'ribeiro', 'joaosoft@gmail.com', '698dc19d489c4e4db73e28a713eab07b', '', true);

-- the password is 'teste1'
INSERT INTO "session"."user" VALUES(2, 'luis', 'ribeiro', 'luissoft@gmail.com', 'e959088c6049f1104c84c9bde5560a13', '', true);




-- migrate down
DELETE FROM "session"."user";