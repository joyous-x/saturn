package main

import (
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/govern/wconsul"
	"fmt"
	capi "github.com/hashicorp/consul/api"
)

const (
	hostAddr      = "10.0.2.15"
	serverAddress = hostAddr + ":8500"
)

func testConsulAll() error {
	client, err := wconsul.NewClient(serverAddress)
	if err != nil {
		return err
	}

	testAgent(client)
	testCatalog(client)
	testKV(client)
	testHealth(client)
	testSvcRegistry(client)
	return nil
}

func testSvcRegistry(client *capi.Client) error {
	xlog.Info("======> testSvcRegistry")
	serviceA := &wconsul.SvcInfo{"svcHelloWorldA", "svcHelloWorld", hostAddr, 6671, []string{"version_x_1"}, nil}
	serviceB := &wconsul.SvcInfo{"svcHelloWorldB", "svcHelloWorld", hostAddr, 6671, []string{"version_x_2"}, nil}
	serviceC := &wconsul.SvcInfo{"svcSayHi", "svcSayHi", hostAddr, 6671, []string{"version_c_1"}, nil}

	agent := client.Agent()
	nodeName, err := agent.NodeName()
	if err != nil {
		xlog.Error("testSvcRegistry agent.NodeName error: %v", err)
		return err
	} else {
		xlog.Debug("testSvcRegistry agent.NodeName nodename=%v", nodeName)
	}

	serviceCheck := func(svc *wconsul.SvcInfo) *capi.AgentServiceCheck {
		url := fmt.Sprintf("http://%v:%v/health", svc.Host, svc.Port)
		check := wconsul.NewHTTPSvcCheck(svc.Name, url, map[string][]string{"x-foo": []string{"bar", "baz"}}, "POST", "", "")
		return check
	}
	svcRegistration := &capi.AgentServiceRegistration{
		Kind:              capi.ServiceKindTypical,
		ID:                serviceA.ID,
		Name:              serviceA.Name,
		Tags:              serviceA.Tags,
		Port:              serviceA.Port,
		Address:           serviceA.Host, //> 默认agent所连接的server地址
		EnableTagOverride: false,
		Meta:              map[string]string{},
		Weights:           &capi.AgentWeights{Passing: 1, Warning: 1},
		Check:             serviceCheck(serviceA), //> capi.AgentServiceCheck {}
		Checks:            []*capi.AgentServiceCheck{},
		Proxy:             nil, //> capi.AgentServiceConnectProxyConfig{}
		Connect:           nil, //> capi.AgentServiceConnect{}
	}
	err = agent.ServiceRegister(svcRegistration)
	if err != nil {
		xlog.Error("testSvcRegistry agent.ServiceRegister error: %v", err)
		return err
	}
	svcRegistration.ID, svcRegistration.Name, svcRegistration.Tags, svcRegistration.Port, svcRegistration.Check = serviceB.ID, serviceB.Name, serviceB.Tags, serviceB.Port, serviceCheck(serviceB)
	err = agent.ServiceRegister(svcRegistration)
	if err != nil {
		xlog.Error("testSvcRegistry agent.ServiceRegister id=%v name=%v error: %v", svcRegistration.ID, svcRegistration.Name, err)
		return err
	}
	svcRegistration.ID, svcRegistration.Name, svcRegistration.Tags, svcRegistration.Port, svcRegistration.Check = serviceC.ID, serviceC.Name, serviceC.Tags, serviceC.Port, serviceCheck(serviceC)
	err = agent.ServiceRegister(svcRegistration)
	if err != nil {
		xlog.Error("testSvcRegistry agent.ServiceRegister id=%v name=%v error: %v", svcRegistration.ID, svcRegistration.Name, err)
		return err
	}

	svcDeregistrationID := serviceB.ID
	err = agent.ServiceDeregister(svcDeregistrationID)
	if err != nil {
		xlog.Error("testSvcRegistry agent.ServiceDeregister error: %v", err)
	} else {
		xlog.Debug("testSvcRegistry agent.ServiceDeregister svc=%v ok", svcDeregistrationID)
	}

	agentSvc, queryMeta, err := agent.Service(svcDeregistrationID, nil)
	if err != nil {
		xlog.Error("testSvcRegistry agent.Service svc=%v error: %v", svcDeregistrationID, err)
	} else {
		xlog.Debug("testSvcRegistry agent.Service svc=%v info=%v queryMeta=%v ok", svcDeregistrationID, agentSvc, queryMeta)
	}

	svcRegistration.ID, svcRegistration.Name, svcRegistration.Tags, svcRegistration.Port, svcRegistration.Check = serviceB.ID, serviceB.Name, serviceB.Tags, serviceB.Port, serviceCheck(serviceB)
	err = agent.ServiceRegister(svcRegistration)
	if err != nil {
		xlog.Error("testSvcRegistry agent.ServiceRegister id=%v name=%v error: %v", svcRegistration.ID, svcRegistration.Name, err)
	}

	mapSvcs, err := agent.Services()
	if err != nil {
		xlog.Error("testSvcRegistry agent.Services error: %v", err)
	} else {
		for i, v := range mapSvcs {
			xlog.Debug("testSvcRegistry agent.Services index=%v svc=%v ", i, v)
		}
	}

	return err
}

func testHealth(client *capi.Client) error {
	xlog.Info("======> testHealth")

	health := client.Health()
	checks, queryMeta, err := health.State("any", nil)
	if err != nil {
		xlog.Error("testHealth health.State error: %v", err)
		return err
	} else {
		for i, v := range checks {
			xlog.Debug("testHealth health.State index=%v val=%v queryMeta=%v", i, v, queryMeta)
		}
	}

	return nil
}

func testKV(client *capi.Client) error {
	xlog.Info("======> testKV")
	//> reference: https://www.consul.io/api/kv.html
	//>      note: Values in the KV store cannot be larger than 512kb.

	kvTestKey := "kv-key-a"
	kvTestKeyb := "kv-key-b"
	kvTestVala := "kv-val:helloworld-a"
	kvTestValb := "kv-val:helloworld-b"

	kv := client.KV()
	writeOpt := &capi.WriteOptions{
		Datacenter:  "",
		Token:       "",
		RelayFactor: 0,
	}
	putKVPair := &capi.KVPair{
		Key:         kvTestKey,
		Value:       []byte(kvTestVala),
		Session:     "",
		CreateIndex: 0,
		ModifyIndex: 0,
		LockIndex:   0,
		Flags:       888888,
	}
	writeMeta, err := kv.Put(putKVPair, writeOpt) // writeOpt can be nil
	if err != nil {
		xlog.Error("testKV kv.Put error: %v", err)
		return err
	} else {
		xlog.Debug("testAgent kv.Put rst: %v", *writeMeta)
	}
	putKVPair.Key, putKVPair.Value = kvTestKeyb, []byte(kvTestValb)
	writeMeta, err = kv.Put(putKVPair, nil)
	if err != nil {
		xlog.Error("testKV kv.Put error: %v", err)
		return err
	} else {
		xlog.Debug("testAgent kv.Put rst: %v", *writeMeta)
	}

	queryOpt := &capi.QueryOptions{
		RequireConsistent: true,
		UseCache:          false,
	}

	prefix := "kv"
	separator := ""
	keys, queryMeta, err := kv.Keys(prefix, separator, queryOpt)
	if err != nil {
		xlog.Error("testAgent kv.Keys error: %v", err)
		return err
	} else {
		xlog.Debug("testAgent kv.Keys val=%v queryMeta=%v", keys, queryMeta)
	}

	pairs, queryMeta, err := kv.List(prefix, queryOpt)
	if err != nil {
		xlog.Error("testAgent kv.Keys error: %v", err)
		return err
	} else {
		for i, v := range pairs {
			xlog.Debug("testAgent kv.Keys index=%v val=%v queryMeta=%v", i, v, queryMeta)
		}
	}

	val, queryMeta, err := kv.Get(kvTestKey, queryOpt)
	if err != nil {
		xlog.Error("testAgent kv.Get error: %v", err)
		return err
	} else {
		xlog.Debug("testAgent kv.Get val=%v queryMeta=%v", val, queryMeta)
	}

	// cas: key val 都不变
	putKVPair.ModifyIndex = val.ModifyIndex
	putKVPair.Key, putKVPair.Value = kvTestKey, []byte(kvTestVala)
	rst, writeMeta, err := kv.CAS(putKVPair, writeOpt)
	if err != nil {
		xlog.Error("testAgent kv.CAS a error: %v", err)
	} else {
		xlog.Debug("testAgent kv.CAS a rst=%v key=%v modifyIndex=%v writeMeta=%v", rst, putKVPair.Key, val.ModifyIndex, writeMeta)
	}
	// cas: val 改变
	val, _, _ = kv.Get(kvTestKey, queryOpt)
	putKVPair.ModifyIndex, putKVPair.Value = val.ModifyIndex, []byte(kvTestValb+"_cas")
	rst, writeMeta, err = kv.CAS(putKVPair, writeOpt)
	if err != nil {
		xlog.Error("testAgent kv.CAS b error: %v", err)
	} else {
		xlog.Debug("testAgent kv.CAS b rst=%v key=%v modifyIndex=%v writeMeta=%v", rst, putKVPair.Key, val.ModifyIndex, writeMeta)
	}

	writeMeta, err = kv.DeleteTree(prefix, nil)
	if err != nil {
		xlog.Error("testAgent kv.DeleteTree error: %v", err)
	} else {
		xlog.Debug("testAgent kv.DeleteTree prefix=%v writeMeta=%v", prefix, writeMeta)
	}

	return nil
}

func testAgent(client *capi.Client) error {
	xlog.Info("======> testAgent")

	agent := client.Agent()
	infos, err := agent.Host()
	if err != nil {
		xlog.Error("testAgent agent.Host() error: %v", err)
		return err
	} else {
		xlog.Debug("testAgent agent.Host() infos: %v", infos)
	}

	members, err := agent.Members(true)
	if err != nil {
		xlog.Error("testAgent agent.Members() error: %v", err)
		return err
	} else {
		for i, v := range members {
			xlog.Debug("testAgent agent.Members() index=%v data=%v", i, *v)
		}
	}

	return nil
}

func testCatalog(client *capi.Client) error {
	xlog.Info("======> testCatalog")
	catalog := client.Catalog()
	nodes, queryMeta, err := catalog.Nodes(nil)
	if err != nil {
		xlog.Error("testCatalog catalog.Nodes error: %v", err)
		return err
	} else {
		for i, v := range nodes {
			xlog.Debug("testCatalog catalog.Nodes index=%v node=%v queryMeta=%v", i, *v, queryMeta)
		}
	}
	return nil
}
