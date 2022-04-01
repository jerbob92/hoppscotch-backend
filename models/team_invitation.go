package models

import "gorm.io/gorm"

type TeamInvitation struct {
	gorm.Model
	TeamID       uint
	Team         Team
	UserID       uint
	User         User
	InviteeRole  TeamMemberRole
	InviteeEmail string
	Code         string
}
