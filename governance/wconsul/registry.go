package wconsul

import (
	"github.com/joyous-x/saturn/common/xlog"
	"fmt"

	capi "github.com/hashicorp/consul/api"
)

var (
	errInvalidConsulClient = fmt.Errorf("invalid consul client")
)

// SvcInfo ...
type SvcInfo struct {
	ID      string
	Name    string
	Host    string
	Port    int
	Tags    []string
	Checker *capi.AgentServiceCheck
}

// SvcRegistration ...
func SvcRegistration(svc *SvcInfo) error {
	return innerSvcRegistration(DefaultClient(), svc)
}

// SvcDeregistration ...
func SvcDeregistration(serviceID string) error {
	if nil == DefaultClient() {
		return errInvalidConsulClient
	}
	agent := DefaultClient().Agent()
	err := agent.ServiceDeregister(serviceID)
	if err != nil {
		xlog.Error("SvcDeregistration agent.ServiceDeregister serviceID=%v error: %v", serviceID, err)
		return err
	}
	xlog.Debug("SvcDeregistration success serviceID=%v", serviceID)
	return err
}

func innerSvcRegistration(client *capi.Client, svc *SvcInfo) error {
	if nil == client {
		return errInvalidConsulClient
	}
	agent := client.Agent()
	svcRegistration := &capi.AgentServiceRegistration{
		Kind:              capi.ServiceKindTypical,
		ID:                svc.ID,
		Name:              svc.Name,
		Tags:              svc.Tags,
		Port:              svc.Port,
		Address:           svc.Host, //> 默认agent所连接的server地址
		EnableTagOverride: false,
		Meta:              map[string]string{},
		Weights:           &capi.AgentWeights{Passing: 1, Warning: 1},
		Check:             svc.Checker, //> capi.AgentServiceCheck {}
		Checks:            []*capi.AgentServiceCheck{},
		Proxy:             nil, //> capi.AgentServiceConnectProxyConfig{}
		Connect:           nil, //> capi.AgentServiceConnect{}
	}
	err := agent.ServiceRegister(svcRegistration)
	if err != nil {
		xlog.Error("SvcRegistration agent.ServiceRegister error: %v", err)
		return err
	}
	xlog.Debug("SvcRegistration success id=%v name=%v host=%v:%v", svc.ID, svc.Name, svc.Host, svc.Port)
	return err
}
