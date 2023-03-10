package db

import (
	"context"

	"github.com/mosiur404/goserver/util"
)

const insertImage = `-- name: InsertImage :exec
INSERT INTO images (image_title, image_url, thumbnail_url, height, width) 
VALUES (?, ?, ?, ?, ?);
`

type InsertImageParams struct {
	ImageTitle   string
	ImageUrl     string
	ThumbnailUrl string
	Height       int
	Width        int
}

func (q *Queries) InsertImage(ctx context.Context, params *InsertImageParams) (*int64, error) {
	result, err := q.db.ExecContext(ctx, insertImage,
		params.ImageTitle,
		params.ImageUrl,
		params.ThumbnailUrl,
		params.Height,
		params.Width,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()

	return &id, err
}

const insertPost = `-- name: InsertPost :exec
INSERT INTO posts (post_title, post_slug, post_content, post_author, post_image) 
VALUES (?, ?, ?, ?, ?);
`

type InsertPostParams struct {
	PostSlug    string
	PostTitle   string
	PostContent string
	Authorid    int64
	ImageID     int64
}

func (q *Queries) InsertPost(ctx context.Context, params *InsertPostParams) (*int64, error) {
	result, err := q.db.ExecContext(ctx, insertPost,
		params.PostTitle,
		params.PostSlug,
		params.PostContent,
		params.Authorid,
		params.ImageID,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()

	return &id, err
}

const getPosts = `-- name: GetPosts :many
SELECT id, post_title, post_slug, post_content, post_image, post_author, created_at, updated_at
FROM posts ORDER BY ? LIMIT ? OFFSET ?;
`

type GetPostsParams struct {
	OrderBy string
	Limit   int
	Offset  int
}

func (q *Queries) GetPosts(ctx context.Context, arg *GetPostsParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPosts, arg.OrderBy, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID, &i.PostTitle,
			&i.PostSlug, &i.PostContent,
			&i.PostImage, &i.PostAuthor,
			&i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

const getImage = `-- name: GetImage :one
SELECT id, user_name, user_email, user_pass, user_activation_key, user_role, display_name, user_image, created_at, updated_at
FROM users
WHERE ID = ?
LIMIT 1;
`

func (q *Queries) GetImage(ctx context.Context, id int64) (Image, error) {
	row := q.db.QueryRowContext(ctx, getImage, id)
	var i Image
	err := row.Scan(
		&i.ID, &i.ImageTitle,
		&i.ImageUrl, &i.ThumbnailUrl,
		&i.Height, &i.Width,
		&i.CreatedAt, &i.UpdatedAt,
	)
	return i, err
}

const getPost = `-- name: GetPost :one
SELECT id, post_title, post_slug, post_content, post_image, post_author, created_at, updated_at
FROM posts WHERE ID = ?
LIMIT 1;
`

func (q *Queries) GetPost(ctx context.Context, id int64) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, id)
	var i Post
	err := row.Scan(
		&i.ID, &i.PostTitle,
		&i.PostSlug, &i.PostContent,
		&i.PostImage, &i.PostAuthor,
		&i.CreatedAt, &i.UpdatedAt,
	)
	return i, err
}

const getPostCount = `-- name: GetPostCount :exec
SELECT COUNT(*) FROM posts;
`

func (q *Queries) GetPostCount(ctx context.Context) (int, error) {
	row := q.db.QueryRowContext(ctx, getPostCount)
	var count int
	err := row.Scan(&count)
	return count, err
}

const getPostBySlug = `-- name: GetPost :one
SELECT id, post_title, post_slug, post_content, post_image, post_author, created_at, updated_at
FROM posts WHERE post_slug LIKE ?
LIMIT 1;
`

func (q *Queries) GetPostBySlug(ctx context.Context, slug string) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostBySlug, slug)
	var i Post
	err := row.Scan(
		&i.ID, &i.PostTitle,
		&i.PostSlug, &i.PostContent,
		&i.PostImage, &i.PostAuthor,
		&i.CreatedAt, &i.UpdatedAt,
	)
	return i, err
}

func (q *Queries) GetUsersByIds(ctx context.Context, ids []string) ([]User, error) {
	getUsersByIds := `SELECT * FROM users WHERE ID` + util.MakeInQuery(len(ids)) + `;`
	// making array for args
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	// exec
	rows, err := q.db.QueryContext(ctx, getUsersByIds, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID, &i.UserName,
			&i.UserEmail, &i.UserPass,
			&i.UserActivationKey, &i.UserRole,
			&i.DisplayName, &i.UserImage,
			&i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmailOrUsername = `-- name: GetUser :one
SELECT * FROM users WHERE user_email = ? OR user_name = ? LIMIT 1;
`

func (q *Queries) GetUserByEmailOrUsername(ctx context.Context, usernameOrEmail string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmailOrUsername, usernameOrEmail, usernameOrEmail)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserName,
		&i.UserEmail,
		&i.UserPass,
		&i.UserActivationKey,
		&i.UserRole,
		&i.DisplayName,
		&i.UserImage,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

func (q *Queries) GetImagesByIds(ctx context.Context, ids []string) ([]Image, error) {
	getImagesByIds := `SELECT * FROM images WHERE ID` + util.MakeInQuery(len(ids)) + `;`

	// making array for args
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	// exec
	rows, err := q.db.QueryContext(ctx, getImagesByIds, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ID, &i.ImageTitle,
			&i.ImageUrl, &i.ThumbnailUrl,
			&i.Height, &i.Width,
			&i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCategories = `-- name: GetCategories :many
SELECT * FROM categories;
`

func (q *Queries) GetCategories(ctx context.Context) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, getCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID, &i.CategoryTitle,
			&i.CategorySlug, &i.Description,
			&i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

const getPostsByCat = `-- name: GetPosts :many
SELECT posts.* from posts
INNER JOIN post_category pc ON posts.ID = pc.post_id
INNER JOIN categories c ON c.ID = pc.category_id
WHERE c.category_slug = ? ORDER BY ? LIMIT ? OFFSET ?;
`

type GetPostsByCatParams struct {
	Category string
	OrderBy  string
	Limit    int
	Offset   int
}

func (q *Queries) GetPostsByCategory(ctx context.Context, arg *GetPostsByCatParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByCat, arg.Category, arg.OrderBy, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID, &i.PostTitle,
			&i.PostSlug, &i.PostContent,
			&i.PostImage, &i.PostAuthor,
			&i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
