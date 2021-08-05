package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username *string `gorm:"unique;not null"`
	Password *string `gorm:"not null"`
	Role     *string `gorm:"not null;check:role in ('member', 'mod', 'admin')"`
	Age      *uint
	Gender   *string `gorm:"check:gender in ('male', 'female', 'gay', 'les', 'other')"`
	Rooms    []*Room `gorm:"many2many:detailed_rooms"`
}
