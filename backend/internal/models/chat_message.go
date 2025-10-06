package models

import (
	"gorm.io/gorm"
)

type ChatMessage struct {
	gorm.Model
	SenderID   uint   `gorm:"not null;index"`
	ReceiverID uint   `gorm:"not null;index"`
	Content    string `gorm:"not null"`
	Sender     User   `gorm:"foreignKey:SenderID"`
	Receiver   User   `gorm:"foreignKey:ReceiverID"`
}
