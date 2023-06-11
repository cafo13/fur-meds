package cors

import "github.com/gin-gonic/gin"

type CORSMiddleware interface {
	Middleware() gin.HandlerFunc
}

type AllowingCORSMiddleware struct{}

func NewAllowingCORSMiddleware() CORSMiddleware {
	return AllowingCORSMiddleware{}
}

func (c AllowingCORSMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
