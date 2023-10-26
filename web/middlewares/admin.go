package middlewares

import (
	"buding-job/common/utils"
	"buding-job/web/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, _ := c.Get("claims")
		claims := value.(*utils.CustomClaims)
		if claims.Role == 2 {
			c.JSON(http.StatusOK, dto.NewResponse(http.StatusForbidden, "暂无权限", ""))
			c.Abort()
		}
		c.Next()
	}
}
