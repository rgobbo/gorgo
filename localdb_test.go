package gorgo

import "testing"

func TestLocalDialect_InitDB(t *testing.T) {
	config := ConfigDB{}
	config.ModelFile = "model.json"
	config.Type = "localdb"
	config.Server = "localtest.db"
	config.WatchModel = true

	DB, err := NewDB(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("DB created : %s", config.Server)
	err = DB.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("DB close")

}