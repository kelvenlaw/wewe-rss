package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置API路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	
	// 解决跨域问题
	r.Use(CORSMiddleware())
	
	// API路由组
	api := r.Group("/api/v2")
	
	// 登录相关接口
	api.GET("/login/platform", CreateLoginURL)
	api.GET("/login/platform/:id", GetLoginResult)
	
	// 公众号文章相关接口
	api.GET("/platform/mps/:mpId/articles", GetMpArticles)
	api.POST("/platform/wxs2mp", GetMpInfo)
	
	return r
}

// CORSMiddleware 处理跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, xid")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
} 