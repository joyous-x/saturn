package xlog

import (
	"time"
	"testing"
	"strings"
)

func Test_ZapWrite(t *testing.T) {
	strJsonA := ` {"t_str": "helloworld", "t_int": 1} `
	strJsonB := "{\"HIDMD5\":\"292eb733582806043e0d9328c28a0f12\",\"Class\":\"system\",\"LastUpdateTime\":\"2020/3/31 17:25:27\"}"
	InitZapLogger(true)
	AsyncZapWriteJson(strJsonA)
	AsyncZapWriteJson(strJsonB)

	data := map[string]interface{}{
		"uuid": "test-uuid",
		"os": "test-os",
		"subids": strings.Join([]string{"a", "b", "c"}, ","),
	}
	AsyncZapWriteMap(data)

	time.Sleep(1 * time.Second)
}
