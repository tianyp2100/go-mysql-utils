package tsgmysqlutils

/*
 string utils
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
	db, err := db.Open(strings.ToLower(MySQL), dbConnString)
	consume := tsgutils.Millisecond() - start
	if consume > ConnDBTimeoutMillisecond {
		PrintSlowConn(MySQL, config.DbHost, config.DbName, consume)
	}
	tsgutils.CheckAndPrintError(MySQL+" connection failed, db conn string: \n"+dbConnString, err)
	return db
}

/*
  Get a MySQL statement
 */
func (client *DBClient) GetStmt(sql string) *db.Stmt {
	stmt, err := client.Db.Prepare(sql)
	tsgutils.CheckAndPrintError(MySQL+" prepare stmt failed", err)
	PrintErrorSql(MySQL, err, sql)
	return stmt
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
  Get MySQL database table metadata
 */
func (client *DBClient) QueryMetaData(tabName string) *db.Rows {
	rows, err := client.Db.Query("SELECT * FROM " + tabName + " WHERE 1=1 LIMIT 1;")
	tsgutils.CheckAndPrintError(MySQL+" Query table meta data failed", err)
	return rows
}

/*
  Get MySQL database table a row data
 */
func (client *DBClient) QueryRow(sql string, args []interface{}, orm ORMBase) {
	start := tsgutils.Millisecond()
	row := client.forkQuery(sql, args...)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	orm.RowToStruct(row)
}

func (client *DBClient) QueryAggregate(sql string, args []interface{}) int64 {
	start := tsgutils.Millisecond()
	row := client.forkQuery(sql, args...)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	var result int64
	err := row.Scan(&result)
	tsgutils.CheckAndPrintError("MySQL query aggregate scan error", err)
	return result
}

/*
  Get MySQL database table multiple rows data
 */
func (client *DBClient) QueryList(sql string, args []interface{}, orm ORMBase) {
	start := tsgutils.Millisecond()
	rows := client.forkQueryList(sql, args...)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	orm.RowsToStruct(rows)
}

/*
  Begin the transaction
 */
func (client *DBClient) TxBegin() *db.Tx {
	tx, err := client.Db.Begin()
	tsgutils.CheckAndPrintError(MySQL+" begin tx failed", err)
	return tx
}

/*
  Commit the transaction
 */
func (client *DBClient) TxCommit(tx *db.Tx) {
	err := tx.Commit()
	tsgutils.CheckAndPrintError(MySQL+" commit tx failed", err)
}

/*
  Rollback the transaction
 */
func (client *DBClient) TxRollback(tx *db.Tx) {
	err := tx.Rollback()
	tsgutils.CheckAndPrintError(MySQL+" rollback tx failed", err)
}

/*
  Modify MySQL database table info or data
 */
func (client *DBClient) Exec(sql string, args []interface{}) int64 {
	stmt := client.GetStmt(sql)
	start := tsgutils.Millisecond()
	result, err := stmt.Exec(args...)
	tsgutils.CheckAndPrintError(MySQL+" exec sql failed", err)
	PrintErrorSql(MySQL, err, sql, args)
	client.slowSql(tsgutils.Millisecond()-start, sql, args...)
	var intResult int64
	var flag string
	if tsgutils.NewString(sql).ContainsIgnoreCase("INSERT") {
		intResult, err = result.LastInsertId()
		flag = "get last insert id"
	} else {
		intResult, err = result.RowsAffected()
		flag = "get rows affected"
	}
	tsgutils.CheckAndPrintError(MySQL+" exec and "+flag+" failed", err)
	PrintErrorSql(MySQL, err, sql, args...)
	client.CloseStmt(stmt)
	return intResult
}

func (client *DBClient) forkQuery(sql string, args ...interface{}) *db.Row {
	stmt := client.GetStmt(sql)
	var row *db.Row
	if len(args) > 0 {
		row = stmt.QueryRow(args)
	} else {
		row = stmt.QueryRow()
	}
	client.CloseStmt(stmt)
	return row
}

func (client *DBClient) forkQueryList(sql string, args ...interface{}) *db.Rows {
	stmt := client.GetStmt(sql)
	var rows *db.Rows
	var err error
	if len(args) > 0 {
		rows, err = stmt.Query(args)
	} else {
		rows, err = stmt.Query()
	}
	tsgutils.CheckAndPrintError(MySQL+" query rows list failed", err)
	PrintErrorSql(MySQL, err, sql, args)
	client.CloseStmt(stmt)
	return rows
}

func (client *DBClient) slowSql(consume int64, sql string, args ...interface{}) {
	if consume > SlowSqlTimeoutMillisecond {
		PrintSlowSql(MySQL, client.Config.DbHost, client.Config.DbName, consume, sql, args...)
	}
}
