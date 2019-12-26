package base

type ServiceConfig struct {
	Protocal string `yaml:"protocal"`
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
	Debug    bool   `yaml:"debug"`
}

type ClientConfig struct {
	Protocal string `yaml:"protocal"`
	Address  string `yaml:"address"`
	Name     string `yaml:"name"`
	Debug    bool   `yaml:"debug"`
}

type IMServer interface {
	Route(method, relativePath string, routes ...interface{}) error
	Run() error
	Stop() error
}

type IMClient interface {
	Call(method, relativePath string, req, resp interface{}) error
}
