package tsgmysqlutils

import (
	"time"
	db "database/sql"
	"github.com/timespacegroup/go-utils"
)

/*
 string utils
 @author Tony Tian
 @date 2018-04-17
 @version 1.0.0
*/

func TestDbClient() *DBClient {
	var dbConfig DBConfig

	dbConfig.DbHost = "127.0.0.1"
	dbConfig.DbUser = "root"
	dbConfig.DbPass = "123456"
	dbConfig.IsLocalTime = true
	dbConfig.DbName = "test"

	return NewDbClient(dbConfig)
}

func TestORMConfig() *ORMConfig {
	var dbConfig DBConfig

	dbConfig.DbHost = "127.0.0.1"
	dbConfig.DbUser = "root"
	dbConfig.DbPass = "123456"
	dbConfig.IsLocalTime = true
	dbConfig.DbName = "test"

	var ormConfig ORMConfig
	ormConfig.DbConfig = dbConfig
	ormConfig.TabName = []string{"we_test_tab1", "we_test_tab2"}

	return &ormConfig
}

type WeTestTab1 struct {
	Id           int64     `column:"id"`
	Name         string    `column:"name"`
	Gender       int64     `column:"gender"`
	Birthday     time.Time `column:"birthday"`
	Stature      float64   `column:"stature"`
	Weight       float64   `column:"weight"`
	CreatedTime  time.Time `column:"created_time"`
	ModifiedTime time.Time `column:"modified_time"`
	IsDeleted    int64     `column:"is_deleted"`
	WeTestTab1s  [] WeTestTab1
}

func (weTestTab1 *WeTestTab1) RowToStruct(row *db.Row) {
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
	tsgutils.CheckAndPrintError("MySQL query row scan error", err)
}

func (weTestTab1 *WeTestTab1) RowsToStruct(rows *db.Rows) {
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
		tsgutils.CheckAndPrintError("MySQL query rows scan error", err)
		weTestTab1s = append(weTestTab1s, *weTestTab1)
	}
	if rows != nil {
		defer rows.Close()
	}
	weTestTab1.WeTestTab1s = weTestTab1s
}

type WeTestTab2 struct {
	Id                 int64     `column:"id"`
	UserId             int64     `column:"user_id"`
	AreaCode           int64     `column:"area_code"`
	Phone              int64     `column:"phone"`
	Email              string    `column:"email"`
	Postcode           int64     `column:"postcode"`
	AdministrationCode int64     `column:"administration_code"`
	Address            string    `column:"address"`
	CreatedTime        time.Time `column:"created_time"`
	ModifiedTime       time.Time `column:"modified_time"`
	IsDeleted          int64     `column:"is_deleted"`
	WeTestTab2s        [] WeTestTab2
}

func (weTestTab2 *WeTestTab2) RowToStruct(row *db.Row) {
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
	tsgutils.CheckAndPrintError("MySQL query row scan error", err)
}

func (weTestTab2 *WeTestTab2) RowsToStruct(rows *db.Rows) {
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
		tsgutils.CheckAndPrintError("MySQL query rows scan error", err)
		weTestTab2s = append(weTestTab2s, *weTestTab2)
	}
	if rows != nil {
		defer rows.Close()
	}
	weTestTab2.WeTestTab2s = weTestTab2s
}
