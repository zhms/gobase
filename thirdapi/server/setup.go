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
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	sql := `CREATE TABLE IF NOT EXISTS ex_transfer_in  (
		Id bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
		OrderId bigint(20) NULL DEFAULT NULL COMMENT '订单id',
		UserId int(11) NULL DEFAULT NULL COMMENT '玩家id',
		SellerId int(11) NULL DEFAULT NULL COMMENT '玩家所属运营商',
		AssetType int(255) NULL DEFAULT NULL COMMENT '资产类型',
		Symbol varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '币种',
		Side int(255) NULL DEFAULT NULL COMMENT '转入,转出 1转入 2转出',
		Amount bigint(255) NULL DEFAULT NULL COMMENT '划转金额',
		State int(255) NULL DEFAULT NULL COMMENT '订单状态',
		Memo varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备忘录',
		Extra varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '订单额外信息',
		CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '订单创建时间',
		PRIMARY KEY (Id) USING BTREE,
		UNIQUE INDEX OrderId(OrderId) USING BTREE,
		INDEX UserId(UserId) USING BTREE
		) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '其他系统转入转出到本系统的订单' ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)
	sql = `CREATE TABLE IF NOT EXISTS ex_transfer_out  (
			Id bigint(20) NOT NULL AUTO_INCREMENT COMMENT '订单状态 1新建',
			OrderId bigint(255) NULL DEFAULT NULL COMMENT '订单编号',
			UserId int(11) NULL DEFAULT NULL COMMENT '玩家id',
			SellerId int(11) NULL DEFAULT NULL COMMENT '玩家所属运营商',
			AssetType int(255) NULL DEFAULT NULL COMMENT '资产类型',
			Symbol varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '币种',
			Side int(255) NULL DEFAULT NULL COMMENT '1转入 2转出',
			Amount bigint(255) NULL DEFAULT NULL COMMENT '划转金额',
			State int(255) NULL DEFAULT NULL COMMENT '订单状态 1新建 2完成 3失败',
			Memo varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备忘录',
			Extra varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '订单额外信息',
			CreateTime datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '订单创建时间',
			PRIMARY KEY (Id) USING BTREE,
			INDEX UserId(UserId) USING BTREE,
			INDEX OrderId(OrderId) USING BTREE
		  ) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '本系统转入转出到其他系统的订单' ROW_FORMAT = DYNAMIC;`
	sql = replace_sql(sql)
	db.QueryNoResult(sql)

	sql = `CREATE PROCEDURE ex_third_transfer_out_update(p_OrderId BIGINT,p_State INT,p_Extra VARCHAR(1024),p_Reason INT ,p_Memo VARCHAR(1024))
proc:BEGIN
	##############################################################################################
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET CURRENT DIAGNOSTICS CONDITION 1	@errcode = MYSQL_ERRNO, @errmsg = MESSAGE_TEXT;
		ROLLBACK;
		INSERT INTO ex_error(FunName,ErrCode,ErrMsg)VALUES('ex_third_transfer_out_update',@errcode,@errmsg);
		SELECT @errcode AS errcode,@errmsg AS errmsg;
	END;
	##############################################################################################
	SET @ErrCode = 10;
	IF p_State <> 2 OR p_State <> 3 THEN #p_State = 2 成功 p_State = 3 失败
		LEAVE proc;
	END IF;
	START TRANSACTION;
		SET @UserId = NULL;
		SET @AssetType = NULL;
		SET @Symbol = NULL;
		SET @Side = NULL;
		SET @State = NULL;
		SET @Amount = NULL;
		SELECT UserId,AssetType,Symbol,Side,State,Amount INTO @UserId,@AssetType,@Symbol,@Side,@State,@Amount FROM ex_transfer_out WHERE id = p_OrderId FOR UPDATE;
		IF ROW_COUNT() = 0 THEN
			ROLLBACK;
			SELECT @ErrCode AS errcode,'订单不存在' AS errmsg;
			LEAVE proc;
		END IF;
		SET @ErrCode = @ErrCode + 1;
		IF @State <> 1 THEN
			ROLLBACK;
			SELECT @ErrCode AS errcode,'订单状态不正确' AS errmsg;
			LEAVE proc;
		END IF;
		SET @ErrCode = @ErrCode + 1;
		IF p_State = 2 AND @Side = 1 THEN #转入成功,加钱
			CALL ex_db_asset_alter(p_UserId,p_SellerId,@AssetType,@Symbol,@Amount,p_Reason,p_Memo,@AfterAssetAmt);
		END IF;
		IF p_State = 3 AND @Side = 2 THEN #转出失败,退钱
			CALL ex_db_asset_alter(p_UserId,p_SellerId,@AssetType,@Symbol,@Amount,p_Reason,p_Memo,@AfterAssetAmt);
			UPDATE ex_asset SET FrozenAmt = FrozenAmt - p_Amount WHERE UserId = p_UserId AND AssetAmt = p_AssetType AND Symbol = p_Symbol;
		END IF;
		UPDATE ex_transfer_out SET State = p_State WHERE Id = p_OrderId;
	COMMIT;
END`
	sql = replace_sql(sql)
	_, err := Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	sql = `CREATE PROCEDURE ex_third_transfer_out_out(p_UserId INT,p_SellerId INT,p_AssetType INT,p_Symbol VARCHAR(32),p_Amount BIGINT,p_Extra VARCHAR(1024),p_Reason INT,p_Memo VARCHAR(1024))
proc:BEGIN
	##############################################################################################
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET CURRENT DIAGNOSTICS CONDITION 1	@errcode = MYSQL_ERRNO, @errmsg = MESSAGE_TEXT;
		ROLLBACK;
		INSERT INTO ex_error(FunName,ErrCode,ErrMsg)VALUES('ex_third_transfer_out_out',@errcode,@errmsg);
		SELECT @errcode AS errcode,@errmsg AS errmsg;
	END;
	##############################################################################################
	SET @ErrCode = 10;
	IF p_Amount <= 0 THEN
		SELECT @ErrCode AS errcode,"参数错误" AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF NOT EXISTS(SELECT Id FROM ex_asset_change_reason WHERE Id = p_Reason) THEN
		SELECT @ErrCode AS errcode,"变化原因不存在" AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	SET @OrderId = NULL;
	SET @whilecount = 0;
	WHILE @whilecount < 10 AND @OrderId IS NULL DO
		SET @whilecount = @whilecount + 1;
		SET @tmpid = FLOOR(100000000000000 + RAND() * (999999999999999 - 100000000000000));
		IF NOT EXISTS(SELECT id FROM ex_transfer_out WHERE OrderId = @tmpid) THEN
			SET @OrderId = @tmpid;
			SET @whilecount = 100;
		END IF;
	END WHILE;
	IF @OrderId IS NULL THEN
		SELECT @ErrCode AS errcode,"分配OrderId失败" AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	START TRANSACTION;
		SET @AfterAssetAmt  = 0;
		CALL ex_db_asset_alter(OrderId,p_UserId,p_SellerId,p_AssetType,p_Symbol,-p_Amount,p_Reason,p_Memo,@AfterAssetAmt);
		IF @AfterAssetAmt < 0 THEN
			ROLLBACK;
			SELECT @ErrCode AS errcode,"资产不足" AS errmsg;
			LEAVE proc;
		END IF;
		SET @ErrCode = @ErrCode + 1;
		UPDATE ex_asset SET FrozenAmt = FrozenAmt + p_Amount WHERE UserId = p_UserId AND AssetType = p_AssetType AND Symbol = p_Symbol;
		INSERT INTO ex_transfer_out(UserId,SellerId,Symbol,Side,AssetType,Amount,State,Extra,Memo)
		VALUES(@OrderId,p_UserId,p_SellerId,p_Symbol,2,p_AssetType,p_Amount,1,p_Extra,"新建订单");
		SELECT @OrderId AS OrderId;
	COMMIT;
END`
	sql = replace_sql(sql)
	_, err = Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	sql = `CREATE PROCEDURE ex_third_transfer_out_in(p_UserId INT,p_SellerId INT,p_AssetType INT,p_Symbol VARCHAR(32),p_Amount BIGINT,p_Extra VARCHAR(1024),p_Reason INT,p_Memo VARCHAR(1024))
proc:BEGIN
	##############################################################################################
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET CURRENT DIAGNOSTICS CONDITION 1	@errcode = MYSQL_ERRNO, @errmsg = MESSAGE_TEXT;
		ROLLBACK;
		INSERT INTO ex_error(FunName,ErrCode,ErrMsg)VALUES('ex_third_transfer_out_in',@errcode,@errmsg);
		SELECT @errcode AS errcode,@errmsg AS errmsg;
	END;
	##############################################################################################
	SET @ErrCode = 10;
	IF p_Amount <= 0 THEN
		SELECT @ErrCode AS errcode,"参数错误" AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF NOT EXISTS(SELECT Id FROM ex_asset_change_reason WHERE Id = p_Reason) THEN
		SELECT @ErrCode AS errcode,"变化原因不存在" AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF NOT EXISTS(SELECT Id FROM ex_user WHERE UserId = p_UserId) THEN
		SELECT @ErrCode AS errcode,"玩家不存在" AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	START TRANSACTION;
		SET @OrderId = NULL;
		SET @whilecount = 0;
		WHILE @whilecount < 10 AND @OrderId IS NULL DO
			SET @whilecount = @whilecount + 1;
			SET @tmpid = FLOOR(100000000000000 + RAND() * (999999999999999 - 100000000000000));
			IF NOT EXISTS(SELECT id FROM ex_transfer_out WHERE OrderId = @tmpid) THEN
				SET @OrderId = @tmpid;
				SET @whilecount = 100;
			END IF;
		END WHILE;
		IF @OrderId IS NULL THEN
			SELECT @ErrCode AS errcode,"分配OrderId失败" AS errmsg;
			LEAVE proc;
		END IF;
		SET @ErrCode = @ErrCode + 1;
		INSERT INTO ex_transfer_out(OrderId,UserId,SellerId,Symbol,Side,AssetType,Amount,State,Extra,Memo)
		VALUES(@OrderId,p_UserId,p_SellerId,p_Symbol,1,p_AssetType,p_Amount,1,p_Extra,"新建订单");
		SELECT @OrderId AS OrderId;
	COMMIT;
END`
	sql = replace_sql(sql)
	_, err = Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	sql = `CREATE PROCEDURE ex_third_transfer_in_out(p_OrderId BIGINT,p_UserId INT,p_SellerId INT,p_AssetType INT,p_Symbol VARCHAR(32),p_Amount BIGINT,p_Extra VARCHAR(1024),p_Reason INT,p_Memo VARCHAR(1024))
proc:BEGIN
	##############################################################################################
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET CURRENT DIAGNOSTICS CONDITION 1	@errcode = MYSQL_ERRNO, @errmsg = MESSAGE_TEXT;
		ROLLBACK;
		INSERT INTO ex_error(FunName,ErrCode,ErrMsg)VALUES('ex_third_transfer_in_out',@errcode,@errmsg);
		SELECT @errcode AS errcode,@errmsg AS errmsg;
	END;
	##############################################################################################
	SET @ErrCode = 10;
	IF EXISTS (SELECT id FROM ex_transfer_in WHERE OrderId = p_OrderId) THEN
		SELECT @ErrCode AS errcode,"订单已存在" AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	INSERT INTO ex_transfer_in(OrderId,UserId,SellerId,AssetType,Symbol,Amount,Side,State,Extra,Memo)
	VALUES(p_OrderId,p_UserId,p_SellerId,p_AssetType,p_Symbol,p_Amount,2,1,p_Extra,p_Memo);
	IF p_Amount <= 0 THEN
		SELECT @ErrCode AS errcode,"参数错误" AS errmsg;
		UPDATE ex_transfer_in SET State = 3,Memo = '订单失败,参数错误' WHERE OrderId = p_OrderId;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF NOT EXISTS(SELECT Id FROM ex_asset_change_reason WHERE Id = p_Reason) THEN
		SELECT @ErrCode AS errcode,"变化原因不存在" AS errmsg;
		UPDATE ex_transfer_in SET State = 3,Memo = '订单失败,变化原因不存在' WHERE OrderId = p_OrderId;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF NOT EXISTS(SELECT Id FROM ex_user WHERE UserId = p_UserId) THEN
		SELECT @ErrCode AS errcode,"玩家不存在" AS errmsg;
		UPDATE ex_transfer_in SET State = 3,Memo = '订单失败,玩家不存在' WHERE OrderId = p_OrderId;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	START TRANSACTION;
		SET @AfterAssetAmt = 0;
		CALL ex_db_asset_alter(p_UserId,p_SellerId,p_AssetType,p_Symbol,-p_Amount,p_Reason,p_Memo,@AfterAssetAmt);
		IF @AfterAssetAmt < 0 THEN
			ROLLBACK;
			SELECT @ErrCode AS errcode,"资产不足" AS errmsg;
			UPDATE ex_transfer_in SET State = 3,Memo = '订单失败,资产不足' WHERE OrderId = p_OrderId;
			LEAVE proc;
		END IF;
		UPDATE ex_transfer_in SET State = 2,Memo = '订单成功' WHERE OrderId = p_OrderId;
		SELECT @AfterAssetAmt AS Balance;
	COMMIT;
END`
	sql = replace_sql(sql)
	_, err = Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	sql = `CREATE PROCEDURE ex_third_transfer_in_in(p_OrderId BIGINT,p_UserId INT,p_SellerId INT,p_AssetType INT,p_Symbol VARCHAR(32),p_Amount BIGINT,p_Extra VARCHAR(1024),p_Reason INT,p_Memo VARCHAR(1024))
proc:BEGIN
	##############################################################################################
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET CURRENT DIAGNOSTICS CONDITION 1	@errcode = MYSQL_ERRNO, @errmsg = MESSAGE_TEXT;
		ROLLBACK;
		INSERT INTO ex_error(FunName,ErrCode,ErrMsg)VALUES('ex_third_transfer_in_in',@errcode,@errmsg);
		SELECT @errcode AS errcode,@errmsg AS errmsg;
	END;
	##############################################################################################
	SET @ErrCode = 10;
	IF EXISTS (SELECT id FROM ex_transfer_in WHERE OrderId = p_OrderId) THEN
		SELECT @ErrCode AS errcode,"订单已存在" AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	INSERT INTO ex_transfer_in(OrderId,UserId,SellerId,AssetType,Symbol,Amount,Side,State,Extra,Memo)
	VALUES(p_OrderId,p_UserId,p_SellerId,p_AssetType,p_Symbol,p_Amount,1,1,p_Extra,p_Memo);
	IF p_Amount <= 0 THEN
		SELECT @ErrCode AS errcode,"参数错误" AS errmsg;
		UPDATE ex_transfer_in SET State = 3,Memo = '订单失败,参数错误' WHERE OrderId = p_OrderId;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF NOT EXISTS(SELECT Id FROM ex_asset_change_reason WHERE Id = p_Reason) THEN
		SELECT @ErrCode AS errcode,"变化原因不存在" AS errmsg;
		UPDATE ex_transfer_in SET State = 3,Memo = '订单失败,变化原因不存在' WHERE OrderId = p_OrderId;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF NOT EXISTS(SELECT Id FROM ex_user WHERE UserId = p_UserId) THEN
		SELECT @ErrCode AS errcode,"玩家不存在" AS errmsg;
		UPDATE ex_transfer_in SET State = 3,Memo = '订单失败,玩家不存在' WHERE OrderId = p_OrderId;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	START TRANSACTION;
		SET @AfterAssetAmt = 0;
		CALL ex_db_asset_alter(p_UserId,p_SellerId,p_AssetType,p_Symbol,p_Amount,p_Reason,p_Memo,@AfterAssetAmt);
		UPDATE ex_transfer_in SET State = 2,Memo = '订单成功' WHERE OrderId = p_OrderId;
		SELECT @AfterAssetAmt AS Balance;
	COMMIT;
END`
	sql = replace_sql(sql)
	_, err = Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
	sql = `CREATE PROCEDURE ex_third_register(p_ThirdId VARCHAR(32),p_SellerId INT,p_Password VARCHAR(64),p_ExtraData VARCHAR(10240))
proc:BEGIN
	##############################################################################################
	DECLARE EXIT HANDLER FOR SQLEXCEPTION
	BEGIN
		GET CURRENT DIAGNOSTICS CONDITION 1	@errcode = MYSQL_ERRNO, @errmsg = MESSAGE_TEXT;
		ROLLBACK;
		INSERT INTO ex_error(FunName,ErrCode,ErrMsg)VALUES('ex_third_register',@errcode,@errmsg);
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
	SET @UserId = 0;
	SELECT UserId INTO @UserId FROM ex_user WHERE ThirdId = p_ThirdId AND SellerId = p_SellerId;
	#账号已存在,更新账号密码
	IF FOUND_ROWS() > 0 THEN
		UPDATE ex_user SET 2416796325297210Password2416796325297210 = p_Password WHERE ThirdId = p_ThirdId AND SellerId = p_SellerId;
		SELECT @UserId AS UserId;
		LEAVE proc;
	END IF;
	#账号不存在,创建账号
	CALL ex_db_get_userid(@UserId);
	IF @UserId = 0 THEN
		SELECT @ErrCode AS errcode,'分配玩家Id失败' AS errmsg;
		LEAVE proc;
	END IF;
	SET @ErrCode = @ErrCode + 1;
	IF p_ExtraData = NULL OR LENGTH(p_ExtraData) = 0 THEN
		SET p_ExtraData = '{}';
	END IF;
	INSERT INTO ex_user(UserId,SellerId,Account,ThirdId,2416796325297210Password2416796325297210,NickName)
	VALUES(@UserId,p_SellerId,p_ThirdId,p_ThirdId,p_Password,CONCAT(@UserId));
	SELECT @UserId AS UserId;
END`
	sql = replace_sql(sql)
	_, err = Db().Conn().Exec(sql)
	if err != nil && strings.Index(err.Error(), "1304") <= 0 {
		fmt.Println(err)
	}
}
