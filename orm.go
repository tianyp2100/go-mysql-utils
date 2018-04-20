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
	db "database/sql"
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
	RowToStruct(row *db.Row)
	RowsToStruct(rows *db.Rows)
}

func GenerateORM(config ORMConfig) {
	db := NewDbClient(config.DbConfig)
	// builder go file import
	importBuilder := tsgutils.NewStringBuilder()
	importBuilder.Append("import (\n")
	importBuilder.Append("\t").Append("\"time\"\n")
	importBuilder.Append("\t").Append("db \"database/sql\"\n")
	importBuilder.Append("\t").Append("\"github.com/timespacegroup/go-utils\"\n")
	importBuilder.Append(")\n")
	tsgutils.Stdout(importBuilder.ToString())
	for ti := range config.TabName {
		tabName := config.TabName[ti]
		rows := db.QueryMetaData(tabName)
		cols, err := rows.Columns()
		tsgutils.CheckAndPrintError(MySQL+" Query table meta data columns failed", err)
		types, err := rows.ColumnTypes()
		tsgutils.CheckAndPrintError(MySQL+" Query table meta data column types failed", err)
		len := len(cols)
		// builder this table's struct
		structBuilder := importBuilder.Clear()
		structName := tsgutils.FirstCaseToUpper(tabName, true)
		aliasStructName := tsgutils.FirstCaseToUpper(tabName, false)
		structNames := structName + "s"
		aliasStructNames := aliasStructName + "s"
		structBuilder.Append("type ").Append(structName).Append(" struct {\n")
		fieldNames := tsgutils.NewInterfaceBuilder()
		for i := 0; i < len; i++ {
			colName := cols[i]
			colType := DBTypes[types[i].DatabaseTypeName()]
			fieldName := tsgutils.FirstCaseToUpper(colName, true)
			fieldNames.Append(fieldName)
			structBuilder.Append("\t").Append(fieldName)
			structBuilder.Append(" ").Append(colType)
			structBuilder.Append(" `column:\"").Append(colName).Append("\"`")
			structBuilder.Append("\n")
		}
		structBuilder.Append("\t").Append(structName).Append("s").Append(" [] ").Append(structName).Append("\n")
		structBuilder.Append("}\n")
		tsgutils.Stdout(structBuilder.ToString())

		// builder this table's function
		fieldNamesArray := fieldNames.ToInterfaces()
		funcBuilder1 := structBuilder.Clear()
		funcBuilder1.Append("func (").Append(aliasStructName).Append(" *").Append(structName).Append(") RowToStruct(row *db.Row) {\n")
		funcBuilder1.Append("\tbuilder := tsgutils.NewInterfaceBuilder()\n")
		for i := range fieldNamesArray {
			funcBuilder1.Append("\tbuilder.Append(&").Append(aliasStructName).Append(".").Append(tsgutils.InterfaceToString(fieldNamesArray[i])).Append(")\n")
		}
		funcBuilder1.Append("\terr := row.Scan(builder.ToInterfaces()...)\n")
		funcBuilder1.Append("\ttsgutils.CheckAndPrintError(\"MySQL query row scan error\", err)\n")
		funcBuilder1.Append("}\n")
		tsgutils.Stdout(funcBuilder1.ToString())
		funcBuilder2 := funcBuilder1.Clear()

		funcBuilder2.Append("func (").Append(aliasStructName).Append(" *").Append(structName).Append(") RowsToStruct(rows *db.Rows) {\n")
		funcBuilder2.Append("\tvar ").Append(aliasStructNames).Append(" [] ").Append(structName).Append("\n")
		funcBuilder2.Append("\tbuilder := tsgutils.NewInterfaceBuilder()\n")
		funcBuilder2.Append("\tfor rows.Next() {\n")
		funcBuilder2.Append("\t\tbuilder.Clear()\n")
		for i := range fieldNamesArray {
			funcBuilder2.Append("\t\tbuilder.Append(&").Append(aliasStructName).Append(".").Append(tsgutils.InterfaceToString(fieldNamesArray[i])).Append(")\n")
		}
		funcBuilder2.Append("\t\terr := rows.Scan(builder.ToInterfaces()...)\n")
		funcBuilder2.Append("\t\ttsgutils.CheckAndPrintError(\"MySQL query rows scan error\", err)\n")
		funcBuilder2.Append("\t\t").Append(aliasStructNames).Append(" = append(").Append(aliasStructNames).Append(", *").Append(aliasStructName).Append(")\n")
		funcBuilder2.Append("\t}\n")
		funcBuilder2.Append("\tif rows != nil {\n")
		funcBuilder2.Append("\t\tdefer rows.Close()\n")
		funcBuilder2.Append("\t}\n")
		funcBuilder2.Append("\t").Append(aliasStructName).Append(".").Append(structNames).Append(" = ").Append(aliasStructNames).Append("\n")
		funcBuilder2.Append("}\n")
		tsgutils.Stdout(funcBuilder2.ToString())
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
