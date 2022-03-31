package models

import "github.com/jerbob92/hoppscotch-backend/db"

func AutoMigrate() error {
	return db.DB.AutoMigrate(&Shortcode{}, &Team{}, &TeamCollection{}, &TeamInvitation{}, &TeamMember{}, &TeamRequest{}, &User{})
}
