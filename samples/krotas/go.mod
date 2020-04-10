module krotas

go 1.12

replace github.com/joyous-x/saturn => ../../../saturn

require (
	github.com/gin-contrib/static v0.0.0-20191128031702-f81c604d8ac2
	github.com/gin-gonic/gin v1.5.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jinzhu/gorm v1.9.11
	github.com/joyous-x/saturn v0.0.0-00010101000000-000000000000
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.2
)
