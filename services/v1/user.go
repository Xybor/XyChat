package v1

import (
	"errors"

	"github.com/xybor/xychat/models"
	representation "github.com/xybor/xychat/representations/v1"
	"gorm.io/gorm"
)

type userService struct {
	id       *uint
	username *string
	password *string
}

func CreateUserService(id *uint, username, password *string) userService {
	return userService{id: id, username: username, password: password}
}

func (us userService) Register() error {
	var user models.User

	role := "member"
	user.Username = us.username
	user.Password = us.password
	user.Role = &role

	db := models.GetDB()

	err := db.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (us userService) Authenticate() (representation.UserRepresentation, error) {
	response := representation.UserRepresentation{}

	if us.username == nil {
		return response, errors.New("empty username")
	}

	if us.password == nil {
		return response, errors.New("empty password")
	}

	err := models.GetDB().
		Select("ID, username, role, age, gender").
		Where("username = ? and password = ?", us.username, us.password).
		First(&models.User{}).
		Scan(&response).Error

	return response, err
}

func (us userService) GetProfile() (representation.UserRepresentation, error) {
	response := representation.UserRepresentation{}

	if us.id == nil {
		return response, errors.New("empty id")
	}

	err := models.GetDB().
		Select("ID, username, role").
		First(&models.User{Model: gorm.Model{ID: *us.id}}).
		Scan(&response).Error

	return response, err
}
