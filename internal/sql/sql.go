package sql

import (
	"github.com/gangming/sql2struct/internal/dto"
)

type ISqlDriver interface {
	ParseDDL(s string) (*dto.Table, error)
	Execute() error
}
