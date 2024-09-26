-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1
LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;

-- name: ListTransfersByFromAccount :many
SELECT * FROM transfers
WHERE from_account_id = $1
ORDER BY created_at DESC;

-- name: ListTransfersByToAccount :many
SELECT * FROM transfers
WHERE to_account_id = $1
ORDER BY created_at DESC;

-- name: GetTransfersAfter :many
SELECT * FROM transfers
WHERE created_at > $1
ORDER BY created_at DESC;