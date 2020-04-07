package tb

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/joyous-x/saturn/common/xlog"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

// TBRespError 请求失败时的返回信息，json 的根为："error_response"
type TBRespError struct {
	SubMsg  string `json:"sub_msg"`
	SubCode string `json:"sub_code"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
}

// TBRespNorm 请求的返回信息，json 的根为："[method]_response" 的下一层
type TBRespNorm struct {
	RequestId   string                 `json:"request_id"`
	ResultList  map[string]interface{} `json:"result_list"`
	ResultTotal float64                `json:"total_results"`
	ResultData  map[string]interface{} `json:"data"`
}

// TBItemInfo 淘宝商品信息
type TBItemInfo struct {
	CouponStartTime        string              `json:"coupon_start_time"` // 2017-10-29	优惠券信息-优惠券开始时间
	CouponEndTime          string              `json:"coupon_end_time"`   // 2017-10-29	优惠券信息-优惠券结束时间
	InfoDxjh               string              `json:"info_dxjh"`         // 商品信息-定向计划信息, json格式：如，{"19013551":"2850","74510538":"2550"}
	TkTotalSales           string              `json:"tk_total_sales"`    // 商品信息-淘客30天推广量
	CouponId               string              `json:"coupon_id"`         // 优惠券信息-优惠券id
	Title                  string              `json:"title"`
	PictUrl                string              `json:"pict_url"`                  // 商品信息-商品主图
	SmallImages            map[string][]string `json:"small_images"`              // 商品信息-商品小图列表
	ReservePrice           string              `json:"reserve_price"`             // 商品信息-商品一口价格
	ZkFinalPrice           string              `json:"zk_final_price"`            // 折扣价（元） 若属于预售商品，付定金时间内，折扣价=预售价
	UserType               float64             `json:"user_type"`                 // 店铺信息-卖家类型。0表示集市，1表示天猫
	ProvCity               string              `json:"provcity"`                  // 商品信息-宝贝所在地
	ItemUrl                string              `json:"item_url"`                  // 链接-宝贝地址
	IncludeMkt             string              `json:"include_mkt"`               // 商品信息-是否包含营销计划
	IncludeDxjh            string              `json:"include_dxjh"`              // 商品信息-是否包含定向计划
	CommissionRate         string              `json:"commission_rate"`           // 1550表示15.5%	商品信息-佣金比率。1550表示15.5%
	Volume                 float64             `json:"volume"`                    // 商品信息-30天销量
	SellerId               float64             `json:"seller_id"`                 // 店铺信息-卖家id
	CouponTotalCount       float64             `json:"coupon_total_count"`        // 优惠券信息-优惠券总量
	CouponRemainCount      float64             `json:"coupon_remain_count"`       // 优惠券信息-优惠券剩余量
	CouponInfo             string              `json:"coupon_info"`               // 优惠券信息-优惠券满减信息
	CommissionType         string              `json:"commission_type"`           // MKT表示营销计划，SP表示定向计划，COMMON表示通用计划	商品信息-佣金类型。MKT表示营销计划，SP表示定向计划，COMMON表示通用计划
	ShopTitle              string              `json:"shop_title"`                // 店铺信息-店铺名称
	ShopDsr                float64             `json:"shop_dsr"`                  // 店铺信息-店铺dsr评分
	CouponShareUrl         string              `json:"coupon_share_url"`          // 链接-宝贝+券二合一页面链接
	Url                    string              `json:"url"`                       // 链接-宝贝推广链接
	LevelOneCategoryName   string              `json:"level_one_category_name"`   // 商品信息-一级类目名称
	LevelOneCategoryId     float64             `json:"level_one_category_id"`     // 商品信息-一级类目ID
	CategoryName           string              `json:"category_name"`             // 商品信息-叶子类目名称
	CategoryId             float64             `json:"category_id"`               // 商品信息-叶子类目id
	ShortTitle             string              `json:"short_title"`               // 商品信息-商品短标题
	WhiteImage             string              `json:"white_image"`               // 商品信息-商品白底图
	Oetime                 string              `json:"oetime"`                    // 2018-08-21 11:23:43	拼团专用-拼团结束时间
	Ostime                 string              `json:"ostime"`                    // 2018-08-21 11:23:43	拼团专用-拼团开始时间
	JddNum                 float64             `json:"jdd_num"`                   // 拼团专用-拼团几人团
	JddPrice               string              `json:"jdd_price"`                 // 拼团专用-拼团拼成价，单位元
	UvSumPreSale           float64             `json:"uv_sum_pre_sale"`           // 预售专用-预售数量
	Xid                    string              `json:"x_id"`                      // 链接-物料块id(测试中请勿使用)
	CouponStartFee         string              `json:"coupon_start_fee"`          // 优惠券信息-优惠券起用门槛，满X元可用。如：满299元减20元
	CouponAmount           string              `json:"coupon_amount"`             // 优惠券（元） 若属于预售商品，该优惠券付尾款可用，付定金不可用
	ItemDescription        string              `json:"item_description"`          // 商品信息-宝贝描述(推荐理由)
	Nick                   string              `json:"nick"`                      // 店铺信息-卖家昵称
	OrigPrice              string              `json:"orig_price"`                // 拼团专用-拼团一人价（原价)，单位元
	TotalStock             float64             `json:"total_stock"`               // 拼团专用-拼团库存数量
	SellNum                float64             `json:"sell_num"`                  // 拼团专用-拼团已售数量
	Stock                  float64             `json:"stock"`                     // 拼团专用-拼团剩余库存
	TmallPlayActivityInfo  string              `json:"tmall_play_activity_info"`  // 营销-天猫营销玩法: 如，前n件x折
	ItemId                 float64             `json:"item_id"`                   // 商品信息-宝贝id
	RealPostFee            string              `json:"real_post_fee"`             // 商品邮费，单位元
	LockRate               string              `json:"lock_rate"`                 // 锁住的佣金率
	LockRateEndTime        float64             `json:"lock_rate_end_time"`        // 锁佣结束时间, 1567440000000
	LockRateStartTime      float64             `json:"lock_rate_start_time"`      // 锁佣开始时间, 1567440000000
	PresaleDiscountFeeText string              `json:"presale_discount_fee_text"` // 预售商品-优惠: 如，付定金立减5元
	PresaleTailEndTime     float64             `json:"presale_tail_end_time"`     // 预售商品-付尾款结束时间（毫秒）
	PresaleTailStartTime   float64             `json:"presale_tail_start_time"`   // 预售商品-付尾款开始时间（毫秒）
	PresaleEndTime         float64             `json:"presale_end_time"`          // 预售商品-付定金结束时间（毫秒）
	PresaleStartTime       float64             `json:"presale_start_time"`        // 预售商品-付定金开始时间（毫秒）
	PresaleDeposit         string              `json:"presale_deposit"`           // 预售商品-定金（元）
}

// TBMyShareItemInfo 淘宝商品分享时需要的信息
type TBMyShareItemInfo struct {
	TBItemInfo
	Tpwd string `json:"tpwd"`
}

func parseAliUrl(s string) (short bool, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}
	host := u.Host
	if host == "www.dwntme.com" || host == "dwntme.com" || strings.Index(s, ".tb.cn/h.") >= 4 {
		short = true
	} else if host == "detail.tmall.com" || host == "item.taobao.com" {
		short = false
	} else {
		err = fmt.Errorf("not ali url")
	}
	return
}

func getOriginalFromdwntme(s string) string {
	getOriginalUrl := func(s string) string {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req, err := http.NewRequest("GET", s, nil)
		resp, err := client.Do(req)
		if err != nil {
			return ""
		}
		defer resp.Body.Close()
		return resp.Header.Get("Location")
	}
	rurl := ""
	client := &http.Client{}
	req, err := http.NewRequest("GET", s, nil)
	resp, err := client.Do(req)
	if err != nil {
		return rurl
	}
	defer resp.Body.Close()
	httpbyte, _ := ioutil.ReadAll(resp.Body)
	reg := regexp.MustCompile(`var url = '(http.+)';`)
	matches := reg.FindStringSubmatch(string(httpbyte))
	if len(matches) == 2 {
		location := getOriginalUrl(matches[1])
		if len(location) < 1 {
			rurl = matches[1]
		} else {
			rurl = location
		}
	}
	return rurl
}

func makeSignHmac(values url.Values, secret string) string {
	data := ""
	keys := make([]string, 0, len(values))
	for k, _ := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := values.Get(k)
		if len(k) > 0 && len(v) > 0 {
			data += k + v
		}
	}
	mac := hmac.New(md5.New, []byte(secret))
	mac.Write([]byte(data))
	signature := strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
	return signature
}

// ParseSharedContent 解析分享内容
//    二合一链接
//       二合一链接是可以将优惠券链接和商品链接合并为一条券品二合一链接,为淘客推广提高转化率。
//       二合一链接是以uland.taobao.com开头的，二合一链接的参数还多种形式，包含三个信息：
//           activityId：优惠券id, itemId：商品id, pid：您的淘客推广pid
func ParseSharedContent(c string, debug bool) (surl, rurl, title, prodid string, err error) {
	title, surl, rurl, prodid, err = "", "", "", "", nil
	// 1. 解析title 和 url
	c = strings.Trim(c, " ")
	reg := regexp.MustCompile(`【+(.+?)】.*?\s*((https|http|ftp|rtsp|mms):\/\/\S+)?`)
	if c[:1] != "【" && strings.Count(c, "【") > 2 && strings.Count(c, "【") == strings.Count(c, "】") {
		// 解析：淘宝 ----> 分享 ---> 复制链接
		reg = regexp.MustCompile(`(.+?)\s(.|\s)*?【.*?】((https|http|ftp|rtsp|mms):\/\/\S+)?\s*-+.*?`)
	}
	matches := reg.FindStringSubmatch(c)
	if debug {
		for i := range matches {
			xlog.Debug("===> ParseSharedContent index=%v %v \n", i, matches[i])
		}
	}
	if len(matches) >= 1 {
		title = matches[1]
	}
	tmpUrl := ""
	if len(matches) == 4 {
		tmpUrl = matches[2]
	} else if len(matches) == 5 {
		tmpUrl = matches[3]
	}
	if len(tmpUrl) < 1 {
		return
	}
	// 2. 解析短url获取原始链接
	short, err := parseAliUrl(tmpUrl)
	if err == nil {
		if short {
			surl = tmpUrl
		} else {
			rurl = tmpUrl
		}
	} else {
		//> logs
	}
	// 3. 解析长url获取更多信息
	if len(rurl) < 1 && len(surl) > 0 {
		rurl = getOriginalFromdwntme(surl)
	}
	if len(rurl) < 1 {
		return
	}
	// 4. 获取精准商品id
	u, err := url.Parse(rurl)
	if err != nil {
		return
	}
	args, _ := url.ParseQuery(u.RawQuery)
	if info, ok := args["trackInfo"]; ok {
		fields := strings.Split(info[0], "_")
		for i := range fields {
			k := strings.Split(fields[i], ":")
			if k[0] == "itemId" && len(k) == 2 {
				prodid = k[1]
			}
		}
	}
	if u.Host == "uland.taobao.com" {
		if v, ok := args["e"]; ok {
			// TODO: 二合一连接，解析 e 参数
			xlog.Debug("===> uland - e : %v\n", v)
		}
	}
	if u.Host == "item.taobao.com" || u.Host == "detail.tmall.com" {
		if v, ok := args["id"]; ok {
			prodid = v[0]
		}
	}
	return
}

// DoTKRequest 请求接口
// 	  note: md5 必须大写
func DoTKRequest(appKey, appSecret, method string, data map[string]string) ([]byte, error) {
	host := "http://gw.api.taobao.com/router/rest"
	args := url.Values{}
	// 公共参数
	args.Add("app_key", appKey)
	args.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	args.Add("format", "json")
	args.Add("v", "2.0")
	args.Add("sign_method", "hmac")
	// 请求参数
	for k, v := range data {
		args.Add(k, v)
	}
	args.Add("method", method)
	// 签名参数
	args.Add("sign", makeSignHmac(args, appSecret))
	// Do Request
	httpClient := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}
	req, err := http.NewRequest("POST", host, strings.NewReader(args.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	//req.Header.Set("Accept", "text/xml,text/javascript")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	return d, err
}

// ParseTKResponse 解析 DoTKRequest 的请求结果
func ParseTKResponse(respBytes []byte) (*TBRespNorm, error) {
	var err error
	var rst = &TBRespNorm{}
	datas := make(map[string]interface{}, 0)
	if err = json.Unmarshal(respBytes, &datas); err != nil {
		return rst, err
	}
	if m, ok := datas["error_response"]; ok {
		if v, ok := m.(map[string]interface{}); ok {
			subCode, subMsg := "", ""
			if d, ok := v["sub_code"].(string); ok {
				subCode = d
			}
			if d, ok := v["sub_msg"].(string); ok {
				subMsg = d
			}
			err = fmt.Errorf("parseResp error: code=%v msg=%v sub_code=%v sub_msg=%v", v["code"].(float64), v["msg"].(string), subCode, subMsg)
		} else {
			err = fmt.Errorf("parseResp error: unknown")
		}
		return rst, err
	}
	if len(datas) != 1 {
		err = fmt.Errorf("parseResp somethings error : len(datas) != 1")
		return rst, err
	}
	for _, data := range datas {
		if d, ok := data.(map[string]interface{}); ok {
			for k, v := range d {
				switch k {
				case "request_id":
					rst.RequestId = v.(string)
				case "result_list":
					rst.ResultList = v.(map[string]interface{})
				case "total_results":
					rst.ResultTotal = v.(float64)
				case "data":
					rst.ResultData = v.(map[string]interface{})
				default:
					err = fmt.Errorf("parseResp invalid key(%v)", k)
				}
			}
		} else {
			err = fmt.Errorf("parseResp datas error : invalid data type")
		}
		break
	}
	return rst, err
}

// GenerateTpwd 生成淘口令
// 			user_id	String	false	123	生成口令的淘宝用户ID
// 			text	String	true	长度大于5个字符	口令弹框内容
// 			url		String	true	https://uland.taobao.com/	口令跳转目标页
// 			logo	String	false	https://uland.taobao.com/	口令弹框logoURL
// 			ext		String	false	{}	扩展字段JSON格式
func GenerateTpwd(appKey, appSecret, text, url, logo, userId, ext string) (string, error) {
	rst := ""
	if len(text) < 1 || len(url) < 1 {
		return rst, fmt.Errorf("invalid param: text or url")
	}
	reqArgs := map[string]string{
		"text": text,
		"url":  url,
	}
	if len(logo) > 0 {
		reqArgs["logo"] = logo
	}
	if len(ext) > 0 {
		reqArgs["ext"] = ext
	}
	if len(userId) > 0 {
		reqArgs["user_id"] = userId
	}
	respBytes, err := DoTKRequest(appKey, appSecret, "taobao.tbk.tpwd.create", reqArgs)
	if err != nil {
		return rst, err
	} else {
		xlog.Debug("GenerateTpwd.DoTKRequest ===> args=%+v\n", reqArgs)
	}
	respData, err := ParseTKResponse(respBytes)
	if err != nil {
		return rst, err
	}
	if _, ok := respData.ResultData["model"]; !ok {
		return rst, fmt.Errorf("response don't have model")
	} else {
		rst = respData.ResultData["model"].(string)
	}
	return rst, err
}

// MakeMyTKShare 解析 DoTKRequest 的请求结果
func MakeMyTKShare(appKey, appSecret, userId string, genTpwd bool, info map[string]interface{}) (*TBMyShareItemInfo, error) {
	prodInfo := &TBMyShareItemInfo{}
	infoStr, err := json.Marshal(info)
	if err != nil {

		return prodInfo, err
	}
	err = json.Unmarshal(infoStr, prodInfo)
	if err != nil {
		return prodInfo, err
	}
	if genTpwd {
		urls := prodInfo.Url
		if len(prodInfo.CouponShareUrl) > 0 {
			urls = "https:" + prodInfo.CouponShareUrl
		}
		pictUrl := prodInfo.PictUrl
		pictUrl = ""
		tpwd, err := GenerateTpwd(appKey, appSecret, prodInfo.Title, urls, pictUrl, userId, "{}")
		if err != nil {
			fmt.Printf("MakeMyTKShare error: %v %v \n", err, fmt.Sprintf("%v%v", "https:", prodInfo.CouponShareUrl))
			return prodInfo, fmt.Errorf("GenerateTpwd error(%v) url=%v", err, urls)
		}
		prodInfo.Tpwd = tpwd
	}
	return prodInfo, nil
}
