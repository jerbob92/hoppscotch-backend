package models

import (
	"gorm.io/gorm"
	"time"
)

type Shortcode struct {
	gorm.Model
	Code      string
	Request   string
	UserID    uint
	User      User
	CreatedOn time.Time `gorm:"autoCreateTime"`
}
