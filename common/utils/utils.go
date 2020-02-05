package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	guuid "github.com/google/uuid"
	"hash/crc64"
	"strconv"
	"strings"
	"time"
)

// CalMD5 ...
func CalMD5(value string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(value))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// NewUUID ...
func NewUUID(appname, uniqueID string) string {
	return CalMD5(strings.Replace(guuid.New().String(), "-", "", -1) + appname + uniqueID)
}

// NewToken ...
func NewToken(appname, uniqueID string) string {
	suffix := strconv.FormatInt(time.Now().UnixNano(), 10) + guuid.New().String()
	hash64 := crc64.Checksum([]byte(appname+uniqueID+suffix), crc64.MakeTable(crc64.ISO))
	return fmt.Sprintf("%x", hash64)
}
