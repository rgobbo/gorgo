package gorgo

import (
	"fmt"
)

type DB struct {
	dialectDB Dialect
	showSQL  bool
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
		return &DB{dialectDB: dialect}, nil
	case "localdb":
		dialect := &LocalDialect{}
		err := dialect.InitDB(config)
		if err != nil {
			return nil, err
		}
		return &DB{dialectDB: dialect}, nil
	default:
		return nil, fmt.Errorf("[WARNING] dbtype not found", nil)
	}

	return nil, nil
}

func (d *DB)NewSession() *Session {
	session := &Session{db: d}
	session.Init()
	return session
}

func (d *DB) Table(tbl string) *Session  {
	session := d.NewSession()
	session.tableName = tbl
	return session
}

func (d *DB) Limit(i int) *Session  {
	session := d.NewSession()
	session.limit = i
	return session
}

func (d *DB) Offset(i int) *Session  {
	session := d.NewSession()
	session.offset = i
	return session
}

func (d *DB) Where(query string, params ...interface{}) *Session {
	session := d.NewSession()
	session.where = query
	session.params = params
	return session
}

// Query a raw sql and return records as []map[string][]byte
func (d *DB) Get() ([]JSONDoc, error) {
	session := d.NewSession()
	return session.Get()
}

func (d *DB) GetByID( id string) (JSONDoc, error) {
	session := d.NewSession()
	return session.GetByID(id)

}

func (d *DB) Insert(data JSONDoc) (JSONDoc, error) {
	session := d.NewSession()
	return session.Insert(data)
}

func (d *DB) Update( data JSONDoc) error {
	session := d.NewSession()
	return session.Update(data)
}

func (d *DB) DeleteByID( id string) error {
	session := d.NewSession()
	return session.DeleteByID(id)
}

func (d *DB) DeleteByWhere( ) error {
	session := d.NewSession()
	return session.DeleteByWhere()
}

func (d *DB) Count( ) (int, error) {
	session := d.NewSession()
	return session.Count()
}

func (d *DB) Close() error {
	return d.dialectDB.CloseDB()
}


