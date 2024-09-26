-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1
LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListEntriesByAccount :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY created_at DESC;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;

-- name: GetEntriesAfter :many
SELECT * FROM entries
WHERE created_at > $1
ORDER BY created_at DESC;

-- name: ListEntriesByAmount :many
SELECT * FROM entries
WHERE amount = $1
ORDER BY created_at DESC;

-- name: ListEntriesBetweenDates :many
SELECT * FROM entries
WHERE created_at BETWEEN $1 AND $2
ORDER BY created_at DESC;