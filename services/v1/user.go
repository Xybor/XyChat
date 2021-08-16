package v1

import (
	"errors"
	"log"
	"strings"

	"github.com/xybor/xychat/models"
	"github.com/xybor/xychat/xyerrors"

	resources "github.com/xybor/xychat/resources/v1"
	"gorm.io/gorm"
)

const (
	RoleMember   = "member"
	RoleModifier = "mod"
	RoleAdmin    = "admin"
)

var roleMap = map[string]int{
	RoleMember:   1,
	RoleModifier: 10,
	RoleAdmin:    100,
}

const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderGay    = "gay"
	GenderLes    = "les"
	GenderOther  = "other"
)

type userService struct {
	user        *models.User
	isLoaded    bool
	isValidated bool
}

// compareRole comapares two roles and return:
//
// + -1 if role1 < role2
//
// + 0  if role1 = role2
//
// + 1  if role1 > role2
//
// @note: role1 and role2 must be a valid role
func compareRole(role1 string, role2 string) int {
	roleLevel1 := roleMap[role1]
	roleLevel2 := roleMap[role2]

	if roleLevel1 < roleLevel2 {
		return -1
	} else if roleLevel1 == roleLevel2 {
		return 0
	} else {
		return 1
	}
}

// isValidRole checks if a role is in role list or not.
func isValidRole(role string) bool {
	for key := range roleMap {
		if role == key {
			return true
		}
	}

	return false
}

// CreateUserService create a userService with the subject is the user having
// given id.  If the parameter is nil, there is no subject in this service.
//
// @note: id must be nil or a valid id
func CreateUserService(id *uint, valid bool) userService {
	if id == nil {
		return userService{user: nil, isLoaded: false}
	}

	user := models.User{BaseModel: models.BaseModel{ID: *id}}
	return userService{
		user:        &user,
		isLoaded:    false,
		isValidated: valid,
	}
}

// The load method gets user data from database by id.  If the userService is
// signed as valid, but it can't be loaded from database, PANIC.  If it hasn't
// been valid yet and also can't be loaded from database, the user is assigned
// by nil.
func (us *userService) load() {
	if us.isLoaded {
		return
	}

	us.isLoaded = true

	if us.user == nil {
		return
	}

	err := models.GetDB().First(us.user, us.user.ID).Error

	if err != nil {
		if us.isValidated {
			log.Panicln(err)
		} else {
			us.user = nil
		}
	}

	us.isValidated = true
}

// The validate method check if a userService is validated or not.  It it
// hasn't been validated yet, it calls load() method instead.
func (us *userService) validate() {
	if us.isValidated {
		return
	}

	us.load()
}

// Register creates a user with given username and password.
//
// If the role is not member, a subject in service is required.
//
// @error: ErrorUnknownRole, ErrorPermission, ErrorExistedUsername, ErrorUnknown
func (us *userService) Register(username, password, role string) xyerrors.XyError {
	if !isValidRole(role) {
		return xyerrors.ErrorUnknownInput.New("Unknown role %s", role)
	}

	if role != RoleMember {
		us.load()

		if us.user == nil {
			return xyerrors.ErrorPermission.New("You must login before")
		}

		if compareRole(*us.user.Role, role) != 1 {
			return xyerrors.ErrorPermission.New("You aren't entitled to register a/an %s user", role)
		}
	}

	var user = models.User{
		Username: &username,
		Password: &password,
		Role:     &role,
	}

	err := models.GetDB().Create(&user).Error
	if err != nil {
		// If the error is duplicated record error, return a specific error and
		// doesn't log.
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return xyerrors.ErrorExistedUsername.New("Existed username, please choose another")
		}

		return xyerrors.ErrorUnknown
	}

	return xyerrors.NoError
}

// Remove deletes a user with a given id.  It needs a subject to determine the
// permission.
//
// @error: ErrorPermission, ErrorUnknown
func (us *userService) Remove(id uint) xyerrors.XyError {
	us.validate()

	if us.user == nil {
		return xyerrors.ErrorPermission.New("You must login before")
	}

	if us.user.ID != id {
		us.load()

		userResponse, xerr := us.Select(id)
		if xerr.Errno() != 0 {
			return xerr
		}

		if compareRole(*us.user.Role, userResponse.Role) != 1 {
			return xyerrors.ErrorPermission.New("You aren't entitled to remove the %s user", userResponse.Role)
		}
	}

	removedUser := models.User{BaseModel: models.BaseModel{ID: id}}

	err := models.GetDB().Delete(&removedUser).Error
	if err != nil {
		return xyerrors.ErrorUnknown
	}

	return xyerrors.NoError
}

// RemoveByUsername deletes a user with a given name.  It needs a subject to
// determine the permission.
//
// @error: ErrorPermission, ErrorUnknown
func (us *userService) RemoveByUsername(username string) xyerrors.XyError {
	us.validate()

	if us.user == nil {
		return xyerrors.ErrorPermission.New("You must login before")
	}

	userResponse, xerr := us.SelectByName(username)
	if xerr.Errno() != 0 {
		return xerr
	}

	if userResponse.ID != us.user.ID &&
		compareRole(*us.user.Role, userResponse.Role) != 1 {
		return xyerrors.ErrorPermission.New(
			"You aren't entitled to remove the %s user", userResponse.Role)
	}

	err := models.GetDB().Where("username=?", username).Delete(&models.User{}).Error
	if err != nil {
		return xyerrors.ErrorUnknown
	}

	return xyerrors.NoError
}

// Authenticate checks if the given username and password belongs to a user.
// If yes, it returns a userResponse of that user.
//
// @error: ErrorFailedAuthentication, ErrorUnknown
func (us *userService) Authenticate(
	username, password string,
) (*resources.UserResponse, xyerrors.XyError) {
	us.validate()

	userResponse := resources.UserResponse{}

	err := models.GetDB().
		Select("id, username, role, age, gender").
		Where("username = ? and password = ?", username, password).
		First(&models.User{}).
		Scan(&userResponse).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xyerrors.ErrorFailedAuthentication.New("Invalid username or password")
		}

		return nil, xyerrors.ErrorUnknown
	}

	return &userResponse, xyerrors.NoError
}

// AuthenticateById is the same as Authenticate, except it uses id instead of
// username.  It doesn't return userResponse.
//
// @error: ErrorFailedAuthentication, ErrorUnknown
func (us *userService) AuthenticateById(
	id uint, password string,
) (*resources.UserResponse, xyerrors.XyError) {
	us.validate()

	userResponse := resources.UserResponse{}

	err := models.GetDB().
		Select("id, username, role, age, gender").
		Where("id = ? and password = ?", id, password).
		First(&models.User{}).
		Scan(&userResponse).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xyerrors.ErrorFailedAuthentication.New("Invalid id or password")
		}

		return nil, xyerrors.ErrorUnknown
	}

	return &userResponse, xyerrors.NoError
}

// Select gets the userResponse of a given userid from database.  It
// needs a subject in service to determine the permission.
//
// @error: ErrorPermission, ErrorUnknown
func (us *userService) Select(
	id uint,
) (*resources.UserResponse, xyerrors.XyError) {
	us.validate()

	if us.user == nil {
		return nil, xyerrors.ErrorPermission.New("You must login before")
	}

	if us.user.ID != id {
		us.load()

		if *us.user.Role != RoleModifier && *us.user.Role != RoleAdmin {
			return nil, xyerrors.ErrorPermission.New(
				"You aren't entitled to view the profile of another user")
		}
	}

	userResponse := resources.UserResponse{}
	err := models.GetDB().
		Select("ID, username, role, age, gender").
		First(&models.User{}, id).
		Scan(&userResponse).Error

	if err != nil {
		return nil, xyerrors.ErrorUnknown
	}

	if us.user.ID != id && compareRole(*us.user.Role, userResponse.Role) != 1 {
		return nil, xyerrors.ErrorPermission.New(
			"You aren't entitled to view the profile of a/an %s user", userResponse.Role)
	}

	return &userResponse, xyerrors.NoError
}

// SelectByUsername gets the userResponse of a given name from database.  It
// needs a subject in service to determine the permission.
//
// @error: ErrorPermission, ErrorUnknown
func (us *userService) SelectByName(
	username string,
) (*resources.UserResponse, xyerrors.XyError) {
	us.load()

	if us.user == nil {
		return nil, xyerrors.ErrorPermission.New("You must login before")
	}

	userResponse := resources.UserResponse{}
	err := models.GetDB().
		Select("ID, username, role, age, gender").
		Where("username=?", username).
		First(&models.User{}).
		Scan(&userResponse).Error

	if err != nil {
		return nil, xyerrors.ErrorUnknown
	}

	if us.user.ID != userResponse.ID &&
		compareRole(*us.user.Role, userResponse.Role) != 1 {
		return nil, xyerrors.ErrorPermission.New(
			"You aren't entitled to view the profile of a/an %s user", userResponse.Role)
	}

	return &userResponse, xyerrors.NoError
}

// SelfSelect is a shortcut of us.Select(us.user.ID).
//
// @error: ErrorPermission, ErrorUnknown
func (us *userService) SelfSelect() (*resources.UserResponse, xyerrors.XyError) {
	us.validate()

	if us.user == nil {
		return nil, xyerrors.ErrorPermission.New("You must login before")
	}

	return us.Select(us.user.ID)
}

// UpdateInfo updates age and gender for a specific user determined by id.  It
// needs a subject in the service to determine permission.
//
// @error: ErrorPermission, ErrorUnknown
func (us *userService) UpdateInfo(
	id uint,
	age *uint,
	gender *string,
) xyerrors.XyError {
	us.validate()

	if us.user == nil {
		return xyerrors.ErrorPermission.New("You must login before")
	}

	if us.user.ID != id {
		us.load()

		userResponse, xerr := us.Select(id)
		if xerr.Errno() != 0 {
			return xerr
		}

		if compareRole(*us.user.Role, userResponse.Role) != 1 {
			return xyerrors.ErrorPermission.New(
				"You aren't entitled to update the profile of a/an %s user", userResponse.Role)
		}
	}

	user := models.User{BaseModel: models.BaseModel{ID: id}}
	err := models.GetDB().Model(&user).Updates(
		models.User{
			Age:    age,
			Gender: gender,
		},
	).Error

	if err != nil {
		return xyerrors.ErrorUnknown
	}

	return xyerrors.NoError
}

// UpdateRole updates role for a specific user determined by id.  It needs a
// subject in the service to determine permission.
//
// @error: ErrorUnknownRole, ErrorPermission, ErrorUnknown
func (us *userService) UpdateRole(id uint, role string) xyerrors.XyError {
	if !isValidRole(role) {
		return xyerrors.ErrorUnknownInput.New("Unknown role %s", role)
	}

	us.load()

	if us.user == nil {
		return xyerrors.ErrorPermission.New("You must login before")
	}

	if compareRole(*us.user.Role, role) != 1 {
		return xyerrors.ErrorPermission.New("You aren't entitled to update a user to %s", role)
	}

	userResponse, xerr := us.Select(id)
	if xerr.Errno() != 0 {
		return xerr
	}

	if compareRole(*us.user.Role, userResponse.Role) != 1 {
		return xyerrors.ErrorPermission.New(
			"You aren't entitled to update the role of a/an %s user", userResponse.Role)
	}

	user := models.User{BaseModel: models.BaseModel{ID: id}}
	err := models.GetDB().Model(&user).Updates(
		models.User{
			Role: &role,
		},
	).Error

	if err != nil {
		return xyerrors.ErrorUnknown
	}

	return xyerrors.NoError
}

// UpdatePassword updates password for a specific user determined by id.  It
// needs a subject in the service to determine permission.  If the user change
// his password, he needs to provide oldpwd.  If the admin or mod changes the
// password of another, he doesn't need to provide oldpwd.
//
// @error: ErrorInvalidOldPassword, ErrorPermission, ErrorUnknown
func (us *userService) UpdatePassword(id uint, oldpwd *string, newpwd string) xyerrors.XyError {
	us.validate()

	if us.user == nil {
		return xyerrors.ErrorPermission.New("You must login before")
	}

	if us.user.ID == id {
		if oldpwd == nil {
			return xyerrors.ErrorFailedAuthentication.New("Please provide the old password")
		}

		_, xerr := us.AuthenticateById(id, *oldpwd)
		if xerr.Errno() != 0 {
			return xyerrors.ErrorFailedAuthentication.New("Old password is incorrect")
		}
	} else {
		us.load()
		userResponse, xerr := us.Select(id)
		if xerr.Errno() != 0 {
			return xerr
		}

		if compareRole(*us.user.Role, userResponse.Role) != 1 {
			return xyerrors.ErrorPermission.New(
				"You aren't entitled to update password of a/an %s user", userResponse.Role)
		}
	}

	user := models.User{BaseModel: models.BaseModel{ID: id}}
	err := models.GetDB().Model(&user).Updates(
		models.User{
			Password: &newpwd,
		},
	).Error

	if err != nil {
		return xyerrors.ErrorUnknown
	}

	return xyerrors.NoError
}

// LoadRooms retrieves all rooms which this user joined,  these rooms will be
// store in us.user.Rooms
//
// @error: ErrorPermission, ErrorUnknown
func (us *userService) LoadRooms() xyerrors.XyError {
	us.load()

	if us.user == nil {
		return xyerrors.ErrorPermission.New("You must login before")
	}

	err := models.GetDB().Model(us.user).Association("Rooms").Find(&us.user.Rooms)
	if err != nil {
		return xyerrors.ErrorUnknown
	}

	return xyerrors.NoError
}
