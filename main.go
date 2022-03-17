package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/helpers"
	"github.com/xybor/xychat/models"
	"github.com/xybor/xychat/routers"
	"github.com/xybor/xychat/seeds"
	servicev1 "github.com/xybor/xychat/services/v1"
)

func main() {
	reset := flag.Bool("reset", false, "Drop all tables before auto-migrating")
	admin := flag.String(
		"admin",
		"",
		"Create an admin user using an environment variable. Set the name of that variable here.",
	)
	run := flag.Bool("run", false, "Run the server")
	dotenv := flag.Bool("dotenv", false, "Load environment variables from .env file")

	flag.Parse()

	if *dotenv {
		// Load environment variables in .env file
		helpers.LoadEnv()
	}

	models.InitializeDB()
	models.CreateTables(*reset)

	if *admin != "" {
		data := helpers.MustReadEnv(*admin)

		credentials := strings.Split(data, ":")

		if len(credentials) != 2 {
			log.Fatalln("Invalid admin credentials")
		}

		seeds.SeedAdminUser(credentials[0], credentials[1])
	}

	if *run {
		// xychat = test, debug or release
		xychat := helpers.ReadEnvDefault("XYCHAT", "release")
		gin.SetMode(xychat)

		servicev1.InitializeMatchQueue(60 * time.Second)

		router := routers.Route()

		addr := ":" + helpers.MustReadEnv("PORT")

		TLS, err := helpers.ReadEnv("TLS")

		if err != nil {
			log.Fatal(router.Run(addr))
		} else {
			crt := TLS + ".crt"
			key := TLS + ".key"
			log.Fatal(router.RunTLS(addr, crt, key))
		}
	}
}
