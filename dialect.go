package gorgo

import "encoding/json"

//DialectSQL : interface for sql database dialect
type Dialect interface {
	InitDB(ConfigDB) error
	CloseDB() error
	Count(string) (int, error)
	Create(string, JSONDoc) (JSONDoc,error)
	GetById(string, string) (JSONDoc, error)
	GetOneByQuery(string, map[string]interface{}) (JSONDoc, error)
	GetManyByQuery(string, map[string]interface{}) ([]JSONDoc, error)
	GetAll(string, int, int, string) ([]JSONDoc, error)
	GetAllBySearch(string, string, string, int, int, string) ([]JSONDoc, error)
	Update(string, JSONDoc) error
	Delete(string, string) error
	GetByGroup(string, map[string]interface{}) (JSONDoc, error)
}

//JSONDoc map string for interfaces like json
type JSONDoc map[string]interface{}
func (j *JSONDoc) ToString() string{
	encoded, e := json.Marshal(j)
	if e != nil {
		return ""
	}
	return string(encoded)
}