package errors

import (
	comerrs "github.com/joyous-x/saturn/common/errors"
)

var (
	ErrIp2regionMemSearch = comerrs.NewError(11100001, "ip mem search error")
	ErrUserLogin          = comerrs.NewError(11100002, "user login error")
	ErrLoginByWxMiniApp   = comerrs.NewError(11100003, "wx miniapp login error")
	ErrGetUserInfo        = comerrs.NewError(11100004, "get user info error")
	ErrUpdateUserInfo     = comerrs.NewError(11100005, "update user info error")
	ErrDecryptUserInfo    = comerrs.NewError(11100006, "decrypt user info error")
	ErrGetAccessToken     = comerrs.NewError(11100007, "get access token error")
)
