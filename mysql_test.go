package tsgmysqlutils

import (
	"testing"
	"errors"
	"github.com/timespacegroup/go-utils"
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
	tabSql := tsgutils.NewStringBuilder()
	tabSql.Append("CREATE TABLE `we_test_tab` (")
	tabSql.Append("`id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'The primary key id',")
	tabSql.Append("`name` varchar(64) NOT NULL DEFAULT '' COMMENT 'The user name',")
	tabSql.Append("`gender` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'The user gerder, 1:male 2:female 0:default',")
	tabSql.Append("`birthday` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'The user birthday, eg: 2018-04-16',")
	tabSql.Append("`stature` decimal(16, 2) NOT NULL DEFAULT '0.00' COMMENT 'The user stature, eg: 172.22cm',")
	tabSql.Append("`weight` decimal(16, 2) NOT NULL DEFAULT '0.00' COMMENT 'The user weight, eg: 21.77kg',")
	tabSql.Append("`modified_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'record time',")
	tabSql.Append("`is_deleted` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'Logic to delete(0:normal 1:deleted)',")
	tabSql.Append("PRIMARY KEY (`id`),")
	tabSql.Append("UNIQUE KEY `name` (`name`)")
	tabSql.Append(") ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE= utf8_bin COMMENT = 'test table';")

	db := TestDbClient()
	result := db.Exec(tabSql.ToString(), nil)
	tsgutils.FmtPrintln(result)
}
