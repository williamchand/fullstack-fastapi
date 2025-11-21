-- name: GetUserByID :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO "user" (
    email, phone_number, full_name, hashed_password
) VALUES (
    $1, $2, $3, $4
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

-- name: GetUserRoles :many
SELECT r.* FROM role r
JOIN user_role ur ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: AssignRoleToUser :exec
INSERT INTO user_role (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- name: CreateOAuthAccount :one
INSERT INTO oauth_account (
    user_id, provider, provider_user_id, access_token, refresh_token, token_expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOAuthAccount :one
SELECT * FROM oauth_account
WHERE provider = $1 AND provider_user_id = $2 LIMIT 1;