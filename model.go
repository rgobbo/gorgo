package gorgo

import (
	"github.com/rgobbo/fileutils"
	"strings"
	"fmt"
	"github.com/spf13/cast"
)

type model struct {
	Schema string
	Tables map[string]table
}

type table struct {
	Name string
	Fields []*field
}

type field struct {
	Name string
	Type string
	Autoincrement bool
	Unique bool
	Default interface{}
	Alias string
	Maxlen int
	Minlen int
	Required bool
	Validation string
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
	tables := make(map[string]table)
	for _ ,t := range conf.Tables{
		newTable := table{}
		newTable.Name = t.Name
		fields := []*field{}
		for _, s := range t.Fields {
			newField, err := parseField(s)
			if err != nil {
				return err
			}
			fields = append(fields,newField)
		}
		newTable.Fields = fields
		tables[t.Name] = newTable
	}
	m.Tables = tables

	return nil
}

func parseField (str string) (*field, error) {
	newField := &field{}
	parts := strings.Split(str,",")

	if len(parts) < 2 {
		return nil, fmt.Errorf("Field must have two properties (name,type)")
	}
	for i, s := range parts {
		if i == 0 {
			newField.Name = s
		} else if i == 1 {
			fieldType := strings.ToLower(s)
			fieldType = strings.Trim(s," ")
			switch fieldType {
			case "string" :
				newField.Type = fieldType
			case "date" :
				newField.Type = fieldType
			case "int" :
				newField.Type = fieldType
			case "bigint" :
				newField.Type = fieldType
			case "float" :
				newField.Type = fieldType
			case "double" :
				newField.Type = fieldType
			case "varchar" :
				newField.Type = fieldType
			default:
				newField.Type = "string"
			}
		} else {
			slc := strings.ToLower(s)
			slc = strings.Trim(slc," ")
			if slc == "autoincrement" {
				newField.Autoincrement = true
			}
			if slc == "unique" {
				newField.Unique = true
			}
			if strings.HasPrefix(slc,"minlen") {
				sints := strings.Split(slc,"=")
				i, err := cast.ToIntE(strings.Trim(sints[1]," "))
				if err != nil {
					return nil, err
				}
				newField.Minlen = i
			}
			if strings.HasPrefix(s,"maxlen") {
				sints := strings.Split(slc,"=")
				i, err := cast.ToIntE(strings.Trim(sints[1]," "))
				if err != nil {
					return nil, err
				}
				newField.Maxlen = i
			}
			if slc == "required"  {
				newField.Required = true
			}
			if strings.HasPrefix(slc,"alias") {
				ss := strings.Split(strings.Trim(s," "),"=")
				newField.Alias = strings.Trim(ss[1]," ")
			}
			if strings.HasPrefix(slc,"default") {
				ss := strings.Split(slc,"=")
				newField.Default = strings.Trim(ss[1]," ")
			}
			if strings.HasPrefix(slc,"validation") {
				ss := strings.Split(slc,"=")
				newField.Validation = strings.Trim(ss[1]," ")
			}
		}

	}

	return newField, nil
}

