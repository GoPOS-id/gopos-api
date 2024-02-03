package main

import "github.com/GoPOS-id/gopos-api/api/server"

// @title           GoPOS API
// @version         1.0
// @description     Rest API Endpoint for GoPOS.
// @termsOfService  http://gopos.web.id/terms/
// @contact.name   Muhamad Fadli Aqil
// @contact.url    http://instagram.com/fadliiaqil_
// @contact.email  fadli.aqil12@gmail.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      api.gopos.web.id
// @BasePath  /
// @securityDefinitions.apiKey  Bearer Token
// @in Authorization
// @name Bearer Token
func main() {
	server.Init()
}
