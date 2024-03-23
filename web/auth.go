/**
* @author: gongquanlin
* @since: 2024/3/23
* @desc:
 */

package web

import (
	"easywebhook/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var tokenMap = map[string]int{}

// AuthMiddleware 认证中间件
func AuthMiddleware(whiteList []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		for _, s := range whiteList {
			if c.Request.URL.Path == s {
				c.Next()
				return
			}
		}

		// 检查是否包含认证信息
		authHeader, err := c.Cookie("session")
		if err != nil || authHeader == "" {
			// c.AbortWithStatus(http.StatusUnauthorized)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		userId, ok := tokenMap[authHeader]
		if !ok {
			// c.AbortWithStatus(http.StatusUnauthorized)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}

// LoginToken 获取登录token
func LoginToken(userId int) string {
	// 随机生成token
	token := utils.GenerateRandomString(32)
	tokenMap[token] = userId
	return token
}
