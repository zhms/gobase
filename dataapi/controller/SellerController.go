package controller

import (
	"fmt"
	"xserver/abugo"
	"xserver/server"

	"github.com/beego/beego/logs"
)

type SellerController struct {
}

func (c *SellerController) Init() {
	group := server.Http().NewGroup("/seller")
	{
		group.PostNoAuth("/flush", seller_flush)
		group.PostNoAuth("/get", seller_get)
	}
	seller_flush(nil)
}

func seller_flush(ctx *abugo.AbuHttpContent) {
	type SellerData struct {
		SellerId              int
		SellerName            string
		State                 int
		ApiPublicKey          string
		ApiPrivateKey         string
		ApiThirdPublicKey     string
		ApiRiskPublicKey      string
		ApiRiskPrivateKey     string
		ApiThirdRiskPublicKey string
	}
	rediskey := fmt.Sprint(server.SystemName(), ":seller")
	sql := "select * from x_seller"
	dbresult, err := server.Db().Conn().Query(sql)
	if err != nil {
		return
	}
	keys := server.Redis().HKeys(rediskey)
	for dbresult.Next() {
		sellerdata := SellerData{}
		abugo.GetDbResult(dbresult, &sellerdata)
		if sellerdata.State != 1 {
			server.Redis().HDel(rediskey, fmt.Sprint(sellerdata.SellerId))
		} else {
			server.Redis().HSet(rediskey, fmt.Sprint(sellerdata.SellerId), sellerdata)
		}
		for i := 0; i < len(keys); i++ {
			if keys[i] == fmt.Sprint(sellerdata.SellerId) {
				keys = append(keys[:i], keys[i+1:]...)
			}
		}
	}
	for i := 0; i < len(keys); i++ {
		server.Redis().HDel(rediskey, keys[i])
	}
	logs.Debug("刷新运营商")
}

func seller_get(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId int `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	ctx.RespOK()
}
