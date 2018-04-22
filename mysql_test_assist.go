package tsgmysqlutils

/*
 string utils
 @author Tony Tian
 @date 2018-04-17
 @version 1.0.0
*/

import (
	db "database/sql"
	"errors"
	"github.com/timespacegroup/go-utils"
	"reflect"
	"time"
)

func TestDbClient() *DBClient {
	var dbConfig DBConfig

	dbConfig.DbHost = "127.0.0.1"
	dbConfig.DbUser = "root"
	dbConfig.DbPass = "123456"
	dbConfig.IsLocalTime = true
	dbConfig.DbName = "test"

	return NewDbClient(dbConfig)
}

/*
	test table1
*/
type WeTestTab1 struct {
	Id           int64     `column:"id"`            // The primary key id
	Name         string    `column:"name"`          // The user name
	Gender       int64     `column:"gender"`        // The user gerder, 1:male 2:female 0:default
	Birthday     time.Time `column:"birthday"`      // The user birthday, eg: 2018-04-16
	Stature      float64   `column:"stature"`       // The user stature, eg: 172.22cm
	Weight       float64   `column:"weight"`        // The user weight, eg: 21.77kg
	CreatedTime  time.Time `column:"created_time"`  // created time
	ModifiedTime time.Time `column:"modified_time"` // record time
	IsDeleted    int64     `column:"is_deleted"`    // Logic to delete(0:normal 1:deleted)
	WeTestTab1s  [] WeTestTab1                      // This value is used for batch queries and inserts.
}

func (weTestTab1 *WeTestTab1) RowToStruct(row *db.Row) error {
	builder := tsgutils.NewInterfaceBuilder()
	builder.Append(&weTestTab1.Id)
	builder.Append(&weTestTab1.Name)
	builder.Append(&weTestTab1.Gender)
	builder.Append(&weTestTab1.Birthday)
	builder.Append(&weTestTab1.Stature)
	builder.Append(&weTestTab1.Weight)
	builder.Append(&weTestTab1.CreatedTime)
	builder.Append(&weTestTab1.ModifiedTime)
	builder.Append(&weTestTab1.IsDeleted)
	err := row.Scan(builder.ToInterfaces()...)
	if err != nil {
		return err
	}
	return nil
}

func (weTestTab1 *WeTestTab1) RowsToStruct(rows *db.Rows) error {
	var weTestTab1s [] WeTestTab1
	builder := tsgutils.NewInterfaceBuilder()
	for rows.Next() {
		builder.Clear()
		builder.Append(&weTestTab1.Id)
		builder.Append(&weTestTab1.Name)
		builder.Append(&weTestTab1.Gender)
		builder.Append(&weTestTab1.Birthday)
		builder.Append(&weTestTab1.Stature)
		builder.Append(&weTestTab1.Weight)
		builder.Append(&weTestTab1.CreatedTime)
		builder.Append(&weTestTab1.ModifiedTime)
		builder.Append(&weTestTab1.IsDeleted)
		err := rows.Scan(builder.ToInterfaces()...)
		if err != nil {
			return err
		}
		weTestTab1s = append(weTestTab1s, *weTestTab1)
	}
	if rows != nil {
		defer rows.Close()
	}
	weTestTab1.WeTestTab1s = weTestTab1s
	return nil
}

func (weTestTab1 *WeTestTab1) Insert(client *DBClient, idSet bool) (int64, error) {
	structParam := *weTestTab1
	sql := tsgutils.NewStringBuilder()
	qSql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("INSERT INTO ")
	sql.Append("we_test_tab1")
	sql.Append(" (")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(col).Append("`,")
		qSql.Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	qSql.RemoveLast()
	sql.Append(") VALUES (").Append(qSql.ToString()).Append(");")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (weTestTab1 *WeTestTab1) UpdateWeTestTab1ById(client *DBClient) (int64, error) {
	structParam := *weTestTab1
	sql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("UPDATE ")
	sql.Append("we_test_tab1")
	sql.Append(" SET ")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	var id interface{}
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" {
			id = v
			continue
		}
		sql.Append(col).Append("=").Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	params.Append(id)
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (weTestTab1 *WeTestTab1) DeleteWeTestTab1ById(client *DBClient) (int64, error) {
	structParam := weTestTab1
	sql := tsgutils.NewStringBuilder()
	sql.Append("DELETE FROM ")
	sql.Append("we_test_tab1")
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), structParam.Id)
}

func (weTestTab1 *WeTestTab1) BatchInsert(client *DBClient, idSet, returnIds bool) ([]int64, error) {
	structParam := *weTestTab1
	list := structParam.WeTestTab1s
	var result []int64
	listLen := len(list)
	if listLen == 0 {
		return result, errors.New("no data needs to be inserted")
	}
	sql := tsgutils.NewStringBuilder()
	oneQSql := tsgutils.NewStringBuilder()
	batchQSql := tsgutils.NewStringBuilder()
	ks := reflect.TypeOf(structParam)
	fieldsNum := ks.NumField() - 1
	sql.Append("INSERT INTO ")
	sql.Append("we_test_tab1")
	sql.Append(" (")
	for i := 0; i < fieldsNum; i++ {
		iCol := ks.Field(i).Tag.Get("column")
		if iCol == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(iCol).Append("`,")
	}
	sql.RemoveLast().Append(") VALUES ")
	batchInsertColsLen := tsgutils.InterfaceToInt(tsgutils.IIIInterfaceOperator(idSet, fieldsNum, fieldsNum-1))
	oneQSql.Append("(")
	for j := 0; j < batchInsertColsLen; j++ {
		oneQSql.Append("?,")
	}
	oneQSql.RemoveLast().Append(")")
	if !returnIds {
		for j := 0; j < listLen; j++ {
			batchQSql.Append(oneQSql.ToString()).Append(",")
		}
		batchQSql.RemoveLast()
		batchSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(batchQSql.ToString()).Append(";").ToString()
		batchParams := tsgutils.NewInterfaceBuilder()
		for k := range list {
			item := list[k]
			kItem := reflect.ValueOf(item)
			for l := 0; l < fieldsNum; l++ {
				lCol := ks.Field(l).Tag.Get("column")
				if lCol == "id" && !idSet {
					continue
				}
				batchParams.Append(kItem.Field(l).Interface())
			}
		}
		id, err := client.Exec(batchSql, batchParams.ToInterfaces()...)
		if err != nil {
			return result, err
		}
		result = append(result, id)
	} else {
		oneSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(oneQSql.ToString()).Append(";").ToString()
		oneParams := tsgutils.NewInterfaceBuilder()
		tx, err := client.TxBegin()
		if err != nil {
			return result, err
		}
		for m := range list {
			oneParams.Clear()
			item := list[m]
			mItem := reflect.ValueOf(item)
			for n := 0; n < fieldsNum; n++ {
				nCol := ks.Field(n).Tag.Get("column")
				if nCol == "id" && !idSet {
					continue
				}
				oneParams.Append(mItem.Field(n).Interface())
			}
			id, err := client.TxExec(tx, oneSql, oneParams.ToInterfaces()...)
			if err != nil {
				client.TxRollback(tx)
				var resultTxRollback []int64
				return resultTxRollback, err
			}
			result = append(result, id)
		}
		if !client.TxCommit(tx) {
			return result, errors.New("batch insert (returnIds=true) tx commit failed")
		}
	}
	defer client.CloseConn()
	return result, nil
}

/*
	test table2
*/
type WeTestTab2 struct {
	Id                 int64     `column:"id"`                  // The primary key id
	UserId             int64     `column:"user_id"`             // The user id
	AreaCode           int64     `column:"area_code"`           // The user area code
	Phone              int64     `column:"phone"`               // The user phone
	Email              string    `column:"email"`               // The user email
	Postcode           int64     `column:"postcode"`            // The user postcode
	AdministrationCode int64     `column:"administration_code"` // The user administration code
	Address            string    `column:"address"`             // The user address
	CreatedTime        time.Time `column:"created_time"`        // created time
	ModifiedTime       time.Time `column:"modified_time"`       // modified time
	IsDeleted          int64     `column:"is_deleted"`          // Logic to delete(0:normal 1:deleted)
	WeTestTab2s        [] WeTestTab2                            // This value is used for batch queries and inserts.
}

func (weTestTab2 *WeTestTab2) RowToStruct(row *db.Row) error {
	builder := tsgutils.NewInterfaceBuilder()
	builder.Append(&weTestTab2.Id)
	builder.Append(&weTestTab2.UserId)
	builder.Append(&weTestTab2.AreaCode)
	builder.Append(&weTestTab2.Phone)
	builder.Append(&weTestTab2.Email)
	builder.Append(&weTestTab2.Postcode)
	builder.Append(&weTestTab2.AdministrationCode)
	builder.Append(&weTestTab2.Address)
	builder.Append(&weTestTab2.CreatedTime)
	builder.Append(&weTestTab2.ModifiedTime)
	builder.Append(&weTestTab2.IsDeleted)
	err := row.Scan(builder.ToInterfaces()...)
	if err != nil {
		return err
	}
	return nil
}

func (weTestTab2 *WeTestTab2) RowsToStruct(rows *db.Rows) error {
	var weTestTab2s [] WeTestTab2
	builder := tsgutils.NewInterfaceBuilder()
	for rows.Next() {
		builder.Clear()
		builder.Append(&weTestTab2.Id)
		builder.Append(&weTestTab2.UserId)
		builder.Append(&weTestTab2.AreaCode)
		builder.Append(&weTestTab2.Phone)
		builder.Append(&weTestTab2.Email)
		builder.Append(&weTestTab2.Postcode)
		builder.Append(&weTestTab2.AdministrationCode)
		builder.Append(&weTestTab2.Address)
		builder.Append(&weTestTab2.CreatedTime)
		builder.Append(&weTestTab2.ModifiedTime)
		builder.Append(&weTestTab2.IsDeleted)
		err := rows.Scan(builder.ToInterfaces()...)
		if err != nil {
			return err
		}
		weTestTab2s = append(weTestTab2s, *weTestTab2)
	}
	if rows != nil {
		defer rows.Close()
	}
	weTestTab2.WeTestTab2s = weTestTab2s
	return nil
}

func (weTestTab2 *WeTestTab2) Insert(client *DBClient, idSet bool) (int64, error) {
	structParam := *weTestTab2
	sql := tsgutils.NewStringBuilder()
	qSql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("INSERT INTO ")
	sql.Append("we_test_tab2")
	sql.Append(" (")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(col).Append("`,")
		qSql.Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	qSql.RemoveLast()
	sql.Append(") VALUES (").Append(qSql.ToString()).Append(");")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (weTestTab2 *WeTestTab2) UpdateWeTestTab2ById(client *DBClient) (int64, error) {
	structParam := *weTestTab2
	sql := tsgutils.NewStringBuilder()
	params := tsgutils.NewInterfaceBuilder()
	sql.Append("UPDATE ")
	sql.Append("we_test_tab2")
	sql.Append(" SET ")
	ks := reflect.TypeOf(structParam)
	vs := reflect.ValueOf(structParam)
	var id interface{}
	for i, ksLen := 0, ks.NumField()-1; i < ksLen; i++ {
		col := ks.Field(i).Tag.Get("column")
		v := vs.Field(i).Interface()
		if col == "id" {
			id = v
			continue
		}
		sql.Append(col).Append("=").Append("?,")
		params.Append(v)
	}
	sql.RemoveLast()
	params.Append(id)
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), params.ToInterfaces()...)
}

func (weTestTab2 *WeTestTab2) DeleteWeTestTab2ById(client *DBClient) (int64, error) {
	structParam := weTestTab2
	sql := tsgutils.NewStringBuilder()
	sql.Append("DELETE FROM ")
	sql.Append("we_test_tab2")
	sql.Append(" WHERE id = ?;")
	defer client.CloseConn()
	return client.Exec(sql.ToString(), structParam.Id)
}

func (weTestTab2 *WeTestTab2) BatchInsert(client *DBClient, idSet, returnIds bool) ([]int64, error) {
	structParam := *weTestTab2
	list := structParam.WeTestTab2s
	var result []int64
	listLen := len(list)
	if listLen == 0 {
		return result, errors.New("no data needs to be inserted")
	}
	sql := tsgutils.NewStringBuilder()
	oneQSql := tsgutils.NewStringBuilder()
	batchQSql := tsgutils.NewStringBuilder()
	ks := reflect.TypeOf(structParam)
	fieldsNum := ks.NumField() - 1
	sql.Append("INSERT INTO ")
	sql.Append("we_test_tab2")
	sql.Append(" (")
	for i := 0; i < fieldsNum; i++ {
		iCol := ks.Field(i).Tag.Get("column")
		if iCol == "id" && !idSet {
			continue
		}
		sql.Append("`").Append(iCol).Append("`,")
	}
	sql.RemoveLast().Append(") VALUES ")
	batchInsertColsLen := tsgutils.InterfaceToInt(tsgutils.IIIInterfaceOperator(idSet, fieldsNum, fieldsNum-1))
	oneQSql.Append("(")
	for j := 0; j < batchInsertColsLen; j++ {
		oneQSql.Append("?,")
	}
	oneQSql.RemoveLast().Append(")")
	if !returnIds {
		for j := 0; j < listLen; j++ {
			batchQSql.Append(oneQSql.ToString()).Append(",")
		}
		batchQSql.RemoveLast()
		batchSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(batchQSql.ToString()).Append(";").ToString()
		batchParams := tsgutils.NewInterfaceBuilder()
		for k := range list {
			item := list[k]
			kItem := reflect.ValueOf(item)
			for l := 0; l < fieldsNum; l++ {
				lCol := ks.Field(l).Tag.Get("column")
				if lCol == "id" && !idSet {
					continue
				}
				batchParams.Append(kItem.Field(l).Interface())
			}
		}
		id, err := client.Exec(batchSql, batchParams.ToInterfaces()...)
		if err != nil {
			return result, err
		}
		result = append(result, id)
	} else {
		oneSql := tsgutils.NewStringBuilder().Append(sql.ToString()).Append(oneQSql.ToString()).Append(";").ToString()
		oneParams := tsgutils.NewInterfaceBuilder()
		tx, err := client.TxBegin()
		if err != nil {
			return result, err
		}
		for m := range list {
			oneParams.Clear()
			item := list[m]
			mItem := reflect.ValueOf(item)
			for n := 0; n < fieldsNum; n++ {
				nCol := ks.Field(n).Tag.Get("column")
				if nCol == "id" && !idSet {
					continue
				}
				oneParams.Append(mItem.Field(n).Interface())
			}
			id, err := client.TxExec(tx, oneSql, oneParams.ToInterfaces()...)
			if err != nil {
				client.TxRollback(tx)
				var resultTxRollback []int64
				return resultTxRollback, err
			}
			result = append(result, id)
		}
		if !client.TxCommit(tx) {
			return result, errors.New("batch insert (returnIds=true) tx commit failed")
		}
	}
	defer client.CloseConn()
	return result, nil
}
