# Go MySQL utils

##### Two main functions: database operation and automatic ORM(object(struct) relational mapping).

---------------------------------------
### Installation
```
$ go get -u github.com/timespacegroup/go-utils
$ go get -u github.com/timespacegroup/go-mysql-utils
```
### Usage

##### 1. Create a new mysql client:
```
func TestDbClient() *DBClient {
	var dbConfig tsgmysqlutils.DBConfig
	dbConfig.DbHost = "127.0.0.1"
	dbConfig.DbUser = "root"
	dbConfig.DbPass = "123456"
	dbConfig.IsLocalTime = true
	dbConfig.DbName = "test"
	return tsgmysqlutils.NewDbClient(dbConfig)
}
```
##### 2. Create Object(struct) Relational Mapping:
```
func TestGenerateORM(t *testing.T) {
	client := TestDbClient()
	orm := tsgmysqlutils.NewORMGenerator(client)
	orm.AddComment = true
	tabNames := []string{"we_test_tab1", "we_test_tab2"}
	orm.DefaultGenerator(tabNames)
	client.CloseConn()
}
```
##### More info see:
###### See the client operation: mysql_test.go
###### See the orm result: mysql_test_assist.go
###### https://blog.csdn.net/typa01_kk/article/category/7629914
