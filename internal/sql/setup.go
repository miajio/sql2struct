package sql

import (
	"github.com/gangming/sql2struct/config"
	mysqlparser "github.com/gangming/sql2struct/internal/service/mysql"
	"github.com/gangming/sql2struct/internal/service/postgresql"
)

var (
	SqlDriver ISqlDriver
)

func Init() {
	SqlDriver = NewSqlDriver()
}

func NewSqlDriver() ISqlDriver {
	if config.Cnf.DBType == "mysql" {
		return mysqlparser.NewMysql()
	}
	return postgresql.NewPsql()
}
