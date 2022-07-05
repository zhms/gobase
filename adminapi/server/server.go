package server

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"xserver/abugo"

	"github.com/beego/beego/logs"
	"github.com/spf13/viper"
	"github.com/valyala/fastjson"
)

var http *abugo.AbuHttp
var redis *abugo.AbuRedis
var db *abugo.AbuDb
var websocket *abugo.AbuWebsocket
var debug bool = false

type SellerData struct {
	SellerId   int
	SellerName string
	State      int
	Remark     string
	CreateTime string
}
type CacheSellerData struct {
	SellerId               int
	SellerName             string
	State                  int
	ApiPublicKey           string
	ApiPrivateKey          string
	ApiThirdPublicKey      string
	ApiRiskPublicKey       string
	ApiRiskPrivateKey      string
	ApiRiskThirdPublicKey  string
	HbcServerPublicKey     string
	HbcLocalPublicKey      string
	HbcLocalPrivateKey     string
	HbcRiskServerPublicKey string
	HbcRiskLocalPublicKey  string
	HbcRiskLocalPrivateKey string
	HbcAppId               string
}

func seller_list(ctx *abugo.AbuHttpContent) {
	defer recover()
	errcode := 0
	type RequestData struct {
		Page     int
		PageSize int
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize == 0 {
		reqdata.PageSize = 10
	}
	token := GetToken(ctx)
	if ctx.RespErrString(token.SellerId != -1, &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(!Auth2(token, "系统管理", "运营商管理", "查"), &errcode, "权限不足") {
		return
	}
	where := abugo.AbuWhere{}
	where.OrderBy = "ASC"
	where.OrderKey = "SellerId"
	var total int
	db.QueryScan(where.CountSql(db_seller_tablename), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := db.Conn().Query(where.Sql(db_seller_tablename, reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if ctx.RespErr(err, &errcode) {
		return
	}
	data := []SellerData{}
	for dbresult.Next() {
		data_element := SellerData{}
		abugo.GetDbResult(dbresult, &data_element)
		data_element.CreateTime = abugo.TimeToUtc(data_element.CreateTime)
		data = append(data, data_element)
	}
	dbresult.Close()
	ctx.Put("data", data)
	ctx.Put("page", reqdata.Page)
	ctx.Put("pagesize", reqdata.PageSize)
	ctx.Put("total", total)
	ctx.RespOK()
}

func seller_add(ctx *abugo.AbuHttpContent) {
	defer recover()
	errcode := 0
	type RequestData struct {
		SellerName string `validate:"required"`
		State      int    `validate:"required"`
		Remark     string
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(token.SellerId != -1, &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(!Auth2(token, "系统管理", "运营商管理", "增"), &errcode, "权限不足") {
		return
	}
	sql := fmt.Sprintf("insert into %s(SellerName,State,Remark)values(?,?,?)", db_seller_tablename)
	db.QueryNoResult(sql, reqdata.SellerName, reqdata.State, reqdata.Remark)
	WriteAdminLog("添加运营商", ctx, reqdata)
	seller_flush(nil)
	ctx.RespOK()
}

func seller_modify(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId   int    `validate:"required"`
		SellerName string `validate:"required"`
		State      int    `validate:"required"`
		Remark     string
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(token.SellerId != -1, &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(!Auth2(token, "系统管理", "运营商管理", "改"), &errcode, "权限不足") {
		return
	}
	sql := fmt.Sprintf("update %s set SellerName = ?,State = ?,Remark = ? where SellerId = ?", db_seller_tablename)
	db.QueryNoResult(sql, reqdata.SellerName, reqdata.State, reqdata.Remark, reqdata.SellerId)
	WriteAdminLog("修改运营商", ctx, reqdata)
	seller_flush(nil)
	ctx.RespOK()
}

func seller_delete(ctx *abugo.AbuHttpContent) {
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
	token := GetToken(ctx)
	if ctx.RespErrString(token.SellerId != -1, &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(!Auth2(token, "系统管理", "运营商管理", "增"), &errcode, "权限不足") {
		return
	}
	sql := fmt.Sprintf("delete from  %s where SellerId = ?", db_seller_tablename)
	db.QueryNoResult(sql, reqdata.SellerId)
	WriteAdminLog("删除运营商", ctx, reqdata)
	redis.HDel("x_hash:seller", fmt.Sprint(reqdata.SellerId))
	seller_flush(nil)
	ctx.RespOK()
}

func Init() {
	abugo.Init()
	debug = viper.GetBool("server.debug")
	http = new(abugo.AbuHttp)
	http.Init("server.http.http.port")
	redis = new(abugo.AbuRedis)
	redis.Init("server.redis")
	db = new(abugo.AbuDb)
	db.Init("server.db")
	SetupDatabase()
	{
		http.PostNoAuth("/admin/user/login", user_login)
		http.Post("/admin/login_log", login_log)
		http.Post("/admin/role/list", role_list)
		http.Post("/admin/role/listall", role_listall)
		http.Post("/admin/role/roledata", role_data)
		http.Post("/admin/role/modify", role_modify)
		http.Post("/admin/role/add", role_add)
		http.Post("/admin/role/delete", role_delete)
		http.Post("/admin/opt_log", opt_log)
		http.Post("/admin/user/list", user_list)
		http.Post("/admin/user/modify", user_modify)
		http.Post("/admin/user/delete", user_delete)
		http.Post("/admin/user/add", user_add)
		http.Post("/admin/user/google", user_google)
		http.Post("seller/name", seller_name)
		http.Post("seller/list", seller_list)
		http.Post("seller/add", seller_add)
		http.Post("seller/delete", seller_delete)
		http.Post("seller/modify", seller_modify)
		http.Post("seller/flush", seller_flush)
	}
	sql := "select RoleData from admin_role where SellerId = -1 and RoleName = '超级管理员'"
	var dbauthdata string
	db.QueryScan(sql, []interface{}{}, &dbauthdata)
	if dbauthdata != AuthDataStr {
		sql = "select Id,SellerId,RoleName,RoleData from admin_role"
		dbresult, _ := db.Conn().Query(sql)
		for dbresult.Next() {
			var roleid int
			var sellerid int
			var rolename string
			var roledata string
			dbresult.Scan(&roleid, &sellerid, &rolename, &roledata)
			if sellerid == -1 && rolename == "超级管理员" {
				continue
			}
			jnewdata := make(map[string]interface{})
			json.Unmarshal([]byte(AuthDataStr), &jnewdata)
			clean_auth(jnewdata)
			jrdata := make(map[string]interface{})
			json.Unmarshal([]byte(roledata), &jrdata)
			for k, v := range jrdata {
				set_auth(k, jnewdata, v.(map[string]interface{}))
			}
			newauthbyte, _ := json.Marshal(&jnewdata)
			sql = "update admin_role set RoleData = ? where id = ?"
			db.QueryNoResult(sql, string(newauthbyte), roleid)
		}
		dbresult.Close()
		sql = "update admin_role set RoleData = ? where RoleName = '超级管理员'"
		db.QueryNoResult(sql, AuthDataStr)
	}
	seller_flush(nil)
}
func clean_auth(node map[string]interface{}) {
	for k, v := range node {
		if strings.Index(reflect.TypeOf(v).Name(), "float") >= 0 {
			node[k] = 0
		} else {
			clean_auth(v.(map[string]interface{}))
		}
	}
}
func set_auth(parent string, newdata map[string]interface{}, node map[string]interface{}) {
	for k, v := range node {
		if strings.Index(reflect.TypeOf(v).Name(), "float") >= 0 {
			if v.(float64) != 1 {
				continue
			}
			path := strings.Split(parent, ".")
			if len(path) == 0 {
				continue
			}
			fk, fok := newdata[path[0]]
			if !fok {
				continue
			}
			var pn *interface{} = &fk
			var finded bool = true
			for i := 1; i < len(path); i++ {
				tk := path[i]
				tv, ok := (*pn).(map[string]interface{})[tk]
				if !ok {
					finded = false
					break
				}
				pn = &tv
			}
			if finded {
				(*pn).(map[string]interface{})[k] = 1
			}
		} else {
			set_auth(parent+"."+k, newdata, v.(map[string]interface{}))
		}
	}
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
func Run() {
	abugo.Run()
}

type TokenData struct {
	Account  string
	SellerId int
	AuthData string
}

func GetToken(ctx *abugo.AbuHttpContent) *TokenData {
	tokendata := TokenData{}
	err := json.Unmarshal([]byte(ctx.TokenData), &tokendata)
	if err != nil {
		return nil
	}
	return &tokendata
}
func WriteAdminLog(opt string, ctx *abugo.AbuHttpContent, data interface{}) {
	token := ctx.Token
	strdata, _ := json.Marshal(&data)
	tokendata := GetToken(ctx)
	Ip := ctx.GetIp()
	go func() {
		sql := "insert into admin_opt_log(Account,Opt,Token,Data,Ip,SellerId)values(?,?,?,?,?,?)"
		db.QueryNoResult(sql, tokendata.Account, opt, token, string(strdata), Ip, tokendata.SellerId)
	}()
}
func Auth2(td *TokenData, m string, s string, o string) bool {
	defer recover()
	authdata := make(map[string]interface{})
	json.Unmarshal([]byte(td.AuthData), &authdata)
	im, imok := authdata[m]
	if !imok {
		return false
	}
	is, isok := im.(map[string]interface{})[s]
	if !isok {
		return false
	}
	io, iook := is.(map[string]interface{})[o]
	if !iook {
		return false
	}
	if strings.Index(reflect.TypeOf(io).Name(), "float64") < 0 {
		return false
	}
	if io.(float64) != 1 {
		return false
	}
	return true
}
func user_login(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Account    string `validate:"required"`
		Password   string `validate:"required"`
		VerifyCode string `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type MenuData struct {
		Icon  string     `json:"icon"`
		Index string     `json:"index"`
		Title string     `json:"title"`
		Subs  []MenuData `json:"subs"`
	}
	sqlstr := "select Id,SellerId,`Password`,RoleName,Token,State,GoogleSecret,LoginTime,LoginCount,RoleSellerId from admin_user where Account = ?"
	type UserData struct {
		Id           int
		Password     string
		Token        string
		Rolename     string
		SellerId     int
		State        int
		GoogleSecret string
		LoginTime    string
		LoginCount   int
		RoleSellerId int
	}
	userdata := UserData{}
	uqsresult, uqserr := db.Conn().Query(sqlstr, reqdata.Account)
	if ctx.RespErr(uqserr, &errcode) {
		return
	}
	if ctx.RespErrString(!uqsresult.Next(), &errcode, "账号不存在") {
		return
	}
	abugo.GetDbResult(uqsresult, &userdata)
	uqsresult.Close()

	if ctx.RespErrString(userdata.State != 1, &errcode, "账号已被禁用") {
		return
	}
	if ctx.RespErrString(userdata.Password != reqdata.Password, &errcode, "密码不正确") {
		return
	}
	if userdata.SellerId != -1 {
		var dbsellerid int
		sqlstr = fmt.Sprintf("select SellerId from %s where SellerId = ?", db_seller_tablename)
		db.QueryScan(sqlstr, []interface{}{userdata.SellerId}, &dbsellerid)
		if ctx.RespErrString(dbsellerid == 0, &errcode, "所属运营商已停用") {
			return
		}
	}
	var authdata string
	sqlstr = "select RoleData from admin_role where RoleName = ? and SellerId = ?"
	sqlparam := []interface{}{userdata.Rolename, userdata.RoleSellerId}
	qserr, qsresult := db.QueryScan(sqlstr, sqlparam, &authdata)
	if ctx.RespErr(qserr, &errcode) {
		return
	}
	if ctx.RespErrString(!qsresult, &errcode, "角色不存在") {
		return
	}
	if ctx.RespErrString(!debug && len(userdata.GoogleSecret) > 0 && !abugo.VerifyGoogleCode(userdata.GoogleSecret, reqdata.VerifyCode), &errcode, "验证码不正确") {
		return
	}
	if len(userdata.Token) > 0 {
		http.DelToken(userdata.Token)
	}
	tokendata := TokenData{}
	tokendata.Account = reqdata.Account
	tokendata.SellerId = userdata.SellerId
	tokendata.AuthData = authdata
	token := abugo.GetUuid()
	http.SetToken(token, tokendata)
	sqlstr = "update admin_user set Token = ?,LoginCount = LoginCount + 1,LoginTime = now(),LoginIp = ? where id = ?"
	err = db.QueryNoResult(sqlstr, token, ctx.GetIp(), userdata.Id)
	if ctx.RespErr(err, &errcode) {
		return
	}
	sqlstr = "insert into admin_login_log(UserId,SellerId,Account,Token,LoginIp)values(?,?,?,?,?)"
	err = db.QueryNoResult(sqlstr, userdata.Id, userdata.SellerId, reqdata.Account, token, ctx.GetIp())
	if ctx.RespErr(err, &errcode) {
		return
	}
	menu := []MenuData{}
	json.Unmarshal([]byte(MenuDataStr), &menu)
	parser := fastjson.Parser{}
	jauthdata, _ := parser.ParseBytes([]byte(authdata))
	//三级菜单
	for i := 0; i < len(menu); i++ {
		for j := 0; j < len(menu[i].Subs); j++ {
			smenu := []MenuData{}
			for k := 0; k < len(menu[i].Subs[j].Subs); k++ {
				open := jauthdata.GetInt(menu[i].Title, menu[i].Subs[j].Title, menu[i].Subs[j].Subs[k].Title, "查")
				if open == 1 {
					smenu = append(smenu, menu[i].Subs[j].Subs[k])
				}
			}
			menu[i].Subs[j].Subs = smenu
		}
	}
	//二级菜单
	for i := 0; i < len(menu); i++ {
		smenu := []MenuData{}
		for j := 0; j < len(menu[i].Subs); j++ {
			open := jauthdata.GetInt(menu[i].Title, menu[i].Subs[j].Title, "查")
			if open == 1 || len(menu[i].Subs[j].Subs) > 0 {
				smenu = append(smenu, menu[i].Subs[j])
			}
		}
		menu[i].Subs = smenu
	}
	//一级菜单
	smenu := []MenuData{}
	for i := 0; i < len(menu); i++ {
		open := jauthdata.GetInt(menu[i].Title, "查")
		if open == 1 || len(menu[i].Subs) > 0 {
			smenu = append(smenu, menu[i])
		}
	}
	menu = smenu
	jauth := make(map[string]interface{})
	json.Unmarshal([]byte(authdata), &jauth)

	ctx.Put("UserId", userdata.Id)
	ctx.Put("SellerId", userdata.SellerId)
	ctx.Put("Account", reqdata.Account)
	ctx.Put("AuthData", jauth)
	ctx.Put("MenuData", menu)
	ctx.Put("Token", token)
	ctx.Put("LoginTime", userdata.LoginTime)
	ctx.Put("Ip", ctx.GetIp())
	ctx.Put("LoginCount", userdata.LoginCount)
	ctx.Put("Version", "1.0.0")
	ctx.RespOK()
}
func login_log(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		Account  string
		SellerId int
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize == 0 {
		reqdata.PageSize = 10
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "登录日志", "查"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	where := abugo.AbuWhere{}
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	where.AddString("and", "Account", reqdata.Account, "")
	var total int
	db.QueryScan(where.CountSql("admin_login_log"), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := db.Conn().Query(where.Sql("admin_login_log", reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Id         int
		UserId     int
		SellerId   int
		Account    string
		LoginIp    string
		CreateTime string
	}
	data := []ReturnData{}
	for dbresult.Next() {
		data_element := ReturnData{}
		abugo.GetDbResult(dbresult, &data_element)
		data_element.CreateTime = abugo.TimeToUtc(data_element.CreateTime)
		data = append(data, data_element)
	}
	dbresult.Close()
	ctx.Put("data", data)
	ctx.Put("page", reqdata.Page)
	ctx.Put("pagesize", reqdata.PageSize)
	ctx.Put("total", total)
	ctx.RespOK()
}
func role_list(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		SellerId int
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize == 0 {
		reqdata.PageSize = 10
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "角色管理", "查"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	where := abugo.AbuWhere{}
	where.OrderBy = "ASC"
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	var total int
	db.QueryScan(where.CountSql("admin_role"), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := db.Conn().Query(where.Sql("admin_role", reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Id             int
		RoleName       string
		SellerId       int
		ParentSellerId int
		Parent         string
		RoleData       string
	}
	data := []ReturnData{}
	for dbresult.Next() {
		data_element := ReturnData{}
		abugo.GetDbResult(dbresult, &data_element)
		data = append(data, data_element)
	}
	dbresult.Close()
	ctx.Put("data", data)
	ctx.Put("page", reqdata.Page)
	ctx.Put("pagesize", reqdata.PageSize)
	ctx.Put("total", total)
	ctx.RespOK()
}
func role_listall(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		SellerId int
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "角色管理", "查"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	sql := "select RoleName from admin_role where SellerId = ?"
	dbresult, err := db.Conn().Query(sql, reqdata.SellerId)
	if ctx.RespErr(err, &errcode) {
		return
	}
	names := []string{}
	for dbresult.Next() {
		var RoleName string
		dbresult.Scan(&RoleName)
		names = append(names, RoleName)
	}
	dbresult.Close()
	ctx.RespOK(names)
}
func role_data(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId int    `validate:"required"`
		RoleName string `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "角色管理", "查"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	sql := "select RoleData from admin_role where SellerId = ? and RoleName = ?"
	var RoleData string
	db.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &RoleData)
	var SuperRoleData string
	sql = "select RoleData from admin_role where SellerId = -1 and RoleName = '超级管理员'"
	db.QueryScan(sql, []interface{}{}, &SuperRoleData)
	ctx.Put("RoleData", RoleData)
	ctx.Put("SuperRoleData", SuperRoleData)
	ctx.RespOK()
}
func role_check(parent string, parentdata map[string]interface{}, data map[string]interface{}, result *string) {
	defer recover()
	for k, v := range data {
		if strings.Index(reflect.TypeOf(v).Name(), "float") >= 0 {
			if v.(float64) != 1 {
				continue
			}
			path := strings.Split(parent, ".")
			if len(path) == 0 {
				continue
			}
			fk, fok := parentdata[path[0]]
			if !fok {
				continue
			}
			var pn *interface{} = &fk
			var finded bool = true
			for i := 1; i < len(path); i++ {
				tk := path[i]
				tv, ok := (*pn).(map[string]interface{})[tk]
				if !ok {
					finded = false
					break
				}
				pn = &tv
			}
			if finded {
				fv := (*pn).(map[string]interface{})[k].(float64)
				if fv != 1 {
					(*result) = "fail"
				}
			} else {
				(*result) = "fail"
			}

		} else {
			role_check(parent+"."+k, parentdata, v.(map[string]interface{}), result)
		}
	}
}
func role_modify(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId int    `validate:"required"`
		RoleName string `validate:"required"`
		RoleData string `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "角色管理", "改"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	var ParentSellerId int
	var ParentRoleName string
	sql := "select ParentSellerId,Parent from admin_role where SellerId = ? and RoleName = ?"
	db.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &ParentSellerId, &ParentRoleName)
	if ctx.RespErrString(len(ParentRoleName) == 0, &errcode, "上级角色不存在") {
		return
	}
	var ParentRoleData string
	sql = "select RoleData from admin_role where SellerId = ? and RoleName = ?"
	db.QueryScan(sql, []interface{}{ParentSellerId, ParentRoleName}, &ParentRoleData)
	if ctx.RespErrString(len(ParentRoleData) == 0, &errcode, "获取上级角色数据失败") {
		return
	}
	jparent := make(map[string]interface{})
	err = json.Unmarshal([]byte(ParentRoleData), &jparent)
	if ctx.RespErr(err, &errcode) {
		return
	}
	jdata := make(map[string]interface{})
	err = json.Unmarshal([]byte(reqdata.RoleData), &jdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	result := ""
	for k, v := range jdata {
		role_check(k, jparent, v.(map[string]interface{}), &result)
	}
	if ctx.RespErrString(len(result) > 0, &errcode, "权限大过上级角色") {
		return
	}
	sql = "update admin_role set  RoleData = ? where SellerId = ? and RoleName = ?"
	err = db.QueryNoResult(sql, reqdata.RoleData, reqdata.SellerId, reqdata.RoleName)
	if ctx.RespErr(err, &errcode) {
		return
	}
	WriteAdminLog("修改角色", ctx, reqdata)
	ctx.RespOK()
}
func role_add(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		ParentSellerId int    `validate:"required"`
		Parent         string `validate:"required"`
		SellerId       int    `validate:"required"`
		RoleName       string `validate:"required"`
		RoleData       string `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "角色管理", "增"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	if ctx.RespErrString(reqdata.SellerId != -1 && reqdata.SellerId != reqdata.ParentSellerId && reqdata.ParentSellerId != -1, &errcode, "上级角色运营商只能是总后台角色或跟自己所属运营商一致,不可以是别的运营商") {
		return
	}
	if ctx.RespErrString(reqdata.SellerId == -1 && reqdata.ParentSellerId != -1, &errcode, "总后台角色上级角色只能是总后台的角色") {
		return
	}
	var rid int = 0
	sql := "select id from admin_role where SellerId = ? and RoleName = ?"
	db.QueryScan(sql, []interface{}{reqdata.ParentSellerId, reqdata.Parent}, &rid)
	if ctx.RespErrString(rid == 0, &errcode, "上级角色不存在") {
		return
	}
	rid = 0
	sql = "select id from admin_role where SellerId = ? and RoleName = ? "
	db.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &rid)
	if ctx.RespErrString(rid > 0, &errcode, "角色已经存在") {
		return
	}
	if reqdata.SellerId != -1 {
		sql = fmt.Sprintf("select SellerId from %s where SellerId = ? and state = 1", db_seller_tablename)
		var sellerid int
		db.QueryScan(sql, []interface{}{reqdata.SellerId}, &sellerid)
		if ctx.RespErrString(sellerid == 0, &errcode, "运营商不存在") {
			return
		}
	}
	sql = "insert into admin_role(RoleName,SellerId,ParentSellerId,Parent,RoleData)values(?,?,?,?,?)"
	param := []interface{}{reqdata.RoleName, reqdata.SellerId, reqdata.ParentSellerId, reqdata.Parent, reqdata.RoleData}
	err = db.QueryNoResult(sql, param...)
	if ctx.RespErr(err, &errcode) {
		logs.Error(err)
		return
	}
	WriteAdminLog("添加角色", ctx, reqdata)
	ctx.RespOK()
}
func role_delete(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId int    `validate:"required"`
		RoleName string `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "角色管理", "删"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	sql := "select id,Parent from admin_role where ParentSellerId = ? and Parent = ?"
	var id int
	var parent string
	db.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &id, &parent)
	if ctx.RespErrString(id > 0, &errcode, "该角色有下级角色,不可删除") {
		return
	}
	if ctx.RespErrString(parent == "god", &errcode, "该角色不可删除") {
		return
	}
	id = 0
	sql = "select id from admin_user where RoleSellerId = ? and RoleName = ?"
	db.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &id)
	if ctx.RespErrString(id > 0, &errcode, "该角色下存在账号,不可删除") {
		return
	}
	sql = "delete from admin_role where SellerId = ? and RoleName = ?"
	db.QueryNoResult(sql, reqdata.SellerId, reqdata.RoleName)
	WriteAdminLog("删除角色", ctx, reqdata)
	ctx.RespOK()
}
func opt_log(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		SellerId int
		Account  string
		Opt      string
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize == 0 {
		reqdata.PageSize = 10
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "操作日志", "查"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	where := abugo.AbuWhere{}
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	where.AddString("and", "Account", reqdata.Account, "")
	where.AddString("and", "Opt", reqdata.Opt, "")
	var total int
	db.QueryScan(where.CountSql("admin_opt_log"), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := db.Conn().Query(where.Sql("admin_opt_log", reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Id         int
		Account    string
		SellerId   int
		Ip         string
		Opt        string
		Data       string
		CreateTime string
	}
	data := []ReturnData{}
	for dbresult.Next() {
		data_element := ReturnData{}
		abugo.GetDbResult(dbresult, &data_element)
		data_element.CreateTime = abugo.TimeToUtc(data_element.CreateTime)
		data = append(data, data_element)
	}
	dbresult.Close()
	ctx.Put("data", data)
	ctx.Put("page", reqdata.Page)
	ctx.Put("pagesize", reqdata.PageSize)
	ctx.Put("total", total)
	ctx.RespOK()
}
func user_list(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		Account  string
		SellerId int
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize == 0 {
		reqdata.PageSize = 10
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "账号管理", "查"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	where := abugo.AbuWhere{}
	where.OrderBy = "ASC"
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	where.AddString("and", "Account", reqdata.Account, "")
	var total int
	db.QueryScan(where.CountSql("admin_user"), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := db.Conn().Query(where.Sql("admin_user", reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Id           int
		Account      string
		SellerId     int
		RoleSellerId int
		RoleName     string
		Remark       string
		State        int
		LoginCount   int
		LoginIp      string
		LoginTime    string
		CreateTime   string
	}
	data := []ReturnData{}
	for dbresult.Next() {
		data_element := ReturnData{}
		abugo.GetDbResult(dbresult, &data_element)
		data_element.CreateTime = abugo.TimeToUtc(data_element.CreateTime)
		data_element.LoginTime = abugo.TimeToUtc(data_element.LoginTime)
		data = append(data, data_element)
	}
	dbresult.Close()
	ctx.Put("data", data)
	ctx.Put("page", reqdata.Page)
	ctx.Put("pagesize", reqdata.PageSize)
	ctx.Put("total", total)
	ctx.RespOK()
}
func user_modify(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Account      string `validate:"required"`
		SellerId     int    `validate:"required"`
		Password     string
		RoleSellerId int    `validate:"required"`
		RoleName     string `validate:"required"`
		State        int    `validate:"required"`
		Remark       string
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "账号管理", "改"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	if ctx.RespErrString(reqdata.RoleSellerId != -1 && reqdata.RoleSellerId != reqdata.SellerId, &errcode, "运营商不正确") {
		return
	}
	sql := "select id from admin_role  where SellerId = ? and RoleName = ?"
	var rid int
	db.QueryScan(sql, []interface{}{reqdata.RoleSellerId, reqdata.RoleName}, &rid)
	if ctx.RespErrString(rid == 0, &errcode, "角色不存在") {
		return
	}
	if len(reqdata.Password) > 0 {
		sql = "update admin_user set RoleSellerId = ?,RoleName = ?,State = ?,Remark = ?,`Password` = ? where Account = ? and SellerId = ?"
		db.QueryNoResult(sql, reqdata.RoleSellerId, reqdata.RoleName, reqdata.State, reqdata.Remark, reqdata.Password, reqdata.Account, reqdata.SellerId)
	} else {
		sql = "update admin_user set RoleSellerId = ?,RoleName = ?,State = ?,Remark = ? where Account = ? and SellerId = ?"
		db.QueryNoResult(sql, reqdata.RoleSellerId, reqdata.RoleName, reqdata.State, reqdata.Remark, reqdata.Account, reqdata.SellerId)
	}
	if reqdata.State != 1 {
		sql = "select Token from admin_user where Account = ? and SellerId = ? "
		var tokenstr string
		db.QueryScan(sql, []interface{}{reqdata.Account, reqdata.SellerId}, &tokenstr)
		if len(tokenstr) > 0 {
			http.DelToken(tokenstr)
		}
	}
	WriteAdminLog("修改管理员", ctx, reqdata)
	ctx.RespOK()
}
func user_delete(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Id       int    `validate:"required"`
		Account  string `validate:"required"`
		SellerId int    `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "账号管理", "改"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	sql := "delete from admin_user where Id = ? and Account = ? and SellerId = ?"
	db.QueryNoResult(sql, reqdata.Id, reqdata.Account, reqdata.SellerId)
	WriteAdminLog("删除管理员", ctx, reqdata)
	ctx.RespOK()
}
func user_add(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Account      string `validate:"required"`
		SellerId     int    `validate:"required"`
		Password     string `validate:"required"`
		RoleSellerId int    `validate:"required"`
		RoleName     string `validate:"required"`
		State        int    `validate:"required"`
		Remark       string
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "账号管理", "增"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	if ctx.RespErrString(reqdata.RoleSellerId != -1 && reqdata.RoleSellerId != reqdata.SellerId, &errcode, "运营商不正确") {
		return
	}
	sql := "select id from admin_role  where SellerId = ? and RoleName = ?"
	var rid int
	db.QueryScan(sql, []interface{}{reqdata.RoleSellerId, reqdata.RoleName}, &rid)
	if ctx.RespErrString(rid == 0, &errcode, "角色不存在") {
		return
	}
	sql = "select id from admin_user where Account = ? and SellerId = ?"
	var uid int
	db.QueryScan(sql, []interface{}{reqdata.Account, reqdata.SellerId}, &uid)
	if ctx.RespErrString(uid > 0, &errcode, "账号已经存在") {
		return
	}
	sql = "insert into admin_user(Account,Password,SellerId,RoleSellerId,RoleName,State)values(?,?,?,?,?,?)"
	db.QueryNoResult(sql, reqdata.Account, reqdata.Password, reqdata.SellerId, reqdata.RoleSellerId, reqdata.RoleName, reqdata.State)
	WriteAdminLog("添加管理员", ctx, reqdata)
	ctx.RespOK()
}

func user_google(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Account  string `validate:"required"`
		SellerId int    `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	token := GetToken(ctx)
	if ctx.RespErrString(!Auth2(token, "系统管理", "账号管理", "改"), &errcode, "权限不足") {
		return
	}
	if ctx.RespErrString(token.SellerId > 0 && reqdata.SellerId != token.SellerId, &errcode, "运营商不正确") {
		return
	}
	verifykey := abugo.GetGoogleSecret()
	verifyurl := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=abugo", reqdata.Account, verifykey)
	sql := "update admin_user set GoogleSecret = ? where Account = ? and SellerId = ?"
	db.QueryNoResult(sql, verifykey, reqdata.Account, reqdata.SellerId)
	ctx.RespOK(verifyurl)
}

func seller_name(ctx *abugo.AbuHttpContent) {
	token := GetToken(ctx)
	if token.SellerId != -1 {
		ctx.Put("data", []interface{}{})
		ctx.RespOK()
		return
	}
	errcode := 0
	sql := fmt.Sprintf("select * from %sseller where State = 1", DbPrefix)
	dbresult, err := db.Conn().Query(sql)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		SellerId   int
		SellerName string
	}
	data := []ReturnData{}
	data = append(data, ReturnData{0, "全部"})
	data = append(data, ReturnData{-1, "总后台"})
	for dbresult.Next() {
		data_element := ReturnData{}
		abugo.GetDbResult(dbresult, &data_element)
		data = append(data, data_element)
	}
	dbresult.Close()
	ctx.RespOK(data)
}

func seller_flush(ctx *abugo.AbuHttpContent) {
	sql := fmt.Sprintf("select * from %sseller", DbPrefix)
	dbresult, err := db.Conn().Query(sql)
	errcode := 0
	if ctx != nil {
		if ctx.RespErr(err, &errcode) {
			return
		}
	} else {
		if err != nil {
			logs.Error(err)
			return
		}
	}
	for dbresult.Next() {
		data := CacheSellerData{}
		abugo.GetDbResult(dbresult, &data)
		if data.State == 1 {
			redis.HSet("x_hash:seller", fmt.Sprint(data.SellerId), data)
		} else {
			redis.HDel("x_hash:seller", fmt.Sprint(data.SellerId))
		}
	}
	dbresult.Close()
	if ctx != nil {
		ctx.RespOK()
	}
}

func GetSeller(SellerId int) *CacheSellerData {
	data := redis.HGet("x_hash:seller", fmt.Sprint(SellerId))
	if data != nil {
		r := CacheSellerData{}
		json.Unmarshal(data.([]uint8), &r)
		return &r
	}
	return nil
}
