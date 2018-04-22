package tsgmysqlutils

import (
	db "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/timespacegroup/go-utils"
	"strings"
)

/*
 MySQL client
  Usage:
    var dbConfig tsgmysqlutils.DBConfig
	dbConfig.DbHost = "127.0.0.1"
	dbConfig.DbUser = "root"
	dbConfig.DbPass = "123456"
	dbConfig.IsLocalTime = true
	dbConfig.DbName = "test"
	client := tsgmysqlutils.NewDbClient(dbConfig)

   @author Tony Tian
   @date 2018-04-16
   @version 1.0.0
*/

type DBClient struct {
	Config DBConfig
	Db     *db.DB
}

const (
	MySQL                     = "MySQL"
	SlowSqlTimeoutMillisecond = 2000
	ConnDBTimeoutMillisecond  = 2000
)

/*
 Get a MySQL client
*/
func NewDbClient(config DBConfig) *DBClient {
	var client DBClient
	client.Config = config
	client.Db = GetConn(config)
	return &client
}

/*
 MySQL database connect configuration
*/
type DBConfig struct {
	DbHost string
	DbUser string
	DbPass string
	DbName string
	// if false, local default UTC +0
	IsLocalTime bool
}

/*
  Get a MySQL connection string
*/
func getDbConnString(config DBConfig) string {
	builder := tsgutils.NewStringBuilder()
	builder.Append(config.DbUser).Append(":").Append(config.DbPass)
	builder.Append("@tcp(").Append(config.DbHost).Append(":").Append("3306").Append(")/")
	builder.Append(config.DbName).Append("?").Append("charset=utf8")
	if config.IsLocalTime {
		builder.Append("&parseTime=true&loc=Local")
	}
	return builder.ToString()
}

/*
  Get a MySQL connection
*/
func GetConn(config DBConfig) *db.DB {
	dbConnString := getDbConnString(config)
	start := tsgutils.Millisecond()
	mysql, err := db.Open(strings.ToLower(MySQL), dbConnString)
	consume := tsgutils.Millisecond() - start
	if consume > ConnDBTimeoutMillisecond {
		PrintSlowConn(MySQL, config.DbHost, config.DbName, consume)
	}
	tsgutils.CheckAndPrintError(MySQL+" connection failed, db conn string: \n"+dbConnString, err)
	return mysql
}

/*
  Get MySQL statement
*/
func (client *DBClient) GetStmt(sql string) (stmt *db.Stmt, err error) {
	stmt, err = client.Db.Prepare(sql)
	if err != nil {
		PrintErrorSql(err, sql, nil)
		return nil, err
	}
	return stmt, nil
}

/*
 Close MySQL connection
*/
func (client *DBClient) CloseConn() {
	if client.Db != nil {
		client.Db.Close()
	}
}

/*
  Close MySQL statement
*/
func (client *DBClient) CloseStmt(stmt *db.Stmt) {
	if stmt != nil {
		stmt.Close()
	}
}

/*
  Get database table a row data
*/
func (client *DBClient) QueryRow(orm ORMBase, sql string, args ...interface{}) (row *db.Row, err error) {
	start := tsgutils.Millisecond()
	stmt, err := client.GetStmt(sql)
	if stmt == nil || err != nil {
		return nil, err
	}
	row, err = client.forkQuery(stmt, orm, sql, args...)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	defer client.CloseStmt(stmt)
	return row, err
}

/*
  Get database table multiple rows data
*/
func (client *DBClient) QueryList(orm ORMBase, sql string, args ...interface{}) (rows *db.Rows, err error) {
	start := tsgutils.Millisecond()
	stmt, err := client.GetStmt(sql)
	if stmt == nil || err != nil {
		return nil, err
	}
	rows, err = client.forkQueryList(stmt, orm, sql, args...)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	defer client.CloseStmt(stmt)
	return rows, err
}

/*
 Database aggregate function, eg: SUM(*),COUNT(*) etc.
*/
func (client *DBClient) QueryAggregate(sql string, args ...interface{}) (aggregate int64, err error) {
	row, err := client.QueryRow(nil, sql, args...)
	if err != nil {
		return 0, err
	}
	var result int64
	err = row.Scan(&result)
	if err != nil {
		return 0, nil
	}
	return result, nil
}

/*
  Modify database table info or data
*/
func (client *DBClient) Exec(sql string, args ...interface{}) (result int64, err error) {
	start := tsgutils.Millisecond()
	stmt, err := client.GetStmt(sql)
	if stmt == nil || err != nil {
		return 0, err
	}
	result, err = client.forkExec(stmt, sql, args...)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	defer client.CloseStmt(stmt)
	return result, err
}

/*
  Begin the transaction
*/
func (client *DBClient) TxBegin() (tx *db.Tx, err error) {
	return client.Db.Begin()
}

/*
  Commit the transaction
*/
func (client *DBClient) TxCommit(tx *db.Tx) bool {
	err := tx.Commit()
	if err != nil {
		tsgutils.CheckAndPrintError(MySQL+" tx commit failed", err)
		return false
	}
	return true
}

/*
  Rollback the transaction
*/
func (client *DBClient) TxRollback(tx *db.Tx) {
	err := tx.Rollback()
	if err != nil {
		tsgutils.CheckAndPrintError(MySQL+" tx rollback failed", err)
	}
}

/*
  Get MySQL statement,transaction
*/
func (client *DBClient) GetTxStmt(tx *db.Tx, sql string) (stmt *db.Stmt, err error) {
	stmt, err = tx.Prepare(sql)
	if err != nil {
		PrintErrorSql(err, sql, nil)
		return nil, err
	}
	return stmt, nil
}

/*
  Get database table a row data,transaction
*/
func (client *DBClient) TxQueryRow(tx *db.Tx, orm ORMBase, sql string, args ...interface{}) (row *db.Row, err error) {
	start := tsgutils.Millisecond()
	stmt, err := client.GetTxStmt(tx, sql)
	if stmt == nil || err != nil {
		return nil, err
	}
	row, err = client.forkQuery(stmt, orm, sql, args...)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	return row, err
}

/*
  Get database table multiple rows data,transaction
*/
func (client *DBClient) TxQueryList(tx *db.Tx, orm ORMBase, sql string, args ...interface{}) (rows *db.Rows, err error) {
	start := tsgutils.Millisecond()
	stmt, err := client.GetTxStmt(tx, sql)
	if stmt == nil || err != nil {
		return nil, err
	}
	rows, err = client.forkQueryList(stmt, orm, sql, args...)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	return rows, err
}

/*
 Database aggregate function, eg: SUM(*),COUNT(*) etc,transaction
*/
func (client *DBClient) TxQueryAggregate(tx *db.Tx, sql string, args ...interface{}) (aggregate int64, err error) {
	row, err := client.TxQueryRow(tx, nil, sql, args...)
	if err != nil {
		return 0, err
	}
	var result int64
	err = row.Scan(&result)
	if err != nil {
		return 0, nil
	}
	return result, nil
}

/*
  Modify database table info or data,transaction
*/
func (client *DBClient) TxExec(tx *db.Tx, sql string, args ...interface{}) (result int64, err error) {
	start := tsgutils.Millisecond()
	stmt, err := client.GetTxStmt(tx, sql)
	if stmt == nil || err != nil {
		return 0, err
	}
	result, err = client.forkExec(stmt, sql, args...)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	defer client.CloseStmt(stmt)
	return result, err
}

/*
  Get database table metadata
*/
func (client *DBClient) QueryMetaData(tabName string) *db.Rows {
	rows, err := client.Db.Query("SELECT * FROM " + tabName + " WHERE 1=1 LIMIT 1;")
	tsgutils.CheckAndPrintError("Query '"+tabName+"' table meta data failed", err)
	return rows
}

/*
  Get database information
*/
func (client *DBClient) QueryDBInfo() *db.Rows {
	sql := "SELECT tab.TABLE_NAME,tab.TABLE_COMMENT,col.COLUMN_NAME,col.COLUMN_TYPE,col.COLUMN_COMMENT " +
		"FROM information_schema.TABLES tab,INFORMATION_SCHEMA.Columns col " +
		"WHERE col.TABLE_NAME=tab.TABLE_NAME AND tab.`TABLE_SCHEMA` = ?"
	rows, err := client.Db.Query(sql, client.Config.DbName)
	tsgutils.CheckAndPrintError("Query '"+client.Config.DbName+"' db info failed", err)
	return rows
}

func (client *DBClient) forkQuery(stmt *db.Stmt, orm ORMBase, sql string, args ...interface{}) (row *db.Row, err error) {
	row = stmt.QueryRow(args...)
	if orm != nil {
		err = orm.RowToStruct(row)
		if err != nil {
			PrintErrorSql(err, sql, args...)
			return nil, err
		}
	}
	return row, nil
}

func (client *DBClient) forkQueryList(stmt *db.Stmt, orm ORMBase, sql string, args ...interface{}) (rows *db.Rows, err error) {
	rows, err = stmt.Query(args...)
	if err != nil {
		PrintErrorSql(err, sql, args...)
		return nil, err
	}
	if orm != nil {
		err = orm.RowsToStruct(rows)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil

}

func (client *DBClient) forkExec(stmt *db.Stmt, sql string, args ...interface{}) (result int64, err error) {
	var results db.Result
	results, err = stmt.Exec(args...)
	if err != nil {
		PrintErrorSql(err, sql, args...)
		return 0, err
	}
	var intResult int64
	if tsgutils.NewString(sql).ContainsIgnoreCase("INSERT") {
		intResult, err = results.LastInsertId()
	} else {
		intResult, err = results.RowsAffected()
	}
	if err != nil {
		PrintErrorSql(err, sql, args...)
		return 0, err
	}
	return intResult, nil
}

func (client *DBClient) slowSql(consume int64, sql string, args ...interface{}) {
	if consume > SlowSqlTimeoutMillisecond {
		PrintSlowSql(client.Config.DbHost, client.Config.DbName, consume, sql, args...)
	}
}
