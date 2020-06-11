package controller

import (
	"encoding/json"
	"fmt"
	"testing"

	"krotas/bizs"

	"github.com/joyous-x/saturn/common/reqresp"
	"github.com/joyous-x/saturn/common/xnet"
	"github.com/joyous-x/saturn/foos/user"
)

var (
	localhost = "127.0.0.1:8000"
	client    = xnet.NewEasyHTTP()
)

func unmarshalResp(t *testing.T, respData []byte, resp interface{}) error {
	if err := json.Unmarshal(respData, resp); err != nil {
		t.Fatalf("unmarshal error: %v, respData=%v", err, string(respData))
	}

	iResp, ok := resp.(reqresp.IResponse)
	if !ok {
		t.Fatalf("error, invalid resp, not IResponse: %#v", resp)
	}

	if iResp.GetCommon().Ret != 0 {
		t.Fatalf("error, resp: %#v", resp)
	}

	t.Logf("ok, resp: %#v", resp)
	return nil
}

func Test_Ip2Region(t *testing.T) {
	req := &bizs.Ip2RegionReq{
		ClientIP: "10.20.13.11",
		Debug:    true,
	}
	resp := &bizs.Ip2RegionResp{}
	respData, err := client.PostJSON(fmt.Sprintf("http://%s/%s", localhost, "c/ip2region"), req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if err := unmarshalResp(t, respData, resp); err != nil {
		t.Fatalf("error, unmarshalResp: %#v", err)
	}
}

func Test_Login(t *testing.T) {
	req := &bizs.UserLoginReq{
		Params: &user.LoginRequest{
			InviterId:     "test_inviter_a",
			InviteScene:   "test_scene",
			InvitePayload: nil,
			LoginType:     user.LoginByWxApp,
			WX:            user.LoginWxParams{Code: "--------------------"},
			QQ:            user.LoginQQParams{},
			Mobile:        user.LoginMobileParams{},
		},
	}
	resp := &bizs.UserLoginResp{}
	respData, err := client.PostJSON(fmt.Sprintf("http://%s/%s", localhost, "c/login"), req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if err := unmarshalResp(t, respData, resp); err != nil {
		t.Fatalf("error, unmarshalResp: %#v", err)
	}
}

func Test_AliPay(t *testing.T) {
	req := &bizs.AliPayReq{}
	resp := &bizs.AliPayResp{}
	respData, err := client.PostJSON(fmt.Sprintf("http://%s/%s", localhost, "c/pay/ali"), req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if err := unmarshalResp(t, respData, resp); err != nil {
		t.Fatalf("error, unmarshalResp: %#v", err)
	}
}
