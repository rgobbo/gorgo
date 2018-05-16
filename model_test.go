package gorgo

import "testing"

func TestModel_LoadFile(t *testing.T) {
	m := new(model)
	err := m.LoadFile("./model.json")
	if err != nil {
		t.Fatal("Error loading model :", err)
	}

	t.Log("Success loading model.")
	t.Log("Schema=",m.Schema)
	t.Log("-> Tables <-")
	for _,table := range m.Tables {
		t.Log("-->",table.Name)
		t.Log("--> FIELDS <--",)
		for _, ff := range table.Fields {
			t.Log("----->",ff.Name , "  validation:", ff.Validation, "  unique:", ff.Unique, "  alias:", ff.Alias)
		}
	}
}