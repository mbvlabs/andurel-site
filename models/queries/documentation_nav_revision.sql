-- Fetches one DocumentationNavRevision by primary key.
-- name: GetDocumentationNavRevision :one
SELECT *
FROM documentation_nav_revisions
WHERE id = $1
LIMIT 1;

-- Lists DocumentationNavRevision rows.
-- name: ListDocumentationNavRevisions :many
-- @order id
SELECT *
FROM documentation_nav_revisions;

-- Counts DocumentationNavRevision rows.
-- name: CountDocumentationNavRevisions :one
SELECT count(*)
FROM documentation_nav_revisions;

-- Inserts one DocumentationNavRevision and returns the row.
-- name: CreateDocumentationNavRevision :one
INSERT INTO documentation_nav_revisions (
	created_at,
	version_id,
	status,
	tree,
	created_by,
	published_at
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
RETURNING *;

-- Updates one DocumentationNavRevision by primary key and returns the row.
-- name: UpdateDocumentationNavRevision :one
UPDATE documentation_nav_revisions
SET
	version_id = $2,
	status = $3,
	tree = $4,
	created_by = $5,
	published_at = $6
WHERE id = $1
RETURNING *;

-- Deletes one DocumentationNavRevision by primary key.
-- name: DeleteDocumentationNavRevision :exec
DELETE FROM documentation_nav_revisions
WHERE id = $1;

-- Inserts or updates one DocumentationNavRevision and returns the row.
-- name: UpsertDocumentationNavRevision :one
INSERT INTO documentation_nav_revisions (
	id,
	created_at,
	version_id,
	status,
	tree,
	created_by,
	published_at
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7
)
ON CONFLICT (id) DO UPDATE SET
	version_id = excluded.version_id,
	status = excluded.status,
	tree = excluded.tree,
	created_by = excluded.created_by,
	published_at = excluded.published_at
RETURNING *;

-- name: GetDocumentationNavRevisionDraftByVersion :one
SELECT *
FROM documentation_nav_revisions
WHERE version_id = $1 AND status = 'draft'
ORDER BY created_at DESC
LIMIT 1;

-- name: ListDocumentationNavRevisionsByVersion :many
SELECT *
FROM documentation_nav_revisions
WHERE version_id = $1
ORDER BY created_at DESC, id DESC;
