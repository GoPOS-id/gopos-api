package server

import (
	"github.com/GoPOS-id/gopos-api/config"
	"github.com/GoPOS-id/gopos-api/database"
	"github.com/GoPOS-id/gopos-api/module"
)

func Init() {
	database.Init()
	app := fiberConfigrutaion()
	hostname := config.APP_HOST + ":" + config.APP_PORT
	module.Init(app)
	app.Listen(hostname)
}
