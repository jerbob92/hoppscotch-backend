package models

import "gorm.io/gorm"

type TeamCollection struct {
	gorm.Model
	TeamID   uint
	Team     Team
	Title    string
	ParentID uint
}
