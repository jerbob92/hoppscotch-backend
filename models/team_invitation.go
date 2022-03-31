package models

type TeamInvitation struct {
	ID           int64
	CreatorUid   string
	InviteeRole  TeamMemberRole
	inviteeEmail string
	TeamID       uint64
}
