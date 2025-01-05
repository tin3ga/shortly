-- +goose Up
-- +goose StatementBegin
INSERT INTO "shortly" ("id", "short_link", "long_link") VALUES ('45348b11-6984-4b73-96e2-5722c445ba07', 'test', 'www.test.com');
INSERT INTO "shortly" ("id", "short_link", "long_link") VALUES ('6309f7cc-aba7-49b1-9d9d-c7a4fafd6238', 'google', 'www.google.com');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "shortly" WHERE "id" IN ('45348b11-6984-4b73-96e2-5722c445ba07', '6309f7cc-aba7-49b1-9d9d-c7a4fafd6238');
-- +goose StatementEnd


