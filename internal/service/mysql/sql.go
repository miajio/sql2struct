package mysqlparser

import (
	"fmt"
	"github.com/gangming/sql2struct/internal/service"
	"github.com/gangming/sql2struct/internal/sql"
	"strings"

	"github.com/gangming/sql2struct/config"
	"github.com/gangming/sql2struct/internal/dto"
	"github.com/gangming/sql2struct/internal/infra"
	"github.com/gangming/sql2struct/utils"
)

var MysqlType2GoType = map[string]string{
	"int":       "int64",
	"tinyint":   "uint8",
	"bigint":    "int64",
	"varchar":   "string",
	"text":      "string",
	"date":      "time.Time",
	"time":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
	"json":      "string",
}

type mysql struct{}

func NewMysql() sql.ISqlDriver {
	return &mysql{}
}

func (m *mysql) ParseDDL(s string) (*dto.Table, error) {
	lines := strings.Split(s, "\n")
	var table dto.Table
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "CREATE TABLE") {
			tableName := strings.Split(line, "`")[1]
			table.Name = tableName
			table.UpperCamelCaseName = utils.Underline2UpperCamelCase(tableName)
			continue
		}
		if strings.Contains(line, "ENGINE") && strings.Contains(line, "COMMENT=") {
			table.Comment = strings.Trim(strings.Split(line, "COMMENT='")[1], "'")
			fmt.Println(table.Comment)
			continue
		}
		if line[0] == '`' {
			field := dto.Field{}
			field.Name = strings.Split(line, "`")[1]
			field.UpperCamelCaseName = utils.Underline2UpperCamelCase(field.Name)
			field.Type = strings.TrimRightFunc(strings.Split(line, " ")[1], func(r rune) bool {
				return r < 'a' || r > 'z'
			})
			field.Type = MysqlType2GoType[field.Type]
			if strings.Contains(field.Type, "time") {
				table.ContainsTimeField = true
			}
			if strings.Contains(line, "COMMENT") {
				field.Comment = strings.Trim(strings.Split(line, "COMMENT '")[1], "',")
			}
			if strings.Contains(line, "DEFAULT'") {
				field.DefaultValue = strings.Split(line, "DEFAULT ")[1]
			}
			if strings.Contains(line, "PRIMARY KEY") {
				field.IsPK = true
			}

			table.Fields = append(table.Fields, field)

		}
	}
	return &table, nil
}

func GetDDLs() ([]string, error) {
	var result []string
	tables := GetTables()
	for _, tableName := range tables {
		rows, err := infra.GetDB().Query("show create table " + tableName)
		if err != nil {
			panic(err)
		}

		if rows.Next() {
			var r string
			err := rows.Scan(&tableName, &r)
			if err != nil {
				panic(err)
			}
			result = append(result, r)
		}
	}

	return result, nil
}
func GetTables() []string {
	if len(config.Cnf.Tables) > 0 {
		return config.Cnf.Tables
	}
	var result []string
	rows, err := infra.GetDB().Query("show tables")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var r string
		err := rows.Scan(&r)
		if err != nil {
			panic(err)
		}
		result = append(result, r)
	}
	return result
}
func (m *mysql) Execute() error {
	ddls, err := GetDDLs()
	if err != nil {
		return err
	}
	for _, ddl := range ddls {
		table, err := m.ParseDDL(ddl)
		if err != nil {
			return err
		}
		service.Generater.Generate(table)
	}
	return nil
}
