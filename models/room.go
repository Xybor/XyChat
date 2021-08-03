package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Users []*User `gorm:"many2many:detailed_rooms"`
}

type DetailedRoom struct {
	UserID uint `gorm:"primaryKey"`
	RoomID uint `gorm:"primaryKey"`

	User User
	Room Room
}
