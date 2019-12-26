package wconsul

import (
	"fmt"
	capi "github.com/hashicorp/consul/api"
	"github.com/joyous-x/saturn/common/xlog"
	"math/rand"
	"time"
)

const (
	// LBPolicyRandom lb policy: random
	LBPolicyRandom = "csl_random"
)

// NewHTTPSvcCheck  interval and timeout are in Go time format
//       reference: https://www.consul.io/docs/agent/checks.html
func NewHTTPSvcCheck(name, url string, headers map[string][]string, method, interval, timeout string) *capi.AgentServiceCheck {
	if len(interval) < 1 {
		interval = "10s"
	}
	if len(timeout) < 1 {
		timeout = "3s"
	}
	if len(method) < 1 {
		method = "GET"
	}
	newHTTPSvcCheck := func() *capi.AgentServiceCheck {
		return &capi.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("checkid_%v", time.Now().Unix()),
			Name:                           fmt.Sprintf("checker_%v", name),
			Args:                           []string{},
			DockerContainerID:              "",
			Shell:                          "", // Only supported for Docker.
			Interval:                       interval,
			Timeout:                        timeout,
			TTL:                            "",
			HTTP:                           url,
			Header:                         headers,
			Method:                         method, //> default: GET
			TCP:                            "",
			Status:                         "",
			Notes:                          "",
			TLSSkipVerify:                  true,
			GRPC:                           "",
			GRPCUseTLS:                     false,
			AliasNode:                      "",
			AliasService:                   "",
			DeregisterCriticalServiceAfter: "0h3m",
		}
	}
	return newHTTPSvcCheck()
}

// NewGrpcSvcCheck ...
func NewGrpcSvcCheck(host string, port int, interval, timeout string, useTLS bool) *capi.AgentServiceCheck {
	if len(interval) < 1 {
		interval = "5s"
	}
	if len(timeout) < 1 {
		timeout = "3s"
	}
	newSvcCheck := func() *capi.AgentServiceCheck {
		return &capi.AgentServiceCheck{
			Interval:                       interval,
			Timeout:                        timeout,
			TLSSkipVerify:                  true,
			GRPC:                           fmt.Sprintf("%v:%v", host, port),
			GRPCUseTLS:                     useTLS,
			DeregisterCriticalServiceAfter: "0h10m",
		}
	}
	return newSvcCheck()
}

// OneHealthSvcRandom ...
func OneHealthSvcRandom(svcName, tag string) (*capi.AgentService, error) {
	return getHealthSvcs(svcName, tag, LBPolicyRandom)
}

func getHealthSvcs(svcName, tag, policy string) (*capi.AgentService, error) {
	var agentService *capi.AgentService

	if nil == DefaultClient() {
		return agentService, errInvalidConsulClient
	}

	// health.State : capi.HealthPassing
	health := DefaultClient().Health()
	entrys, _, err := health.Service(svcName, tag, true, nil)
	if err != nil {
		return agentService, err
	}

	if len(entrys) < 1 {
		err = fmt.Errorf("no service:%v available", svcName)
	}

	if err != nil {
		xlog.Error("getHealthSvcs error: %v", err)
		return agentService, err
	}

	switch policy {
	case LBPolicyRandom:
		fallthrough
	default:
		index := rand.New(rand.NewSource(time.Now().Unix())).Intn(len(entrys))
		agentService = entrys[index].Service
		xlog.Debug("getHealthSvcs svcName=%v policy=%v index(len:%v)=%v service=%+v", svcName, policy, len(entrys), index, *agentService)
	}

	return agentService, nil
}
