-- name: UpsertAICredential :one
INSERT INTO ai_credential (user_id, provider, api_key_enc)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, provider)
DO UPDATE SET api_key_enc = EXCLUDED.api_key_enc, updated_at = now()
RETURNING *;

-- name: GetAICredential :one
SELECT * FROM ai_credential WHERE user_id = $1 AND provider = $2 LIMIT 1;
