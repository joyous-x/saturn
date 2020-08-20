package dbs

// RedisConfExts ...
type RedisConfExts struct {
	MaxIdle      int    `yaml:"maxIdle"`
	MaxActive    int    `yaml:"maxActive"`
	IdleTimeout  string `yaml:"idleTimeout"`
	ConnTimeout  string `yaml:"connTimeout"`
	MaxConnLife  string `yaml:"maxConnLifetime"`
	TestOnBorrow string `yaml:"testOnBorrow"`
}

// RedisConf ...
type RedisConf struct {
	Name   string        `yaml:"name"`
	Host   string        `yaml:"host"`
	Passwd string        `yaml:"passwd"`
	Db     int           `yaml:"db"`
	Exts   RedisConfExts `yaml:"exts"`
}

// MysqlConfExts ...
type MysqlConfExts struct {
	MaxIdle     int    `yaml:"maxIdle"`
	MaxOpen     int    `yaml:"maxOpen"`
	MaxPoolSize int    `yaml:"maxPoolSize"`
	MaxConnLife string `yaml:"maxConnLifeTime"`
}

// MysqlConf ...
type MysqlConf struct {
	Name string        `yaml:"name"`
	Type string        `yaml:"type"`
	Dsn  string        `yaml:"dsn"`
	Exts MysqlConfExts `yaml:"exts"`
}
