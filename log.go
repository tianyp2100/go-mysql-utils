package ts

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

func PrintErrorSql(driverName string, err error, sql string, args ...interface{}) {
	if err != nil {
		log.Println(driverName+" Error Sql: ", sql)
		if len(args) > 0 {
			log.Println(driverName+" Error Sql Args: ", args[0])
		}
	}
}

func PrintSlowSql(driverName, host, dbName string, consume int64, sql string, args ...interface{}) {
	builder := tsgutils.NewStringBuilder()
	builder.Append(driverName)
	builder.Append(" Slow Sql: ")
	builder.Append("Host: ")
	builder.Append(host)
	builder.Append(", DBName: ")
	builder.Append(dbName)
	builder.Append(", Consume time: ")
	builder.AppendInt64(consume)
	builder.Append("ms")
	log.Println(builder.ToString())

	builder.Clear()
	builder.Append(driverName)
	builder.Append(" Slow Sql: ")
	builder.Append(sql)
	log.Println(builder.ToString())

	if len(args) > 0 {
		builder.Clear()
		builder.Append(driverName)
		builder.Append(" Slow Sql Args: ")
		log.Println(builder.ToString(), args[0])
	}
}
