package server

import (
	"fmt"
	"strings"
)

var db_asset_tablename string
var db_config_tablename string
var db_error_tablename string
var db_seller_tablename string
var db_user_tablename string
var db_verify_tablename string
var db_transfer_in_tablename string
var db_transfer_out_tablename string
var db_asset_change_reason_tablename string
var db_asset_log_tablename string
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
	sql = strings.Replace(sql, "ex_", DbPrefix(), -1)
	return sql
}

func SetupDatabase() {
	db_asset_tablename = fmt.Sprintf("%sasset", DbPrefix())
	db_config_tablename = fmt.Sprintf("%sconfig", DbPrefix())
	db_error_tablename = fmt.Sprintf("%serror", DbPrefix())
	db_seller_tablename = fmt.Sprintf("%sseller", DbPrefix())
	db_user_tablename = fmt.Sprintf("%suser", DbPrefix())
	db_verify_tablename = fmt.Sprintf("%sverify", DbPrefix())
	db_transfer_in_tablename = fmt.Sprintf("%stransfer_in", DbPrefix())
	db_transfer_out_tablename = fmt.Sprintf("%stransfer_out", DbPrefix())
	db_asset_change_reason_tablename = fmt.Sprintf("%sasset_change_reason", DbPrefix())
	db_asset_log_tablename = fmt.Sprintf("%sasset_log", DbPrefix())
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql := `CREATE TABLE IF NOT EXISTS ex_user  (
			Id int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
			UserId int(11) NOT NULL COMMENT '玩家',
			SellerId int(11) NULL DEFAULT NULL COMMENT '运营商',
			Account varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '账号',
			Password varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录密码',
			ThirdId varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '第三方id',
			Email varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '绑定邮箱',
			NickName varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '昵称',
			PhoneNum varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '绑定手机',
			Token varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录token',
			Agents varchar(10240) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'json数组,所有上级id,下标0是直属代理,越往后,代理等级越高',
			TopAgentId int(11) NULL DEFAULT NULL COMMENT '顶级代理id',
			AgentId int(11) NULL DEFAULT NULL COMMENT '直属代理',
			RegisterIp varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '注册ip',
			RegisterTime datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '注册时间',
			PRIMARY KEY (UserId) USING BTREE,
			UNIQUE INDEX Account(Account,SellerId) USING BTREE,
			INDEX ThirdId(ThirdId) USING BTREE,
			INDEX Id(Id) USING BTREE
		  ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE TABLE IF NOT EXISTS ex_seller  (
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
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE TABLE IF NOT EXISTS ex_verify  (
			Account varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '账号',
			SellerId int(11) NOT NULL COMMENT '运营商',
			UseType int(255) NOT NULL COMMENT '使用途径 1注册 2登录',
			VerifyCode varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '验证码',
			CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0),
			PRIMARY KEY (Account, SellerId,UseType) USING BTREE
		  ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	db.QueryNoResult(sql)
	sql = `CREATE TABLE IF NOT EXISTS ex_config  (
			SellerId int(11) NOT NULL COMMENT '运营商',
			ConfigName varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配置名称',
			ConfigValue varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '配置值',
			Remark varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '注释',
			PRIMARY KEY (SellerId, ConfigName) USING BTREE
		  ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	//
	sql = `INSERT IGNORE INTO ex_config VALUES (1, 'SystemOpen', '1', '系统是否开放 1开放 2关闭');`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	//
	sql = `INSERT IGNORE INTO ex_config VALUES (1, 'Verify', '0', '是否开启验证码 1开启 2关闭');`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE TABLE IF NOT EXISTS ex_asset  (
				UserId int(11) NOT NULL COMMENT '玩家',
				SellerId int(11) NOT NULL COMMENT '运营商',
				AssetType int(11) NOT NULL COMMENT '钱包类型 1平台 ',
				Symbol varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '币种',
				AssetAmt bigint(20) NOT NULL DEFAULT 0 COMMENT '钱包余额',
				FrozenAmt bigint(20) NOT NULL DEFAULT 0 COMMENT '冻结余额',
				CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
				PRIMARY KEY (UserId, AssetType, Symbol) USING BTREE
			  ) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE TABLE IF NOT EXISTS ex_asset_change_reason  (
			Id int(11) NOT NULL COMMENT 'id',
			Description varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '描述',
			CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
			PRIMARY KEY (Id) USING BTREE
		   ) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)

	sql = "INSERT IGNORE INTO `x_asset_change_reason` VALUES (1, '钱包转入', now());"
	db.QueryNoResult(sql)

	sql = "INSERT IGNORE INTO `x_asset_change_reason` VALUES (2, '钱包转出', now());"
	db.QueryNoResult(sql)

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE TABLE IF NOT EXISTS ex_asset_log  (
			Id int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
			UserId int(11) NOT NULL COMMENT '玩家id',
			BeforeAmount bigint(255) NOT NULL COMMENT '变化前',
			ChangeAmount bigint(255) NOT NULL COMMENT '变化值',
			AfterAmount bigint(255) NOT NULL COMMENT '变化后',
			Reason int(255) NOT NULL COMMENT '变化原因',
			Extra varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '额外信息',
			CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '变化时间',
			PRIMARY KEY (Id) USING BTREE,
			INDEX UserId(UserId) USING BTREE
		) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE TABLE IF NOT EXISTS ex_error  (
			Id bigint(11) NOT NULL AUTO_INCREMENT,
			FunName varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
			ErrCode int(255) NOT NULL,
			ErrMsg varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
			CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
			PRIMARY KEY (Id) USING BTREE
			) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE  PROCEDURE ex_db_get_userid(OUT p_UserId INT)
BEGIN
	SET p_UserId = 0;
	#创建玩家Id
	SET @whilecount = 0;
	SET @UserId = NULL;
	WHILE @whilecount < 10 AND @UserId IS NULL DO
		SET @whilecount = @whilecount + 1;
		SET @tmpid = 0;
		SELECT FLOOR(10000000 + RAND() * (99999999 - 10000000)) INTO @tmpid;
		IF NOT EXISTS(SELECT UserId FROM ex_user WHERE UserId = @tmpid) THEN
			SET @UserId = @tmpid;
		END IF;
	END WHILE;
	IF @UserId IS NULL THEN
		SET @UserId = 0;
	END IF;
	SET p_UserId = @UserId;
END`
	sql = replace_sql(sql)
	_, err := Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE PROCEDURE ex_db_verify(p_Account VARCHAR(64),p_SellerId INT,p_UseType INT,p_VerifyCode VARCHAR(64),OUT p_Result INT)
proc:BEGIN
/*
	验证码验证
	返回值:
		0:成功
		1:验证码不存在
		2:验证码已过期
		3:验证码不正确
*/
	SET p_Result = 0;
	SET @Verify = NULL;
	SELECT ConfigValue INTO @Verify FROM ex_config WHERE SellerId = p_SellerId AND ConfigName = 'Verify';
	IF @Verify <> '1' THEN
		LEAVE proc;
	END IF;
	SET @VerifyCode = NULL;
	SET @CreateTime = NULL;
	SELECT VerifyCode,CreateTime INTO @VerifyCode,@CreateTime FROM ex_verify WHERE Account = p_Account AND SellerId = p_SellerId AND UseType = p_UseType;
	IF ROW_COUNT() = 0 THEN
		SET p_Result = 1;
		LEAVE proc;
	END IF;
	IF DATE_ADD(@CreateTime, interval 10 MINUTE) < NOW() THEN
		SET p_Result = 2;
		LEAVE proc;
	END IF;
	IF @VerifyCode <> p_VerifyCode THEN
		SET p_Result = 3;
		LEAVE proc;
	END IF;
	DELETE FROM ex_verify  WHERE Account = p_Account AND UseType = p_UseType;
END`
	sql = replace_sql(sql)
	_, err = Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE PROCEDURE ex_api_user_register(p_Account VARCHAR(64),p_SellerId INT,p_Password VARCHAR(64),p_VerifyCode VARCHAR(10),p_ExtraData VARCHAR(10240))
proc:BEGIN
	##############################################################################################
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET CURRENT DIAGNOSTICS CONDITION 1	@errcode = MYSQL_ERRNO, @errmsg = MESSAGE_TEXT;
		ROLLBACK;
		INSERT INTO ex_error(FunName,ErrCode,ErrMsg)VALUES('ex_api_user_register',@errcode,@errmsg);
		SELECT @errcode AS errcode,@errmsg AS errmsg;
	END;
	##############################################################################################
	SET @ErrCode = 10;
	SET @SystemOpen = NULL;
	SELECT ConfigValue INTO @SystemOpen FROM ex_config WHERE SellerId = p_SellerId AND ConfigName = 'SystemOpen';
	IF @SystemOpen <> '1' THEN
		SELECT @ErrCode AS errcode,'系统维护,请稍后再试' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	SET @SellerState = NULL;
	SELECT State INTO @SellerState FROM ex_seller WHERE SellerId = p_SellerId;
	IF ROW_COUNT() = 0 THEN
		SELECT @ErrCode AS errcode,'运营商不存在' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF @SellerState <> 1 THEN
		SELECT @ErrCode AS errcode,'运营商已被禁用' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	SET @Email = NULL;
	SET @PhoneNum = NULL;
	#IF LOCATE('@',p_Account) > 0 THEN
	#	SET @Email = p_Account;
	#ELSE
	#	SET @PhoneNum = p_Account;
	#END IF;
	#DELETE FROM ex_user WHERE Account = p_Account;
	IF EXISTS(SELECT UserId FROM ex_user WHERE Account = p_Account) THEN
		SELECT @ErrCode AS errcode,'账号已经存在' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF p_ExtraData = NULL OR LENGTH(p_ExtraData) = 0 THEN
		SET p_ExtraData = '{}';
	END IF;
	SET @Ip = JSON_UNQUOTE(JSON_EXTRACT(p_ExtraData,'$.Ip'));
	SET @AgentId = CAST(JSON_UNQUOTE(JSON_EXTRACT(p_ExtraData,'$.AgentId')) AS UNSIGNED);
	SET @Agents = NULL;
	SET @TopAgent = NULL;
	IF @AgentId > 0 THEN
		SELECT Agents INTO @Agents FROM ex_user WHERE UserId = @AgentId;
		IF FOUND_ROWS() = 0 THEN
			SELECT @ErrCode AS errcode,'代理不存在' AS errmsg;
			LEAVE proc;
		END IF;
		IF @Agents IS NULL THEN
			SET @Agents = '[]';
		END IF;
		SET @Agents = JSON_ARRAY_INSERT(@Agents, '$[0]', @AgentId);
		SET @TopAgent = JSON_EXTRACT(@Agents, CONCAT('$[',JSON_LENGTH(@Agents) - 1,']'));
	END IF;
	SET @VerifyResult = 0;
	CALL ex_db_verify(p_Account,p_SellerId,1,p_VerifyCode,@VerifyResult);
	IF @VerifyResult = 1 THEN
		SELECT @ErrCode + 0 AS errcode, '验证码不存在' AS errmsg;
		LEAVE proc;
	ELSEIF @VerifyResult = 2 THEN
		SELECT @ErrCode + 1 AS errcode, '验证码已过期' AS errmsg;
		LEAVE proc;
	ELSEIF @VerifyResult = 3 THEN
		SELECT @ErrCode + 2 AS errcode, '验证码不正确' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 3;
	SET @UserId = 0;
	CALL ex_db_get_userid(@UserId);
	IF @UserId = 0 THEN
		SELECT @ErrCode AS errcode,'分配玩家Id失败' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	START TRANSACTION;
		INSERT INTO ex_user(UserId,SellerId,Account,2416796325297210Password2416796325297210,Email,PhoneNum,NickName,RegisterIp,Agents,TopAgentId,AgentId)
		VALUES(@UserId,p_SellerId,p_Account,p_Password,@Email,@PhoneNum,CONCAT(@UserId),@Ip,@Agents,@TopAgent,@AgentId);
		SET @ChildLevel = JSON_LENGTH(@Agents) - 1;
		WHILE @ChildLevel >= 0 DO
			SET @Parentid = JSON_EXTRACT(@Agents, CONCAT('$[',@ChildLevel,']'));
			INSERT INTO ex_agent_child(UserId,Child,ChildLevel)VALUES(@Parentid,@UserId,@ChildLevel);
			SET @ChildLevel = @ChildLevel - 1;
		END WHILE;
	COMMIT;
	SELECT @UserId AS UserId;
END`
	sql = replace_sql(sql)
	_, err = Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE PROCEDURE ex_api_user_login_verifycode(p_Account VARCHAR(64),p_SellerId INT,p_Password VARCHAR(64),p_VerifyCode VARCHAR(10),p_ExtraData VARCHAR(10240))
proc:BEGIN
	##############################################################################################
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET CURRENT DIAGNOSTICS CONDITION 1	@errcode = MYSQL_ERRNO, @errmsg = MESSAGE_TEXT;
		ROLLBACK;
		INSERT INTO ex_error(FunName,ErrCode,ErrMsg)VALUES('ex_api_user_login_verifycode',@errcode,@errmsg);
		SELECT @errcode AS errcode,@errmsg AS errmsg;
	END;
	##############################################################################################
	SET @ErrCode = 10;
	SET @2416796325297210Password2416796325297210 = NULL;
	SET @OldToken = NULL;
	SET @UserId = NULL;
	SELECT UserId,2416796325297210Password2416796325297210,Token INTO @UserId,@2416796325297210Password2416796325297210,@OldToken FROM ex_user WHERE Account = p_Account AND SellerId = p_SellerId;
	IF ROW_COUNT() = 0 THEN
		SELECT @ErrCode AS errcode,'账号不存在' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF @2416796325297210Password2416796325297210 <> p_Password THEN
		SELECT @ErrCode AS errcode,'密码不正确' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	SET @VerifyResult = 0;
	CALL ex_db_verify(p_Account,p_SellerId,2,p_VerifyCode,@VerifyResult);
	IF @VerifyResult = 1 THEN
		SELECT @ErrCode + 1 AS errcode, '验证码不存在' AS errmsg;
		LEAVE proc;
	ELSEIF @VerifyResult = 2 THEN
		SELECT @ErrCode + 2 AS errcode, '验证码已过期' AS errmsg;
		LEAVE proc;
	ELSEIF @VerifyResult = 3 THEN
		SELECT @ErrCode AS errcode, '验证码不正确' AS errmsg;
		LEAVE proc;
	END IF;
	SET @NewToken = UUID();
	UPDATE ex_user SET Token = @NewToken WHERE UserId = @UserId;
	SELECT @UserId AS UserId,p_SellerId AS SellerId,@OldToken AS OldToken,@NewToken AS NewToken;
END`
	sql = replace_sql(sql)
	_, err = Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE PROCEDURE ex_api_user_login_password(p_Account VARCHAR(64),p_SellerId INT,p_Password VARCHAR(64),p_ExtraData VARCHAR(10240))
proc:BEGIN
	##############################################################################################
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET CURRENT DIAGNOSTICS CONDITION 1	@errcode = MYSQL_ERRNO, @errmsg = MESSAGE_TEXT;
		ROLLBACK;
		INSERT INTO ex_error(FunName,ErrCode,ErrMsg)VALUES('ex_api_user_login_password',@errcode,@errmsg);
		SELECT @errcode AS errcode,@errmsg AS errmsg;
	END;
	##############################################################################################
	SET @ErrCode = 10;
	SET @2416796325297210Password2416796325297210 = NULL;
	SELECT 2416796325297210Password2416796325297210 INTO @2416796325297210Password2416796325297210 FROM ex_user WHERE Account = p_Account AND SellerId = p_SellerId;
	IF ROW_COUNT() = 0 THEN
		SELECT @ErrCode AS errcode,'账号不存在' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF @2416796325297210Password2416796325297210 <> p_Password THEN
		SELECT @ErrCode AS errcode,'密码不正确' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
END`
	sql = replace_sql(sql)
	_, err = Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql = `CREATE PROCEDURE ex_db_asset_alter(p_UserId INT,p_SellerId INT,p_AssetType INT,p_Symbol VARCHAR(32),p_Amount BIGINT,p_Reason INT,p_Memo VARCHAR(1024),OUT p_ResultAmount BIGINT)
proc:BEGIN
	IF p_Amount = 0 THEN
		LEAVE proc;
	END IF;
	SET @AssetAmt = NULL;
	SELECT AssetAmt INTO @AssetAmt from ex_asset WHERE UserId = p_UserId AND AssetType = p_AssetType AND Symbol = p_Symbol FOR UPDATE;
	IF FOUND_ROWS() = 0 THEN
		SET @AssetAmt = 0;
		INSERT INTO ex_asset(UserId,SellerId,AssetType,Symbol,AssetAmt,FrozenAmt)
		VALUES(p_UserId,p_SellerId,p_AssetType,p_Symbol,0,0);
	END IF;
	UPDATE ex_asset SET AssetAmt = AssetAmt + p_Amount WHERE UserId = p_UserId AND AssetType = p_AssetType AND Symbol = p_Symbol;
	INSERT INTO ex_asset_log(UserId,BeforeAmount,ChangeAmount,AfterAmount,Reason,Extra)VALUES(p_UserId,@AssetAmt,p_Amount,@AssetAmt + p_Amount,p_Reason,p_Memo);
	SET p_ResultAmount = @AssetAmt + p_Amount;
END`
	sql = replace_sql(sql)
	_, err = Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
}
