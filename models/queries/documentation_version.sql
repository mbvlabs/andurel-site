-- Fetches one DocumentationVersion by primary key.
-- name: GetDocumentationVersion :one
SELECT *
FROM documentation_versions
WHERE id = $1
LIMIT 1;

-- Lists DocumentationVersion rows.
-- name: ListDocumentationVersions :many
-- @order id
SELECT *
FROM documentation_versions;

-- Counts DocumentationVersion rows.
-- name: CountDocumentationVersions :one
SELECT count(*)
FROM documentation_versions;

-- Inserts one DocumentationVersion and returns the row.
-- name: CreateDocumentationVersion :one
INSERT INTO documentation_versions (
	created_at,
	updated_at,
	slug,
	label,
	is_latest,
	position,
	published_nav_id
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7
)
RETURNING *;

-- Updates one DocumentationVersion by primary key and returns the row.
-- name: UpdateDocumentationVersion :one
UPDATE documentation_versions
SET
	updated_at = $2,
	slug = $3,
	label = $4,
	is_latest = $5,
	position = $6,
	published_nav_id = $7
WHERE id = $1
RETURNING *;

-- Deletes one DocumentationVersion by primary key.
-- name: DeleteDocumentationVersion :exec
DELETE FROM documentation_versions
WHERE id = $1;

-- Inserts or updates one DocumentationVersion and returns the row.
-- name: UpsertDocumentationVersion :one
INSERT INTO documentation_versions (
	id,
	created_at,
	updated_at,
	slug,
	label,
	is_latest,
	position,
	published_nav_id
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8
)
ON CONFLICT (id) DO UPDATE SET
	updated_at = excluded.updated_at,
	slug = excluded.slug,
	label = excluded.label,
	is_latest = excluded.is_latest,
	position = excluded.position,
	published_nav_id = excluded.published_nav_id
RETURNING *;

-- name: GetDocumentationVersionBySlug :one
SELECT *
FROM documentation_versions
WHERE slug = $1
LIMIT 1;

-- Lists versions by sidebar position.
-- name: ListDocumentationVersionsByPosition :many
SELECT *
FROM documentation_versions
ORDER BY position ASC, id ASC;

-- Clears the latest flag on every version.
-- name: ClearDocumentationVersionLatest :exec
UPDATE documentation_versions
SET is_latest = false, updated_at = $1
WHERE is_latest = true;
