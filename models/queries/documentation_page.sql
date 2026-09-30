-- Fetches one DocumentationPage by primary key.
-- name: GetDocumentationPage :one
SELECT *
FROM documentation_pages
WHERE id = $1
LIMIT 1;

-- Lists DocumentationPage rows.
-- name: ListDocumentationPages :many
-- @order id
SELECT *
FROM documentation_pages;

-- Counts DocumentationPage rows.
-- name: CountDocumentationPages :one
SELECT count(*)
FROM documentation_pages;

-- Inserts one DocumentationPage and returns the row.
-- name: CreateDocumentationPage :one
INSERT INTO documentation_pages (
	created_at,
	updated_at,
	version_id,
	slug,
	published_revision_id
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
RETURNING *;

-- Updates one DocumentationPage by primary key and returns the row.
-- name: UpdateDocumentationPage :one
UPDATE documentation_pages
SET
	updated_at = $2,
	version_id = $3,
	slug = $4,
	published_revision_id = $5
WHERE id = $1
RETURNING *;

-- Deletes one DocumentationPage by primary key.
-- name: DeleteDocumentationPage :exec
DELETE FROM documentation_pages
WHERE id = $1;

-- Inserts or updates one DocumentationPage and returns the row.
-- name: UpsertDocumentationPage :one
INSERT INTO documentation_pages (
	id,
	created_at,
	updated_at,
	version_id,
	slug,
	published_revision_id
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
ON CONFLICT (id) DO UPDATE SET
	updated_at = excluded.updated_at,
	version_id = excluded.version_id,
	slug = excluded.slug,
	published_revision_id = excluded.published_revision_id
RETURNING *;

-- name: ListDocumentationPagesByVersion :many
SELECT *
FROM documentation_pages
WHERE version_id = $1
ORDER BY id ASC;

-- name: GetDocumentationPageByVersionSlug :one
SELECT *
FROM documentation_pages
WHERE version_id = $1 AND slug = $2
LIMIT 1;
