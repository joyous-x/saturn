package biz

import(
	"github.com/gin-gonic/gin"
)

// GetUserInfo ...
func GetUserInfo(appname, uuid string) (token string, err error) {
	return "token_test", nil
}

// UserLogin user login via phone and third account 
// such as wechat, qq, and so on
func UserLogin(c *gin.Context) {

}


func MiniappWxLogin(c *gin.Context) {
	
}