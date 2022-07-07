package abucache

import (
	"encoding/json"
	"fmt"
	"xserver/server"

	"github.com/beego/beego/logs"
)

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

func GetSeller(SellerId int) *SellerData {
	rediskey := fmt.Sprint(server.SystemName(),":seller")
	bdata := server.Redis().HGet(rediskey,fmt.Sprint(SellerId))
	if bdata == nil {
		return nil
	}
	sellerdata := SellerData{}
	err := json.Unmarshal(bdata.([]byte),&sellerdata)
	if err != nil{
		logs.Error(err)
		return nil
	}
	return &sellerdata
}

func GetSystemConfig(SellerId int,ConfigName string) string{
	rediskey := fmt.Sprint(server.SystemName(), ":systemconfig:",SellerId)
	bdata := server.Redis().HGet(rediskey,ConfigName)
	if bdata == nil{
		return ""
	}
	return string(bdata.([]byte))
}