package tsgmysqlutils

/*
 string utils
 @author Tony Tian
 @date 2018-04-16
 @version 1.0.0
*/

import (
	"log"
	"github.com/timespacegroup/go-utils"
)

func PrintSlowConn(driverName, host, dbName string, consume int64) {
	builder := tsgutils.NewStringBuilder()
	builder.Append(driverName)
	builder.Append(" Conn Timeout: ")
	builder.Append("Host: ")
	builder.Append(host)
	builder.Append(", DBName: ")
	builder.Append(dbName)
	builder.Append(", Consume time: ")
	builder.AppendInt64(consume)
	builder.Append("ms")
	log.Println(builder.ToString())
}

func PrintErrorSql(err error, sql string, args ...interface{}) {
	if err != nil {
		log.Println("Error Sql: ", sql)
		if ArgsIsNotNil(args...) {
			log.Print("Error Sql Args: ")
			log.Println(args...)
		}
	}
}

func PrintSlowSql(host, dbName string, consume int64, sql string, args ...interface{}) {
	builder := tsgutils.NewStringBuilder()
	builder.Append("Slow Sql: ")
	builder.Append("Host: ")
	builder.Append(host)
	builder.Append(", DBName: ")
	builder.Append(dbName)
	builder.Append(", Consume time: ")
	builder.AppendInt64(consume)
	builder.Append("ms")
	log.Println(builder.ToString())

	builder.Clear()
	builder.Append("Slow Sql: ")
	builder.Append(sql)
	log.Println(builder.ToString())

	if ArgsIsNotNil(args...) {
		builder.Clear()
		builder.Append(" Slow Sql Args: ")
		log.Print(builder.ToString())
		log.Println(args...)
	}
}
