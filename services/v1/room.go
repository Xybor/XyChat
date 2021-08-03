package v1

import (
	"github.com/xybor/xychat/models"
)

type roomService struct {
	id *uint
}

func CreateRoomService(id *uint) roomService {
	return roomService{id: id}
}

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
