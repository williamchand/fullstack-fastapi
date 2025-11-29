-- name: CreateDataSource :one
INSERT INTO data_source (
    user_id, name, type, host, port, database_name, username, password_enc, options
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetDataSourceByID :one
SELECT * FROM data_source WHERE id = $1 AND user_id = $2 LIMIT 1;
