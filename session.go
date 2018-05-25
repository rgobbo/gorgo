package gorgo

import "fmt"

type Session struct {
	tableName string
	limit int
	offset int
	where string
	order string
	params []interface{}
	colums string
	pk  string
	join     string
	groupBy  string
	db *DB
}

func (s *Session) Init() {
	s.limit = 0
	s.offset = 0
	s.order = ""
}

func (s *Session) Where(query string, params ...interface{}) *Session {
	s.where = query
	s.params = params
	return s
}


func (s *Session) Limit(i int) *Session  {
	s.limit = i
	return s
}

func (s *Session) Offset(i int) *Session  {
	s.offset = i
	return s
}

func (s *Session) Get()([]JSONDoc, error) {
	if s.tableName == "" {
		return []JSONDoc{}, fmt.Errorf("need to set a tablename")
	}
	if s.where != "" {
		s.db.dialectDB.GetManyByQuery(s.tableName,s.where, s.params...)
	}
	return s.db.dialectDB.GetAll(s.tableName, s.offset, s.limit, s.order)
}

func (s *Session) GetByID(id string) (JSONDoc, error) {
	if s.tableName == "" {
		return JSONDoc{}, fmt.Errorf("need to set a tablename")
	}
	return s.db.dialectDB.GetById(s.tableName, id)
}

func (s *Session) Insert( data JSONDoc) (JSONDoc, error) {
	if s.tableName == "" {
		return JSONDoc{}, fmt.Errorf("need to set a tablename")
	}
	return s.db.dialectDB.Create(s.tableName, data)
}

func (s *Session) Update( data JSONDoc) error {
	if s.tableName == "" {
		return fmt.Errorf("need to set a tablename")
	}
	return s.db.dialectDB.Update(s.tableName, data)
}

func (s *Session) DeleteByID(id string) error {
	return s.db.dialectDB.Delete(s.tableName, id)

}

func (s *Session) DeleteByWhere() error {
	if s.where == "" {
		return fmt.Errorf("need where clause")
	}
	return s.db.dialectDB.Delete(s.tableName, s.where)

}


func (s *Session) Count() (int, error) {
	if s.tableName == "" {
		return 0,fmt.Errorf("need to set a tablename")
	}
	i, err := s.db.dialectDB.Count(s.tableName)
	if err != nil {
		return 0, err
	}
	return i, nil
}
