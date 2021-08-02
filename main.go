package main

import (
	_ "github.com/godror/godror"

	"github.com/xybor/xychat/helpers"
	"github.com/xybor/xychat/models"
	"github.com/xybor/xychat/routers"
)

func main() {
	helpers.LoadEnv()

	models.Initialize()
	models.CreateTables(false)

	router := routers.Route()
	router.Run(":1999")
}
