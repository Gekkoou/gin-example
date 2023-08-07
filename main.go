package main

import (
	"gin-example/initialize"
)

// @title gin-example api
// @version 1.0.0
// @description gin-exampl api.
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        token
// @BasePath
func main() {
	initialize.Initialize(false)
}
