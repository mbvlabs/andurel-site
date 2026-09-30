-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS documentation_versions (
    id bigserial PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    slug VARCHAR(64) NOT NULL UNIQUE,
    label VARCHAR(255) NOT NULL,
    is_latest BOOLEAN NOT NULL DEFAULT false,
    position INTEGER NOT NULL DEFAULT 0,
    published_nav_id BIGINT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS documentation_versions;
-- +goose StatementEnd
