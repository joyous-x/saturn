package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"

	"github.com/joyous-x/saturn/common/xlog"
)

// MakeSign 将key按字典序排列后拼接keyvalue最后再添加appSecret值，形成的字符串计算md5
func MakeSign(appid, appSecret string, values url.Values) string {
	values.Add("appid", appid)
	keys := make([]string, 0, len(values))
	for k, _ := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	raw := ""
	for _, v := range keys {
		raw += v + values.Get(v)
	}
	raw += appSecret
	hashed := md5.Sum([]byte(raw))
	xlog.Debug("makeSign raw=%v token=%v sign=%x", raw, appSecret, hashed)
	return fmt.Sprintf("%x", hashed)
}

// MakeHMac ...
func MakeHMac(token string, body []byte) string {
	mac := hmac.New(sha1.New, []byte(token))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))
	return signature
}
