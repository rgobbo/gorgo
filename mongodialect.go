package gorgo

import (
	"crypto/tls"
	"net"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/rgobbo/watchfy"
	"log"
)

//MySQLDialect - dialect for mysql database
type MongoDialect struct {
	Session *mgo.Session
	DBName  string
	Model  *model
}

func (m *MongoDialect) InitDB(config ConfigDB) error {

	if config.ModelFile != "" {
		err := m.Model.LoadFile(config.ModelFile)
		if err != nil {
			return err
		}
		if config.WatchModel == true {
			go watchfy.NewWatcher([]string{config.ModelFile}, true, func(filename string){
				m.Model = new(model)
				err := m.Model.LoadFile(config.ModelFile)
				if err != nil {
					log.Println(err)
				}
			})
		}
	}

	servers := config.Servers
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

func (m *MongoDialect) Create(collection string, json JSONDoc) error {
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	//json["ID"] = bson.NewObjectId()
	return c.Insert(json)
}

func (m *MongoDialect) GetById(collection string, id string) (JSONDoc, error) {
	var data JSONDoc
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&data)
	return data, err
}

func (m *MongoDialect) GetOneByQuery(collection string, query map[string]interface{}) (JSONDoc, error) {
	var data JSONDoc
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	err := c.Find(query).One(&data)
	return data, err
}

func (m *MongoDialect) GetManyByQuery(collection string, query map[string]interface{}) ([]JSONDoc, error) {
	var data []JSONDoc
	ss := m.Session.Copy()
	defer ss.Close()
	c := ss.DB(m.DBName).C(collection)
	err := c.Find(query).All(&data)
	return data, err
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
