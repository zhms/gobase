package server

import (
	"fmt"
	"strings"
)

/*
	超级管理员
	账号:admin
	密码:admin
*/
var DbPrefix = "x_"
var db_asset_tablename = fmt.Sprintf("%sasset", DbPrefix)
var db_config_tablename = fmt.Sprintf("%sconfig", DbPrefix)
var db_error_tablename = fmt.Sprintf("%serror", DbPrefix)
var db_seller_tablename = fmt.Sprintf("%sseller", DbPrefix)
var db_user_tablename = fmt.Sprintf("%suser", DbPrefix)
var db_verify_tablename = fmt.Sprintf("%sverify", DbPrefix)
var db_transfer_in_tablename = fmt.Sprintf("%stransfer_in", DbPrefix)
var db_transfer_out_tablename = fmt.Sprintf("%stransfer_out", DbPrefix)
var db_asset_change_reason_tablename = fmt.Sprintf("%sasset_change_reason", DbPrefix)
var db_asset_log_tablename = fmt.Sprintf("%sasset_log", DbPrefix)
var replace_symbol = "2416796325297210"

func replace_sql(sql string) string {
	sql = strings.Replace(sql, "2416796325297210", "`", -1)
	sql = strings.Replace(sql, "ex_asset", db_asset_tablename, -1)
	sql = strings.Replace(sql, "ex_config", db_config_tablename, -1)
	sql = strings.Replace(sql, "ex_error", db_error_tablename, -1)
	sql = strings.Replace(sql, "ex_seller", db_seller_tablename, -1)
	sql = strings.Replace(sql, "ex_user", db_user_tablename, -1)
	sql = strings.Replace(sql, "ex_verify", db_verify_tablename, -1)
	sql = strings.Replace(sql, "ex_transfer_in", db_transfer_in_tablename, -1)
	sql = strings.Replace(sql, "ex_transfer_out", db_transfer_out_tablename, -1)
	sql = strings.Replace(sql, "ex_asset_change_reason", db_asset_change_reason_tablename, -1)
	sql = strings.Replace(sql, "ex_asset_log", db_asset_log_tablename, -1)
	sql = strings.Replace(sql, "ex_", DbPrefix, -1)
	return sql
}

func SetupDatabase() {
	sql := `CREATE TABLE IF NOT EXISTS ex_seller  (
				SellerId int(11) NOT NULL AUTO_INCREMENT COMMENT '运营商',
				SellerName varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '运营名称',
				State int(255) NULL DEFAULT 1 COMMENT '状态 1启用 2禁用',
				Remark varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
			  	ApiPublicKey text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
				ApiPrivateKey text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
				ApiThirdPublicKey text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
				ApiRiskPublicKey text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
				ApiRiskPrivateKey text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
				ApiRiskThirdPublicKey text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
				CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
				PRIMARY KEY (SellerId) USING BTREE
			  ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	sql = `CREATE TABLE IF NOT EXISTS admin_login_log  (
			Id int(11) NOT NULL AUTO_INCREMENT,
			UserId int(11) NULL DEFAULT NULL COMMENT '管理员id',
			SellerId int(11) NULL DEFAULT NULL COMMENT '运营商',
			Account varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '管理员账号',
			Token varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '当次登录token',
			LoginIp varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录ip',
			CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '登录时间',
			PRIMARY KEY (Id) USING BTREE,
			INDEX Account(Account) USING BTREE,
			INDEX SellerId(SellerId) USING BTREE
		) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	sql = `CREATE TABLE IF NOT EXISTS admin_opt_log  (
			Id int(11) NOT NULL AUTO_INCREMENT,
			Account varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '操作账号',
			SellerId int(11) NOT NULL DEFAULT -1 COMMENT '账号所属运营商',
			Opt varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '操作类型',
			Ip varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作ip',
			Token varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '请求token',
			Data text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '请求数据',
			CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
			PRIMARY KEY (Id) USING BTREE,
			INDEX Account(Account) USING BTREE,
			INDEX Opt(Opt) USING BTREE
		) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	sql = `CREATE TABLE IF NOT EXISTS admin_role  (
			Id int(11) NOT NULL AUTO_INCREMENT,
			RoleName varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
			SellerId int(11) NOT NULL,
			Parent varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '上级角色',
			RoleData text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色数据',
			CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
			PRIMARY KEY (RoleName, SellerId) USING BTREE,
			UNIQUE INDEX id(Id) USING BTREE
		) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	sql = "INSERT IGNORE INTO `admin_role` VALUES (1, '超级管理员', -1, 'god', '{}', now());"
	db.QueryNoResult(sql)
	sql = `CREATE TABLE IF NOT EXISTS admin_user  (
			Id int(11) NOT NULL AUTO_INCREMENT,
			Account varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '账号',
			Password varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
			SellerId int(11) NOT NULL COMMENT '运营商',
			RoleName varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色名',
			State int(255) NULL DEFAULT 1 COMMENT '状态 1启用 2禁用',
			Token varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'token',
			GoogleSecret varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '谷歌验证码',
			Remark varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
			LoginCount int(255) NULL DEFAULT 0 COMMENT '登录次数',
			LoginTime datetime(0) NULL DEFAULT NULL COMMENT '最后登录时间',
			LoginIp varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '最后登录Ip',
			CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
			PRIMARY KEY (Account) USING BTREE,
			UNIQUE INDEX Id(Id) USING BTREE
		) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	sql = `INSERT IGNORE INTO admin_user VALUES (1, 'admin', '21232f297a57a5a743894a0e4a801fc3', -1, -1, '超级管理员', 1, '', '', '超级管理员,不可删除,编辑', 0, now(), '', now());`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
}
