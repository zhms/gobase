package controller

import (
	"fmt"
	"strings"
	"time"
	"xserver/abugo"
	"xserver/server"
)

type ThirdController struct {
}

func (c *ThirdController) Init() {
	group := server.Http().NewGroup("/third/v1")
	{
		group.PostNoAuth("/user_register", third_user_register)
		group.PostNoAuth("/get_balance", third_get_balance)
		group.PostNoAuth("/transfer_in", third_transfer_in)
		group.PostNoAuth("/transfer_out", third_transfer_out)
		group.PostNoAuth("/server_login", third_server_login)
	}
}

func third_user_register(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UniqueId  string `validate:"required"`
		Password  string `validate:"required"`
		TimeStamp int64  `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	seller := server.GetSeller(reqdata.SellerId)
	if ctx.RespErrString(seller == nil, &errcode, "商户不存在") {
		return
	}
	if ctx.RespErrString(!server.Debug() && !abugo.RsaVerify(reqdata, seller.ApiThirdPublicKey), &errcode, "签名不正确") {
		return
	}
	sql := fmt.Sprintf("call %sthird_register(?,?,?,?)", server.DbPrefix)
	dbresult, err := server.Db().Conn().Query(sql, reqdata.UniqueId, reqdata.SellerId, reqdata.Password, "{}")
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		UserId    int
		Timestamp int64
		Sign      string
	}
	retdata := ReturnData{}
	if dbresult.Next() {
		dberr := abugo.GetDbResult(dbresult, &retdata)
		if ctx.RespDbErr(dberr) {
			return
		}
	}
	retdata.Timestamp = time.Now().Unix()
	if !server.Debug() {
		retdata.Sign = abugo.RsaSign(retdata, seller.ApiPrivateKey)
	}
	ctx.RespOK(retdata)
}

func third_get_balance(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UserId    int32  `validate:"required"`
		Symbol    string `validate:"required"`
		AssetType int    `validate:"required"`
		TimeStamp int64  `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	reqdata.Symbol = strings.ToLower(reqdata.Symbol)
	seller := server.GetSeller(reqdata.SellerId)
	if ctx.RespErrString(seller == nil, &errcode, "商户不存在") {
		return
	}
	if ctx.RespErrString(!server.Debug() && !abugo.RsaVerify(reqdata, seller.ApiThirdPublicKey), &errcode, "签名不正确") {
		return
	}
	sql := fmt.Sprintf("select AssetAmt from %sasset where UserId = ? and AssetType = ? and Symbol = ?", server.DbPrefix)
	dbresult, err := server.Db().Conn().Query(sql, reqdata.UserId, reqdata.AssetType, reqdata.Symbol)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Balance   int64
		Timestamp int64
		Sign      string
	}
	retdata := ReturnData{}
	if dbresult.Next() {
		dbresult.Scan(&retdata.Balance)
	}
	dbresult.Close()
	retdata.Timestamp = time.Now().Unix()
	if !server.Debug(){
		retdata.Sign = abugo.RsaSign(retdata, seller.ApiPrivateKey)
	}
	ctx.RespOK(retdata)
}

func third_transfer_in(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UserId    int32 `validate:"required"`
		Symbol    string `validate:"required"`
		AssetType int    `validate:"required"`
		OrderId   int64  `validate:"required"`
		Amount    int64  `validate:"required"`
		TimeStamp int64  `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	reqdata.Symbol = strings.ToLower(reqdata.Symbol)
	seller := server.GetSeller(reqdata.SellerId)
	if ctx.RespErrString(seller == nil, &errcode, "商户不存在") {
		return
	}
	if ctx.RespErrString(!server.Debug() && !abugo.RsaVerify(reqdata, seller.ApiThirdPublicKey), &errcode, "签名不正确") {
		return
	}
	sql := fmt.Sprintf("call %sthird_transfer_in_in(?,?,?,?,?,?,?,?,?)", server.DbPrefix)
	dbresult, err := server.Db().Conn().Query(sql, reqdata.OrderId, reqdata.UserId, reqdata.SellerId, reqdata.AssetType, reqdata.Symbol, reqdata.Amount, "{}", 1, "钱包转入")
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Balance   int64
		Timestamp int64
		Sign      string
	}
	retdata := ReturnData{}
	if dbresult.Next() {
		dberr := abugo.GetDbResult(dbresult, &retdata)
		if ctx.RespDbErr(dberr) {
			return
		}
	}
	retdata.Timestamp = time.Now().Unix()
	if !server.Debug(){
		retdata.Sign = abugo.RsaSign(retdata, seller.ApiPrivateKey)
	}
	ctx.Put("balance", retdata.Balance)
	ctx.RespOK(retdata)
}

func third_transfer_out(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UserId    string `validate:"required"`
		Symbol    string `validate:"required"`
		AssetType int    `validate:"required"`
		OrderId   int64  `validate:"required"`
		Amount    int64  `validate:"required"`
		TimeStamp int64  `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	reqdata.Symbol = strings.ToLower(reqdata.Symbol)
	seller := server.GetSeller(reqdata.SellerId)
	if ctx.RespErrString(seller == nil, &errcode, "商户不存在") {
		return
	}
	if ctx.RespErrString(!server.Debug() && !abugo.RsaVerify(reqdata, seller.ApiThirdPublicKey), &errcode, "签名不正确") {
		return
	}
	sql := fmt.Sprintf("call %sapi_transfer_in_out(?,?,?,?,?,?,?,?,?)", server.DbPrefix)
	dbresult, err := server.Db().Conn().Query(sql, reqdata.OrderId, reqdata.UserId, reqdata.SellerId, reqdata.AssetType, reqdata.Symbol, reqdata.Amount, "{}", 1, "钱包转入")
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Balance   int64
		Sign      string
		TimeStamp int64
	}
	retdata := ReturnData{}
	if dbresult.Next() {
		dberr := abugo.GetDbResult(dbresult, &retdata)
		if ctx.RespDbErr(dberr) {
			return
		}
	}
	retdata.TimeStamp = time.Now().Unix()
	retdata.Sign = abugo.RsaSign(retdata, seller.ApiPrivateKey)
	ctx.RespOK(retdata)
}

type third_server_token struct {
	UserId   int
	SellerId int
}

func third_server_login(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UserId    int    `validate:"required"`
		Password  string `validate:"required"`
		TimeStamp int64  `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	seller := server.GetSeller(reqdata.SellerId)
	if ctx.RespErrString(seller == nil, &errcode, "商户不存在") {
		return
	}
	if ctx.RespErrString(!server.Debug() && !abugo.RsaVerify(reqdata, seller.ApiThirdPublicKey), &errcode, "签名不正确") {
		return
	}
	sql := fmt.Sprintf("select Password from %suser where UserId = ? and SellerId = ?", server.DbPrefix)
	dbresult, err := server.Db().Conn().Query(sql, reqdata.UserId, reqdata.SellerId)
	if ctx.RespErr(err, &errcode) {
		return
	}
	if dbresult.Next() {
		var password string
		dbresult.Scan(&password)
		if reqdata.Password != password {
			ctx.RespErrString(true, &errcode, "密码不正确")
			return
		}
	} else {
		ctx.RespErrString(true, &errcode, "账号不存在")
		return
	}
	type ReturnData struct {
		Token     string
		Sign      string
		TimeStamp int64
	}
	retdata := ReturnData{}
	retdata.Token = abugo.GetUuid()
	tokendata := third_server_token{}
	tokendata.UserId = reqdata.UserId
	tokendata.SellerId = reqdata.SellerId
	server.Redis().SetEx(fmt.Sprint("exchange:third:server:login:", retdata.Token), 60*5, tokendata)
	retdata.TimeStamp = time.Now().Unix()
	retdata.Sign = abugo.RsaSign(retdata, seller.ApiPrivateKey)
	ctx.RespOK(retdata)
}
