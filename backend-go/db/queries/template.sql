-- name: ListTemplates :many
SELECT * FROM template ORDER BY created_at DESC;

-- name: GetTemplateByID :one
SELECT * FROM template WHERE id = $1 LIMIT 1;
