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
}

