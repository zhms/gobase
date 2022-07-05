package controller

import (
	"fmt"
	"math/rand"
	"xserver/abugo"
	"xserver/server"

	"github.com/beego/beego/logs"
)

type VerifyController struct {
}

func (c *VerifyController) Init() {
	gropu := server.Http().NewGroup("/verify")
	{
		gropu.PostNoAuth("/send", c.send)
	}
}

////////////////////////////////////////////////////////////////////////
//发送验证码
///////////////////////////////////////////////////////////////////////

func (c *VerifyController) send(ctx *abugo.AbuHttpContent) {
	type RequestData struct {
		Account  string `validate:"required"` //账号
		SellerId int    `validate:"required"` //运营商
		UseType  int    `validate:"required"` //用途
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	VerifyCode := fmt.Sprint(rand.Intn(999999-100000) + 100000)
	logs.Debug(reqdata.Account)
	sql := "replace into x_verify(Account,SellerId,UseType,VerifyCode)values(?,?,?,?)"
	_ , err = server.Db().Conn().Query(sql, reqdata.Account, reqdata.SellerId, reqdata.UseType, VerifyCode)
	if ctx.RespErr(err,&errcode) {
		return
	}
	if server.Debug() {
		ctx.Put("VerifyCode", VerifyCode)
	}
	ctx.RespOK()
}

/////////////////////////////////////////////////////////////////////////
