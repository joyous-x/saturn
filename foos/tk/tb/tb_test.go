package tb

import (
	"fmt"
	"testing"
)

func Test_ParseContent(t *testing.T) {
	config := &TkTbConfig{
		UserId:    "488285073",
		AppKey:    "27995011",
		AppSecret: "b59cff7b0cb4f904fb5a393ebca9bac5",
		AdZoneId:  "109586600288",
		SiteId:    "933150027",
	}

	content := func(index string) string {
		contentA := `南极人被子冬被芯冬天加厚保暖宿舍单人学生双人空调春秋冬季棉被 【包邮】
		【在售价】79.00元
		【券后价】59.00元
		【下单链接】https://m.tb.cn/h.epk7MKO 
		----------------- 
		复制这条信息，$5vLKYKRwjvA$，到【手机淘宝】即可查看`
		contentB := `【Ciate丝绒亮泽口红哑光雾面丝绒口红不脱色枫叶红保湿防水西柚色】https://m.tb.cn/h.epGgMsj?sm=53c47a 嚸↑↓擊鏈ㄣ接，再选择瀏覽嘂..咑№亓；或復zんíゞ这句话€DUxNYKswah4€后咑閞淘灬寳`
		contentC := `【秋冬季棉拖鞋女厚底室内保暖防滑家居家用毛拖鞋男士托鞋女士冬天】https://c.tb.cn/h.epkwXjX?sm=962ae4 點￡擊☆鏈バ接，再选择瀏覽→噐咑№亓；或復ず■淛这句话₤25bxYqgePw4₤后咑閞綯℡寳`
		switch index {
		case "a":
			return contentA
		case "b":
			return contentB
		case "c":
			return contentC
		default:
			return ""
		}
	}
	resp, err := SearchForCoupon(content("c"), true, config)
	if err != nil {
		t.Log(fmt.Sprintf("--------- ParseForCoupon err=%v", err))
	} else {
		if resp.Found != nil {
			t.Log(fmt.Sprintf("--------- found: %+v", *resp.Found))
		}
		if resp.Relative != nil && len(resp.Relative) > 0 {
			for i := range resp.Relative {
				t.Log(fmt.Sprintf("--------- relative: index=%v %+v", i, *resp.Relative[i]))
			}
		}
	}
}
