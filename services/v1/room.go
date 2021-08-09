package v1

import (
	"github.com/xybor/xychat/models"
)

type roomService struct {
	room *models.Room
}

// CreateRoomService creates a roomService struct with given roomid.
func CreateRoomService(id *uint) roomService {
	if id == nil {
		return roomService{nil}
	}

	return roomService{room: &models.Room{BaseModel: models.BaseModel{ID: *id}}}
}

// create creates a room and assigns roomid to this roomService.
func (rs *roomService) create() error {
	rs.room = &models.Room{}
	err := models.GetDB().Create(rs.room).Error

	if err != nil {
		return err
	}

	return nil
}
