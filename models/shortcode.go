package models

import (
	"gorm.io/gorm"
)

type Shortcode struct {
	gorm.Model
	Code    string
	Request string
	UserID  uint
	User    User
}
