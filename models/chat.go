package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	RoomID  uint `gorm:"not null"`
	Room    Room
	UserID  uint `gorm:"not null"`
	User    User
	Message string `gorm:"not null"`
}
