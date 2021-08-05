package seeds

import (
	"log"

	"github.com/xybor/xychat/models"
)

// SeedAdminUser creates an admin account with given username and password
func SeedAdminUser(username, password string) {
	role := "admin"

	err := models.GetDB().Create(
		&models.User{
			Username: &username,
			Password: &password,
			Role:     &role,
		},
	).Error

	if err != nil {
		log.Fatal("Cannot create admin user: " + err.Error())
	}

	log.Println("[Xychat] Successfully create user " + username + ":" + password)
}
