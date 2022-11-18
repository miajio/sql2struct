package infra

import (
	"database/sql"
	"time"

	"github.com/gangming/sql2struct/config"
	_ "github.com/go-sql-driver/mysql"
)

var pool *sql.DB

func InitDB() {
	var err error
	pool, err = sql.Open(string(config.Cnf.DBType), config.Cnf.DSN)
	if err != nil {
		panic(err)
	}
	pool.SetMaxOpenConns(100)
	pool.SetMaxIdleConns(20)
	pool.SetConnMaxLifetime(100 * time.Second)

}
func GetDB() *sql.DB {
	return pool
}
func Init() {
	InitDB()
}
