package v1

import (
	"errors"

	"github.com/xybor/xychat/models"
	resp "github.com/xybor/xychat/responses/api/v1"
	"github.com/jinzhu/gorm"
)

type UserService struct {
	ID       *uint
	Username *string
	Password *string
}

func (us UserService) RegisterUser() error {
	var user models.User

	role := "normal"
	user.Username = us.Username
	user.Password = us.Password
	user.Role = &role

	db := models.GetDB()

	err := db.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (us UserService) Authenticate() (resp.UserResponse, error) {
	response := resp.UserResponse{}

	if us.Username == nil {
		return response, errors.New("empty username")
	}

	if us.Password == nil {
		return response, errors.New("empty password")
	}

	err := models.GetDB().
		Select("ID, username, role").
		Where("username = ? and password = ?", us.Username, us.Password).
		First(&models.User{}).
		Scan(&response).Error

	return response, err
}

func (us UserService) GetProfile() (resp.UserResponse, error) {
	response := resp.UserResponse{}

	if us.ID == nil {
		return response, errors.New("empty id")
	}

	err := models.GetDB().
		Select("ID, username, role").
		First(&models.User{Model: gorm.Model{ID: *us.ID}}).
		Scan(&response).Error

	return response, err
}
