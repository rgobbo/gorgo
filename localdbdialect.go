package gorgo

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rgobbo/watchfy"
	"github.com/spf13/cast"
	"github.com/tidwall/buntdb"
	"gopkg.in/mgo.v2/bson"
)

type LocalDialect struct {
	DB     *buntdb.DB
	Model  *model
	Config ConfigDB
}

func (s *LocalDialect) InitDB(config ConfigDB) error {
	if config.ModelFile != "" {
		s.Model = new(model)
		err := s.Model.LoadFile(config.ModelFile)
		if err != nil {
			return err
		}
		if config.WatchModel == true {
			go watchfy.NewWatcher([]string{config.ModelFile}, true, func(filename string) {
				s.Model = new(model)
				err := s.Model.LoadFile(config.ModelFile)
				if err != nil {
					log.Println(err)
				}
			})
		}

	} else {
		s.Model = new(model)
	}

	server := config.Server
	db, err := buntdb.Open(server)
	if err != nil {
		return err
	}
	s.DB = db
	err = s.generateIndexes()

	s.Config = config

	return err
}

func (s *LocalDialect) generateIndexes() error {

	err := s.DB.Update(func(tx *buntdb.Tx) error {

		indexes, err := tx.Indexes()
		if err != nil {
			return err
		}
		if !contains(indexes, "buckets") {
			err = tx.CreateIndex("buckets", "BUCKETS:*", buntdb.IndexString)
			if err != nil {
				return err
			}
		}

		var buckets []string
		tx.Ascend("buckets", func(key, val string) bool {
			buckets = append(buckets, val)
			return true
		})

		for _, s := range buckets {
			if !contains(indexes, "idx"+s) {
				err = tx.CreateIndex("idx"+s, s+":*", buntdb.IndexString)
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
	err := s.DB.View(func(tx *buntdb.Tx) error {
		count := 0
		err := tx.Ascend("idx"+collection, func(key, value string) bool {
			if strings.HasPrefix(key, collection+":") {
				count++
			}
			return true

		})
		ret = count
		return err
	})
	return ret, err

}

func (s *LocalDialect) Create(collection string, data JSONDoc) (JSONDoc, error) {
	id := bson.NewObjectId()
	sid := id.Hex()
	data["_id"] = sid
	key := collection + ":" + sid
	uniques := []string{}
	var newDoc JSONDoc

	data["_created"] = time.Now()

	err := validateFields(collection, data, s.Model, s.Config.Validations)
	if err != nil {
		return newDoc, err
	}

	if val, ok := s.Model.Tables[collection]; ok {
		for _, f := range val.Fields {
			if f.Unique == true {
				if data[f.Name] == nil {
					return newDoc, fmt.Errorf("Unique field %s, could not be null", f.Name)
				}
				uniques = append(uniques, "unique_"+collection+":"+cast.ToString(data[f.Name]))
			}
		}
	}

	err = s.DB.Update(func(tx *buntdb.Tx) error {

		_, _, err := tx.Set("BUCKETS:"+collection, collection, nil)
		if err != nil {
			return err
		}

		indexes, err := tx.Indexes()
		if err != nil {
			return err
		}

		if !contains(indexes, "idx"+collection) {
			err = tx.CreateIndex("idx"+collection, collection+":*", buntdb.IndexString)
			if err != nil {
				return err
			}
		}
		if !contains(indexes, "idx_unique"+collection) {
			err = tx.CreateIndex("idx_unique"+collection, "unique_"+collection+":*", buntdb.IndexString)
			if err != nil {
				return err
			}
		}

		for _, s := range uniques {
			v, _ := tx.Get(s, true)
			if v != "" {
				return fmt.Errorf("Unique key violated - key[%s] ", s)
			} else {
				_, _, err = tx.Set(s, sid, nil)
				if err != nil {
					return err
				}
			}
		}
		_, _, err = tx.Set(key, data.ToString(), nil)
		if err != nil {
			return err
		}

		return err
	})

	newDoc = data
	return newDoc, err
}

func (s *LocalDialect) GetById(collection string, id string) (JSONDoc, error) {
	var data JSONDoc
	err := s.DB.View(func(tx *buntdb.Tx) error {
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

func (s *LocalDialect) Update(collection string, data JSONDoc) error {
	sid := ""
	if data["_id"] != nil {
		sid = cast.ToString(data["_id"])
	} else {
		return fmt.Errorf("Field _id could not be null")
	}
	key := collection + ":" + sid

	err := validateFields(collection, data, s.Model, s.Config.Validations)
	if err != nil {
		return err
	}

	err = s.DB.Update(func(tx *buntdb.Tx) error {
		item, err := tx.Get(key)
		if err != nil {
			return err
		}
		if item == "" {
			return fmt.Errorf("Item id[%s] not found", sid)
		}
		var olddata JSONDoc
		e := json.Unmarshal([]byte(item), &olddata)
		if e != nil {
			return e
		}

		if val, ok := s.Model.Tables[collection]; ok {
			for _, f := range val.Fields {
				if f.Unique == true {
					if data[f.Name] == nil {
						return fmt.Errorf("Unique field %s, could not be null", f.Name)
					}
					s := "unique_" + collection + ":" + cast.ToString(data[f.Name])
					v, _ := tx.Get(s, true)

					if v != "" && v != cast.ToString(data["_id"]) {
						return fmt.Errorf("Unique key violated - %s ", v)
					} else {
						oldunique := "unique_" + collection + ":" + cast.ToString(olddata[f.Name])
						_, _ = tx.Delete(oldunique)

						_, _, err = tx.Set(s, sid, nil)
						if err != nil {
							return err
						}
					}
				}
			}
		}

		_, _, err = tx.Set(key, data.ToString(), nil)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *LocalDialect) Delete(collection string, id string) error {
	err := s.DB.Update(func(tx *buntdb.Tx) error {

		key := collection + ":" + id

		_, err := tx.Delete(key)
		if err != nil {
			return err
		}

		err = tx.Ascend("idx_unique"+collection, func(key, value string) bool {
			if value == id {
				_, err := tx.Delete(key)
				if err != nil {
					return false
				}
			}
			return true
		})

		return err
	})

	return err
}

func (s *LocalDialect) GetAll(tableName string, skip int, limit int, sorted string) ([]JSONDoc, error) {
	var result []JSONDoc
	count := 0
	skipLimit := skip * limit
	initLimit := skipLimit - limit
	err := s.DB.View(func(tx *buntdb.Tx) error {

		err := tx.Ascend("idx"+tableName, func(key, value string) bool {
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
func (s *LocalDialect) GetOneByQuery(collection string, query string) (JSONDoc, error) {
	var data JSONDoc
	err := s.DB.View(func(tx *buntdb.Tx) error {

		err := tx.Ascend("idx"+collection, func(key, value string) bool {
			res := strings.Contains(value, query)
			if res {
				err := json.Unmarshal([]byte(value), &data)
				if err != nil {
					log.Println("Error when unmarshal data")
				}
				return false
			}
			return true

		})
		return err

	})

	return data, err
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
