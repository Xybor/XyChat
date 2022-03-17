package v1

import (
	"github.com/xybor/xychat/models"
	"github.com/xybor/xychat/xyerrors"
)

type roomService struct {
	room *models.Room
}

// CreateRoomService creates a roomService struct with given roomid.
func CreateRoomService(id *uint) roomService {
	if id == nil {
		return roomService{room: nil}
	}

	return roomService{room: &models.Room{BaseModel: models.BaseModel{ID: *id}}}
}

// Create creates a room with some users.  It need at least two users to create
// a room.  The room information is saved in the called object.
//
// @error: ErrorNotEnoughUserToCreateRoom, ErrorUnknown
func (rs *roomService) Create(us ...*userService) xyerrors.XyError {
	if len(us) < 2 {
		return xyerrors.ErrorNotEnoughUserToCreateRoom.New(
			"It needs at least two users to create a room")
	}

	db := models.GetDB().Begin()

	rs.room = &models.Room{}
	err := models.GetDB().Create(rs.room).Error

	if err != nil {
		db.Rollback()
		return xyerrors.ErrorUnknown
	}

	users := make([]*models.User, len(us))
	for i, u := range us {
		u.validate()

		if u.user == nil {
			db.Rollback()
			return xyerrors.ErrorUnknown.New("There is an invalid user")
		}

		u.load()
		users[i] = u.user
	}

	err = models.GetDB().Model(rs.room).Association("Users").Append(users)

	if err != nil {
		db.Rollback()
		return xyerrors.ErrorUnknown
	}

	db.Commit()

	return xyerrors.NoError
}

// Admit adds a user to the room
//
// @error: ErrorPermission, ErrorUnknown
func (rs *roomService) Admit(us *userService) xyerrors.XyError {
	if rs.room == nil {
		return xyerrors.ErrorInvalidService.New("Invalid room")
	}

	us.validate()

	if us.user == nil {
		return xyerrors.ErrorInvalidService.New("Invalid user")
	}

	err := models.GetDB().Model(rs.room).Association("Users").Append(us.user)

	if err != nil {
		return xyerrors.ErrorUnknown
	}

	return xyerrors.NoError
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
func (rs *roomService) LoadUsers() xyerrors.XyError {
	if rs.room == nil {
		return xyerrors.ErrorInvalidService.New("Invalid room")
	}

	err := models.GetDB().Model(rs.room).
		Association("Users").
		Find(&rs.room.Users)

	if err != nil {
		return xyerrors.ErrorUnknown
	}

	return xyerrors.NoError
}
