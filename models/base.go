package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xybor/xychat/helpers"
)

var db *gorm.DB

func Initialize() {
	var err error

	postgres_host := helpers.ReadEnv("postgres_host", "localhost")
	postgres_user := helpers.ReadEnv("postgres_user", "postgres")
	postgres_dbname := helpers.ReadEnv("postgres_dbname", "xychat")
	postgres_password := helpers.ReadEnv("postgres_password", "pwd")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		postgres_host,
		postgres_user,
		postgres_dbname,
		postgres_password,
	)
	fmt.Println(dsn)

	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[CHAT] Connecting to database success")
}

func GetDB() *gorm.DB {
	return db
}

func CreateTables(drop_if_exists bool) {
	if drop_if_exists {
		db.DropTableIfExists(&User{})
	}
	db.CreateTable(&User{})
}
