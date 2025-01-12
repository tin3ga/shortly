-- +goose Up
-- +goose StatementBegin

INSERT INTO "users"
("id", "username","email","password_hash") 
VALUES ('00008b11-6984-4b73-96e2-5722c445ba07', 'test1', 'test1@example.com', 'hash_example1');

INSERT INTO "users"
("id", "username","email","password_hash")
VALUES ('0000f7cc-aba7-49b1-9d9d-c7a4fafd6238', 'test2', 'test2@example.com', 'hash_example2');

INSERT INTO "shortly" 
("id", "user_id","short_link", "long_link", "click_count")
VALUES ('45348b11-6984-4b73-96e2-5722c445ba07', '00008b11-6984-4b73-96e2-5722c445ba07', 'test', 'https://www.test.com', '0');

INSERT INTO "shortly" 
("id", "user_id","short_link", "long_link", "click_count") 
VALUES ('6309f7cc-aba7-49b1-9d9d-c7a4fafd6238', '0000f7cc-aba7-49b1-9d9d-c7a4fafd6238', 'google', 'https://www.google.com', '0');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "shortly" WHERE "id" IN ('45348b11-6984-4b73-96e2-5722c445ba07', '6309f7cc-aba7-49b1-9d9d-c7a4fafd6238');
DELETE FROM "users" WHERE "id" IN('00008b11-6984-4b73-96e2-5722c445ba07', '0000f7cc-aba7-49b1-9d9d-c7a4fafd6238');
-- +goose StatementEnd


