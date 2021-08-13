package models

type ChatMessage struct {
	BaseModel
	RoomID  uint `gorm:"not null"`
	Room    Room
	UserID  uint `gorm:"not null"`
	User    User
	Message string `gorm:"not null"`
}
