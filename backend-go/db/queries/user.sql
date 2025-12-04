-- name: GetUserByID :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1 LIMIT 1;

-- name: GetUserByPhone :one
SELECT * FROM "user"
WHERE phone_number = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO "user" (
    email, phone_number, full_name, hashed_password, is_active, is_email_verified
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateUser :one
UPDATE "user"
SET 
    email = COALESCE($2, email),
    phone_number = COALESCE($3, phone_number),
    full_name = COALESCE($4, full_name),
    hashed_password = COALESCE($5, hashed_password),
    is_active = COALESCE($6, is_active),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: SetPhoneVerified :one
UPDATE "user"
SET is_phone_verified = TRUE,
    is_active = TRUE,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: SetEmailVerified :one
UPDATE "user"
SET is_email_verified = TRUE,
    is_active = TRUE,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: GetUserRole :many
SELECT r.* FROM role r
INNER JOIN user_role ur ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: GetRole :many
SELECT r.id
FROM role r
WHERE r.name = ANY($1::text[]);

-- name: AssignRoleToUser :exec
INSERT INTO user_role (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- name: DeleteUserRole :exec
DELETE FROM user_role
WHERE user_id = $1;
