package user

import (
	"sync"

	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/dbs"
	"gopkg.in/redis.v5"
)

// SessHelper with this object, we can handle the uuid and token conveniently
type SessHelper struct {
	redisCli *redis.Client
}

// GetToken get the token for uuid
func (s *SessHelper) GetToken(appid, uuid string) (string, error) {
	return "", nil
}

// UpdateToken update the token for uuid
func (s *SessHelper) UpdateToken(appid, uuid, token string, expireMS int) error {
	return nil
}

// DeleteToken delete the current token of uuid
func (s *SessHelper) DeleteToken(appid, uuid string) error {
	return nil
}

// SetRedis ...
func (s *SessHelper) SetRedis(cli *redis.Client) error {
	s.redisCli = cli
	return nil
}

func (s *SessHelper) makeTokenKey(appid, uuid string) (string, error) {
	return "", nil
}

var gSessHelper *SessHelper
var gSessHelperOnce sync.Once

// SessHelperInst the global instance for SessHelper
func SessHelperInst() *SessHelper {
	gSessHelperOnce.Do(func() {
		gSessHelper = &SessHelper{}
		gSessHelper.SetRedis(dbs.RedisInst().Client("default"))
		xlog.Info("init => SessHelper : ok")
	})
	return gSessHelper
}
