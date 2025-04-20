package app

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Port            int
	Host            string
	WeReadProxyURL  string // 微信读书API代理URL，可选
	DefaultArticles int    // 默认文章数量
	Debug           bool
	QRCodeSize      int    // 二维码尺寸
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 尝试加载.env文件，如果存在的话
	_ = godotenv.Load()

	port, _ := strconv.Atoi(getEnv("PORT", "8080"))
	defaultArticles, _ := strconv.Atoi(getEnv("DEFAULT_ARTICLES", "20"))
	qrCodeSize, _ := strconv.Atoi(getEnv("QR_CODE_SIZE", "200"))
	debug := getEnv("DEBUG", "false") == "true"

	return &Config{
		Port:            port,
		Host:            getEnv("HOST", "0.0.0.0"),
		WeReadProxyURL:  getEnv("WEREAD_PROXY_URL", ""),
		DefaultArticles: defaultArticles,
		Debug:           debug,
		QRCodeSize:      qrCodeSize,
	}
}

// 从环境变量获取配置，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if appConfig == nil {
		log.Println("Loading configuration...")
		appConfig = LoadConfig()
	}
	return appConfig
}

var appConfig *Config 