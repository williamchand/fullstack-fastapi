-- name: GetEmailTemplateByName :one
SELECT *
FROM email_template
WHERE name = $1
  AND is_active = TRUE
LIMIT 1;
