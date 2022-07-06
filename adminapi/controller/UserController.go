package controller

import (
	"fmt"
	"xserver/abugo"
	"xserver/server"
)

type UserController struct {
}

func (c *UserController) Init() {
	group := server.Http().NewGroup("/user")
	{
		group.Post("/list", user_list)
	}
}

func user_list(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		SellerId int `validate:"required"`
		UserId   int
		Account  string
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := server.GetToken(ctx)
	if ctx.RespErrString(!server.Auth2(token, "玩家管理", "账号管理", "查"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	if reqdata.SellerId == -1 {
		reqdata.SellerId = 0
	}
	where := abugo.AbuWhere{}
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	where.AddInt("and", "UserId", reqdata.UserId, 0)
	where.AddString("and", "Account", reqdata.Account, "")
	var total int
	sql := fmt.Sprintf("%suser", server.DbPrefix())
	server.Db().QueryScan(where.CountSql(sql), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	sql = fmt.Sprintf("%suser", server.DbPrefix())
	dbresult, err := server.Db().Conn().Query(where.Sql(sql, reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Id           int
		Account      string
		SellerId     int
		UserId       int
		NickName     string
		RegisterTime string
	}
	data := []ReturnData{}
	for dbresult.Next() {
		data_element := ReturnData{}
		abugo.GetDbResult(dbresult, &data_element)
		data = append(data, data_element)
	}
	dbresult.Close()
	ctx.Put("data", data)
	ctx.Put("page", reqdata.Page)
	ctx.Put("pagesize", reqdata.PageSize)
	ctx.Put("total", total)
	ctx.RespOK()
}
