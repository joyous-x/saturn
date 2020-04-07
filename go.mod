module github.com/joyous-x/saturn

go 1.12

replace github.com/joyous-x/saturn => ../saturn

require (
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.5.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/hashicorp/consul/api v1.3.0
	github.com/jhump/protoreflect v1.6.0
	github.com/jinzhu/gorm v1.9.11
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/opentracing/opentracing-go v1.1.0
	github.com/uber/jaeger-client-go v2.21.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.uber.org/zap v1.10.0
	google.golang.org/grpc v1.26.0
)
