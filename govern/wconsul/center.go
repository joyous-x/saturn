package wconsul

import (
	capi "github.com/hashicorp/consul/api"
	"github.com/joyous-x/saturn/common/xlog"
	"net/http"
	"sync"
	"time"
)

var defaultConsulClient *capi.Client
var defaultOnce sync.Once

// InitDefaultClient init default consul client
func InitDefaultClient(consulServerAddr string) *capi.Client {
	if len(consulServerAddr) > 0 {
		client, err := NewClient(consulServerAddr)
		if err == nil {
			_, err = client.Status().Leader()
			if err == nil {
				defaultOnce.Do(func() {
					defaultConsulClient = client
				})
			}
		}
		if err != nil {
			xlog.Error("InitDefaultClient addr=%v error=%v", consulServerAddr, err)
		}
	}
	return defaultConsulClient
}

// DefaultClient get default consul client instance
func DefaultClient() *capi.Client {
	return defaultConsulClient
}

// NewClient new consul client with (URI)Scheme='http'
func NewClient(consulServerAddr string) (*capi.Client, error) {
	consulServerURIScheme := "http"
	consulCfg := &capi.Config{
		Address:    consulServerAddr,
		Scheme:     consulServerURIScheme, // The URL scheme of the agent to use ("http" or "https"). Defaults to "http"
		Datacenter: "",                    // If not provided, the default agent datacenter is used.
		Transport:  &http.Transport{},     // use for the http client.
		HttpClient: nil,                   // the client to use. Default will be used if not provided.
		HttpAuth:   nil,                   //the auth info to use for http access.
		WaitTime:   time.Duration(0),      // limits how long a Watch will block. If not provided, the agent default values will be used.
		Token:      "",
		TokenFile:  "",
		TLSConfig:  capi.TLSConfig{},
	}
	client, err := capi.NewClient(consulCfg)
	if err != nil {
		xlog.Error("NewClient error: %v", err)
		return nil, err
	}

	return client, nil
}
