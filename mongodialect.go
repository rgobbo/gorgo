package gorgo

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/spf13/cast"

	"log"

	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/rgobbo/fsmodify"
)

//MySQLDialect - dialect for mysql database
type MongoDialect struct {
	Session *mgo.Session
	DBName  string
	Model   *model
}

func (m *MongoDialect) InitDB(config ConfigDB) error {

	if config.ModelFile != "" {
		m.Model = new(model)
		err := m.Model.LoadFile(config.ModelFile)
		if err != nil {
			return err
		}
		if config.WatchInterval > 0 {
			go fsmodify.NewWatcher(config.ModelFile, "", config.WatchInterval, func(filename string) {
				m.Model = new(model)
				err := m.Model.LoadFile(config.ModelFile)
				if err != nil {
					log.Println(err)
				}
			})
		}
	} else {
		m.Model = new(model)
	}

	var servers []string
	if len(config.Servers) > 0 {
		servers = config.Servers
	} else {
		servers = []string{config.Server}
	}

	dbname := config.Database
	user := config.User
	pass := config.Password
	ssl := config.UseSSL

	var info mgo.DialInfo
	info.Database = dbname
	info.Username = user
	info.Password = pass
	info.Addrs = servers
	info.Timeout = 60 * time.Second

	if ssl == true {
		info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		}
	}
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return err
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	m.Session = session

	m.DBName = dbname
	return nil
}

//CloseDB  - close database
func (m *MongoDialect) CloseDB() error {
	m.Session.Close()
	return nil
}

func (m *MongoDialect) Count(collection string) (int, error) {
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	return c.Count()
}

func (m *MongoDialect) CountByWhere(collection string, query string) (int, error) {
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return 0, err
	}
	return c.Find(q).Count()
}

func (m *MongoDialect) Create(collection string, json JSONDoc) (JSONDoc, error) {
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	return json, c.Insert(json)
}

func (m *MongoDialect) GetById(collection string, id string) (JSONDoc, error) {
	var data JSONDoc
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	if !bson.IsObjectIdHex(id) {
		return data, fmt.Errorf("Mongo ObjectID is invalid")
	}
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&data)
	return data, err
}

func (m *MongoDialect) GetOneByQuery(collection string, query string) (JSONDoc, error) {
	var data JSONDoc
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)

	var qjson map[string]interface{}
	err := json.Unmarshal([]byte(query), &qjson)
	if err != nil {
		return data, fmt.Errorf("Query parsing error:%v", err)
	}

	err = c.Find(qjson).One(&data)
	return data, err
}

func (m *MongoDialect) GetManyByQuery(collection string, query string, params ...interface{}) ([]JSONDoc, error) {
	var data []JSONDoc
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)

	squery := m.parseQuery(query, params...)
	log.Println(squery)

	var qjson map[string]interface{}
	err := json.Unmarshal([]byte(squery), &qjson)
	if err != nil {
		return data, fmt.Errorf("Query parsing error:%v", err)
	}
	err = c.Find(qjson).All(&data)
	return data, err
}

func (m *MongoDialect) parseQuery(query string, params ...interface{}) string {
	for _, p := range params {
		if reflect.TypeOf(p).String() == "string" {
			ps := `"` + p.(string) + `"`
			query = strings.Replace(query, "?", ps, 1)
		} else {
			query = strings.Replace(query, "?", cast.ToString(p), 1)
		}

	}
	return query
}

func (m *MongoDialect) GetAll(collection string, page int, qtd int, sorted string) ([]JSONDoc, error) {
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)

	var err error
	var result []JSONDoc
	if sorted != "" {
		err = c.Find(bson.M{}).Sort(sorted).Skip((page - 1) * qtd).Limit(qtd).All(&result)
	} else {
		err = c.Find(bson.M{}).Skip((page - 1) * qtd).Limit(qtd).All(&result)
	}

	return result, err
}

func (m *MongoDialect) GetAllBySearch(collection string, searchtext string, field string, page int, qtd int, sorted string) ([]JSONDoc, error) {
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)

	var err error
	var result []JSONDoc
	if sorted != "" {
		err = c.Find(bson.M{field: bson.RegEx{searchtext, ""}}).Sort(sorted).Skip((page - 1) * qtd).Limit(qtd).All(&result)
	} else {
		err = c.Find(bson.M{field: bson.RegEx{searchtext, ""}}).Skip((page - 1) * qtd).Limit(qtd).All(&result)
	}

	return result, err
}

func (m *MongoDialect) Update(collection string, json JSONDoc) error {
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	return c.Update(bson.M{"_id": json["_id"]}, json)
}

func (m *MongoDialect) Delete(collection string, id string) error {
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	return c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}

func (m *MongoDialect) DeleteByWhere(collection string, query string) error {
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return err
	}
	return c.Remove(q)
}

func (m *MongoDialect) GetByGroup(collection string, query map[string]interface{}) (JSONDoc, error) {
	//db.empresa.aggregate( [ { $group: { _id: null, total: { $sum: "$InteresseEmprestimo" } } } ] )
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)

	//	query := []bson.M{{
	//		"$group": bson.M{
	//			"_id":   bson.M{},
	//			"total": bson.M{"$sum": "$interesseemprestimo"},
	//			"count": bson.M{"$sum": 1},
	//		}},
	//	}
	var result JSONDoc
	err := c.Pipe(query).One(&result)
	return result, err
}
