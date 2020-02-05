package pubacc

import (
	"ceggs/wechat"
	"ceggs/wechat/pubacc/message"
	"ceggs/wechat/pubacc/util"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/enceladus/common/xlog"
	"io"
	"io/ioutil"
	"reflect"
	"sort"
	"strconv"
	"time"
)

type MsgHandler func(*message.MixMessage) (*message.Reply, error)

type MsgRequestData struct {
	OpenID    string `json:"openid"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
	// echostr 随机字符串, 用于校验开发者接入是否成功
	//     若确认此次GET请求来自微信服务器，请原样返回echostr参数内容，则接入生效，否则接入失败
	EchoStr      string `json:"echostr"`
	EchoStrExist bool   `json:"echostr_exist"`
	// encrypt_type == aes
	SafeMode bool
	Random   []byte
	// input message data
	InputMixMsg *message.MixMessage
}

func makeSign(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func checkSign(signature string, params ...string) bool {
	return signature == makeSign(params...)
}

func ParseRequestFromGin(c *gin.Context, wxcfg *wechat.WxConfig, needCheckSign bool) (*MsgRequestData, error) {
	xlog.Debug("ParseRequestFromGin ===> Params = %+v", c.Params)
	xlog.Debug("ParseRequestFromGin ===> Request.URL = %+v", c.Request.URL.Query())
	resp := &MsgRequestData{}
	resp.OpenID = c.Query("openid")
	resp.Nonce = c.Query("nonce")
	resp.Timestamp = c.Query("timestamp")
	resp.Signature = c.Query("signature")
	resp.EchoStr, resp.EchoStrExist = c.GetQuery("echostr")
	resp.SafeMode = c.Query("encrypt_type") == "aes"
	resp.InputMixMsg = &message.MixMessage{}
	token := wxcfg.PubAccToken
	xlog.Debug("ParseRequestFromGin token=%v signature=%v openid=%v nonce=%v timestamp=%v echostr=%v encrypt_type=%v",
		token, resp.Signature, resp.OpenID, resp.Nonce, resp.Timestamp, resp.EchoStr, c.Query("encrypt_type"))
	// appName = c.Query("appname")
	// pubName = c.Query("pubname")
	// unionID = c.Query("unionid")
	mysign := makeSign(token, resp.Timestamp, resp.Nonce)
	if needCheckSign && resp.Signature != mysign {
		xlog.Error("ParseRequestFromGin error: expectSign=%v sign=%v", mysign, resp.Signature)
		return resp, fmt.Errorf("invalid signature(%v)", resp.Signature)
	}
	if resp.EchoStrExist {
		return resp, nil
	}

	var err error
	var rawXMLMsgBytes []byte
	if resp.SafeMode {
		var encryptedXMLMsg message.EncryptedXMLMsg
		if err := xml.NewDecoder(c.Request.Body).Decode(&encryptedXMLMsg); err != nil {
			return resp, fmt.Errorf("parse body error(%v)", err)
		}

		msgSignature := c.Query("msg_signature")
		if !checkSign(msgSignature, token, resp.Timestamp, resp.Nonce, encryptedXMLMsg.EncryptedMsg) {
			return resp, fmt.Errorf("invalid msg signature")
		}

		resp.Random, rawXMLMsgBytes, err = util.DecryptMsg(wxcfg.AppID, encryptedXMLMsg.EncryptedMsg, wxcfg.EncodingAESKey)
		if err != nil {
			return resp, fmt.Errorf("decrypt msg error(%v)", err)
		}
	} else {
		rawXMLMsgBytes, err = ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return resp, fmt.Errorf("request body error(%v)", err)
		}
	}
	xlog.Debug("ParseRequestFromGin ===> Request.Body = %+v", string(rawXMLMsgBytes))

	err = xml.Unmarshal(rawXMLMsgBytes, resp.InputMixMsg)
	if err != nil {
		return resp, fmt.Errorf("xml unmarshal error(%v)", err)
	} else {
		xlog.Debug("ParseRequestFromGin inputmixmsg: %v", resp.InputMixMsg)
	}
	return resp, nil
}

// InteractEntry the entry of user and public account server interactive with each other
//     return reply message
func InteractEntry(wxcfg *wechat.WxConfig, req *MsgRequestData, handler MsgHandler) (interface{}, error) {
	reply, err := handler(req.InputMixMsg)
	if err != nil || reply == nil {
		if reply == nil {
			err = fmt.Errorf("reply is nil")
		}
		return nil, fmt.Errorf("replay error(%v)", err)
	}

	msgType := reply.MsgType
	switch msgType {
	case message.MsgTypeText:
	case message.MsgTypeImage:
	case message.MsgTypeVoice:
	case message.MsgTypeVideo:
	case message.MsgTypeMusic:
	case message.MsgTypeNews:
	case message.MsgTypeTransfer:
	default:
		err = fmt.Errorf("invalid MsgType(%v)", msgType)
		return nil, err
	}

	msgData := reply.MsgData
	value := reflect.ValueOf(msgData)
	kind := value.Kind().String() //msgData must be a ptr
	if "ptr" != kind {
		err = fmt.Errorf("invalid MsgDataType(%v)", kind)
		return nil, err
	}

	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(req.InputMixMsg.FromUserName)
	value.MethodByName("SetToUserName").Call(params)
	params[0] = reflect.ValueOf(req.InputMixMsg.ToUserName)
	value.MethodByName("SetFromUserName").Call(params)
	params[0] = reflect.ValueOf(msgType)
	value.MethodByName("SetMsgType").Call(params)
	params[0] = reflect.ValueOf(time.Now().Unix())
	value.MethodByName("SetCreateTime").Call(params)

	replyMsg := msgData
	if req.SafeMode {
		respRawXMLMsg, err := xml.Marshal(replyMsg)
		if err != nil {
			return replyMsg, err
		}
		var encryptedMsg []byte
		encryptedMsg, err = util.EncryptMsg(req.Random, respRawXMLMsg, wxcfg.AppID, wxcfg.EncodingAESKey)
		if err != nil {
			return replyMsg, err
		}
		timestamp := func() int64 {
			ret, err := strconv.ParseInt(req.Timestamp, 10, 32)
			if err != nil {
				ret = time.Now().Unix()
			}
			return ret
		}()
		msgSignature := makeSign(wxcfg.PubAccToken, strconv.FormatInt(timestamp, 10), req.Nonce, string(encryptedMsg))
		replyMsg = message.ResponseEncryptedXMLMsg{
			EncryptedMsg: string(encryptedMsg),
			MsgSignature: msgSignature,
			Timestamp:    timestamp,
			Nonce:        req.Nonce,
		}
	}
	return replyMsg, nil
}
