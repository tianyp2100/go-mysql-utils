package tsgmysqlutils

func TestDbClient() *DBClient {
	var dbConfig DBConfig

	dbConfig.DbHost = "127.0.0.1"
	dbConfig.DbUser = "root"
	dbConfig.DbPass = "123456"
	dbConfig.DbName = "test"

	return NewDbClient(dbConfig)
}
