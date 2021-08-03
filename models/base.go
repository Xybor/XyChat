package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xybor/xychat/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Initialize() {
	var err error

	postgres_host := helpers.ReadEnv("postgres_host", "localhost")
	postgres_user := helpers.ReadEnv("postgres_user", "postgres")
	postgres_dbname := helpers.ReadEnv("postgres_dbname", "xychat")
	postgres_port := helpers.ReadEnv("postgres_port", "5432")
	postgres_password := helpers.ReadEnv("postgres_password", "pwd")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable password=%s",
		postgres_host,
		postgres_user,
		postgres_dbname,
		postgres_port,
		postgres_password,
	)
	fmt.Println(dsn)

	_, err = os.Stat("logs")
	if os.IsNotExist(err) {
		os.Mkdir("logs", 0600)
	}

	out, err := os.OpenFile(
		"logs/db.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0600,
	)

	if err != nil {
		log.Fatal(err)
	}

	newLogger := logger.New(
		log.New(out, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	db, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: newLogger,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	sqldb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	err = sqldb.Ping()
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
		err := db.Migrator().DropTable(
			&User{},
			&Room{},
			&DetailedRoom{},
			&Chat{},
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	err := db.AutoMigrate(
		&User{},
		&Room{},
		&DetailedRoom{},
		&Chat{},
	)

	if err != nil {
		log.Fatal(err)
	}
}
