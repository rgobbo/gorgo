package gorgo

import (
	"github.com/tidwall/buntdb"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"time"
	"encoding/json"
	"github.com/spf13/cast"
	"fmt"
)

type LocalDialect struct {
	DB *buntdb.DB
}

func (s *LocalDialect) InitDB(config ConfigDB) error {
	server := config.Server
	db, err := buntdb.Open(server)
	if err != nil {
		return err
	}
	s.DB = db
	err = s.generateIndexes()
	return err
}

func (s *LocalDialect) generateIndexes() (error) {

	err := s.DB.Update(func(tx *buntdb.Tx) error {

		indexes, err := tx.Indexes()
		if err != nil {
			return err
		}
		if !contains(indexes, "buckets") {
			err = tx.CreateIndex("buckets",  "BUCKETS:*", buntdb.IndexString)
			if err != nil {
				return err
			}
		}

		var buckets []string
		tx.Ascend("buckets", func(key, val string) bool {
			buckets = append(buckets,val)
			return true
		})

		for _, s := range buckets {
			if ! contains(indexes, "idx" + s) {
				err = tx.CreateIndex("idx" + s,  s + ":*", buntdb.IndexString)
				if err != nil {
					return err
				}
			}
		}


		return err
	})


	return err
}

func contains(slice []string, elem string) bool {
	for _, a := range slice {
		if a == elem {
			return true
		}
	}
	return false
}

//CloseDB  - close database
func (s *LocalDialect) CloseDB() error {
	return s.DB.Close()
}

func (s *LocalDialect) Count(collection string) (int, error) {
	var ret int
	err := s.DB.View(func(tx *buntdb.Tx) error  {
		count := 0
		err := tx.Ascend( "idx" + collection , func(key, value string) bool{
			if strings.HasPrefix(key , collection + ":") {
				count++
			}
			return true

		})
		ret = count
		return err
	})
	return ret, err

}

func (s *LocalDialect) Create(collection string, data JSONDoc) error {
	id := bson.NewObjectId()
	sid := id.Hex()
	data["_id"] = sid
	key := collection + ":" + sid

	data["_created"] = time.Now()

	err := s.DB.Update(func(tx *buntdb.Tx) error {

		_,_,err := tx.Set("BUCKETS:" + collection,collection,nil)
		if err != nil {
			return err
		}

		indexes, err := tx.Indexes()
		if err != nil {
			return err
		}

		if ! contains(indexes, "idx" + collection) {
			err = tx.CreateIndex("idx" + collection,  collection + ":*", buntdb.IndexString)
			if err != nil {
				return err
			}
		}
		_ , _ ,err = tx.Set(key,data.ToString(),nil)
		if err != nil {
			return err
		}



		return err
	})


	return err
}

func (s *LocalDialect) GetById(collection string, id string) (JSONDoc, error) {
	var data JSONDoc
	err := s.DB.View(func(tx *buntdb.Tx) error  {
		key := collection + ":" + id
		// Retrieve the record
		item, err := tx.Get(key)
		if err != nil {
			return err
		}
		// Decode the record
		e := json.Unmarshal([]byte(item), &data)
		if e != nil {
			return e
		}

		return nil
	})

	return data, err
}


func (s *LocalDialect) Update(tableName string, data JSONDoc) error {
	sid := ""
	if data["_id"] != nil {
		sid = cast.ToString(data["_id"])
	} else {
		return fmt.Errorf("Field _id could not be null")
	}
	key := tableName + ":" + sid

	err := s.DB.Update(func(tx *buntdb.Tx) error {

		_, _ ,err := tx.Set(key, data.ToString(),nil)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *LocalDialect) Delete(tableName string, id string) error {
	err := s.DB.Update(func(tx *buntdb.Tx) error {

		key := tableName + ":" + id
		//KeyUnique := "unique:" + tableName + ":"

		_, err := tx.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *LocalDialect) GetAll(tableName string, skip int, limit int, sorted string) ([]JSONDoc, error) {
	var result []JSONDoc
	count := 0
	skipLimit := skip * limit
	initLimit := skipLimit - limit
	err := s.DB.View(func(tx *buntdb.Tx) error {

		err := tx.Ascend( "idx" + tableName,  func(key, value string) bool{
			if count >= initLimit && count < skipLimit {

				var single JSONDoc
				err := json.Unmarshal([]byte(value), &single)
				if err != nil {
					return false
				}
				result = append(result, single)
				count++
				return true

			} else if count < initLimit {
				count++
				return true
			} else {
				return false
			}

		})
		return err

	})
	return result, err
}
func (s *LocalDialect) GetOneByQuery(collection string, query map[string]interface{}) (JSONDoc, error) {
	var data JSONDoc
	return data, nil
}
func (s *LocalDialect) GetManyByQuery(collection string, query map[string]interface{}) ([]JSONDoc, error) {
	var data []JSONDoc
	return data, nil
}
func (s *LocalDialect) GetAllBySearch(collection string, text string, field string, page int, qtd int, sorted string) ([]JSONDoc, error) {
	var data []JSONDoc
	return data, nil

}
func (s *LocalDialect) GetByGroup(collection string, query map[string]interface{}) (JSONDoc, error) {
	var data JSONDoc
	return data, nil

}
