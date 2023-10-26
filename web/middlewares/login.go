package middlewares

import (
	"buding-job/common/utils"
	"buding-job/web/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusOK, dto.NewResponse(http.StatusUnauthorized, "请登录", ""))
			c.Abort()
			return
		}
		j := utils.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == utils.TokenExpired {
				if err == utils.TokenExpired {
					c.JSON(http.StatusOK, dto.NewResponse(http.StatusUnauthorized, "授权已过期", ""))
					c.Abort()
					return
				}
			}
			c.JSON(http.StatusOK, dto.NewResponse(http.StatusUnauthorized, "未登陆", ""))
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
