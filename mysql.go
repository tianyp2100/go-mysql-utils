package tsgmysqlutils

/*
 MySQL utils
 @author Tony Tian
 @date 2018-04-16
 @version 1.0.0
*/

import (
	"strings"
	db "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/timespacegroup/go-utils"
)

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
  Get a MySQL statement
 */
func (client *DBClient) GetStmt(sql string) (stmt *db.Stmt, err error) {
	stmt, err = client.Db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	PrintErrorSql(err, sql, nil)
	return stmt, nil
}

/*
 Close MySQL connection
*/
func (client *DBClient) CloseConn() {
	if client.Db != nil {
		defer client.Db.Close()
	}
}

/*
  Close MySQL statement
 */
func (client *DBClient) CloseStmt(stmt *db.Stmt) {
	if stmt != nil {
		defer stmt.Close()
	}
}

/*
  Get database table a row data
 */
func (client *DBClient) QueryRow(orm ORMBase, sql string, args ...interface{}) error {
	start := tsgutils.Millisecond()
	row, err := client.forkQuery(sql, args...)
	if err != nil {
		return err
	}
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	err = orm.RowToStruct(row)
	if err != nil {
		return err
	}
	return nil
}

/*
 Database aggregate function, eg: SUM(*),COUNT(*) etc.
 */
func (client *DBClient) QueryAggregate(sql string, args ...interface{}) (aggregate int64, err error) {
	start := tsgutils.Millisecond()
	row, err := client.forkQuery(sql, args...)
	if err != nil {
		return 0, nil
	}
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	var result int64
	err = row.Scan(&result)
	if err != nil {
		return 0, nil
	}
	return result, nil
}

/*
  Get database table multiple rows data
 */
func (client *DBClient) QueryList(orm ORMBase, sql string, args ...interface{}) error {
	start := tsgutils.Millisecond()
	rows, err := client.forkQueryList(sql, args...)
	if err != nil {
		return err
	}
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	err = orm.RowsToStruct(rows)
	if err != nil {
		return err
	}
	return nil
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
		tsgutils.Println(MySQL+" tx commit failed", err)
		return false
	}
	return true
}

/*
  Rollback the transaction
 */
func (client *DBClient) TxRollback(tx *db.Tx) {
	err := tx.Rollback()
	tsgutils.CheckAndPrintError(MySQL+" tx rollback failed", err)
}

/*
  Modify database table info or data
 */
func (client *DBClient) Exec(sql string, args ...interface{}) (result int64, err error) {
	stmt, err := client.GetStmt(sql)
	if stmt == nil || err != nil {
		return 0, err
	}
	start := tsgutils.Millisecond()
	var results db.Result
	if ArgsIsNotNil(args...) {
		results, err = stmt.Exec(args...)
	} else {
		results, err = stmt.Exec()
	}
	client.CloseStmt(stmt)
	PrintErrorSql(err, sql, args...)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	var intResult int64
	if tsgutils.NewString(sql).ContainsIgnoreCase("INSERT") {
		intResult, err = results.LastInsertId()
	} else {
		intResult, err = results.RowsAffected()
	}
	PrintErrorSql(err, sql, args...)
	return intResult, nil
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

func (client *DBClient) forkQuery(sql string, args ...interface{}) (row *db.Row, err error) {
	stmt, err := client.GetStmt(sql)
	if stmt == nil || err != nil {
		return nil, err
	}
	if ArgsIsNotNil(args...) {
		row = stmt.QueryRow(args...)
	} else {
		row = stmt.QueryRow()
	}
	client.CloseStmt(stmt)
	if err != nil {
		return nil, err
	}
	PrintErrorSql(err, sql, args...)
	return row, nil
}

func (client *DBClient) forkQueryList(sql string, args ...interface{}) (rows *db.Rows, err error) {
	stmt, err := client.GetStmt(sql)
	if stmt == nil || err != nil {
		return nil, err
	}
	if ArgsIsNotNil(args...) {
		rows, err = stmt.Query(args...)
	} else {
		rows, err = stmt.Query()
	}
	client.CloseStmt(stmt)
	if err != nil {
		return nil, err
	}
	PrintErrorSql(err, sql, args...)
	return rows, nil
}

func ArgsIsNotNil(args ...interface{}) bool {
	return args[0] != nil
}

func (client *DBClient) slowSql(consume int64, sql string, args ...interface{}) {
	if consume > SlowSqlTimeoutMillisecond {
		PrintSlowSql(client.Config.DbHost, client.Config.DbName, consume, sql, args...)
	}
}
