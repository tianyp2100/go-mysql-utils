package tsgmysqlutils

import "time"

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
