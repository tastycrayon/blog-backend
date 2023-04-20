package db

import (
	"time"
)

type Category struct {
	ID            int64     `json:"ID"`
	CategoryTitle string    `json:"category_title"`
	CategorySlug  string    `json:"category_slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Image struct {
	ID           int64
	ImageTitle   string
	ImageUrl     string
	ThumbnailUrl string
	Height       int32
	Width        int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Post struct {
	ID          int64
	PostTitle   string
	PostSlug    string
	PostContent string
	PostImage   int64
	PostAuthor  int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PostCat struct {
	Post       Post
	Categories []Category
}

type PostCategory struct {
	ID         int64
	PostID     int64
	CategoryID int64
}

type User struct {
	ID                int64
	UserName          string
	UserEmail         string
	UserPass          string
	UserActivationKey string
	UserRole          int32
	DisplayName       string
	UserImage         int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
