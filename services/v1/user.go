package v1

import (
	"errors"
	"log"
	"strings"

	"github.com/xybor/xychat/models"
	representation "github.com/xybor/xychat/representations/v1"
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

var (
	ErrorPermission           = errors.New("you don't have the permission")
	ErrorUnknownRole          = errors.New("unknown role")
	ErrorExistedUsername      = errors.New("username existed")
	ErrorInvalidOldPassword   = errors.New("invalid old password")
	ErrorFailedAuthentication = errors.New("invalid username or password")
	ErrorUnknown              = errors.New("some error occurs")
)

type userService struct {
	user         *models.User
	isAuthorized bool
}

// compareRole comapares two roles and return:
//
// + -1 if role1 < role2
//
// + 0  if role1 = role2
//
// + 1  if role1 > role2
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

// authorize reads the subject's data in this service from database and assign
// it to userService.user
func (us *userService) authorize() {
	if !us.isAuthorized {
		us.isAuthorized = true

		if us.user == nil {
			return
		}

		models.GetDB().First(us.user, us.user.ID)
	}
}

// CreateUserService create a userService with the subject is the user having
// given id.  If the parameter is nil, there is no subject in this service.
func CreateUserService(id *uint) userService {
	if id == nil {
		return userService{user: nil, isAuthorized: false}
	}

	user := models.User{BaseModel: models.BaseModel{ID: *id}}
	return userService{
		user:         &user,
		isAuthorized: false,
	}
}

// Register creates a user with given username and password.
//
// If the role is not member, a subject in service is required.
func (us *userService) Register(username, password, role string) error {
	var user models.User

	if !isValidRole(role) {
		return ErrorUnknownRole
	}

	if role != RoleMember {
		us.authorize()

		if us.user == nil || compareRole(*us.user.Role, role) != 1 {
			return ErrorPermission
		}
	}

	user.Username = &username
	user.Password = &password
	user.Role = &role

	err := models.GetDB().Create(&user).Error
	if err != nil {
		// If the error is duplicated record error, return a specific error and
		// doesn't log.
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return ErrorExistedUsername
		}

		log.Println(err)
		return ErrorUnknown
	}

	return nil
}

// Remove deletes a user with a given id.  It needs a subject to determine the
// permission.
func (us *userService) Remove(id uint) error {
	if us.user == nil {
		return ErrorPermission
	}

	if us.user.ID != id {
		us.authorize()

		userRepresentation, err := us.Select(id)
		if err != nil {
			return err
		}

		if compareRole(*us.user.Role, userRepresentation.Role) != 1 {
			return ErrorPermission
		}
	}

	removedUser := models.User{BaseModel: models.BaseModel{ID: id}}

	err := models.GetDB().Delete(&removedUser).Error
	if err != nil {
		log.Println(err)
		return ErrorUnknown
	}

	return nil
}

// RemoveByUsername deletes a user with a given name.  It needs a subject to
// determine the permission.
func (us *userService) RemoveByUsername(username string) error {
	log.Println(*us)
	if us.user == nil {
		return ErrorPermission
	}

	us.authorize()

	userRepresentation, err := us.SelectByName(username)
	if err != nil {
		return err
	}

	if userRepresentation.ID != us.user.ID &&
		compareRole(*us.user.Role, userRepresentation.Role) != 1 {
		return ErrorPermission
	}

	err = models.GetDB().Where("username=?", username).Delete(&models.User{}).Error
	if err != nil {
		log.Println(err)
		return ErrorUnknown
	}

	return nil
}

// Authenticate checks if the given username and password belongs to a user.
// If yes, it returns a userRepresentation of that user.
func (us *userService) Authenticate(
	username, password string,
) (*representation.UserRepresentation, error) {
	userRepresentation := representation.UserRepresentation{}

	err := models.GetDB().
		Select("id, username, role, age, gender").
		Where("username = ? and password = ?", username, password).
		First(&models.User{}).
		Scan(&userRepresentation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorFailedAuthentication
		}

		log.Println(err)
		return nil, ErrorUnknown
	}

	return &userRepresentation, nil
}

// AuthenticateById is the same as Authenticate, except it uses id instead of
// username.  It doesn't return userRepresentation.
func (us *userService) AuthenticateById(
	id uint, password string,
) (*representation.UserRepresentation, error) {
	userRepresentation := representation.UserRepresentation{}

	err := models.GetDB().
		Select("id, username, role, age, gender").
		Where("id = ? and password = ?", id, password).
		First(&models.User{}).
		Scan(&userRepresentation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorFailedAuthentication
		}

		log.Println(err)
		return nil, ErrorUnknown
	}

	return &userRepresentation, nil
}

// Select gets the userRepresentation of a given userid from database.  It
// needs a subject in service to determine the permission.
func (us *userService) Select(id uint) (*representation.UserRepresentation, error) {
	if us.user == nil {
		return nil, ErrorPermission
	}

	if us.user.ID != id {
		us.authorize()

		if *us.user.Role != RoleModifier && *us.user.Role != RoleAdmin {
			return nil, ErrorPermission
		}
	}

	userRepresentation := representation.UserRepresentation{}
	err := models.GetDB().
		Select("ID, username, role, age, gender").
		First(&models.User{}, id).
		Scan(&userRepresentation).Error

	if err != nil {
		log.Println(err)
		return nil, ErrorUnknown
	}

	if us.user.ID != id && compareRole(*us.user.Role, userRepresentation.Role) != 1 {
		return nil, ErrorPermission
	}

	return &userRepresentation, nil
}

// SelectByUsername gets the userRepresentation of a given name from database.  It
// needs a subject in service to determine the permission.
func (us *userService) SelectByName(username string) (*representation.UserRepresentation, error) {
	if us.user == nil {
		return nil, ErrorPermission
	}

	userRepresentation := representation.UserRepresentation{}
	err := models.GetDB().
		Select("ID, username, role, age, gender").
		Where("username=?", username).
		First(&models.User{}).
		Scan(&userRepresentation).Error

	if err != nil {
		log.Println(err)
		return nil, ErrorUnknown
	}

	if us.user.ID != userRepresentation.ID &&
		compareRole(*us.user.Role, userRepresentation.Role) != 1 {
		return nil, ErrorPermission
	}

	return &userRepresentation, nil
}

// SelfSelect is a shortcut of us.Select(us.user.ID).
func (us *userService) SelfSelect() (*representation.UserRepresentation, error) {
	if us.user == nil {
		return nil, ErrorPermission
	}

	return us.Select(us.user.ID)
}

// UpdateInfo updates age and gender for a specific user determined by id.  It
// needs a subject in the service to determine permission.
func (us *userService) UpdateInfo(
	id uint,
	age *uint,
	gender *string,
) error {
	if us.user == nil {
		return ErrorPermission
	}

	if us.user.ID != id {
		us.authorize()

		r, err := us.Select(id)
		if err != nil {
			return err
		}

		if compareRole(*us.user.Role, r.Role) != 1 {
			return ErrorPermission
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
		log.Println(err)
		return ErrorUnknown
	}

	return nil
}

// UpdateRole updates role for a specific user determined by id.  It needs a
// subject in the service to determine permission.
func (us *userService) UpdateRole(id uint, role string) error {
	if !isValidRole(role) {
		return ErrorUnknownRole
	}

	us.authorize()
	if us.user == nil || compareRole(*us.user.Role, role) != 1 {
		return ErrorPermission
	}

	r, err := us.Select(id)
	if err != nil {
		return err
	}

	if compareRole(*us.user.Role, r.Role) != 1 {
		return ErrorPermission
	}

	user := models.User{BaseModel: models.BaseModel{ID: id}}
	err = models.GetDB().Model(&user).Updates(
		models.User{
			Role: &role,
		},
	).Error

	if err != nil {
		log.Println(err)
		return ErrorUnknown
	}

	return nil
}

// UpdatePassword updates password for a specific user determined by id.  It
// needs a subject in the service to determine permission.  If the user change
// his password, he needs to provide oldpwd.  If the admin or mod changes the
// password of another, he doesn't need to provide oldpwd.
func (us *userService) UpdatePassword(id uint, oldpwd *string, newpwd string) error {
	if us.user == nil {
		return ErrorPermission
	}

	if us.user.ID == id {
		if oldpwd == nil {
			return ErrorInvalidOldPassword
		}

		_, err := us.AuthenticateById(id, *oldpwd)
		if err != nil {
			return ErrorInvalidOldPassword
		}
	} else {
		us.authorize()
		userRepresentation, err := us.Select(id)
		if err != nil {
			return err
		}

		if compareRole(*us.user.Role, userRepresentation.Role) != 1 {
			return ErrorPermission
		}
	}

	user := models.User{BaseModel: models.BaseModel{ID: id}}
	err := models.GetDB().Model(&user).Updates(
		models.User{
			Password: &newpwd,
		},
	).Error

	if err != nil {
		log.Println(err)
		return ErrorUnknown
	}

	return nil
}
