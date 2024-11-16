package middleware

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"

	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URL
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头，相当于Authorization: Bearer xxx.xxx.xxx
		// 🌟这里的具体实现方式要依据你的实际业务情况决定，因为前端返回的数据可能有差异，需要和前端做一个沟通
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割 只要token部分
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			zap.L().Debug("JWT中间件判断errors类型，用来refresh",zap.Error(err))
			if err.Error()=="access_token expired"{
				controller.ResponseError(c, controller.CodeExpiredAToken) //access_token过期
			} else if err.Error()=="refresh_token expired" {
				controller.ResponseError(c, controller.CodeExpiredRToken) //refresh_token过期
			} else if err.Error()=="refresh_token not expired" { //refresh_token没过期
				refAToken, err := jwt.GenAToken(mc.UserID, mc.Username)
				if err != nil {
					controller.ResponseError(c, controller.CodeInvalidToken)
				}
				controller.ResponseErrorWithMsg(c, controller.CodeNotExpiredRToken, refAToken)
			} else {
				controller.ResponseError(c, controller.CodeInvalidToken)
			}
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set(controller.ContextUserIDKey, mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get(ContextUserIDKey)来获取当前请求的用户信息
	}
}
