package config

var Configor = new(Config)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Log    LogConfig    `yaml:"log"`
	Mysql  MysqlConfig  `yaml:"mysql"`
}

type ServerConfig struct {
	AppName string `yaml:"app_name"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	IsDebug bool   `yaml:"is_debug"`
}

type LogConfig struct {
	FileName string `yaml:"file_name"`
	Level    string `yaml:"level"`
	MaxAge   int    `yaml:"max_age"`
}

type MysqlConfig struct {
	Host                  string `yaml:"host"`
	Port                  int    `yaml:"port"`
	User                  string `yaml:"user"`
	Password              string `yaml:"password"`
	Name                  string `yaml:"name"`
	MaxOpenConn           int    `yaml:"max_open_conn"`
	MaxIdleConn           int    `yaml:"max_idle_conn"`
	ConnMaxLifetimeSecond int    `yaml:"conn_max_lifetime_second"`
}
