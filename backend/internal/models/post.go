package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Content    string `gorm:"not null" json:"content"`
	UserID     uint   `gorm:"not null" json:"-"`       // Foreign key to User
	User       User   `gorm:"foreignKey:UserID" json:"author"`

	// Fields populated by custom queries, not stored in DB
	LikeCount int64 `gorm:"-" json:"likeCount"`
	IsLiked   bool  `gorm:"-" json:"isLiked"`
}
