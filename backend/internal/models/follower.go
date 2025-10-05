package models

import (
	"gorm.io/gorm"
)

// Follower model represents a follow relationship
// A user (FollowerID) follows another user (FollowingID)
type Follower struct {
	gorm.Model
	FollowerID  uint `gorm:"not null;uniqueIndex:idx_follow"` // The user who is following
	FollowingID uint `gorm:"not null;uniqueIndex:idx_follow"` // The user who is being followed
}
