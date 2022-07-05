package main

import (
	"xserver/controller"
	"xserver/server"
)

func main() {
	server.Init()
	new(controller.VerifyController).Init()
	new(controller.UserController).Init()
	new(controller.AssetController).Init()
	server.Run()
}
