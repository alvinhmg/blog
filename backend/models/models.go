package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"size:50;not null;unique" json:"username"`
	Email     string         `gorm:"size:100;not null;unique" json:"email"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Role      string         `gorm:"size:20;default:'user'" json:"role"` // admin 或 user
	Posts     []Post         `gorm:"foreignKey:AuthorID" json:"-"`
	Comments  []Comment      `json:"-"`
}

// Category 文章分类
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:50;not null;unique" json:"name"`
	Slug        string         `gorm:"size:50;not null;unique" json:"slug"`
	Description string         `gorm:"size:255" json:"description"`
	Posts       []Post         `gorm:"many2many:post_categories" json:"-"`
}

// Tag 文章标签
type Tag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size:50;not null;unique" json:"name"`
	Slug      string         `gorm:"size:50;not null;unique" json:"slug"`
	Posts     []Post         `gorm:"many2many:post_tags" json:"-"`
}

// Post 博客文章
type Post struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Slug       string         `gorm:"size:200;not null;unique" json:"slug"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Excerpt    string         `gorm:"size:500" json:"excerpt"`
	CoverImage string         `gorm:"size:255" json:"cover_image"`
	Status     string         `gorm:"size:20;default:'draft'" json:"status"` // draft, published
	ViewCount  int            `gorm:"default:0" json:"view_count"`
	LikeCount  int            `gorm:"default:0" json:"like_count"`
	AuthorID   uint           `json:"author_id"`
	Author     User           `gorm:"foreignKey:AuthorID" json:"author"`
	Categories []Category     `gorm:"many2many:post_categories" json:"categories"`
	Tags       []Tag          `gorm:"many2many:post_tags" json:"tags"`
	Comments   []Comment      `json:"comments"`
}

// Comment 评论
type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	PostID    uint           `json:"post_id"`
	Post      Post           `json:"-"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user"`
	ParentID  *uint          `json:"parent_id"` // 父评论ID，用于回复功能
	Parent    *Comment       `gorm:"foreignKey:ParentID" json:"-"`
	Replies   []Comment      `gorm:"foreignKey:ParentID" json:"replies"`
	Status    string         `gorm:"size:20;default:'pending'" json:"status"` // pending, approved, rejected
}
