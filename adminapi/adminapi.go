package main

import (
	"xserver/controller"
	"xserver/server"
)

func main() {
	server.Init()
	new(controller.UserController).Init()
	server.Run()
}
