-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS documentation_revisions (
    id bigserial PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    page_id BIGINT NOT NULL REFERENCES documentation_pages (id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    body_markdown TEXT NOT NULL DEFAULT '',
    headings JSONB NOT NULL DEFAULT '[]',
    created_by UUID NOT NULL REFERENCES users (id),
    published_at TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS documentation_revisions;
-- +goose StatementEnd
