package models

type TeamMember struct {
	ID   int64
	Role TeamMemberRole
}

type TeamMemberRole string

const (
	Editor TeamMemberRole = "EDITOR"
	Owner  TeamMemberRole = "OWNER"
	Viewer TeamMemberRole = "VIEWER"
)
