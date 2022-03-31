package models

import "gorm.io/gorm"

type TeamMember struct {
	gorm.Model
	TeamID uint
	Team   Team
	UserID uint
	User   User
	Role   TeamMemberRole
}

type TeamMemberRole string

const (
	Editor TeamMemberRole = "EDITOR"
	Owner  TeamMemberRole = "OWNER"
	Viewer TeamMemberRole = "VIEWER"
)
