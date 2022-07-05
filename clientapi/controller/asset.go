package controller

import (
	"xserver/abugo"
	"xserver/server"
)

type AssetController struct {
}

func (c *AssetController) Init() {
	gropu := server.Http().NewGroup("/asset")
	{
		gropu.Post("/list", c.list)
	}
}

////////////////////////////////////////////////////////////////////////
//资产列表
///////////////////////////////////////////////////////////////////////
func (c *AssetController) list(ctx *abugo.AbuHttpContent) {
	token := server.GetToken(ctx)
	errcode := 0
	queryresult, err := server.Db().Conn().Query("select Symbol,AssetType,AssetAmt,FrozenAmt from ex_asset where userid = ?", token.UserId)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Symbol    string
		AssetType int
		AssetAmt  int64
		FrozenAmt int64
	}
	assets := []ReturnData{}
	for queryresult.Next() {
		dbresult := ReturnData{}
		abugo.GetDbResult(queryresult, &dbresult)
		assets = append(assets, dbresult)
	}
	queryresult.Close()
	ctx.RespOK(assets)
}
