package utils

import (
	"fmt"
	"testing"
)

func TestGetLocalFreePort(t *testing.T) {
	port, err := GetLocalFreePort()
	if err != nil {
		t.Errorf("GetLocalFreePort error: %v \n", err)
	} else {
		t.Logf("GetLocalFreePort port: %v \n", port)
	}
}

func TestHttp(t *testing.T) {
	HttpPostJson("http://111.230.250.123:8848/test", map[string]string{"a": "xxx"})
	HttpPostWwwForm("http://111.230.250.123:8848/test", map[string]string{"b": "yyy"})
}

func TestErrEqual(t *testing.T) {
	scene := "test"
	fmtstr := "scene:%v control"
	err := fmt.Errorf(fmtstr, scene)
	if err.Error() == fmt.Sprintf(fmtstr, scene) {
		t.Logf("TestError ok %v", err)
	} else {
		t.Errorf("TestError err %v", err)
	}
}
