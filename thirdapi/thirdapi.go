package main

import (
	"xserver/cacheserver"
	"xserver/controller"
	"xserver/server"
)

func main() {
	server.Init()
	cacheserver.Init()
	new(controller.ThirdController).Init()
	new(controller.UserController).Init()
	server.Run()
}
