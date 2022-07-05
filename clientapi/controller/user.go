package controller

import (
	"encoding/json"
	"xserver/abugo"
	"xserver/server"
)

type UserController struct {
}

func (c *UserController) Init() {
	gropu := server.Http().NewGroup("/user")
	{
		gropu.PostNoAuth("/register", c.register)
		gropu.PostNoAuth("/login_password", c.login_password)
		gropu.PostNoAuth("/login_verifycode", c.login_verifycode)
	}
}

////////////////////////////////////////////////////////////////////////
//注册
///////////////////////////////////////////////////////////////////////
func (c *UserController) register(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId   int    `validate:"required"` //运营商
		Account    string `validate:"required"` //账号
		Password   string `validate:"required"` //密码
		VerifyCode string `validate:"required"` //验证码
		AgentId    int
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	dbconn := server.Db().Conn()
	type ExtraData struct {
		Ip      string
		AgentId int
	}
	extra := ExtraData{}
	extra.Ip = ctx.GetIp()
	extra.AgentId = reqdata.AgentId
	extrabytes, _ := json.Marshal(&extra)
	queryresult, err := dbconn.Query("call x_api_user_register(?,?,?,?,?)", reqdata.Account, reqdata.SellerId, reqdata.Password, reqdata.VerifyCode, string(extrabytes))
	if ctx.RespErr(err, &errcode) {
		return
	}
	queryresult.Next()
	type ReturnData struct {
		UserId int
	}
	dbresult := ReturnData{}
	dberr := abugo.GetDbResult(queryresult, &dbresult)
	if ctx.RespDbErr(dberr) {
		return
	}
	queryresult.Close()
	ctx.RespOK(dbresult)
}

////////////////////////////////////////////////////////////////////////
//玩家登录(密码)
///////////////////////////////////////////////////////////////////////
func (c *UserController) login_password(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId int    `validate:"required"` //运营商
		Account  string `validate:"required"` //账号
		Password string `validate:"required"` //密码
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ExtraData struct {
		Ip string
	}
	extra := ExtraData{}
	extra.Ip = ctx.GetIp()
	strextra, _ := json.Marshal(&extra)
	dbconn := server.Db().Conn()
	queryresult, err := dbconn.Query("call ex_api_user_login_password(?,?,?,?)", reqdata.Account, reqdata.SellerId, reqdata.Password, string(strextra))
	if ctx.RespErr(err, &errcode) {
		return
	}
	if queryresult.Next() {
		type ReturnData struct {
		}
		dbresult := ReturnData{}
		dberr := abugo.GetDbResult(queryresult, &dbresult)
		if ctx.RespDbErr(dberr) {
			return
		}
	}
	queryresult.Close()
	ctx.RespOK()
}

////////////////////////////////////////////////////////////////////////
//玩家登录(密码+验证码)
///////////////////////////////////////////////////////////////////////
func (c *UserController) login_verifycode(ctx *abugo.AbuHttpContent) {
	type RequestData struct {
		SellerId   int    `validate:"required"` //运营商
		Account    string `validate:"required"` //账号
		Password   string `validate:"required"` //密码
		VerifyCode string `validate:"required"` //机器码
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ExtraData struct {
		Ip string
	}
	extra := ExtraData{}
	extra.Ip = ctx.GetIp()
	strextra, _ := json.Marshal(&extra)
	dbconn := server.Db().Conn()
	queryresult, err := dbconn.Query("call ex_api_user_login_verifycode(?,?,?,?,?)", reqdata.Account, reqdata.SellerId, reqdata.Password, reqdata.VerifyCode, string(strextra))
	if ctx.RespErr(err, &errcode) {
		return
	}
	queryresult.Next()
	type ReturnData struct {
		UserId   int
		SellerId int
		OldToken string
		NewToken string
	}
	dbresult := ReturnData{}
	dberr := abugo.GetDbResult(queryresult, &dbresult)
	if ctx.RespDbErr(dberr) {
		return
	}
	queryresult.Close()
	tokendata := server.TokenData{}
	tokendata.UserId = dbresult.UserId
	tokendata.SellerId = dbresult.SellerId
	if len(dbresult.OldToken) > 0 {
		server.Http().DelToken(dbresult.OldToken)
	}
	server.Http().SetToken(dbresult.NewToken, tokendata)
	ctx.Put("UserId", dbresult.UserId)
	ctx.Put("SellerId", dbresult.SellerId)
	ctx.Put("Token", dbresult.NewToken)
	ctx.RespOK()
}

///////////////////////////////////////////////////////////////////////
