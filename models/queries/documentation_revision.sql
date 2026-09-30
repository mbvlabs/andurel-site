-- Fetches one DocumentationRevision by primary key.
-- name: GetDocumentationRevision :one
SELECT *
FROM documentation_revisions
WHERE id = $1
LIMIT 1;

-- Lists DocumentationRevision rows.
-- name: ListDocumentationRevisions :many
-- @order id
SELECT *
FROM documentation_revisions;

-- Counts DocumentationRevision rows.
-- name: CountDocumentationRevisions :one
SELECT count(*)
FROM documentation_revisions;

-- Inserts one DocumentationRevision and returns the row.
-- name: CreateDocumentationRevision :one
INSERT INTO documentation_revisions (
	created_at,
	page_id,
	status,
	title,
	meta_title,
	description,
	body_markdown,
	headings,
	created_by,
	published_at
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10
)
RETURNING *;

-- Updates one DocumentationRevision by primary key and returns the row.
-- name: UpdateDocumentationRevision :one
UPDATE documentation_revisions
SET
	page_id = $2,
	status = $3,
	title = $4,
	meta_title = $5,
	description = $6,
	body_markdown = $7,
	headings = $8,
	created_by = $9,
	published_at = $10
WHERE id = $1
RETURNING *;

-- Deletes one DocumentationRevision by primary key.
-- name: DeleteDocumentationRevision :exec
DELETE FROM documentation_revisions
WHERE id = $1;

-- Inserts or updates one DocumentationRevision and returns the row.
-- name: UpsertDocumentationRevision :one
INSERT INTO documentation_revisions (
	id,
	created_at,
	page_id,
	status,
	title,
	meta_title,
	description,
	body_markdown,
	headings,
	created_by,
	published_at
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10,
	$11
)
ON CONFLICT (id) DO UPDATE SET
	page_id = excluded.page_id,
	status = excluded.status,
	title = excluded.title,
	meta_title = excluded.meta_title,
	description = excluded.description,
	body_markdown = excluded.body_markdown,
	headings = excluded.headings,
	created_by = excluded.created_by,
	published_at = excluded.published_at
RETURNING *;

-- name: GetDocumentationRevisionDraftByPage :one
SELECT *
FROM documentation_revisions
WHERE page_id = $1 AND status = 'draft'
ORDER BY created_at DESC
LIMIT 1;

-- name: ListDocumentationRevisionsByPage :many
SELECT *
FROM documentation_revisions
WHERE page_id = $1
ORDER BY created_at DESC, id DESC;
