package models

import (
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeLike    NotificationType = "like"
	NotificationTypeComment NotificationType = "comment"
)

type Notification struct {
	gorm.Model
	UserID   uint             `gorm:"not null;index"` // The user receiving the notification
	ActorID  uint             `gorm:"not null"`       // The user who performed the action
	Type     NotificationType `gorm:"not null"`
	EntityID uint             // e.g., The ID of the post that was liked/commented on
	IsRead   bool             `gorm:"default:false"`
	Actor    User             `gorm:"foreignKey:ActorID"`
}
