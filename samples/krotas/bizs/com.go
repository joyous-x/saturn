package bizs

import (
	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/reqresp"
)

// SendSMSReq ...
type SendSMSReq struct {
	reqresp.ReqCommon
	PhoneNum string `json:"phone_num"`
}

// SendSMSResp ...
type SendSMSResp struct {
	reqresp.RespCommon
}

// UserFeedbackReq ...
type UserFeedbackReq struct {
	reqresp.ReqCommon
	Contact string `json:"contact"`
	Content string `json:"content"`
}

// UserFeedbackResp ...
type UserFeedbackResp struct {
	reqresp.RespCommon
}

// SendSMS ...
func SendSMS(c *gin.Context) {
	req := SendSMSReq{}
	resp := SendSMSResp{}
	_, err := reqresp.RequestUnmarshal(c, &req)
	if err != nil {
		reqresp.ResponseMarshal(c, errors.ErrUnmarshalReq, &resp)
		return
	}

	//> TODO: send sms

	reqresp.ResponseMarshal(c, errors.OK, &resp)
}

// UserFeedback ...
func UserFeedback(c *gin.Context) {
	req := UserFeedbackReq{}
	resp := UserFeedbackResp{}

	_, err := reqresp.RequestUnmarshal(c, &req)
	if err != nil {
		reqresp.ResponseMarshal(c, errors.ErrUnmarshalReq, &resp)
		return
	}

	//> TODO: send sms

	reqresp.ResponseMarshal(c, errors.OK, &resp)
}
