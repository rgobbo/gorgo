package gorgo

type ConfigDB struct {
	Type string
	Server string
	Servers []string
	User string
	Password string
	Port int
	Database string
	MaxIdle int
	MaxOpen int
	UseSSL bool
	ShowSQL bool
	ModelFile string
	SyncDB string
	WatchInterval int
	Validations FuncMap
}

func (c *ConfigDB) AddValidations (fn FuncMap) {
	for s, f := range fn {
		c.Validations[s] = f
	}
}
