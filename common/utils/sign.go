package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/joyous-x/saturn/common/xlog"
	"net/url"
	"sort"
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
