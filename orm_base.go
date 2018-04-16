package tsgmysqlutils

/*
 string utils
 @author Tony Tian
 @date 2018-04-17
 @version 1.0.0
*/

import (
	mysql "database/sql"
)

type ORMBase interface {
	Row2struct(row *mysql.Row)
	Rows2struct(rows *mysql.Rows)
}
