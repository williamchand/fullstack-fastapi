-- name: CreatePayment :one
INSERT INTO payment (
    user_id, payment_method_id, provider, amount, currency, status, transaction_id, extra_metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetPaymentByTransaction :one
SELECT * FROM payment WHERE transaction_id = $1 LIMIT 1;

-- name: UpdatePaymentStatus :one
UPDATE payment
SET status = $2, amount = COALESCE($3, amount), currency = COALESCE($4, currency), extra_metadata = COALESCE($5, extra_metadata)
WHERE transaction_id = $1
RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM payment WHERE id = $1 LIMIT 1;
