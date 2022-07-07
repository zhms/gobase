package cache

import (
	"fmt"
	"xserver/abugo"
	"xserver/server"

	"github.com/beego/beego/logs"
)

func FlushSeller() {
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
		logs.Error(err)
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
}
