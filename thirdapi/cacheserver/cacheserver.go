package cacheserver

import (
	"encoding/json"

	"github.com/beego/beego/logs"
	"github.com/imroc/req"
	"github.com/spf13/viper"
)

var cacheserver string

func Init() {
	cacheserver = viper.GetString("server.cacheserver")
}

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
	type RequestData struct{
		SellerId int
	}
	reqdata := RequestData{}
	reqdata.SellerId = SellerId
	bdata,_ := json.Marshal(reqdata)
	resp,err := req.Post(cacheserver + "/seller/get",string(bdata))
	if err != nil{
		logs.Error(err)
		return nil
	}
	type ReponseData struct{
		Data string `json:"data"`
	}
	respdata := ReponseData{}
	bytedata,_ := resp.ToBytes()
	err = json.Unmarshal(bytedata,&respdata)
	if err != nil{
		logs.Error(err)
		return nil
	}
	if len(respdata.Data) == 0 {
		return nil
	}
	sellerdata := SellerData{}
	err = json.Unmarshal([]byte(respdata.Data),&sellerdata)
	if err != nil{
		logs.Error(err)
		return nil
	}
	return &sellerdata
}