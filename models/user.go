package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username *string `gorm:"type:varchar(30);unique;not null"`
	Password *string `gorm:"type:varchar(30);not null"`
	Role     *string `gorm:"not null"`
}
