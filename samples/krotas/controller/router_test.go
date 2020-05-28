package controller

import (
	"encoding/json"
	"fmt"
	"testing"

	"krotas/bizs"

	"github.com/joyous-x/saturn/common/reqresp"
	"github.com/joyous-x/saturn/common/xnet"
)

const (
	localhost = "127.0.0.1:8000"
)

func unmarshalResp(t *testing.T, respData []byte, resp interface{}) error {
	if err := json.Unmarshal(respData, resp); err != nil {
		t.Errorf("unmarshal error: %v, respData=%v", err, string(respData))
	}

	iResp, ok := resp.(reqresp.IResponse)
	if !ok {
		t.Errorf("error, invalid resp, not IResponse: %#v", resp)
	}

	if iResp.GetCommon().Ret != 0 {
		t.Errorf("error, resp: %#v", resp)
	}

	t.Logf("ok, resp: %#v", resp)
	return nil
}

func Test_Ip2Region(t *testing.T) {
	client := xnet.NewEasyHTTP()

	req := &bizs.Ip2RegionReq{
		ClientIP: "10.20.13.11",
		Debug:    true,
	}
	resp := &bizs.Ip2RegionResp{}
	respData, err := client.PostJSON(fmt.Sprintf("http://%s/%s", localhost, "c/ip2region"), req)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if err := unmarshalResp(t, respData, resp); err != nil {
		t.Errorf("error, unmarshalResp: %#v", err)
	}
}

func Test_Login(t *testing.T) {

}
