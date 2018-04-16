package tsgmysqlutils

/*
  Usage:
	var dbConfig ts.DBConfig
	dbConfig.DbHost = "127.0.0.1"
	dbConfig.DbUser = "root"
	dbConfig.DbPass = "123456"
	dbConfig.DbName = "treasure"
	var ormConfig ts.ORMConfig
	ormConfig.DbConfig = dbConfig
	ormConfig.TabName = []string{"super_lotto", "super_lotto_bonus"}
	tsgmysqlutils.GenerateORM(ormConfig)
 */

import (
	"github.com/timespacegroup/go-utils"
)

type ORMConfig struct {
	DbConfig DBConfig
	TabName  [] string
}

func GenerateORM(config ORMConfig) {
	db := NewDbClient(config.DbConfig)
	for ti := range config.TabName {
		tabName := config.TabName[ti]
		rows := db.QueryMetaData(tabName)
		cols, err := rows.Columns()
		tsgutils.CheckAndPrintError(MySQL+" Query table meta data columns failed", err)
		types, err := rows.ColumnTypes()
		tsgutils.CheckAndPrintError(MySQL+" Query table meta data column types failed", err)
		len := len(cols)
		var builder tsgutils.StringBuilder
		structName := tsgutils.FirstCaseToUpper(tabName, true)
		builder.Append("type ").Append(structName).Append(" struct {\n")
		for i := 0; i < len; i++ {
			colName := cols[i]
			colType := TYPES[types[i].DatabaseTypeName()]
			builder.Append("	").Append(tsgutils.FirstCaseToUpper(colName, true))
			builder.Append(" ").Append(colType)
			builder.Append(" `column:\"").Append(colName).Append("\"`")
			builder.Append("\n")
		}
		builder.Append("	").Append(structName).Append("s").Append(" [] ")
		builder.Append(structName).Append("\n")
		builder.Append("}\n")
		println(builder.ToString())
	}
}

var TYPES = map[string]string{
	"TINYINT":   "int64",
	"SMALLINT":  "int64",
	"MEDIUMINT": "int64",
	"INT":       "int64",
	"BIGINT":    "int64",
	"DECIMAL":   "float64",
	"VARCHAR":   "float64",
	"DATE":      "time.Time",
	"TIMESTAMP": "time.Time",
}
