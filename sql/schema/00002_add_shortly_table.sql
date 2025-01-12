-- +goose Up
-- +goose StatementBegin
CREATE TABLE shortly(
    id UUID PRIMARY KEY, 
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    short_link TEXT NOT NULL,
    long_link TEXT NOT NULL,
    click_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP  DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP  DEFAULT NOW() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shortly;
-- +goose StatementEnd
