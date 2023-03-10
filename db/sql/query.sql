-- name: GetUsers :many
SELECT *
FROM users;
-- name: GetUser :one
SELECT *
FROM users
WHERE ID = sqlc.arg(id)
LIMIT 1;
-- name: GetPostsByCategoy :many
SELECT post_id FROM post_category WHERE category_id=sqlc.arg(id)