package v1

import (
	"github.com/xybor/xychat/models"
)

type roomService struct {
	id *uint
}

// CreateRoomService creates a roomService struct with given roomid.
func CreateRoomService(id *uint) roomService {
	return roomService{id: id}
}

// Create creates a room and assigns roomid to this roomService.
func (rs *roomService) Create() error {
	db := models.GetDB()

	room := models.Room{}
	err := db.Create(&room).Error

	if err != nil {
		return err
	}

	rs.id = &room.ID

	return nil
}
