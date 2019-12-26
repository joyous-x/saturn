package xpprof

import (
	"context"
	"fmt"
	"github.com/joyous-x/saturn/common/xlog"
	"net"
	"net/http"
	"net/http/pprof"
)

// reference
//   https://blog.golang.org/profiling-go-programs
// note:
//   var pprofAddr = flag.String("pprof_addr", ":8181", "address for pprof http service")

// Start run goroutine for pprof
func Start(opts ...int) {
	port := func() int {
		if len(opts) > 0 {
			return opts[0]
		}
		return 0
	}()
	startWebPProf(context.Background(), port)
}

func freeLocalPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func startWebPProf(ctx context.Context, port int) error {
	if port == 0 {
		tmpPort, err := freeLocalPort()
		if err != nil {
			xlog.Error("===> pprof ### findValidPort err: %v \n", err)
			return err
		}
		port = tmpPort
	}
	addr := fmt.Sprintf(":%v", port)

	go func() {
		xlog.Info("===> pprof ### start addr=%v \n", addr)
		httpServer := http.NewServeMux()
		httpServer.HandleFunc("/debug/pprof/", pprof.Index)
		httpServer.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpServer.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpServer.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpServer.HandleFunc("/debug/pprof/trace", pprof.Trace)
		err := http.ListenAndServe(addr, httpServer)
		if err != nil {
			xlog.Error("===> pprof ### error %v \n", err)
		} else {
			xlog.Info("===> pprof ### end \n")
		}
	}()

	return nil
}
