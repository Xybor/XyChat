package v1

import (
	"github.com/xybor/xychat/models"
)

type roomService struct {
	room        *models.Room
	isValidated bool
}

// CreateRoomService creates a roomService struct with given roomid.
func CreateRoomService(id *uint) roomService {
	if id == nil {
		return roomService{room: nil, isValidated: false}
	}

	return roomService{room: &models.Room{BaseModel: models.BaseModel{ID: *id}}, isValidated: false}
}

// Create creates a room with some users.  It need at least two users to create
// a room.  The room information is saved in the called object.
//
// @error: ErrorNotEnoughUserToCreateRoom, ErrorUnknown
func (rs *roomService) Create(us ...*userService) error {
	if len(us) < 2 {
		return ErrorNotEnoughUserToCreateRoom
	}

	db := models.GetDB().Begin()

	rs.room = &models.Room{}
	err := models.GetDB().Create(rs.room).Error

	if err != nil {
		db.Rollback()
		return err
	}

	users := make([]*models.User, len(us))
	for i, u := range us {
		u.load()

		if u.user == nil {
			db.Rollback()
			return ErrorPermission
		}

		users[i] = u.user
	}

	err = models.GetDB().Model(rs.room).Association("Users").Append(users)

	if err != nil {
		db.Rollback()
		return ErrorUnknown
	}

	db.Commit()

	return nil
}

// Admit adds a user to the room
//
// @error: ErrorPermission, ErrorUnknown
func (rs *roomService) Admit(us *userService) error {
	if rs.room == nil {
		return ErrorPermission
	}

	us.load()

	if us.user == nil {
		return ErrorPermission
	}

	err := models.GetDB().Model(rs.room).Association("Users").Append(us.user)

	if err != nil {
		return ErrorUnknown
	}

	return nil
}

// The contain method checks if the user with uid is in the room or not.
func (rs *roomService) contain(uid uint) bool {
	for _, user := range rs.room.Users {
		if user.ID == uid {
			return true
		}
	}

	return false
}

// LoadUsers selects all users in the room from database and saves it in
// rs.room.Users
//
// @error: ErrorPermission
func (rs *roomService) LoadUsers() error {
	if rs.room == nil {
		return ErrorPermission
	}

	err := models.GetDB().Model(rs.room).
		Association("Users").
		Find(&rs.room.Users)

	return err
}
