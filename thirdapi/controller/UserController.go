package controller

import (
	"xserver/abugo"
	"xserver/server"
)

type UserController struct {
}

func (c *UserController) Init() {
	group := server.Http().NewGroup("/user")
	{
		group.PostNoAuth("/login", user_login)
	}
}

func user_login(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Token int `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	ctx.RespOK()
}
