package base

type ServiceConfig struct {
	Protocal string `yaml:"protocal"` // Protocal, transport protocal
	Encoding string `yaml:"encoding"` // Encoding, the serilization of datas
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
	Debug    bool   `yaml:"debug"`
}

type ClientConfig struct {
	Protocal string `yaml:"protocal"`
	Encoding string `yaml:"encoding"`
	Scheme   string `yaml:"scheme"`
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
