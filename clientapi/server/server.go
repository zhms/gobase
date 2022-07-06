package server

import (
	"encoding/json"
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
