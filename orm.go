package gorgo

import (
	"fmt"
	"log"
)

type DB struct {
	DialectDB Dialect
}

type FuncMap map[string]interface{}

func NewDB(config ConfigDB) (*DB, error) {
	dbtype := config.Type
	if dbtype == "" {
		return nil, fmt.Errorf("DB_TYPE not defined.")
	}

	config.Validations = GetFunctions()

	switch dbtype {
	case "mongo":
		dialect := &MongoDialect{}
		err := dialect.InitDB(config)
		if err != nil {
			return nil, err
		}
		return &DB{DialectDB: dialect}, nil
	case "localdb":
		dialect := &LocalDialect{}
		err := dialect.InitDB(config)
		if err != nil {
			return nil, err
		}
		return &DB{DialectDB: dialect}, nil
	default:
		return nil, fmt.Errorf("[WARNING] dbtype not found", nil)
	}

	return nil, nil
}

func (d *DB) Create(collection string, data JSONDoc) (JSONDoc, error) {
	return d.DialectDB.Create(collection, data)
}

func (d *DB) Close() error {
	return d.DialectDB.CloseDB()
}

func (d *DB) Update(collection string, data JSONDoc) error {
	return d.DialectDB.Update(collection, data)
}

func (d *DB) Count(collection string) int {

	i, err := d.DialectDB.Count(collection)
	if err != nil {
		log.Println(err)
	}
	return i
}

func (d *DB) Delete(collection string, id string) error {
	return d.DialectDB.Delete(collection, id)

}

func (d *DB) GetAll(collection string, page int, qtd int, sorted string) ([]JSONDoc, error) {
	return d.DialectDB.GetAll(collection, page, qtd, sorted)

}

func (d *DB) GetByID(collection string, id string) (JSONDoc, error) {
	return d.DialectDB.GetById(collection, id)

}

func (d *DB) GetOneByQuery(collectin string, query string) (JSONDoc, error) {
	return d.DialectDB.GetOneByQuery(collectin, query)
}
