package gorgo

import (
	"testing"
	"time"
	"github.com/spf13/cast"
)


func TestLocalDialect_Create(t *testing.T) {
	config := ConfigDB{}
	config.ModelFile = "model.json"
	config.Type = "localdb"
	config.Server = "localtest.db"
	config.WatchModel = true

	DB, err := NewDB(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Open DB success !!")
	defer DB.Close()


	// Creating user
	user := JSONDoc{}
	user["name"] = "john Doe"
	user["email"] = "jd@test.com"
	user["cpf"] = "23749817030"
	user["cnpj"] = "78470985000106"
	user["age"] = 25
	user["teste"] = "asdfgrreqw653ter"
	user["created"] = time.Now()
	user["updated"] = time.Now()

	userRet, err := DB.Table("user").Insert(user)
	if err != nil {
		t.Fatal("DB Create Error : ", err)
	}
	t.Logf("User created id[ %v ] ", userRet["_id"])


	list, err := DB.Table("user").Limit(6).Offset(1).Get()
	if err != nil {
		t.Fatal("DB GetAll Error : ", err)
	}
	t.Log("User list:", len(list))

	userUpdate := list[0]
	userUpdate["name"] = "upd john doe"

	err = DB.Table("user").Update( userUpdate)
	if err != nil {
		t.Fatal("DB Update Error : ", err)
	}
	t.Logf("User updated !!")

	err = DB.Table("user").DeleteByID( cast.ToString(userUpdate["_id"]))
	if err != nil {
		t.Fatal("DB Delete Error : ", err)
	}
	t.Logf("User deleted !!")

}
