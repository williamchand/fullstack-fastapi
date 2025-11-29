-- name: CreateWedding :one
INSERT INTO wedding (user_id, template_id, payment_id, status, custom_domain, slug, config_data)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetWeddingByID :one
SELECT * FROM wedding WHERE id = $1 LIMIT 1;

-- name: GetWeddingsByUser :many
SELECT * FROM wedding WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: UpdateWeddingConfig :one
UPDATE wedding SET config_data = $2, updated_at = now() WHERE id = $1 RETURNING *;

-- name: SetWeddingTemplate :one
UPDATE wedding SET template_id = $2, updated_at = now() WHERE id = $1 RETURNING *;

-- name: SetWeddingPayment :one
UPDATE wedding SET payment_id = $2, updated_at = now() WHERE id = $1 RETURNING *;

-- name: SetWeddingDomain :one
UPDATE wedding SET custom_domain = $2, updated_at = now() WHERE id = $1 RETURNING *;

-- name: SetWeddingSlug :one
UPDATE wedding SET slug = $2, updated_at = now() WHERE id = $1 RETURNING *;

-- name: PublishWedding :one
UPDATE wedding SET status = 'active', updated_at = now() WHERE id = $1 RETURNING *;
