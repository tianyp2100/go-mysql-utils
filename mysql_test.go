package tsgmysqlutils

/*
 string utils
 @author Tony Tian
 @date 2018-04-16
 @version 1.0.0
*/

import (
	"testing"
	"errors"
	"github.com/timespacegroup/go-utils"
	"time"
)

func TestPrintLog(t *testing.T) {
	driverName := MySQL
	PrintSlowConn(driverName, "127.0.0.1", "mysql", 5000)
	sql := "SELECT * FROM user WHERE User = ? AND Host = ?"
	params := []interface{}{"root", "127.0.0.1"}
	PrintSlowSql(driverName, "127.0.0.1", "mysql", 5000, sql, params)
	err := errors.New("test sql error")
	PrintErrorSql(driverName, err, sql, params)
}

/*
	You must execute this SQL in your MySQL database:

	CREATE DATABASE IF NOT EXISTS test DEFAULT CHARACTER SET utf8 COLLATE utf8_bin;
 */

func TestCreateTable(t *testing.T) {

	db := TestDbClient()

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

	db.Exec(tabSql.ToString(), nil)

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

	db.Exec(tabSql.ToString(), nil)
}

func TestGenerateORM(t *testing.T) {
	config := TestORMConfig()
	GenerateORM(*config)
}

func TestInsertSQL(t *testing.T) {

	db := TestDbClient()

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

		result := db.Exec(sql, params.ToInterfaces())
		tsgutils.Stdout("Insert result: last insert id:", result)
	}
}

func TestUpdateSQL(t *testing.T) {
	db := TestDbClient()

	sql := "UPDATE we_test_tab1 SET is_deleted = 1,modified_time=NOW() WHERE id > 3;"

	result := db.Exec(sql, nil)
	tsgutils.Stdout("Update result: rows affected: ", result)
}

func TestSelectSQL(t *testing.T) {
	db := TestDbClient()

	sql := "SELECT * FROM we_test_tab1 WHERE id = 1;"
	weTestTab1 := new(WeTestTab1)
	var orm ORMBase = weTestTab1
	db.QueryRow(sql, nil, orm)
	tsgutils.Stdout("Select row result: ", tsgutils.StructToJson(weTestTab1))
}

func TestSelectListSQL(t *testing.T) {
	db := TestDbClient()

	sql := "SELECT * FROM we_test_tab1 WHERE is_deleted <> 1;"
	weTestTab1 := new(WeTestTab1)
	var orm ORMBase = weTestTab1
	db.QueryList(sql, nil, orm)
	tsgutils.Stdout("Select row result: ", tsgutils.StructToJson(weTestTab1.WeTestTab1s))
}

func TestSelectAggregateSQL(t *testing.T) {
	db := TestDbClient()

	sql := "SELECT COUNT(*) FROM we_test_tab1 WHERE is_deleted <> 1;"

	result := db.QueryAggregate(sql, nil)
	tsgutils.Stdout("Select aggregate result: ", result)
}

func TestDeleteSQL(t *testing.T) {
	db := TestDbClient()

	sql := "DELETE FROM we_test_tab1 WHERE id = 5;"

	result := db.Exec(sql, nil)
	tsgutils.Stdout("Delete result: rows affected: ", result)
}

func TestDbTx()  {
	db := TestDbClient()
	tx := db.TxBegin()
	sql1 := "INSERT INTO `we_test_tab1` (`name`, `gender`, `birthday`, `stature`, `weight`, `created_time`, `modified_time`, `is_deleted`) VALUES('tony', 2, '1991-01-01', 171.31, 41.11, '2018-04-19 13:20:09', '2018-04-19 13:20:09', 0);"
	sql2 := ""
}
