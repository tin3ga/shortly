-- +goose Up
-- +goose StatementBegin
ALTER TABLE shortly
ADD CONSTRAINT unique_short_link UNIQUE (short_link);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE shortly
DROP CONSTRAINT unique_short_link;
-- +goose StatementEnd
