package config

var Configor = new(Config)

type Config struct {
	Server struct {
		AppName string `yaml:"app_name"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		IsDebug bool   `yaml:"is_debug"`
	} `yaml:"server"`

	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`

	Database struct {
		Host                  string `yaml:"host"`
		Port                  int    `yaml:"port"`
		User                  string `yaml:"user"`
		Password              string `yaml:"password"`
		Name                  string `yaml:"name"`
		MaxOpenConn           int    `yaml:"max_open_conn"`
		MaxIdleConn           int    `yaml:"max_idle_conn"`
		ConnMaxLifetimeSecond int    `yaml:"conn_max_lifetime_second"`
	} `yaml:"database"`
}
