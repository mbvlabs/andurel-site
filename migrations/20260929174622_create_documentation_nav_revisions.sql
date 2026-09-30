-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS documentation_nav_revisions (
    id bigserial PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    version_id BIGINT NOT NULL REFERENCES documentation_versions (id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL,
    tree JSONB NOT NULL DEFAULT '[]',
    created_by UUID NOT NULL REFERENCES users (id),
    published_at TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS documentation_nav_revisions;
-- +goose StatementEnd
