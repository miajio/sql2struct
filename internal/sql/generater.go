package sql

import (
	"bytes"
	"github.com/gangming/sql2struct/internal/service"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gangming/sql2struct/config"
	"github.com/gangming/sql2struct/internal/dto"
	"github.com/gangming/sql2struct/utils"
)

type IGenerater interface {
	Generate(t *dto.Table) string
}

type generater struct {
}

func NewGenerater() IGenerater {
	return &generater{}
}

func (g *generater) Generate(t *dto.Table) string {
	tableName := config.Cnf.TablePrefix + t.Name
	fromTmpl := dto.Tmpl
	fromTmpl = strings.Replace(fromTmpl, "{{.Package}}", "model", -1)
	fromTmpl = strings.Replace(fromTmpl, "{{.Table.Name}}", tableName, -1)
	fromTmpl = strings.Replace(fromTmpl, "{{.Table.Comment}}", t.Comment, -1)
	for i, field := range t.Fields {
		tag := "`" + config.Cnf.DBTag + ":\"column:" + field.Name + "\""
		if field.IsPK {
			tag += ";primary_key\" "
		}
		if field.DefaultValue != "" {
			tag += ";default:" + field.DefaultValue + "\""
		}
		if config.Cnf.WithJsonTag {
			tag += " json:\"" + field.Name + "\""
		}
		tag += "`"
		t.Fields[i].Tag = tag
	}
	tl := template.Must(template.New("tmpl").Parse(fromTmpl))
	var res bytes.Buffer
	err := tl.Execute(&res, t)
	if err != nil {
		panic(err)
	}
	return utils.CommonInitialisms(res.String())
}
func (g *generater) GenerateFile(ddl string) error {
	c, _ := SqlDriver.ParseDDL(ddl)
	dir := config.Cnf.OutputDir
	os.MkdirAll(dir, 0755)
	fileName := filepath.Join(dir, strings.ToLower(c.Name)+".go")
	fd, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}
	defer fd.Close()
	fd.Write([]byte(service.Generater.Generate(c)))

	_, err = exec.Command("goimports", "-l", "-w", dir).Output()
	if err != nil {
		utils.PrintRed(err.Error())
	}
	_, err = exec.Command("gofmt", "-l", "-w", dir).Output()
	if err != nil {
		utils.PrintRed(err.Error())
	}
	utils.PrintGreen(fileName + " generate success")
	return nil
}
