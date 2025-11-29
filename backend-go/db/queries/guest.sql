-- name: AddGuest :one
INSERT INTO guest (wedding_id, name, contact, rsvp_status, message)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateGuest :one
UPDATE guest SET name = $2, contact = $3, rsvp_status = $4, message = $5, deleted_at = NULL
WHERE id = $1 RETURNING *;

-- name: DeleteGuest :exec
UPDATE guest SET deleted_at = now() WHERE id = $1;

-- name: ListGuestsByWedding :many
SELECT * FROM guest WHERE wedding_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC;
