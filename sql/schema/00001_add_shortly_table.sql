-- +goose Up
-- +goose StatementBegin
CREATE TABLE shortly(
    id UUID PRIMARY KEY, 
    short_link TEXT NOT NULL,
    long_link TEXT NOT NULL,
    created_at TIMESTAMP  DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP  DEFAULT NOW() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shortly;
-- +goose StatementEnd
