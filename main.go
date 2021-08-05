package main

import (
	"flag"
	"log"
	"strings"

	"github.com/xybor/xychat/helpers"
	"github.com/xybor/xychat/models"
	"github.com/xybor/xychat/routers"
	"github.com/xybor/xychat/seeds"
	servicev1 "github.com/xybor/xychat/services/v1"
)

func main() {
	reset := flag.Bool("reset", false, "Drop all tables before auto-migrating")
	admin := flag.String("admin", "", "Create an admin user with format username:password")
	run := flag.Bool("run", false, "Run the server")

	flag.Parse()

	helpers.LoadEnv()

	models.InitializeDB()
	models.CreateTables(*reset)

	if *admin != "" {
		credentials := strings.Split(*admin, ":")
		if len(credentials) != 2 {
			log.Fatal("Invalid admin credentials")
		}
		seeds.SeedAdminUser(credentials[0], credentials[1])
	}

	if *run {
		servicev1.InitializeMatchQueue()

		router := routers.Route()
		router.Run(":" + helpers.MustReadEnv("port"))
	}
}
