package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// CalMD5 ...
func CalMD5(value string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(value))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
