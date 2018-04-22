package tsgmysqlutils

/*
 MySQL utils test
 @author Tony Tian
 @date 2018-04-16
 @version 1.0.0
*/

import (
	"errors"
	"github.com/timespacegroup/go-utils"
	"testing"
	"time"
)

/*
	You must execute this SQL in your MySQL database:

	CREATE DATABASE IF NOT EXISTS test DEFAULT CHARACTER SET utf8 COLLATE utf8_bin;
*/

func TestCreateTable(t *testing.T) {

	client := TestDbClient()

	tabSql := tsgutils.NewStringBuilder()
	tabSql.Append("CREATE TABLE `we_test_tab1` (")
	tabSql.Append("`id` int(10) unsigned AUTO_INCREMENT NOT NULL COMMENT 'The primary key id',")
	tabSql.Append("`name` varchar(64) NOT NULL DEFAULT '' COMMENT 'The user name',")
	tabSql.Append("`gender` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'The user gerder, 1:male 2:female 0:default',")
	tabSql.Append("`birthday` date NOT NULL COMMENT 'The user birthday, eg: 2018-04-16',")
	tabSql.Append("`stature` decimal(16, 2) NOT NULL DEFAULT '0.00' COMMENT 'The user stature, eg: 172.22cm',")
	tabSql.Append("`weight` decimal(16, 2) NOT NULL DEFAULT '0.00' COMMENT 'The user weight, eg: 21.77kg',")
	tabSql.Append("`created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',")
	tabSql.Append("`modified_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'record time',")
	tabSql.Append("`is_deleted` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'Logic to delete(0:normal 1:deleted)',")
	tabSql.Append("PRIMARY KEY (`id`),")
	tabSql.Append("UNIQUE KEY `name` (`name`)")
	tabSql.Append(") ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE= utf8_bin COMMENT = 'test table1';")

	client.Exec(tabSql.ToString())

	tabSql = tabSql.Clear()
	tabSql.Append("CREATE TABLE `we_test_tab2` (")
	tabSql.Append("`id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'The primary key id',")
	tabSql.Append("`user_id` int(10) unsigned NOT NULL COMMENT 'The user id',")
	tabSql.Append("`area_code` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT 'The user area code',")
	tabSql.Append("`phone` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'The user phone',")
	tabSql.Append("`email` varchar(35) COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT 'The user email',")
	tabSql.Append("`postcode` mediumint(8) unsigned NOT NULL DEFAULT '0' COMMENT 'The user postcode',")
	tabSql.Append("`administration_code` mediumint(8) unsigned NOT NULL DEFAULT '0' COMMENT 'The user administration code',")
	tabSql.Append("`address` varchar(150) COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT 'The user address',")
	tabSql.Append("`created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',")
	tabSql.Append("`modified_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'modified time',")
	tabSql.Append("`is_deleted` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'Logic to delete(0:normal 1:deleted)',")
	tabSql.Append("PRIMARY KEY (`id`)")
	tabSql.Append(") ENGINE =InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_bin COMMENT ='test table2';")

	_, err := client.Exec(tabSql.ToString())

	if err != nil {
		tsgutils.Stdout("TestCreateTable failed", err)
	} else {
		tsgutils.Stdout("we_test_tab1 and we_test_tab2 tables created.")
	}
	client.CloseConn()
}

func TestGenerateORM(t *testing.T) {
	client := TestDbClient()
	orm := NewORMGenerator(client)
	orm.AddComment = true
	tabNames := []string{"we_test_tab1", "we_test_tab2"}
	orm.DefaultGenerator(tabNames)
	client.CloseConn()
}

func TestInsertSQL(t *testing.T) {

	client := TestDbClient()

	sql := "INSERT INTO we_test_tab1(`name`,`gender`,`birthday`,`stature`,`weight`,`created_time`,`modified_time`,`is_deleted`) VALUES(?,?,?,?,?,?,?,?)"
	params := tsgutils.NewInterfaceBuilder()
	for i := 1; i < 6; i++ {
		params.Clear()
		name := tsgutils.NewStringBuilder().Append("可可").AppendInt(i).ToString()
		params.Append(name).Append(i%2 + 1)
		birthDayStr := tsgutils.NewStringBuilder().Append("199").AppendInt(i).Append("-0").AppendInt(i).Append("-0").AppendInt(i).ToString()
		birthDay, err := tsgutils.StringToTime(birthDayStr, 1)
		tsgutils.CheckAndPrintError("birthDay", err)
		curTime := time.Now()
		stature, err := tsgutils.NewString("17").AppendInt(i).AppendString(".3").AppendInt(i).ToFloat()
		tsgutils.CheckAndPrintError("stature", err)
		weight, err := tsgutils.NewString("4").AppendInt(i).AppendString(".1").AppendInt(i).ToFloat()
		tsgutils.CheckAndPrintError("weight", err)
		params.Append(birthDay).Append(stature).Append(weight)
		params.Append(curTime).Append(curTime).Append(0)

		result, err := client.Exec(sql, params.ToInterfaces()...)
		if err != nil {
			tsgutils.Stdout("Insert failed", err)
		} else {
			tsgutils.Stdout("Insert result: last insert id:", result)
		}
	}
	client.CloseConn()
}

func TestUpdateSQL(t *testing.T) {
	client := TestDbClient()

	sql := "UPDATE we_test_tab1 SET is_deleted = ?,modified_time=NOW() WHERE id > 3;"

	result, err := client.Exec(sql, 1)
	if err != nil {
		tsgutils.Stdout("Update failed", err)
	} else {
		tsgutils.Stdout("Update result: rows affected: ", result)
	}
	client.CloseConn()
}

func TestSelectSQL(t *testing.T) {
	client := TestDbClient()

	sql := "SELECT * FROM we_test_tab1 WHERE id = 1;"
	weTestTab1 := new(WeTestTab1)
	var orm ORMBase = weTestTab1
	_, err := client.QueryRow(orm, sql)
	if err != nil {
		tsgutils.Stdout("Select row failed", err)
	} else {
		tsgutils.Stdout("Select row result: ", tsgutils.StructToJson(weTestTab1))
	}
	client.CloseConn()
}

func TestSelectListSQL(t *testing.T) {
	client := TestDbClient()

	sql := "SELECT * FROM we_test_tab1 WHERE is_deleted <> 1;"
	weTestTab1 := new(WeTestTab1)
	var orm ORMBase = weTestTab1
	_, err := client.QueryList(orm, sql)
	if err != nil {
		tsgutils.Stdout("Select rows failed", err)
	} else {
		tsgutils.Stdout("Select rows result: ", tsgutils.StructToJson(weTestTab1.WeTestTab1s))
	}
	client.CloseConn()
}

func TestSelectAggregateSQL(t *testing.T) {
	client := TestDbClient()

	sql := "SELECT COUNT(*) FROM we_test_tab1 WHERE is_deleted <> 1;"

	result, err := client.QueryAggregate(sql)
	if err != nil {
		tsgutils.Stdout("Select aggregate failed", err)
	} else {
		tsgutils.Stdout("Select aggregate result: ", result)
	}
	client.CloseConn()
}

func TestDeleteSQL(t *testing.T) {
	client := TestDbClient()

	sql := "DELETE FROM we_test_tab1 WHERE id = ?;"

	result, err := client.Exec(sql, 5)
	if err != nil {
		tsgutils.Stdout("Delete failed", err)
	} else {
		tsgutils.Stdout("Delete result: rows affected: ", result)
	}
	client.CloseConn()
}

func TestDbTxSuccessful(t *testing.T) {
	client := TestDbClient()
	sql1 := "INSERT INTO `we_test_tab1` (`name`, `gender`, `birthday`, `stature`, `weight`, `created_time`, `modified_time`, `is_deleted`) VALUES('tony', 2, '1991-01-01', 171.31, 41.11, '2018-04-19 13:20:09', '2018-04-19 13:20:09', 0);"
	sql2 := "INSERT INTO `we_test_tab2` (`user_id`, `area_code`, `phone`, `email`, `postcode`, `administration_code`, `address`, `created_time`, `modified_time`, `is_deleted`) VALUES (?, 86, 18212345678, 'tony@timespace.group', 100089, 110108, '北京市海淀区中关村', '2018-04-21 08:39:14', '2018-04-21 08:39:14', 0);"
	tx, err := client.TxBegin()
	if err != nil {
		tsgutils.CheckAndPrintError("TxBegin failed", err)
		return
	}
	id1, err := client.TxExec(tx, sql1)
	if err != nil {
		client.TxRollback(tx)
		tsgutils.Stdout("TxRollback sql1", id1, err)
		return
	}
	id2, err := client.TxExec(tx, sql2, id1)
	if err != nil {
		client.TxRollback(tx)
		tsgutils.Stdout("TxRollback sql2", id2, err)
		return
	}
	if client.TxCommit(tx) {
		tsgutils.Stdout("TestDbTx Successful: ", id1, id2)
	} else {
		tsgutils.Stdout("TestDbTx TxCommit failed")
	}
	client.CloseConn()
}

func TestDbTxFailure(t *testing.T) {
	// The `we_test_tab2` primary key conflict: Duplicate entry '1' for key 'PRIMARY'
	client := TestDbClient()
	sql1 := "INSERT INTO `we_test_tab1` (`name`, `gender`, `birthday`, `stature`, `weight`, `created_time`, `modified_time`, `is_deleted`) VALUES('tian', 2, '1991-01-01', 171.31, 41.11, '2018-04-19 13:20:09', '2018-04-19 13:20:09', 0);"
	sql2 := "INSERT INTO `we_test_tab2` (`id`,`user_id`, `area_code`, `phone`, `email`, `postcode`, `administration_code`, `address`, `created_time`, `modified_time`, `is_deleted`) VALUES (1, ?, 86, 18212345678, 'tony@timespace.group', 100089, 110108, '北京市海淀区中关村', '2018-04-21 08:39:14', '2018-04-21 08:39:14', 0);"
	tx, err := client.TxBegin()
	if err != nil {
		tsgutils.CheckAndPrintError("TxBegin failed", err)
		return
	}
	id1, err := client.TxExec(tx, sql1)
	if err != nil {
		client.TxRollback(tx)
		tsgutils.Stdout("TxRollback sql1", id1, err)
		return
	}
	id2, err := client.TxExec(tx, sql2, id1)
	if err != nil {
		client.TxRollback(tx)
		tsgutils.Stdout("TxRollback sql2", id2, err)
		return
	}
	if client.TxCommit(tx) {
		tsgutils.Stdout("TestDbTx Successful: ", id1, id2)
	} else {
		tsgutils.Stdout("TestDbTx TxCommit failed")
	}
	client.CloseConn()
}

func TestGenerateORM_Insert_true(t *testing.T) {
	client := TestDbClient()
	weTestTab1 := new(WeTestTab1)
	weTestTab1.Id = 7
	weTestTab1.Name = "tina"
	weTestTab1.Gender = 1
	weTestTab1.Birthday = time.Now()
	weTestTab1.Stature = 60.12
	weTestTab1.Weight = 178.34
	weTestTab1.CreatedTime = time.Now()
	weTestTab1.ModifiedTime = time.Now()
	weTestTab1.IsDeleted = 0
	result, err := weTestTab1.Insert(client, true)
	if err != nil {
		tsgutils.Stdout("Update failed", err)
	} else {
		tsgutils.Stdout("Update orm result: rows affected: ", result)
	}
}

func TestGenerateORM_Insert_false(t *testing.T) {
	client := TestDbClient()
	weTestTab1 := new(WeTestTab1)
	weTestTab1.Name = "boss"
	weTestTab1.Gender = 1
	weTestTab1.Birthday = time.Now()
	weTestTab1.Stature = 60.23
	weTestTab1.Weight = 178.45
	weTestTab1.CreatedTime = time.Now()
	weTestTab1.ModifiedTime = time.Now()
	weTestTab1.IsDeleted = 0
	result, err := weTestTab1.Insert(client, false)
	if err != nil {
		tsgutils.Stdout("Update failed", err)
	} else {
		tsgutils.Stdout("Update orm result: last insert id: ", result)
	}
	client.CloseConn()
}

func TestGenerateORM_BatchInsert_returnIds_false(t *testing.T) {
	client := TestDbClient()
	weTestTab1s := new(WeTestTab1)
	for i := 0; i < 3; i++ {
		var weTestTab1 WeTestTab1
		weTestTab1.Name = tsgutils.NewString("Tony").AppendInt(i).ToString()
		weTestTab1.Gender = 1
		weTestTab1.Birthday = time.Now()
		weTestTab1.Stature = 60.88
		weTestTab1.Weight = 178.55
		weTestTab1.CreatedTime = time.Now()
		weTestTab1.ModifiedTime = time.Now()
		weTestTab1.IsDeleted = 0
		weTestTab1s.WeTestTab1s = append(weTestTab1s.WeTestTab1s, weTestTab1)
	}
	result, err := weTestTab1s.BatchInsert(client, false, false)
	if err != nil {
		tsgutils.Stdout("BatchInsert failed", err)
	} else {
		tsgutils.Stdout("BatchInsert orm result: last insert id: ", tsgutils.StructToJson(result))
	}
	client.CloseConn()
}

func TestGenerateORM_BatchInsert_returnIds_true(t *testing.T) {
	client := TestDbClient()
	weTestTab1s := new(WeTestTab1)
	for i := 100; i < 103; i++ {
		var weTestTab1 WeTestTab1
		weTestTab1.Id = int64(i)
		weTestTab1.Name = tsgutils.NewString("Tony").AppendInt64(weTestTab1.Id).ToString()
		weTestTab1.Gender = 1
		weTestTab1.Birthday = time.Now()
		weTestTab1.Stature = 60.99
		weTestTab1.Weight = 178.66
		weTestTab1.CreatedTime = time.Now()
		weTestTab1.ModifiedTime = time.Now()
		weTestTab1.IsDeleted = 0
		weTestTab1s.WeTestTab1s = append(weTestTab1s.WeTestTab1s, weTestTab1)
	}
	result, err := weTestTab1s.BatchInsert(client, true, true)
	if err != nil {
		tsgutils.Stdout("BatchInsert failed", err)
	} else {
		tsgutils.Stdout("BatchInsert orm result: last insert id: ", tsgutils.StructToJson(result))
		tsgutils.Stdout("BatchInsert orm result: inner last insert id: ", weTestTab1s.WeTestTab1s[0].Id)
	}
	client.CloseConn()
}

func TestGenerateORM_Update(t *testing.T) {
	client := TestDbClient()
	weTestTab1 := new(WeTestTab1)
	var orm ORMBase = weTestTab1
	sql := "SELECT * FROM we_test_tab1 WHERE id = 3;"
	client.QueryRow(orm, sql)
	weTestTab1.Name = "可可^_^"
	result, err := weTestTab1.UpdateWeTestTab1ById(client)
	if err != nil {
		tsgutils.Stdout("Update failed", err)
	} else {
		tsgutils.Stdout("Update orm result: rows affected: ", result)
	}
	client.CloseConn()
}

func TestGenerateORM_Delete(t *testing.T) {
	client := TestDbClient()
	weTestTab1 := new(WeTestTab1)
	weTestTab1.Id = 4
	result, err := weTestTab1.DeleteWeTestTab1ById(client)
	if err != nil {
		tsgutils.Stdout("Delete failed", err)
	} else {
		tsgutils.Stdout("Delete orm result: rows affected: ", result)
	}
	client.CloseConn()
}

func TestTxQuery1(t *testing.T) {
	client := TestDbClient()
	sql := "SELECT name FROM we_test_tab1 WHERE id = 1 FOR UPDATE;"
	tx, _ := client.TxBegin()
	stmt, _ := tx.Prepare(sql)
	row := stmt.QueryRow()
	var name string
	row.Scan(&name)
	tsgutils.Stdout(name)
	client.TxCommit(tx)
	client.CloseConn()
}

func TestTxQuery2(t *testing.T) {
	client := TestDbClient()
	sql := "SELECT name FROM we_test_tab1 WHERE id = ? FOR UPDATE;"
	tx, _ := client.TxBegin()
	row, _ := client.TxQueryRow(tx, nil, sql, 1)
	var name string
	row.Scan(&name)
	tsgutils.Stdout(name)
	client.TxCommit(tx)
	client.CloseConn()
}

func TestPrintLog(t *testing.T) {
	driverName := MySQL
	PrintSlowConn(driverName, "127.0.0.1", "mysql", 5000)
	sql := "SELECT * FROM user WHERE User = ? AND Host = ?"
	params := []interface{}{"root", "127.0.0.1"}
	PrintSlowSql("127.0.0.1", "mysql", 5000, sql, params)
	err := errors.New("test sql error")
	PrintErrorSql(err, sql, params)
}
