package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Content string `gorm:"not null" json:"content"`
	UserID  uint   `gorm:"not null" json:"-"` // Foreign key to User
	User    User   `gorm:"foreignKey:UserID" json:"author"`
}
