package gorgo

import (
	"fmt"
)

type ORM struct {
	dialectDB Dialect
	showSQL   bool
}

type FuncMap map[string]interface{}

func NewOrm(config ConfigDB) (*ORM, error) {
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
		return &ORM{dialectDB: dialect}, nil
	case "localdb":
		dialect := &LocalDialect{}
		err := dialect.InitDB(config)
		if err != nil {
			return nil, err
		}
		return &ORM{dialectDB: dialect}, nil
	default:
		return nil, fmt.Errorf("[WARNING] dbtype not found", nil)
	}

	return nil, nil
}

func (d *ORM) NewSession() *Session {
	session := &Session{orm: d}
	session.Init()
	return session
}

func (d *ORM) Table(tbl string) *Session {
	session := d.NewSession()
	session.tableName = tbl
	return session
}

func (d *ORM) Limit(i int) *Session {
	session := d.NewSession()
	session.limit = i
	return session
}

func (d *ORM) Offset(i int) *Session {
	session := d.NewSession()
	session.offset = i
	return session
}

func (d *ORM) Where(query string, params ...interface{}) *Session {
	session := d.NewSession()
	session.where = query
	session.params = params
	return session
}

// Get Query a raw sql and return records as []map[string][]byte
func (d *ORM) Get() ([]JSONDoc, error) {
	session := d.NewSession()
	return session.Get()
}

func (d *ORM) GetByID(id string) (JSONDoc, error) {
	session := d.NewSession()
	return session.GetByID(id)

}

func (d *ORM) Insert(data JSONDoc) (JSONDoc, error) {
	session := d.NewSession()
	return session.Insert(data)
}

func (d *ORM) InsertStruct(i interface{}) error {
	session := d.NewSession()
	return session.InsertStruct(i)
}

func (d *ORM) Update(data JSONDoc) error {
	session := d.NewSession()
	return session.Update(data)
}

func (d *ORM) DeleteByID(id string) error {
	session := d.NewSession()
	return session.DeleteByID(id)
}

func (d *ORM) DeleteByWhere() error {
	session := d.NewSession()
	return session.DeleteByWhere()
}

func (d *ORM) Count() (int, error) {
	session := d.NewSession()
	return session.Count()
}

func (d *ORM) Close() error {
	return d.dialectDB.CloseDB()
}
