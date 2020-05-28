package main

import (
	"enceladus/common/xlog"
	"enceladus/governance/wconsul"
	"enceladus/project/mesh_demo/mserver"
	"enceladus/wgrpc"
	"flag"
	"fmt"
	"time"
)

var (
	server        = "127.0.0.1"
	srvHost       = "10.0.2.15"
	consulHost    = "10.0.2.15"
	consulAddress = consulHost + ":8500"
)

func getHealthSvc() error {
	go func() {
		c := time.Tick(20 * time.Second)
		for x := range c {
			wconsul.OneHealthSvcRandom("mesh_demo_entry", "")
			xlog.Error("---- getHealthSvc %v", x)
		}
	}()
	return nil
}

/*
server---启动server群
	./run.sh run
	./run.sh kill
client---通过grpc client访问server
	./main -t c -p 8026 -svc entry -cmd 100100001 -data 7 //> insert, DoDispatch by middle:DoGrpcMethod
	./main -t c -p 8026 -svc entry -cmd 100100003 -data 7 //> sum, DoDispatch
	./main -t c -p 8026 -svc entry -cmd 100100005         //> max, DoGrpcMethodEx
	./main -t c -p 8001 -svc model -cmd 100300001 -data 7 //> insert
	./main -t c -p 8001 -svc model -cmd 100300007         //> max
client---通过gate发送指令访问server---TODO
   ./main -t g -cmd 100100001 -data 7         	//> insert
   ./main -t g -cmd 100100007          			//> max
*/

/*
TODO:
	godoc
	grpc theory : client/server ---- ing
	conn-pool
	conf-center
	websocket
	client-gate && map: srvname - cmd - svcGrpcName
*/
func main() {
	t := flag.String("t", "c", "c or entry or middle or model : client or servers(entry, middle, model)")
	p := flag.Int("p", 8086, "port for client")
	svc := flag.String("svc", "entry", "entry or middle or model")
	cmd := flag.Int("cmd", 0, "cmd in services")
	dat := flag.Int("data", 0, "numbers input")
	flag.Parse()
	xlog.Debug("------ mode: %v", *t)

	if *t != "c" {
		wgrpc.InitModules(consulAddress, "")
		// getHealthSvc()
	}

	port := fmt.Sprintf("%v:%v", srvHost, *p)
	if *t == "c" {
		clientFuncInsecure(*svc, port, *cmd, *dat)
	} else if *t == "entry" {
		mserver.FirstEntry(port)
	} else if *t == "middle" {
		mserver.MiddleEntry(port)
	} else if *t == "model" {
		mserver.ModelEntry(port)
	} else {
		flag.PrintDefaults()
	}
	return
}
