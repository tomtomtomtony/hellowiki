package middleware

import (
	"github.com/gin-gonic/gin"
	"hellowiki/api/result"
	"hellowiki/common/utils"
	"log"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, code := (*utils.JWT).ParserToken(&utils.JWT{}, c)
		if code != result.SUCCSE {
			log.Printf("解析token失败:{%v}", code)
			c.Abort()
		}
		// 将解析后的有效载荷claims重新写入gin.Context引用对象中
		c.Set("claims", claims)
		c.Next()
	}
}
