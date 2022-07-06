package server

import (
	"encoding/json"
	"fmt"
	"time"
	"xserver/abugo"

	"github.com/spf13/viper"
)

var http *abugo.AbuHttp
var redis *abugo.AbuRedis
var db *abugo.AbuDb
var websocket *abugo.AbuWebsocket
var systemname string
var debug bool = false

func Init() {
	abugo.Init()
	debug = viper.GetBool("server.debug")
	systemname = viper.GetString("server.systemname")
	http = new(abugo.AbuHttp)
	http.Init("server.http.http.port")
	redis = new(abugo.AbuRedis)
	redis.Init("server.redis")
	db = new(abugo.AbuDb)
	db.Init("server.db")
	SetupDatabase()
	go flush_seller()
}

func Http() *abugo.AbuHttp {
	return http
}

func Redis() *abugo.AbuRedis {
	return redis
}

func Db() *abugo.AbuDb {
	return db
}

func Debug() bool {
	return debug
}

func SystemName() string{
	return systemname
}

func Run() {
	abugo.Run()
}

type TokenData struct {
	UserId   int
	SellerId int
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
	ApiRiskThirdPublicKey string
}


func GetToken(ctx *abugo.AbuHttpContent) *TokenData {
	td := TokenData{}
	err := json.Unmarshal([]byte(ctx.TokenData), &td)
	if err != nil {
		return nil
	}
	return &td
}


func GetSeller(SellerId int) *SellerData {
	sql := "select * from x_seller where SellerId = ?"
	dbresult,err := db.Conn().Query(sql,SellerId)
	if err != nil{
		return nil
	}
	if dbresult.Next() {
		r := SellerData{}
		abugo.GetDbResult(dbresult,&r)
		return &r
	}
	return nil
}

func flush_seller(){
	for{
		sql := "select * from x_seller"
		dbresult,err := db.Conn().Query(sql)
		if err != nil{
			return
		}
		for dbresult.Next() {
			r := SellerData{}
			abugo.GetDbResult(dbresult,&r)
			fmt.Println(r)
		}
		time.Sleep(time.Second * 5)
	}
}