package main

import (
	"github.com/xybor/xychat/helpers"
	"github.com/xybor/xychat/models"
	"github.com/xybor/xychat/routers"
	servicev1 "github.com/xybor/xychat/services/v1"
)

func main() {
	helpers.LoadEnv()

	models.Initialize()
	models.CreateTables(false)

	servicev1.InitializeMatchQueue()

	router := routers.Route()
	//fmt.Print(router)
	router.Run(":1999")
}
