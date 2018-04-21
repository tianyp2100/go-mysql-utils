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
	Paste from the console to the IDE, and format, ok.

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

type ORMGenerator struct {
	Client     *DBClient
	addComment bool
}

func NewORMGenerator(client *DBClient) *ORMGenerator {
	var orm ORMGenerator
	orm.Client = client
	return &orm
}

func (orm *ORMGenerator) DefaultGenerator(tabName [] string) {
	orm.getDbInfo()
	warmTips.Append("\n\n")
	orm.DefaultImportAndStruct(tabName)
	orm.DefaultSelectFunc(tabName)
	tsgutils.Stdout(warmTips.ToInterfaces()...)
}

func (orm *ORMGenerator) DefaultImportAndStruct(tabName [] string) {
	// builder go file import
	importBuilder := tsgutils.NewStringBuilder()
	importBuilder.Append("import (\n")
	importBuilder.Append("\t").Append("\"time\"\n")
	importBuilder.Append("\t").Append("db \"database/sql\"\n")
	importBuilder.Append("\t").Append("\"github.com/timespacegroup/go-utils\"\n")
	importBuilder.Append(")\n")
	tsgutils.Stdout(importBuilder.ToString())
	warmTips.Append("// The (table -> struct) generated tabs: ")
	hasComment := orm.addComment
	for i := range ORMTabsCols {
		tabName1 := tabName[i]
		for j := range ORMTabsCols {
			ORMTab := ORMTabsCols[j]
			tabName2 := ORMTab.TName
			if tabName2 == tabName1 {
				tabNameTmp := tabName1
				cols := ORMTab.TColumns
				// builder this table's struct
				structBuilder := importBuilder.Clear()
				structName := tsgutils.FirstCaseToUpper(tabNameTmp, true)
				if hasComment {
					structBuilder.Append("/*\n")
					structBuilder.Append("\t").Append(ORMTab.TComment).Append("\n")
					structBuilder.Append("*/\n")
				}
				structBuilder.Append("type ").Append(structName).Append(" struct {\n")
				fieldNames := tsgutils.NewInterfaceBuilder()
				for k := range cols {
					col := cols[k]
					colName := col.CName
					colType := DBGoTypes[col.CType]
					colComment := col.CComment
					//colComment := DBTypes[col.CComment]
					fieldName := tsgutils.FirstCaseToUpper(colName, true)
					fieldNames.Append(fieldName)
					structBuilder.Append("\t").Append(fieldName)
					structBuilder.Append(" ").Append(colType)
					structBuilder.Append(" `column:\"").Append(colName).Append("\"`")
					if hasComment {
						structBuilder.Append("\t//").Append(colComment)
					}
					structBuilder.Append("\n")
				}
				structBuilder.Append("\t").Append(structName).Append("s").Append(" [] ").Append(structName).Append("\n")
				structBuilder.Append("}\n")
				tsgutils.Stdout(structBuilder.ToString())
				warmTips.Append(tabNameTmp)
			}
		}
	}
	warmTips.Append("\n")
}

func (orm *ORMGenerator) DefaultSelectFunc(tabName [] string) {
	warmTips.Append("// The (table -> func) generated tabs: ")
	for i := range ORMTabsCols {
		tabName1 := tabName[i]
		for j := range ORMTabsCols {
			ORMTab := ORMTabsCols[j]
			tabName2 := ORMTab.TName
			if tabName2 == tabName1 {
				tabNameTmp := tabName1
				cols := ORMTab.TColumns
				structName := tsgutils.FirstCaseToUpper(tabNameTmp, true)
				aliasStructName := tsgutils.FirstCaseToUpper(tabNameTmp, false)
				structNames := structName + "s"
				aliasStructNames := aliasStructName + "s"
				fieldNames := tsgutils.NewInterfaceBuilder()
				for k := range cols {
					col := cols[k]
					colName := col.CName
					fieldName := tsgutils.FirstCaseToUpper(colName, true)
					fieldNames.Append(fieldName)
				}

				// builder this table's select row function
				fieldNamesArray := fieldNames.ToInterfaces()
				funcBuilder1 := tsgutils.NewStringBuilder()
				funcBuilder1.Append("func (").Append(aliasStructName).Append(" *").Append(structName).Append(") RowToStruct(row *db.Row) error {\n")
				funcBuilder1.Append("\tbuilder := tsgutils.NewInterfaceBuilder()\n")
				for i := range fieldNamesArray {
					funcBuilder1.Append("\tbuilder.Append(&").Append(aliasStructName).Append(".").Append(tsgutils.InterfaceToString(fieldNamesArray[i])).Append(")\n")
				}
				funcBuilder1.Append("\terr := row.Scan(builder.ToInterfaces()...)\n")
				funcBuilder1.Append("\tif err != nil{\n")
				funcBuilder1.Append("\t\treturn err\n")
				funcBuilder1.Append("\t}\n")
				funcBuilder1.Append("\treturn nil\n")
				funcBuilder1.Append("}\n")
				tsgutils.Stdout(funcBuilder1.ToString())

				// builder this table's select rows function
				funcBuilder2 := funcBuilder1.Clear()
				funcBuilder2.Append("func (").Append(aliasStructName).Append(" *").Append(structName).Append(") RowsToStruct(rows *db.Rows) error {\n")
				funcBuilder2.Append("\tvar ").Append(aliasStructNames).Append(" [] ").Append(structName).Append("\n")
				funcBuilder2.Append("\tbuilder := tsgutils.NewInterfaceBuilder()\n")
				funcBuilder2.Append("\tfor rows.Next() {\n")
				funcBuilder2.Append("\t\tbuilder.Clear()\n")
				for i := range fieldNamesArray {
					funcBuilder2.Append("\t\tbuilder.Append(&").Append(aliasStructName).Append(".").Append(tsgutils.InterfaceToString(fieldNamesArray[i])).Append(")\n")
				}
				funcBuilder2.Append("\t\terr := rows.Scan(builder.ToInterfaces()...)\n")
				funcBuilder2.Append("\t\tif err != nil{\n")
				funcBuilder2.Append("\t\t\treturn err\n")
				funcBuilder2.Append("\t\t}\n")
				funcBuilder2.Append("\t\t").Append(aliasStructNames).Append(" = append(").Append(aliasStructNames).Append(", *").Append(aliasStructName).Append(")\n")
				funcBuilder2.Append("\t}\n")
				funcBuilder2.Append("\tif rows != nil {\n")
				funcBuilder2.Append("\t\tdefer rows.Close()\n")
				funcBuilder2.Append("\t}\n")
				funcBuilder2.Append("\t").Append(aliasStructName).Append(".").Append(structNames).Append(" = ").Append(aliasStructNames).Append("\n")
				funcBuilder2.Append("\treturn nil\n")
				funcBuilder2.Append("}\n")
				tsgutils.Stdout(funcBuilder2.ToString())

				// builder this table's update by id function

				warmTips.Append(tabNameTmp)
			}
		}
	}
	warmTips.Append("\n")
}

var DBGoTypes = map[string]string{
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

type ORMBase interface {
	RowToStruct(row *db.Row) error
	RowsToStruct(rows *db.Rows) error
}

var ORMTabsCols [] ORMTable
var warmTips tsgutils.InterfaceBuilder

type ORMTable struct {
	TName    string
	TComment string
	TColumns [] ORMColumn
}

type ORMColumn struct {
	CName    string
	CType    string
	CComment string
}

func setORMTabsCols(tName, tComment, cName, cType, cComment string) {
	hasNotTab := false
	for i := range ORMTabsCols {
		table := ORMTabsCols[i]
		if table.TName == tName {
			var column ORMColumn
			column.CName = cName
			column.CType = cType
			column.CComment = cComment
			table.TColumns = append(table.TColumns, column)
			ORMTabsCols[i] = table
			hasNotTab = true
		}
	}
	if !hasNotTab {
		var tab ORMTable
		tab.TName = tName
		tab.TComment = tComment
		var column ORMColumn
		column.CName = cName
		column.CType = cType
		column.CComment = cComment
		tab.TColumns = append(tab.TColumns, column)
		ORMTabsCols = append(ORMTabsCols, tab)
	}
}

func (orm *ORMGenerator) getDbInfo() {
	var tName, tComment, cName, cType, cComment string
	rows := orm.Client.QueryDBInfo()
	for rows.Next() {
		err := rows.Scan(&tName, &tComment, &cName, &cType, &cComment)
		tsgutils.CheckAndPrintError("Get db info rows scan failed", err)
		setORMTabsCols(tName, tComment, cName, getDBType(cType), cComment)
	}
}

func getDBType(cType string) string {
	typeTmp := tsgutils.NewString(cType)
	if typeTmp.Contains(tsgutils.PARENTHESIS_LEFT) {
		return typeTmp.SubstringEnd(typeTmp.Index(tsgutils.PARENTHESIS_LEFT)).ToUpper().ToString()
	}
	return typeTmp.ToUpper().ToString()
}
