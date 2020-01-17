package tb

import (
	"fmt"
	"github.com/joyous-x/enceladus/common/xlog"
)

const (
	debug = true
)

type TkTbConfig struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	UserId    string `json:"user_id"`
	AdZoneId  string `json:"adzone_id"`
	SiteId    string `json:"site_id"`
}

type ParseForCouponResp struct {
	Found    *TBMyShareItemInfo
	Relative []*TBMyShareItemInfo
}

func init() {
}

// SearchForCoupon
//     NOTE: 1、不能处理：福利
//     TODO: 支持拼多多、网站对接淘宝客、每日推广官方活动、返利积分：创建推广位并绑定，查看订单状态，积分变更
func SearchForCoupon(content string, force bool, cfg *TkTbConfig) (*ParseForCouponResp, error) {
	resp, maxRespLen := &ParseForCouponResp{}, 3
	surl, rurl, title, prodid, err := ParseSharedContent(content, debug)
	if err != nil {
		return resp, err
	} else {
		xlog.Debug("ParseForCoupon.ParseSharedContent ===> surl=%v rurl=%v title=%v prodid=%v", surl, rurl, title, prodid)
	}
	if len(title) < 1 && len(prodid) < 1 {
		return resp, fmt.Errorf("invalid content")
	}
	args := map[string]string{
		"adzone_id":  cfg.AdZoneId,             // mm_xxx_xxx_12345678三段式的最后一段数字
		"site_id":    cfg.SiteId,               // mm_xxx_22_xxx三段式的第二段数字
		"has_coupon": fmt.Sprintf("%v", force), // 优惠券筛选-是否有优惠券。true表示该商品有优惠券，false或不设置表示不限
		"page_size":  "10",                     // 页大小，默认20，1~100
		"page_no":    "1",                      // 第几页，默认：１
		"platform":   "2",                      // 链接形式：1：PC，2：无线，默认：１
		"q":          title,
	}
	respBytes, err := DoTKRequest(cfg.AppKey, cfg.AppSecret, "taobao.tbk.dg.material.optional", args)
	if err != nil {
		return resp, err
	} else {
		xlog.Debug("ParseForCoupon.DoTKRequest ===> args=%v resp=%v", args, string(respBytes))
	}
	respData, err := ParseTKResponse(respBytes)
	if err != nil {
		return resp, err
	} else {
		xlog.Debug("ParseForCoupon.ParseTKResponse ===> args=%+v respData=%+v", args, respData)
	}

	rstList := func(datas map[string]interface{}) []interface{} {
		for _, v := range datas {
			if rst, ok := v.([]interface{}); ok {
				return rst
			}
			continue
		}
		return nil
	}(respData.ResultList)
	if nil == rstList {
		return resp, fmt.Errorf("can't find resultlist")
	}

	var rst map[string]interface{}
	if len(prodid) > 0 {
		for i := range rstList {
			item, err := MakeMyTKShare(cfg.UserId, false, rstList[i].(map[string]interface{}))
			if err != nil {
				xlog.Debug("ParseForCoupon.MakeMyTKShare ===> genTpwd=false index=%v err=%v", i, err)
				continue
			}
			if fmt.Sprintf("%.f", item.ItemId) == prodid {
				rst = rstList[i].(map[string]interface{})
				break
			}
		}
	}
	if rst != nil {
		item, err := MakeMyTKShare(cfg.UserId, true, rst)
		if err != nil {
			xlog.Debug("ParseForCoupon.MakeMyTKShare ===> genTpwd=true data=%+v err=%v", rst, err)
			return resp, err
		}
		resp.Found = item
	} else {
		resp.Relative = make([]*TBMyShareItemInfo, 0, maxRespLen)
		for i := range rstList {
			item, err := MakeMyTKShare(cfg.UserId, true, rstList[i].(map[string]interface{}))
			if err != nil {
				xlog.Debug("ParseForCoupon.MakeMyTKShare ===> genTpwd=false index=%v err=%v", i, err)
				return resp, err
			}
			resp.Relative = append(resp.Relative, item)
			if len(resp.Relative) >= maxRespLen {
				break
			}
		}
	}
	xlog.Debug("ParseForCoupon response ===> resp=%+v", *resp)
	return resp, nil
}
