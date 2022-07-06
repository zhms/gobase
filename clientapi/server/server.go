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
var modulename string
var dbprefix string
var debug bool = false

func Init() {
	abugo.Init()
	debug = viper.GetBool("server.debug")
	systemname = viper.GetString("server.systemname")
	modulename = viper.GetString("server.modulename")
	dbprefix = viper.GetString("server.dbprefix")
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

func SystemName() string {
	return systemname
}

func ModuleName() string {
	return modulename
}

func DbPrefix() string {
	return dbprefix
}

func Run() {
	abugo.Run()
}

type TokenData struct {
	UserId   int
	SellerId int
}

func GetToken(ctx *abugo.AbuHttpContent) *TokenData {
	td := TokenData{}
	err := json.Unmarshal([]byte(ctx.TokenData), &td)
	if err != nil {
		return nil
	}
	return &td
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

func GetSeller(SellerId int) *SellerData {
	rediskey := fmt.Sprint(systemname, ":seller")
	bytedata := redis.HGet(rediskey, fmt.Sprint(SellerId))
	if bytedata == nil {
		return nil
	}
	sellerdata := SellerData{}
	json.Unmarshal(bytedata.([]byte), &sellerdata)
	return &sellerdata
}

func flush_seller() {
	for {
		rediskey := fmt.Sprint(systemname, ":seller")
		sql := "select * from x_seller"
		dbresult, err := db.Conn().Query(sql)
		if err != nil {
			return
		}
		keys := redis.HKeys(rediskey)
		for dbresult.Next() {
			sellerdata := SellerData{}
			abugo.GetDbResult(dbresult, &sellerdata)
			if sellerdata.State != 1 {
				redis.HDel(rediskey, fmt.Sprint(sellerdata.SellerId))
			} else {
				redis.HSet(rediskey, fmt.Sprint(sellerdata.SellerId), sellerdata)
			}
			for i := 0; i < len(keys); i++ {
				if keys[i] == fmt.Sprint(sellerdata.SellerId) {
					keys = append(keys[:i], keys[i+1:]...)
				}
			}
		}
		for i := 0; i < len(keys); i++ {
			redis.HDel(rediskey, keys[i])
		}
		time.Sleep(time.Second * 10)
	}
}
