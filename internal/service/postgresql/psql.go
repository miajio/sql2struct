package postgresql

import (
	"github.com/gangming/sql2struct/internal/dto"
	"github.com/gangming/sql2struct/internal/sql"
)

type psql struct{}

func NewPsql() sql.ISqlDriver {
	return &psql{}
}
func (p psql) ParseDDL(s string) (*dto.Table, error) {
	//TODO implement me
	panic("implement me")
}

func (p psql) Execute() error {
	//TODO implement me
	panic("implement me")
}
