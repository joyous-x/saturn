package dbs

import ()

type RedisExts struct {
	MaxIdle      int    `yaml:"maxIdle"`
	MaxActive    int    `yaml:"maxActive"`
	IdleTimeout  string `yaml:"idleTimeout"`
	ConnTimeout  string `yaml:"connTimeout"`
	MaxConnLife  string `yaml:"maxConnLifetime"`
	TestOnBorrow string `yaml:"testOnBorrow"`
}
type RedisDb struct {
	Key string `yaml:"key"`
	Db  int    `yaml:"db"`
}
type RedisConf struct {
	Host   string    `yaml:"host"`
	Passwd string    `yaml:"passwd"`
	Dbs    []RedisDb `yaml:"dbs"`
	Exts   RedisExts `yaml:"exts"`
}

type MysqlExts struct {
	MaxIdle     int    `yaml:"maxIdle"`
	MaxOpen     int    `yaml:"maxOpen"`
	MaxPoolSize int    `yaml:"maxPoolSize"`
	MaxConnLife string `yaml:"maxConnLifeTime"`
}
type MysqlConf struct {
	Key    string    `yaml:"key"`
	Type   string    `yaml:"type"`
	Host   string    `yaml:"host"`
	User   string    `yaml:"user"`
	Passwd string    `yaml:"passwd"`
	DbName string    `yaml:"dbname"`
	Debug  bool      `yaml:"debug"`
	Exts   MysqlExts `yaml:"exts"`
}
