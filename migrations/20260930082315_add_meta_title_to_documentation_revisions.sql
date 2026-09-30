-- +goose Up
-- +goose StatementBegin
ALTER TABLE documentation_revisions
    ADD COLUMN meta_title VARCHAR(255) NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documentation_revisions
    DROP COLUMN IF EXISTS meta_title;
-- +goose StatementEnd
