package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wewe-rss/authserver/internal/api"
	"github.com/wewe-rss/authserver/internal/app"
)

func main() {
	// 加载配置
	config := app.GetConfig()
	
	// 打印启动信息
	log.Printf("=== WeWe-RSS AuthServer ===")
	log.Printf("Configuration: Host=%s, Port=%d, Debug=%v", config.Host, config.Port, config.Debug)
	
	// 设置路由
	router := api.SetupRouter()
	
	// 添加一个简单的健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
	
	// 启动服务器
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Printf("Starting WeWe-RSS AuthServer on %s\n", address)
	
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 