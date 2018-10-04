package gorgo

import (
	"bytes"
	"database/sql"
	"log"
	"sync"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql" // Package for mysql driver
	"github.com/spf13/cast"
)

//MySQLDialect - dialect for mysql database
type MySQLDialect struct {
	DB          *sql.DB
	ShowSQL     bool
	NamedQuerys map[string]string
	CachedMutex sync.Mutex
}

//InitDB  - initialize database
func (m *MySQLDialect) InitDB(config ConfigDB) error {
	user := config.User
	server := config.Server
	pass := config.Password
	port := config.Port
	database := config.Database
	maxIdle := config.MaxIdle
	maxOpen := config.MaxOpen
	sport := cast.ToString(port)
	// "root:@tcp(localhost:3306)/certra"
	var url = user + ":" + pass + "@tcp(" + server + ":" + sport + ")/" + database + "?parseTime=true&timeout=30s"
	//var url = database + "/" + user + "/" + pass
	db, err := sql.Open("mysql", url)
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(int(maxIdle))
	db.SetMaxOpenConns(int(maxOpen))
	m.DB = db
	m.ShowSQL = config.ShowSQL
	m.NamedQuerys = make(map[string]string)

	return nil
}

//CloseDB  - close database
func (m *MySQLDialect) CloseDB() error {
	return m.DB.Close()

}

func (m *MySQLDialect) Create(tableName string, data JSONDoc) (int64, error) {

	m.CachedMutex.Lock()

	keys := sortMap(data)

	var sql = m.NamedQuerys["insert-"+tableName]

	if len(sql) == 0 {
		query := sq.Insert(tableName).Columns(keys...)
		for _, k := range keys {
			v := data[k]
			if k != "id" {
				query.Values(v)
			}
		}
		var err error
		sql, _, err = query.ToSql()
		if err != nil {
			return 0, err
		}
		m.NamedQuerys["insert-"+tableName] = sql
	}
	m.CachedMutex.Unlock()

	if m.ShowSQL == true {
		log.Println("SQL=", sql)
	}

	stmt, err := m.DB.Prepare(sql)
	if err != nil {
		return 0, err
	}

	var params []interface{}

	for _, k := range keys {
		v := data[k]
		if k != "id" {
			params = append(params, v)
		}
	}

	res, err := stmt.Exec(params...)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	stmt.Close()
	return lastId, nil
}

func (s *MySQLDialect) CreateInterface(collection string, i interface{}) error {
	return nil
}

func (m *MySQLDialect) Update(tableName string, data JSONDoc) (int64, error) {

	m.CachedMutex.Lock()
	defer m.CachedMutex.Unlock()

	keys := sortMap(data)

	var sql = m.NamedQuerys["update-"+tableName]
	if len(sql) == 0 {
		builder := sq.Update(tableName)
		for _, k := range keys {
			if k != "id" {
				builder.Set(k, data[k])
			}
		}
		builder.Where("id = ?", 1)
		var err error
		sql, _, err = builder.ToSql()
		if err != nil {
			return 0, err
		}
		m.NamedQuerys["update-"+tableName] = sql
	}

	if m.ShowSQL == true {
		log.Println("SQL=", sql)
	}

	stmt, err := m.DB.Prepare(sql)
	if err != nil {
		return 0, err
	}

	var params []interface{}

	var id interface{}

	for _, k := range keys {
		v := data[k]
		if k != "id" {
			params = append(params, v)
		} else {
			id = v
		}
	}

	params = append(params, id)

	res, err := stmt.Exec(params...)
	if err != nil {
		return 0, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	stmt.Close()
	return rowCnt, nil
}

func (m *MySQLDialect) Delete(tableName string, id int64) (int64, error) {
	m.CachedMutex.Lock()
	defer m.CachedMutex.Unlock()

	var sql = m.NamedQuerys["deletebyid-"+tableName]
	if len(sql) == 0 {
		var err error
		sql, _, err = sq.Delete(tableName).Where("id = ?", id).ToSql()
		if err != nil {
			return 0, err
		}
		m.NamedQuerys["deletebyid-"+tableName] = sql
	}

	if m.ShowSQL == true {
		log.Println("SQL=", sql)
	}
	stmt, err := m.DB.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	// log.Println("delete id=", id)
	res, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}
	i, _ := res.RowsAffected()
	return i, nil
}

func (m *MySQLDialect) Count(tableName string) (int64, error) {
	m.CachedMutex.Lock()
	defer m.CachedMutex.Unlock()

	var sqlbuffer bytes.Buffer
	//hash := calcHash("count-" + tableName + orm.WhereStr)
	sqlbuffer.WriteString(m.NamedQuerys["count-"+tableName])

	if sqlbuffer.Len() == 0 {

		sqlbuffer.WriteString("SELECT COUNT(*) FROM ")
		sqlbuffer.WriteString(tableName)

	}
	m.NamedQuerys["count-"+tableName] = sqlbuffer.String()

	if m.ShowSQL == true {
		log.Println("SQL=", sqlbuffer.String())
	}

	stmt, err := m.DB.Prepare(sqlbuffer.String())
	if err != nil {
		return 0, err
	}
	row := stmt.QueryRow()
	var i int64
	row.Scan(&i)
	//log.Println("interface=", ii)
	stmt.Close()
	return i, nil
}
