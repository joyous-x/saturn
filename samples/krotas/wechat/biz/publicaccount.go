package biz

import (
	"fmt"
	"github.com/joyous-x/saturn/component/wechat/pubacc/message"
	tktb "github.com/joyous-x/saturn/component/tk/tb"
)


func MyMsgHandler(v *message.MixMessage) (*message.Reply, error) {
	resp := &message.Reply{}

	tkcfg := &tktb.TkTbConfig {} // TODO

	switch v.MsgType {
	//文本消息
	case message.MsgTypeText:
		if len(v.Content) < 10 {
			resp = makePubAccHelpRelpy()
		} else {
			coupon, err := tktb.SearchForCoupon(v.Content, true, tkcfg)
			if err == nil {
				resp = makeCouponResp(coupon)
			}
			if err != nil {
				text := message.NewText(fmt.Sprintf("发生了一些错误(%v), 您可以重试下试试看", err))
				resp = &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			}
		}
	//图片消息
	case message.MsgTypeImage:
	//语音消息
	case message.MsgTypeVoice:
	//视频消息
	case message.MsgTypeVideo:
	//小视频消息
	case message.MsgTypeShortVideo:
	//地理位置消息
	case message.MsgTypeLocation:
	//链接消息
	case message.MsgTypeLink:
	//事件推送消息
	case message.MsgTypeEvent:
		resp = eventHandler(v)
	}
	return resp, nil
}

func makeCouponResp(data *tktb.ParseForCouponResp) *message.Reply {
	title := ""
	var target *tktb.TBMyShareItemInfo
	if data != nil && data.Found != nil {
		title = "[玫瑰]找到内部优惠[玫瑰]"
		target = data.Found
	} else if data != nil && data.Relative != nil && len(data.Relative) > 0 {
		title = "[糗大了]您的商品没有内部优惠，[调皮]但是相似产品有~"
		target = data.Relative[0]
	} else {
		title = "没找到您的商品"
	}
	description := fmt.Sprintf("\n【商品】%v\n【店铺】%v\n【折扣价】%v(元)\n【优惠券】[红包][红包][红包]%v \n------\n复制这条信息 %v 到【手机淘宝】领卷下单\n",
		target.Title, target.ShopTitle, target.ZkFinalPrice, target.CouponInfo, target.Tpwd)
	// articleList := []*message.Article{message.NewArticle(title, description, target.PictUrl, "")}
	// return &message.Reply{MsgType: message.MsgTypeNews, MsgData: message.NewNews(articleList)}
	text := message.NewText(fmt.Sprintf("%s\n------\n%s", title, description))
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
}

func makePubAccHelpRelpy() *message.Reply {
	text := message.NewText(fmt.Sprintf("[爱你]欢迎订阅淘优惠(纯净)\n--------\n[右]【使用方法】\n\t在淘宝商品页面 -> 点击'分享' -> '复制链接' -> 此公众号搜索\n\t即可查找内部优惠券"))
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
}

func eventHandler(v *message.MixMessage) *message.Reply {
	switch v.Event {
	//EventSubscribe 订阅
	case message.EventSubscribe:
		return makePubAccHelpRelpy()
	//取消订阅
	case message.EventUnsubscribe:
	//用户已经关注公众号，则微信会将带场景值扫描事件推送给开发者
	case message.EventScan:
	// 上报地理位置事件
	case message.EventLocation:
	// 点击菜单拉取消息时的事件推送
	case message.EventClick:
	// 点击菜单跳转链接时的事件推送
	case message.EventView:
	// 扫码推事件的事件推送
	case message.EventScancodePush:
	// 扫码推事件且弹出“消息接收中”提示框的事件推送
	case message.EventScancodeWaitmsg:
	// 弹出系统拍照发图的事件推送
	case message.EventPicSysphoto:
	// 弹出拍照或者相册发图的事件推送
	case message.EventPicPhotoOrAlbum:
	// 弹出微信相册发图器的事件推送
	case message.EventPicWeixin:
	// 弹出地理位置选择器的事件推送
	case message.EventLocationSelect:
	}
	return nil
}
