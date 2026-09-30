-- name: GetToken :one
SELECT *
FROM tokens
WHERE id = $1
LIMIT 1;

-- name: GetTokenByScopeAndHash :one
SELECT *
FROM tokens
WHERE scope = $1
	AND hash = $2
LIMIT 1;

-- name: ListTokens :many
-- @order id
SELECT *
FROM tokens;

-- name: CountTokens :one
SELECT count(*)
FROM tokens;

-- name: CreateToken :one
INSERT INTO tokens (
	id,
	created_at,
	updated_at,
	scope,
	expires_at,
	hash,
	meta_data
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE id = $1;

-- name: UpdateToken :one
UPDATE tokens
SET
	updated_at = $2,
	scope = $3,
	expires_at = $4,
	hash = $5,
	meta_data = $6
WHERE id = $1
RETURNING *;

-- name: ListTokensByScope :many
SELECT *
FROM tokens
WHERE scope = $1
ORDER BY created_at DESC, id DESC;
