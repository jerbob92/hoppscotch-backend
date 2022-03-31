package models

import "gorm.io/gorm"

type TeamRequest struct {
	gorm.Model
	TeamID           uint
	Team             Team
	TeamCollectionID uint
	TeamCollection   TeamCollection
	Request          string `gorm:"type:MEDIUMTEXT"`
	Title            string
}
