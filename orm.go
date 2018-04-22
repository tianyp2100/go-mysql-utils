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

func (orm *ORMGenerator) DefaultGenerator(tabName []string) {
	orm.getDbInfo()
	warmTips.Append("\n\n")
	orm.buildORMImport()
	orm.ORMBuilder(tabName)
	tsgutils.Stdout(warmTips.ToInterfaces()...)
}

func (orm *ORMGenerator) ORMBuilder(tabNames []string) {
	warmTips.Append("// The generated tabs: ")
	for i := range ORMTabsCols {
		tabName := tabNames[i]
		for j := range ORMTabsCols {
			ORMTab := ORMTabsCols[j]
			if ORMTab.TName == tabName {
				orm.buildORMStruct(tabName, ORMTab, orm.addComment)
				orm.buildORMSqlSelect(tabName, ORMTab.TColumns)
				orm.buildORMSqlInsert(tabName)
				orm.buildORMSqlUpdate(tabName)
				orm.buildORMSqlDelete(tabName)
				orm.buildORMSqlBatchInsert(tabName)
			}
		}
	}
	warmTips.Append("\n")
}

func (orm *ORMGenerator) buildORMImport() {
	importBuilder := tsgutils.NewStringBuilder()
	importBuilder.Append("import (").Append("\n")
	importBuilder.Append("\t").Append("\"time\"").Append("\n")
	importBuilder.Append("\t").Append("\"errors\"").Append("\n")
	importBuilder.Append("\t").Append("\"reflect\"").Append("\n")
	importBuilder.Append("\t").Append("db \"database/sql\"").Append("\n")
	importBuilder.Append("\t").Append("\"github.com/timespacegroup/go-utils\"").Append("\n")
	importBuilder.Append(")").Append("\n")
	tsgutils.Stdout(importBuilder.ToString())
}

func (orm *ORMGenerator) buildORMStruct(tabName string, ORMTab ORMTable, hasComment bool) {
	importBuilder := tsgutils.NewStringBuilder()
	cols := ORMTab.TColumns
	structBuilder := importBuilder.Clear()
	structName := getStructName(tabName)
	structNames := getStructNames(tabName)
	if hasComment {
		structBuilder.Append("/*").Append("\n")
		structBuilder.Append("\t").Append(ORMTab.TComment).Append("\n")
		structBuilder.Append("*/").Append("\n")
	}
	structBuilder.Append("type ").Append(structName).Append(" struct {").Append("\n")
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
			structBuilder.Append("\t").Append("// ").Append(colComment)
		}
		structBuilder.Append("\n")
	}
	structBuilder.Append("\t").Append(structNames).Append(" [] ").Append(structName)
	if hasComment {
		structBuilder.Append("\t").Append("// This value is used for batch queries and inserts.")
	}
	structBuilder.Append("\n")
	structBuilder.Append("}").Append("\n")
	tsgutils.Stdout(structBuilder.ToString())
	warmTips.Append(tabName)
}

func (orm *ORMGenerator) buildORMSqlSelect(tabName string, cols []ORMColumn) {
	structName := getStructName(tabName)
	aliasStructName := getAliasStructName(tabName)
	structNames := getStructNames(tabName)
	aliasStructNames := getAliasStructNames(tabName)
	fieldNames := tsgutils.NewInterfaceBuilder()
	for k := range cols {
		col := cols[k]
		colName := col.CName
		fieldName := tsgutils.FirstCaseToUpper(colName, true)
		fieldNames.Append(fieldName)
	}
	fieldNamesArray := fieldNames.ToInterfaces()
	funcRowBuilder := tsgutils.NewStringBuilder()
	// builder this table's select row function
	funcRowBuilder.Append("func (").Append(aliasStructName).Append(" *").Append(structName).Append(") RowToStruct(row *db.Row) error {").Append("\n")
	funcRowBuilder.Append("\t").Append("builder := tsgutils.NewInterfaceBuilder()").Append("\n")
	for i := range fieldNamesArray {
		funcRowBuilder.Append("\t").Append("builder.Append(&").Append(aliasStructName).Append(".").Append(tsgutils.InterfaceToString(fieldNamesArray[i])).Append(")").Append("\n")
	}
	funcRowBuilder.Append("\t").Append("err := row.Scan(builder.ToInterfaces()...)").Append("\n")
	funcRowBuilder.Append("\t").Append("if err != nil{").Append("\n")
	funcRowBuilder.Append("\t\t").Append("return err").Append("\n")
	funcRowBuilder.Append("\t").Append("}").Append("\n")
	funcRowBuilder.Append("\t").Append("return nil").Append("\n")
	funcRowBuilder.Append("}").Append("\n")
	tsgutils.Stdout(funcRowBuilder.ToString())

	// builder this table's select rows function
	funcRowsBuilder := funcRowBuilder.Clear()
	funcRowsBuilder.Append("func (").Append(aliasStructName).Append(" *").Append(structName).Append(") RowsToStruct(rows *db.Rows) error {").Append("\n")
	funcRowsBuilder.Append("\t").Append("var ").Append(aliasStructNames).Append(" [] ").Append(structName).Append("\n")
	funcRowsBuilder.Append("\t").Append("builder := tsgutils.NewInterfaceBuilder()").Append("\n")
	funcRowsBuilder.Append("\t").Append("for rows.Next() {").Append("\n")
	funcRowsBuilder.Append("\t\t").Append("builder.Clear()").Append("\n")
	for i := range fieldNamesArray {
		funcRowsBuilder.Append("\t\t").Append("builder.Append(&").Append(aliasStructName).Append(".").Append(tsgutils.InterfaceToString(fieldNamesArray[i])).Append(")").Append("\n")
	}
	funcRowsBuilder.Append("\t\t").Append("err := rows.Scan(builder.ToInterfaces()...)").Append("\n")
	funcRowsBuilder.Append("\t\t").Append("if err != nil{").Append("\n")
	funcRowsBuilder.Append("\t\t\t").Append("return err").Append("\n")
	funcRowsBuilder.Append("\t\t").Append("}").Append("\n")
	funcRowsBuilder.Append("\t\t").Append(aliasStructNames).Append(" = append(").Append(aliasStructNames).Append(", *").Append(aliasStructName).Append(")").Append("\n")
	funcRowsBuilder.Append("\t").Append("}").Append("\n")
	funcRowsBuilder.Append("\t").Append("if rows != nil {").Append("\n")
	funcRowsBuilder.Append("\t\t").Append("defer rows.Close()").Append("\n")
	funcRowsBuilder.Append("\t").Append("}").Append("\n")
	funcRowsBuilder.Append("\t").Append("").Append(aliasStructName).Append(".").Append(structNames).Append(" = ").Append(aliasStructNames).Append("\n")
	funcRowsBuilder.Append("\t").Append("return nil").Append("\n")
	funcRowsBuilder.Append("}").Append("\n")
	tsgutils.Stdout(funcRowsBuilder.ToString())
}

func (orm *ORMGenerator) buildORMSqlInsert(tabName string) {
	funcInsertBuilder := tsgutils.NewStringBuilder()
	structName := getStructName(tabName)
	aliasStructName := getAliasStructName(tabName)
	funcInsertBuilder.Append("func (").Append(aliasStructName).Append(" *").Append(structName).Append(") Insert(client *DBClient, idSet bool) (int64, error) {").Append("\n")
	funcInsertBuilder.Append("\t").Append("structParam := *").Append(aliasStructName).Append("\n")
	funcInsertBuilder.Append("\t").Append("sql := tsgutils.NewStringBuilder()").Append("\n")
	funcInsertBuilder.Append("\t").Append("qSql := tsgutils.NewStringBuilder()").Append("\n")
	funcInsertBuilder.Append("\t").Append("params := tsgutils.NewInterfaceBuilder()").Append("\n")
	funcInsertBuilder.Append("\t").Append("sql.Append(\"INSERT INTO \")").Append("\n")
	funcInsertBuilder.Append("\t").Append("sql.Append(\"").Append(tabName).Append("\")\n")
	funcInsertBuilder.Append("\t").Append("sql.Append(\" (\")").Append("\n")
	funcInsertBuilder.Append("\t").Append("ks := reflect.TypeOf(structParam)").Append("\n")
	funcInsertBuilder.Append("\t").Append("vs := reflect.ValueOf(structParam)").Append("\n")
	funcInsertBuilder.Append("\t").Append("for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {").Append("\n")
	funcInsertBuilder.Append("\t\t").Append("col := ks.Field(i).Tag.Get(\"column\")").Append("\n")
	funcInsertBuilder.Append("\t\t").Append("v := vs.Field(i).Interface()").Append("\n")
	funcInsertBuilder.Append("\t\t").Append("if col == \"id\" && !idSet {").Append("\n")
	funcInsertBuilder.Append("\t\t\t").Append("continue").Append("\n")
	funcInsertBuilder.Append("\t\t").Append("}").Append("\n")
	funcInsertBuilder.Append("\t\t").Append("sql.Append(\"`\").Append(col).Append(\"`,\")").Append("\n")
	funcInsertBuilder.Append("\t\t").Append("qSql.Append(\"?,\")").Append("\n")
	funcInsertBuilder.Append("\t\t").Append("params.Append(v)").Append("\n")
	funcInsertBuilder.Append("\t").Append("}").Append("\n")
	funcInsertBuilder.Append("\t").Append("sql.RemoveLast()").Append("\n")
	funcInsertBuilder.Append("\t").Append("qSql.RemoveLast()").Append("\n")
	funcInsertBuilder.Append("\t").Append("sql.Append(\") VALUES (\").Append(qSql.ToString()).Append(\");\")").Append("\n")
	funcInsertBuilder.Append("\t").Append("defer client.CloseConn()").Append("\n")
	funcInsertBuilder.Append("\t").Append("return client.Exec(sql.ToString(), params.ToInterfaces()...)").Append("\n")
	funcInsertBuilder.Append("}").Append("\n")
	tsgutils.Stdout(funcInsertBuilder.ToString())
}

func (orm *ORMGenerator) buildORMSqlUpdate(tabName string) {
	funcUpdateBuilder := tsgutils.NewStringBuilder()
	structName := getStructName(tabName)
	aliasStructName := getAliasStructName(tabName)
	funcUpdateBuilder.Append("func (").Append(aliasStructName).Append(" *").Append(structName).Append(") Update").Append(structName).Append("ById(client *DBClient) (int64, error) {").Append("\n")
	funcUpdateBuilder.Append("\t").Append("structParam := *").Append(aliasStructName).Append("\n")
	funcUpdateBuilder.Append("\t").Append("sql := tsgutils.NewStringBuilder()").Append("\n")
	funcUpdateBuilder.Append("\t").Append("params := tsgutils.NewInterfaceBuilder()").Append("\n")
	funcUpdateBuilder.Append("\t").Append("sql.Append(\"UPDATE \")").Append("\n")
	funcUpdateBuilder.Append("\t").Append("sql.Append(\"").Append(tabName).Append("\")\n")
	funcUpdateBuilder.Append("\t").Append("sql.Append(\" SET \")").Append("\n")
	funcUpdateBuilder.Append("\t").Append("ks := reflect.TypeOf(structParam)").Append("\n")
	funcUpdateBuilder.Append("\t").Append("vs := reflect.ValueOf(structParam)").Append("\n")
	funcUpdateBuilder.Append("\t").Append("var id interface{}").Append("\n")
	funcUpdateBuilder.Append("\t").Append("for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {").Append("\n")
	funcUpdateBuilder.Append("\t\t").Append("col := ks.Field(i).Tag.Get(\"column\")").Append("\n")
	funcUpdateBuilder.Append("\t\t").Append("v := vs.Field(i).Interface()").Append("\n")
	funcUpdateBuilder.Append("\t\t").Append("if col == \"id\" {").Append("\n")
	funcUpdateBuilder.Append("\t\t\t").Append("id = v").Append("\n")
	funcUpdateBuilder.Append("\t\t\t").Append("continue").Append("\n")
	funcUpdateBuilder.Append("\t\t").Append("}").Append("\n")
	funcUpdateBuilder.Append("\t\t").Append("sql.Append(col).Append(\"=\").Append(\"?,\")").Append("\n")
	funcUpdateBuilder.Append("\t\t").Append("params.Append(v)").Append("\n")
	funcUpdateBuilder.Append("\t").Append("}").Append("\n")
	funcUpdateBuilder.Append("\t").Append("sql.RemoveLast()").Append("\n")
	funcUpdateBuilder.Append("\t").Append("params.Append(id)").Append("\n")
	funcUpdateBuilder.Append("\t").Append("sql.Append(\" WHERE id = ?;\")").Append("\n")
	funcUpdateBuilder.Append("\t").Append("defer client.CloseConn()").Append("\n")
	funcUpdateBuilder.Append("\t").Append("return client.Exec(sql.ToString(), params.ToInterfaces()...)").Append("\n")
	funcUpdateBuilder.Append("}").Append("\n")
	tsgutils.Stdout(funcUpdateBuilder.ToString())
}

func (orm *ORMGenerator) buildORMSqlDelete(tabName string) {
	funcDeleteBuilder := tsgutils.NewStringBuilder()
	structName := getStructName(tabName)
	aliasStructName := getAliasStructName(tabName)
	funcDeleteBuilder.Append("func (").Append(aliasStructName).Append(" *").Append(structName).Append(") Delete").Append(structName).Append("ById(client *DBClient) (int64, error) {").Append("\n")
	funcDeleteBuilder.Append("\t").Append("structParam := ").Append(aliasStructName).Append("\n")
	funcDeleteBuilder.Append("\t").Append("sql := tsgutils.NewStringBuilder()").Append("\n")
	funcDeleteBuilder.Append("\t").Append("sql.Append(\"DELETE FROM \")").Append("\n")
	funcDeleteBuilder.Append("\t").Append("sql.Append(\"").Append(tabName).Append("\")\n")
	funcDeleteBuilder.Append("\t").Append("sql.Append(\" WHERE id = ?;\")").Append("\n")
	funcDeleteBuilder.Append("\t").Append("defer client.CloseConn()").Append("\n")
	funcDeleteBuilder.Append("\t").Append("return client.Exec(sql.ToString(), structParam.Id)").Append("\n")
	funcDeleteBuilder.Append("}").Append("\n")
	tsgutils.Stdout(funcDeleteBuilder.ToString())
}

func (orm *ORMGenerator) buildORMSqlBatchInsert(tabName string) {
	funcBatchInsertBuilder := tsgutils.NewStringBuilder()
	structName := getStructName(tabName)
	aliasStructName := getAliasStructName(tabName)
	structNames := getStructNames(tabName)
	funcBatchInsertBuilder.Append("func (").Append(aliasStructName).Append(" *").Append(structName).Append(") BatchInsert(client *DBClient, idSet, returnIds bool) ([]int64, error) {").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("structParam := *").Append(aliasStructName).Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("list := structParam.").Append(structNames).Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("var result []int64").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("listLen := len(list)").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("if listLen == 0 {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("return result, errors.New(\"no data needs to be inserted\")").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("sql := tsgutils.NewStringBuilder()").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("oneQSql := tsgutils.NewStringBuilder()").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("batchQSql := tsgutils.NewStringBuilder()").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("ks := reflect.TypeOf(structParam)").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("fieldsNum := ks.NumField() - 1").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("sql.Append(\"INSERT INTO \")").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("sql.Append(\"").Append(tabName).Append("\")\n")
	funcBatchInsertBuilder.Append("\t").Append("sql.Append(\" (\")").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("for i := 0; i < fieldsNum; i++ {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("iCol := ks.Field(i).Tag.Get(\"column\")").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("if iCol == \"id\" && !idSet {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("continue").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("sql.Append(\"`\").Append(iCol).Append(\"`,\")").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("sql.RemoveLast().Append(\") VALUES \")").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("batchInsertColsLen := tsgutils.InterfaceToInt(tsgutils.IIIInterfaceOperator(idSet, fieldsNum, fieldsNum-1))").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("oneQSql.Append(\"(\")").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("for j := 0; j < batchInsertColsLen; j++ {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("oneQSql.Append(\"?,\")").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("oneQSql.RemoveLast().Append(\")\")").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("if !returnIds {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("for j := 0; j < listLen; j++ {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("batchQSql.Append(oneQSql.ToString()).Append(\",\")").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("batchQSql.RemoveLast()").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("batchSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(batchQSql.ToString()).Append(\";\").ToString()").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("batchParams := tsgutils.NewInterfaceBuilder()").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("for k := range list {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("item := list[k]").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("kItem := reflect.ValueOf(item)").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("for l := 0; l < fieldsNum; l++ {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("lCol := ks.Field(l).Tag.Get(\"column\")").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("if lCol == \"id\" && !idSet {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t\t").Append("continue").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("batchParams.Append(kItem.Field(l).Interface())").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("id, err := client.Exec(batchSql, batchParams.ToInterfaces()...)").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("if err != nil {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("return result, err").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("result = append(result, id)").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("} else {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("oneSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(oneQSql.ToString()).Append(\";\").ToString()").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("oneParams := tsgutils.NewInterfaceBuilder()").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("tx, err := client.TxBegin()").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("if err != nil {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("return result, err").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("for m := range list {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("oneParams.Clear()").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("item := list[m]").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("mItem := reflect.ValueOf(item)").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("for n := 0; n < fieldsNum; n++ {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("nCol := ks.Field(n).Tag.Get(\"column\")").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("if nCol == \"id\" && !idSet {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t\t").Append("continue").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("oneParams.Append(mItem.Field(n).Interface())").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("id, err := client.TxExec(tx, oneSql, oneParams.ToInterfaces()...)").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("if err != nil {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("client.TxRollback(tx)").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("var resultTxRollback []int64").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t\t").Append("return resultTxRollback, err").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("result = append(result, id)").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("if !client.TxCommit(tx) {").Append("\n")
	funcBatchInsertBuilder.Append("\t\t\t").Append("return result, errors.New(\"batch insert (returnIds=true) tx commit failed\")").Append("\n")
	funcBatchInsertBuilder.Append("\t\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("}").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("defer client.CloseConn()").Append("\n")
	funcBatchInsertBuilder.Append("\t").Append("return result, nil").Append("\n")
	funcBatchInsertBuilder.Append("}").Append("\n")
	tsgutils.Stdout(funcBatchInsertBuilder.ToString())
}

func getStructName(tabName string) string {
	return tsgutils.FirstCaseToUpper(tabName, true)
}

func getAliasStructName(tabName string) string {
	return tsgutils.FirstCaseToUpper(tabName, false)
}

func getStructNames(tabName string) string {
	return tsgutils.NewString(getStructName(tabName)).AppendString("s").ToString()
}

func getAliasStructNames(tabName string) string {
	return tsgutils.NewString(getAliasStructName(tabName)).AppendString("s").ToString()
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

var ORMTabsCols []ORMTable
var warmTips tsgutils.InterfaceBuilder

type ORMTable struct {
	TName    string
	TComment string
	TColumns []ORMColumn
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
