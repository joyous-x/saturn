package sms

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// api for qcloud sms functions
const (
	// SENDSMS 发送短信
	SENDSMS string = "sendsms"
	// MULTISMS 群发短信
	MULTISMS string = "sendmultisms2"
	// SENDVOICE 发送语音验证码
	SENDVOICE string = "sendvoice"
	// PROMPTVOICE 发送语音通知
	PROMPTVOICE string = "sendvoiceprompt"
	// ADDSIGN 添加签名
	ADDSIGN string = "add_sign"
)

// url for qcloud sms functions
const (
	// SVR 是腾讯云短信各请求结构的基本 URL
	SVR string = "https://yun.tim.qq.com/v5/"
	// TLSSMSSVR 腾讯云短信业务主URL
	TLSSMSSVR string = "tlssmssvr/"
	// VOICESVR 腾讯云语音URL
	VOICESVR string = "tlsvoicesvr/"
	// TLSSMSSVRAfter 短信业务URL附加内容
	TLSSMSSVRAfter string = "?sdkappid=%s&random=%s"
)

const (
	//SDKName SDK名称，当前主要用于 log 中
	SDKName = "qcloudsms-go-sdk"
	// SDKVersion 版本
	SDKVersion = "0.3.3"
)

// SignReq 请求结构
type SignReq struct {
	Sig    string `json:"sig"`
	Time   int64  `json:"time"`
	Remark string `json:"remark"`
	// 是否为国际短信签名，1=国际，0=国内
	International int `json:"international"`
	// 签名内容，不带"【】"
	Text string `json:"text"`
	// 签名内容相关的证件截图base64格式字符串，非必传参数
	Pic string `json:"pic,omitempty"`
	// 要修改的签名ID
	SignID uint `json:"sign_id,omitempty"`
}

// SignResult 添加/修改/删除 的返回结构
type SignResult struct {
	Result uint   `json:"result"`
	Msg    string `json:"msg"`
	Data   struct {
		ID            uint   `json:"id"`
		International uint   `json:"international,omitempty"`
		Text          string `json:"text"`
		Status        uint   `json:"status"`
	} `json:"data,omitempty"`
}

// SignDelGet 查询/删除签名的请求结构
type SignDelGet struct {
	Sig    string `json:"sig"`
	Time   int64  `json:"time"`
	SignID []uint `json:"sign_id"`
}

// SignStatusResult 短信签名状态返回结构体
type SignStatusResult struct {
	Result uint   `json:"result"`
	Msg    string `json:"msg"`
	Count  uint   `json:"count"`
	Data   []struct {
		ID            uint   `json:"id"`
		Text          string `json:"text"`
		International uint   `json:"international,omitempty"`
		Status        uint   `json:"status"`
		Reply         string `json:"reply"`
		ApplyTime     string `json:"apply_time"`
	} `json:"data"`
}

// SMSTel 国家码，手机号
type SMSTel struct {
	Nationcode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
}

// SMSReq request for SENDSMS
type SMSReq struct {
	Tel    SMSTel   `json:"tel"`
	Type   int      `json:"type,omitempty"`
	Sign   string   `json:"sign,omitempty"`
	TplID  int      `json:"tpl_id,omitempty"`
	Params []string `json:"params"`
	Msg    string   `json:"msg,omitempty"`
	Sig    string   `json:"sig"`
	Time   int64    `json:"time"`
	Extend string   `json:"extend"`
	Ext    string   `json:"ext"`
}

// SMSResult 发送短信返回结构
type SMSResult struct {
	Result       uint   `json:"result"`
	Errmsg       string `json:"errmsg"`
	Ext          string `json:"ext"`
	Sid          string `json:"sid,omitempty"`
	Fee          uint   `json:"fee,omitempty"`
	ActionStatus string `json:"ActionStatus"`
	ErrorCode    uint   `json:"ErrorCode"`
	ErrorInfo    string `json:"ErrorInfo"`
}

// QCloudSms sms sender by qcloud
type QCloudSms struct {
	appID  string
	appKey string
	sign   string
	tplID  int
}

// Init initialize for QCloudSms
func (s *QCloudSms) Init(appID, appKey, sign string, tplID int) error {
	s.appID = appID
	s.appKey = appKey
	s.sign = sign
	s.tplID = tplID
	s.lenRandStr = 6
	s.userAgent = SDKName + "/" + SDKVersion
}

// SendSMS send msg by phonenum
//     nationCode:
//         cn: "86"
func (s *QCloudSms) SendSMS(nationCode, mobile, msg string) error {
	var sm = SMSReq{
		Type:  0,
		Msg:   msg,
		Tel:   SMSTel{Nationcode: nationCode, Mobile: mobile},
		TplID: s.tplID,
	}
	status, err := s.sendSMSSingle(sm)
	if err != nil {
		return err
	}
	if !status {
		return errors.New("send sms failed")
	}
	return nil
}

func (s *QCloudSms) sendSMSSingle(ss SMSSingleReq) error {
	strRand := s.newRandStr(c.lenRandStr)
	urlPath := s.newURL(SENDSMS, s.appID, s.appKey, strRand)
	ss.Time = time.Now().Unix()
	ss.Sig = s.newSig(s.appID, s.appKey, ss.Tel.Mobile, strRand)

	resp, err := c.doRequest(urlPath, s.userAgent, ss, 10*time.Second)
	if err != nil {
		return err
	}
	var res SMSResult
	err = json.Unmarshal(resp, &res)
	if err != nil {
		return err
	}
	if res.Result == 0 && res.ErrorCode == 0 {
		return nil
	}

	return fmt.Errorf("%s,%s", res.Errmsg, res.ErrorInfo)
}

// AddSign function for ADDSIGN
func (s *QCloudSms) AddSign(req SignReq) (SignResult, error) {
	strRand := s.newRandStr(s.lenRandStr)
	urlPath := s.newURL(ADDSIGN, s.appID, s.appKey, strRand)
	req.Time = time.Now().Unix()
	req.Sig = s.newSig(s.appID, s.appKey, "", strRand)

	var res SignResult
	resp, err := s.doRequest(urlPath, s.userAgent, s, 10*time.Second)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *QCloudSms) doRequest(urlPath, userAgent string, params interface{}, timeout time.Duration) ([]byte, error) {
	j, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", urlPath, bytes.NewBuffer([]byte(j)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	httpClient := &http.Client{
		Timeout: timeout,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, ErrRequest
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, err
}

func (s *QCloudSms) newSig(appID, appKey, mobile, strRand string, reqTimeUnix int64) string {
	var t = strconv.FormatInt(reqTimeUnix, 10)
	var sigContent = "appkey=" + appKey + "&random=" + strRand + "&time=" + t
	if len(mobile) > 0 {
		sigContent += "&mobile=" + mobile
	}
	h := sha256.New()
	h.Write([]byte(sigContent))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *QCloudSms) newURL(api, appID, appKey, strRand string) string {
	url := TLSSMSSVR
	if api == SENDVOICE || api == PROMPTVOICE {
		url = VOICESVR
	}
	return SVR + url + api + fmt.Sprintf(TLSSMSSVRAfter, appID, strRand)
}

func (s *QCloudSms) newRandStr(lenRandStr int) string {
	bytes := []byte("0123456789")
	result := make([]byte, lenRandStr)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenRandStr; i++ {
		result[i] = bytes[r.Intn(len(bytes))]
	}
	return string(result)
}
