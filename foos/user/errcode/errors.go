package errcode

import (
	"github.com/joyous-x/saturn/common/errors"
)

func code(id int) int {
	return 10001*1000 + id%1000
}

var (
	OK              = errors.OK
	ErrBadRequest   = errors.NewError(code(001), "bad request")
	ErrAuthForbiden = errors.ErrAuthForbiden
	ErrConfusedUuid = errors.NewError(code(002), "confused user info")
)
