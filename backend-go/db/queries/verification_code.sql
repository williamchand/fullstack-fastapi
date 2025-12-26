-- name: CreateVerificationCode :one
INSERT INTO verification_code (
    user_id,
    verification_code,
    verification_type,
    extra_metadata,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: CreateVerificationCodeNoUser :one
INSERT INTO verification_code (
    verification_code,
    verification_type,
    extra_metadata,
    expires_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetLatestUnusedVerificationCode :one
SELECT * FROM verification_code
WHERE user_id = $1
  AND verification_type = $2
  AND used_at IS NULL
ORDER BY created_at DESC
LIMIT 1;

-- name: GetVerificationCodeByCode :one
SELECT * FROM verification_code
WHERE user_id = $1
  AND verification_type = $2
  AND verification_code = $3
LIMIT 1;

-- name: GetVerificationCodeByCodeOnly :one
SELECT * FROM verification_code
WHERE verification_type = $1
  AND verification_code = $2
LIMIT 1;

-- name: MarkVerificationCodeUsed :one
UPDATE verification_code
SET used_at = now()
WHERE id = $1
RETURNING *;
