package reqresp

type ReqCommonData struct {
	Uid       string `json:"uid,omitempty"`
	AppId     string `json:"appid,omitempty"`
	RequestId string `json:"request_id,omitempty"`
	DeviceID  string `json:"device_id,omitempty"`
	Timestamp int64  `json:"timestamp"`
	EchoToken string `json:"echo_token"`
}

type RespCommonData struct {
	Ret               int    `json:"ret"`
	Msg               string `json:"msg"`
	RequestId string `json:"request_id,omitempty"`
	Timestamp int64  `json:"timestamp"`
	RetryMS   int32  `json:"retry_ms"`
	EchoToken string `json:"echo_token"`
}

// ReqCommon common header in a request
type ReqCommon struct {
	Common ReqCommonData `json:"common"`
}

// RespCommon common header in a response
type RespCommon struct {
	Common RespCommonData `json:"common"`
}

// IRequest the interface of request datas
type IRequest interface {
	GetCommon() *ReqCommonData
}

// IResponse the interface of response datas
type IResponse interface {
	GetCommon() *RespCommonData
}

func (r *ReqCommon) GetCommon() *ReqCommonData {
	return &r.Common
}
func (r *RespCommon) GetCommon() *RespCommonData {
	return &r.Common
}

