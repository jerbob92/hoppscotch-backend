package models

import "gorm.io/gorm"

type TeamEnvironment struct {
	gorm.Model
	TeamID    uint
	Team      Team
	Name      string
	Variables string
}
