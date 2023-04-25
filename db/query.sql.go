package db

import (
	"context"
	"database/sql"

	"github.com/mosiur404/goserver/util"
)

const getPostsByCategoy = `-- name: GetPostsByCategoy :many
SELECT post_id FROM post_category WHERE category_id=?
`

func (q *Queries) GetPostsByCategoy(ctx context.Context, id int64) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByCategoy, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var post_id int64
		if err := rows.Scan(&post_id); err != nil {
			return nil, err
		}
		items = append(items, post_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT id, user_name, user_email, user_pass, user_activation_key, user_role, display_name, user_image, created_at, updated_at
FROM users
WHERE ID = ?
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
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

const getUsers = `-- name: GetUsers :many
SELECT id, user_name, user_email, user_pass, user_activation_key, user_role, display_name, user_image, created_at, updated_at
FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
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
	// var items []Post
	items := make([]Post, 0, 10)
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

const getPostsWithCategory = `-- name: GetPostsWithCategory :many 
SELECT posts.*, categories.* FROM posts 
JOIN post_category ON posts.id = post_category.post_id
JOIN categories ON post_category.category_id = categories.id LIMIT ? OFFSET ?;
`

func (q *Queries) GetPostsWithCategory(ctx context.Context, arg *GetPostsParams) ([]PostCat, error) {
	rows, err := q.db.QueryContext(ctx, getPostsWithCategory, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// var items []PostCat
	items := make([]PostCat, 0, 10)

	temp := make(map[int64]int)

	for rows.Next() {
		var p Post
		var c Category
		if err := rows.Scan(
			&p.ID, &p.PostTitle,
			&p.PostSlug, &p.PostContent,
			&p.PostImage, &p.PostAuthor,
			&p.CreatedAt, &p.UpdatedAt,

			&c.ID, &c.CategoryTitle,
			&c.CategorySlug, &c.Description,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		var index int
		var found bool
		if index, found = temp[p.ID]; !found {
			index = len(items)
			temp[p.ID] = index
			items = append(items, PostCat{Post: p})
		}

		items[index].Categories = append(items[index].Categories, c)
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
const getPostCountWithCat = `-- name: GetPostCount :exec
SELECT COUNT(*) FROM post_category pc
inner join categories c on c.ID=pc.category_id
WHERE c.category_slug = ?;
`

func (q *Queries) GetPostCount(ctx context.Context, cat *string) (int, error) {
	var row *sql.Row
	if cat == nil {
		row = q.db.QueryRowContext(ctx, getPostCount)
	} else {
		row = q.db.QueryRowContext(ctx, getPostCountWithCat, cat)
	}
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

const getPostWithCat = `SELECT
posts.*, categories.*
FROM posts
  JOIN ( SELECT ID FROM posts WHERE post_slug LIKE ? LIMIT 1) d ON posts.ID IN (d.ID)
JOIN post_category ON posts.id = post_category.post_id
JOIN categories ON post_category.category_id = categories.id;
`

func (q *Queries) GetPostWithCatBySlug(ctx context.Context, slug string) (*PostCat, error) {
	var item *PostCat

	rows, err := q.db.QueryContext(ctx, getPostWithCat, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	postAllocated := false
	for rows.Next() {
		var p Post
		var c Category
		if err := rows.Scan(
			&p.ID, &p.PostTitle,
			&p.PostSlug, &p.PostContent,
			&p.PostImage, &p.PostAuthor,
			&p.CreatedAt, &p.UpdatedAt,

			&c.ID, &c.CategoryTitle,
			&c.CategorySlug, &c.Description,
			&c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if !postAllocated {
			item = &PostCat{Post: p}
			postAllocated = true
		}
		item.Categories = append(item.Categories, c)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return item, nil
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
	// var items []Image
	items := make([]Image, 0, 10)
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
	// var items []Category
	items := make([]Category, 0, 10)
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
	// var items []Post
	items := make([]Post, 0, 10)
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
