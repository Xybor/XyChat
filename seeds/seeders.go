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
		log.Panicf("Cannot create admin user: %s\n", err.Error())
	}

	log.Printf("[Xychat] Successfully create user %s:%s\n", username, password)
}
