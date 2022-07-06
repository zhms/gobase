package cacheserver

import (
	"github.com/beego/beego/logs"
	"github.com/imroc/req"
	"github.com/spf13/viper"
)

var cacheserver string

func Init() {
	cacheserver = viper.GetString("server.cacheserver")
}

func FlushSeller() {
	_, err := req.Post(cacheserver + "/seller/flush")
	if err != nil {
		logs.Debug(err)
	}
}
