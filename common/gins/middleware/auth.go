package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joyous-x/saturn/common/errors"
	"github.com/joyous-x/saturn/common/reqresp"
	"github.com/joyous-x/saturn/common/utils"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/foos/user"
)

//
// Authentication & Authorization
//

// AuthToken ...
func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCom := reqresp.ReqCommon{}
		ctx, err := reqresp.RequestUnmarshal(c, &reqCom)
		if err != nil {
			reqresp.ResponseMarshal(c, err, nil)
			c.Abort()
		}
		appid := reqresp.CtxGetStr(ctx, reqresp.CtxKeyAppID)
		uuid := reqresp.CtxGetStr(ctx, reqresp.CtxKeyUuid)
		token, err := user.SessHelperInst().GetToken(appid, uuid)
		if err != nil {
			reqresp.ResponseMarshal(c, err, nil)
			c.Abort()
		}
		if token != reqCom.GetCommon().SessToken {
			reqresp.ResponseMarshal(c, errors.ErrAuthForbiden, nil)
			c.Abort()
		}
		c.Next()
	}
}

// AuthSign ...
func AuthSign() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCom := reqresp.ReqCommon{}
		ctx, err := reqresp.RequestUnmarshal(c, &reqCom)
		if err != nil {
			reqresp.ResponseMarshal(c, err, nil)
			c.Abort()
		}

		appid := reqresp.CtxGetStr(ctx, reqresp.CtxKeyAppID)
		uuid := reqresp.CtxGetStr(ctx, reqresp.CtxKeyUuid)
		reqRawPacket := reqresp.CtxGetRaw(ctx, reqresp.CtxKeyRequestData)

		if err := authSignHandler(appid, uuid, c.GetHeader("Authentication"), reqRawPacket); err != nil {
			reqresp.ResponseMarshal(c, err, nil)
			c.Abort()
		}

		c.Next()
	}
}

// authSignHandler 认证用户信息
func authSignHandler(appid, uuid, authentication string, dataBody []byte) error {
	token, err := user.SessHelperInst().GetToken(appid, uuid)
	if err != nil {
		return err
	}

	serverSign := utils.MakeHMac(token, dataBody)
	if serverSign != authentication {
		xlog.Error("signature client:%s, server:%s, token:%s, body:%s", authentication, serverSign, token, string(dataBody))
		return fmt.Errorf("Authentication check fail")
	}

	return nil
}
