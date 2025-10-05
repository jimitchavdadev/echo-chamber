package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `gorm:"not null" json:"content"`
	UserID  uint   `gorm:"not null" json:"-"`
	PostID  uint   `gorm:"not null" json:"postId"`
	User    User   `gorm:"foreignKey:UserID" json:"author"`
}
