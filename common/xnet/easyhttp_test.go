package xnet

import (
	"testing"
)

func Test_EasyHTTP(t *testing.T) {
	datas := map[string]string{
		"test": "helloworld",
	}
	resp, err := NewEasyHTTP().Options(nil).PostJSON("www.baidu.com", datas)
	if err != nil {
		t.Errorf("postjson error: %v\n", err)
	} else {
		t.Errorf("postjson success: %#v\n", str(resp))
	}
}
