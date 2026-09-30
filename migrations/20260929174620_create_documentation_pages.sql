-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS documentation_pages (
    id bigserial PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    version_id BIGINT NOT NULL REFERENCES documentation_versions (id) ON DELETE CASCADE,
    slug VARCHAR(255) NOT NULL,
    published_revision_id BIGINT,
    UNIQUE (version_id, slug)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS documentation_pages;
-- +goose StatementEnd
