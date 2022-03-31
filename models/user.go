package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FBUID       string `gorm:"column:fb_uid;index"` // Firebase UID
	DisplayName string
	Email       string
	PhotoURL    string
}
