-- name: CreateOAuthAccount :one
INSERT INTO oauth_account (
    user_id, 
    provider, 
    provider_user_id, 
    access_token, 
    refresh_token, 
    token_expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOAuthAccount :one
SELECT * FROM oauth_account 
WHERE provider = $1 AND provider_user_id = $2 
LIMIT 1;

-- name: UpdateOAuthAccountTokens :one
UPDATE oauth_account 
SET 
    access_token = $3,
    refresh_token = $4,
    token_expires_at = $5,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;