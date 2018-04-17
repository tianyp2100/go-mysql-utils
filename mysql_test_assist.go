package tsgmysqlutils

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
