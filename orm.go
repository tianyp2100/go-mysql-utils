package tsgmysqlutils

/*
 string utils
 @author Tony Tian
 @date 2018-04-16
 @version 1.0.0
*/

/*
  ORM: Object(struct) Relational Mapping
  Usage:
	var dbConfig tsgmysqlutils.DBConfig
	dbConfig.DbHost = "127.0.0.1"
	dbConfig.DbUser = "root"
	dbConfig.DbPass = "123456"
	dbConfig.IsLocalTime = true
	dbConfig.DbName = "treasure"

	var ormConfig tsgmysqlutils.ORMConfig
	ormConfig.DbConfig = dbConfig
	ormConfig.TabName = []string{"super_lotto", "super_lotto_bonus"}
	tsgmysqlutils.GenerateORM(ormConfig)
 */

import (
	mysql "database/sql"
	"github.com/timespacegroup/go-utils"
)

/*
  ORM configuration
 */
type ORMConfig struct {
	DbConfig DBConfig
	TabName  [] string
}

type ORMBase interface {
	RowToStruct(row *mysql.Row)
	RowsToStruct(rows *mysql.Rows)
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
		structBuilder := tsgutils.NewStringBuilder()
		// builder this table's struct
		structName := tsgutils.FirstCaseToUpper(tabName, true)
		structBuilder.Append("type ").Append(structName).Append(" struct {\n")
		fieldNames := tsgutils.NewInterfaceBuilder()
		for i := 0; i < len; i++ {
			colName := cols[i]
			colType := DBTypes[types[i].DatabaseTypeName()]
			fieldName := tsgutils.FirstCaseToUpper(colName, true)
			fieldNames.Append(fieldName)
			structBuilder.Append("	").Append(fieldName)
			structBuilder.Append(" ").Append(colType)
			structBuilder.Append(" `column:\"").Append(colName).Append("\"`")
			structBuilder.Append("\n")
		}
		structBuilder.Append("	").Append(structName).Append("s").Append(" [] ")
		structBuilder.Append(structName).Append("\n")
		structBuilder.Append("}\n")
		println(structBuilder.ToString())

		// builder this table's function
		funcBuilder := tsgutils.NewStringBuilder()

		aliasStructName := tsgutils.FirstCaseToUpper(tabName, false)
		funcBuilder.Append("func (").Append(aliasStructName).Append(structName).Append(") RowToStruct(row *mysql.Row) {")
		funcBuilder.Append("  builder := tsgutils.NewInterfaceBuilder()")

		fieldNames
		for i := range fieldNames.ToInterfaces() {
			builder.Append("  builder.Append(&").Append(aliasStructName).Append(".").Append()
		}

		builder.Append("xxxxxsxxxxx")

	}
}

var DBTypes = map[string]string{
	"TINYINT":    "int64",
	"SMALLINT":   "int64",
	"MEDIUMINT":  "int64",
	"INT":        "int64",
	"BIGINT":     "int64",
	"DECIMAL":    "float64",
	"FLOAT":      "float64",
	"DOUBLE":     "float64",
	"NUMERIC":    "float64",
	"CHAR":       "string",
	"VARCHAR":    "string",
	"BINARY":     "string",
	"VARBINARY":  "string",
	"BLOB":       "string",
	"TINYTEXT":   "string",
	"TEXT":       "string",
	"MEDIUMTEXT": "string",
	"LONGTEXT":   "string",
	"ENUM":       "string",
	"SET":        "string",
	"DATE":       "time.Time",
	"DATETIME":   "time.Time",
	"TIMESTAMP":  "time.Time",
}
