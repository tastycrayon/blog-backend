package db

import (
	"time"
)

type Posts struct {
	ID          uint64    `gorm:"column:ID;AUTO_INCREMENT;NOT NULL"`
	PostTitle   string    `gorm:"column:post_title;NOT NULL"`
	PostSlug    string    `gorm:"column:post_slug;NOT NULL"`
	PostContent string    `gorm:"column:post_content;NOT NULL"`
	PostImage   uint64    `gorm:"column:post_image;default:0;NOT NULL"`
	PostAuthor  uint64    `gorm:"column:post_author;default:0;NOT NULL"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"`
}

type Categories struct {
	ID            uint64    `gorm:"column:ID;AUTO_INCREMENT;NOT NULL"`
	CategoryTitle string    `gorm:"column:category_title;NOT NULL"`
	CategorySlug  string    `gorm:"column:category_slug;NOT NULL"`
	Description   string    `gorm:"column:description;NOT NULL"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"`
}

type Users struct {
	ID                uint64    `gorm:"column:ID;AUTO_INCREMENT;NOT NULL"`
	UserName          string    `gorm:"column:user_name;NOT NULL"`
	UserEmail         string    `gorm:"column:user_email;NOT NULL"`
	UserPass          string    `gorm:"column:user_pass;NOT NULL"`
	UserActivationKey string    `gorm:"column:user_activation_key;NOT NULL"`
	UserRole          int       `gorm:"column:user_role;default:2;NOT NULL"`
	DisplayName       string    `gorm:"column:display_name;NOT NULL"`
	UserImage         uint64    `gorm:"column:user_image;default:0;NOT NULL"`
	CreatedAt         time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdatedAt         time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"`
}

type Images struct {
	ID           uint64    `gorm:"column:ID;AUTO_INCREMENT;NOT NULL"`
	ImageTitle   string    `gorm:"column:image_title;NOT NULL"`
	ImageUrl     string    `gorm:"column:image_url;NOT NULL"`
	ThumbnailUrl string    `gorm:"column:thumbnail_url;NOT NULL"`
	Height       uint      `gorm:"column:height;default:0;NOT NULL"`
	Width        uint      `gorm:"column:width;default:0;NOT NULL"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"`
}

type PostCategories struct {
	ID         uint64 `gorm:"column:ID;AUTO_INCREMENT;NOT NULL"`
	PostID     uint64 `gorm:"column:post_id;AUTO_INCREMENT;NOT NULL"`
	CategoryID uint64 `gorm:"column:category_id;AUTO_INCREMENT;NOT NULL"`
}
