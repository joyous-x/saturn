package sv

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/reqresp"
)

// URLParserReq ...
type URLParserReq struct {
	reqresp.ReqCommon
	Type     string `json:"type"`
	ShareURL string `json:"share_url"`
}

// URLParserResp ...
type URLParserResp struct {
	reqresp.RespCommon
	Upper UpperInfo `json:"upper"`
	Item  ItemInfo  `json:"item"`
}

// UpperInfo ...
type UpperInfo struct {
}

// ItemInfo ...
type ItemInfo struct {
}

// URLParser ...
func URLParser(c *gin.Context) {
	req := &URLParserReq{}
	ctx, err := reqresp.RequestUnmarshal(c, req)
	if err != nil {
		reqresp.ResponseMarshal(c, errors.ErrUnmarshalReq, nil)
		return
	}

	resp, err := urlParser(ctx, req)
	if err != nil {
		reqresp.ResponseMarshal(c, err, nil)
		return
	}

	reqresp.ResponseMarshal(c, errors.OK, resp)
	return
}

func urlParser(ctx context.Context, req *URLParserReq) (*URLParserResp, error) {
	return &URLParserResp{}, nil
}
