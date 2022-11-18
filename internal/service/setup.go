package service

import (
	"github.com/gangming/sql2struct/internal/sql"
)

var (
	Generater sql.IGenerater
)

func Init() {
	Generater = sql.NewGenerater()
}
