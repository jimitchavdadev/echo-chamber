package models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserID uint `gorm:"not null;uniqueIndex:idx_like" json:"userId"`
	PostID uint `gorm:"not null;uniqueIndex:idx_like" json:"postId"`
}
