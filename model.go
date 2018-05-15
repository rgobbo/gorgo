package gorgo

import (
	"github.com/rgobbo/fileutils"
	"github.com/spf13/cast"
	"strings"
)

type model struct {
	Schema string
	Tables []table
}

type table struct {
	Name string
	Fields []field
}

type field struct {
	Name string
	Type string
	Autoincrement bool
	Unique bool
	Default interface{}
	Alias string
}

type modelConfig struct {
	Schema string
	Tables []tableConfig
}

type tableConfig struct {
	Name string
	Fields []string
}

func (m* model) LoadFile(path string) error{
	var conf modelConfig
	err := fileutils.LoadJson(path, &conf)
	if err != nil{
		return err
	}

	m.Schema = conf.Schema
	tables := [] table{}
	for _ ,t := range conf.Tables{
		newTable := table{}
		fields := []field{}
		for _, s := range t.Fields {
			newField := field{}
			parts := strings.Split(s," ")
			newField.Name = parts[0]
			newField.Type = parts[1]
			fields = append(fields,newField)
		}
		tables = append(tables,newTable)
	}
	m.Tables = tables

	return nil
}

