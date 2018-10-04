package gorgo

import "encoding/json"

//DialectSQL : interface for sql database dialect
type Dialect interface {
	InitDB(ConfigDB) error
	CloseDB() error
	Count(string) (int, error)
	Create(string, JSONDoc) (JSONDoc, error)
	CreateInterface(string, interface{}) error
	GetById(string, string) (JSONDoc, error)
	GetOneByQuery(string, string) (JSONDoc, error)
	GetManyByQuery(string, string, ...interface{}) ([]JSONDoc, error)
	GetAll(string, int, int, string) ([]JSONDoc, error)
	GetAllBySearch(string, string, string, int, int, string) ([]JSONDoc, error)
	Update(string, JSONDoc) error
	Delete(string, string) error
	DeleteByWhere(string, string) error
	CountByWhere(string, string) (int, error)
	GetByGroup(string, map[string]interface{}) (JSONDoc, error)
}

//JSONDoc map string for interfaces like json
type JSONDoc map[string]interface{}

func (j *JSONDoc) ToString() string {
	encoded, e := json.Marshal(j)
	if e != nil {
		return ""
	}
	return string(encoded)
}
