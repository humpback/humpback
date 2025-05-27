package middleware

import (
	"humpback/types"

	"github.com/gin-gonic/gin"
)

const (
	UserInfoKey   = "userInfo"
	UserCookieKey = "sessionId"
)

func GetUserInfo(c *gin.Context) *types.User {
	info, _ := c.Get(UserInfoKey)
	return info.(*types.User)
}

func SetUserInfo(c *gin.Context, userInfo *types.User) {
	c.Set(UserInfoKey, userInfo)
}

func SetSessionId(c *gin.Context, sessionId string) {
	c.Set(UserCookieKey, sessionId)
}

func GetSessionId(c *gin.Context) string {
	id, _ := c.Get(UserCookieKey)
	return id.(string)
}

func SetCookieSession(c *gin.Context, sessionId string, maxAge int) {
	if maxAge > 0 {
		c.SetCookie(UserCookieKey, sessionId, maxAge, "/webapi", "", false, true)
		return
	}
	c.SetCookie(UserCookieKey, sessionId, -1, "/webapi", "", false, true)
}

func GetCookieSession(c *gin.Context) (string, error) {
	return c.Cookie(UserCookieKey)
}
