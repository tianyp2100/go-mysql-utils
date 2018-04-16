package tsgmysqlutils

import (
	mysql "database/sql"
)

type ORMBase interface {
	Row2struct(row *mysql.Row)
	Rows2struct(rows *mysql.Rows)
}
