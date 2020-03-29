package reqresp

// ReqCommon common header in a request
type ReqCommon struct {
	Uid       string `json:"uid,omitempty"`
	AppId     string `json:"appid,omitempty"`
	RequestId string `json:"request_id,omitempty"`
	DeviceID  string `json:"device_id,omitempty"`
	Timestamp int64  `json:"timestamp"`
	EchoToken string `json:"echo_token"`
}

// RespCommon common header in a response
type RespCommon struct {
	Ret               int    `json:"ret"`
	Msg               string `json:"msg"`
	RequestId string `json:"request_id,omitempty"`
	Timestamp int64  `json:"timestamp"`
	RetryMS   int32  `json:"retry_ms"`
	EchoToken string `json:"echo_token"`
}

// IRequest the interface of request datas
type IRequest interface {
	GetCommon() *ReqCommon
}

// IResponse the interface of response datas
type IResponse interface {
	GetCommon() *RespCommon
}

func (r *ReqCommon) GetCommon() *ReqCommon {
	return r
}
func (r *RespCommon) GetCommon() *RespCommon {
	return r
}

