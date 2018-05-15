package gorgo

import (
	"github.com/rgobbo/fileutils"
	"github.com/spf13/cast"
)

type model struct {
	Schema string
	Tables []table
}

type table struct {
	Name string
	Fields []fields
}

type fields struct {
	Name string
	Type string
	Autoincrement bool
	Unique bool
	Default interface{}
	Alias string
}

func (m* model) LoadFile(path string) error{
	var obj map[string] interface{}
	err := fileutils.LoadJson(path, &obj)
	if err != nil{
		return err
	}

	m.Schema = cast.ToString(obj["schema"])

	tables := cast.ToSlice(obj["tables"])
	for _ ,t := range tables {
		newt := table{}
		newt.Name = cast.ToString(t["name"])

	}

	return nil
}

