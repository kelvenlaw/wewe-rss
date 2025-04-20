package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"github.com/wewe-rss/authserver/internal/app"
	"github.com/wewe-rss/authserver/internal/model"
)

// 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}

// 维护登录会话信息
var loginSessions = make(map[string]*model.LoginResult)

// CreateLoginURL 创建登录URL
func CreateLoginURL(c *gin.Context) {
	// 从微信读书API获取UID
	uid, err := getLoginUIDFromWeRead()
	if err != nil {
		// 如果API调用失败，使用本地生成的UID作为备选方案
		log.Printf("Failed to get UID from WeRead API, using local fallback: %v", err)
		uid = generateFallbackUID()
	}
	
	log.Printf("Login UID123: %s", uid)
	
	// 实际登录确认URL
	confirmURL := fmt.Sprintf("https://weread.qq.com/web/confirm?uid=%s", uid)
	
	// 获取配置
	config := app.GetConfig()
	
	// 生成二维码 - 包含确认URL
	// var qrCode []byte
	// qrCode, err = qrcode.Encode(confirmURL, qrcode.Medium, config.QRCodeSize)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
	// 	return
	// }
	
	// // 转换为base64
	// base64QR := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrCode)
	
	// 初始化登录会话
	loginSessions[uid] = &model.LoginResult{Message: "等待用户扫码"}
	
	// 模拟登录成功（实际应用中需要实现真实的登录逻辑）
	go func() {
		// 延迟30秒后自动模拟登录成功
		time.Sleep(30 * time.Second)
		if result, exists := loginSessions[uid]; exists && result.Vid == 0 {
			// 模拟一个成功的登录结果
			loginSessions[uid] = &model.LoginResult{
				Vid:      12345 + int(time.Now().Unix()%1000),
				Token:    "sample_token_" + uid,
				Username: "WeRead_User_" + uid[0:8],
			}
			log.Printf("Login successful for UID: %s\n", uid)
		}
	}()
	
	// 将新版的响应结构合并到接口中
	c.JSON(http.StatusOK, gin.H{
		"uuid": uid, 
		"scanUrl": confirmURL,  // 提供Base64编码的二维码
	})
}

// 从微信读书API获取登录UID
func getLoginUIDFromWeRead() (string, error) {
	// 配置
	config := app.GetConfig()
	
	// 如果设置了代理URL，使用代理URL
	baseURL := "https://weread.qq.com"
	if config.WeReadProxyURL != "" {
		baseURL = config.WeReadProxyURL
	}
	
	// 构建API请求URL
	apiURL := fmt.Sprintf("%s/api/auth/getLoginUid", baseURL)
	
	// 发起HTTP请求
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()
	
	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body: %v", err)
	}
	
	// 解析JSON
	var result struct {
		UID string `json:"uid"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("Failed to parse JSON: %v", err)
	}
	
	if result.UID == "" {
		return "", fmt.Errorf("Empty UID received from API")
	}
	
	return result.UID, nil
}

// 生成备用UID（当API调用失败时使用）
func generateFallbackUID() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s", 
		generateRandomString(8),
		generateRandomString(4),
		generateRandomString(4),
		generateRandomString(4),
		generateRandomString(12))
}

// GetLoginResult 获取登录结果
func GetLoginResult(c *gin.Context) {
	id := c.Param("id")
	
	// 实际应用中应调用 https://weread.qq.com/api/auth/getLoginInfo?uid={uid}
	// 这里通过本地会话模拟登录结果，实际环境中应替换
	
	result, err := getLoginInfoFromWeRead(id)
	if err != nil {
		// 如果API调用失败，回退到本地会话
		log.Printf("Failed to get login info from WeRead API: %v", err)
		
		// 检查登录会话是否存在
		localResult, exists := loginSessions[id]
		if !exists {
			c.JSON(http.StatusOK, model.LoginResult{
				Message: "登录会话不存在或已过期",
			})
			return
		}
		
		// 返回本地保存的登录状态
		c.JSON(http.StatusOK, localResult)
		return
	}
	
	// 更新本地会话
	loginSessions[id] = result
	
	// 返回API获取的登录状态
	c.JSON(http.StatusOK, result)
}

// 从微信读书API获取登录信息
func getLoginInfoFromWeRead(uid string) (*model.LoginResult, error) {
	// 配置
	config := app.GetConfig()
	
	// 如果设置了代理URL，使用代理URL
	baseURL := "https://weread.qq.com"
	if config.WeReadProxyURL != "" {
		baseURL = config.WeReadProxyURL
	}
	
	// 构建API请求URL
	apiURL := fmt.Sprintf("%s/api/auth/getLoginInfo?uid=%s", baseURL, uid)
	
	// 对于测试环境，返回模拟数据
	if config.Debug {
		// 模拟API响应
		mockResult := &model.LoginResult{
			Message: "模拟数据，仅用于调试",
		}
		
		// 检查本地会话状态
		if localResult, exists := loginSessions[uid]; exists && localResult.Vid != 0 {
			// 如果本地会话已登录，返回成功结果
			mockResult.Vid = localResult.Vid
			mockResult.Token = localResult.Token
			mockResult.Username = localResult.Username
			mockResult.Message = ""
		}
		
		return mockResult, nil
	}
	
	// 实际环境中，这里应实现真实的API调用
	// 临时返回错误以触发本地会话回退
	return nil, fmt.Errorf("API integration not implemented")
}

// GetMpArticles 获取公众号文章列表
func GetMpArticles(c *gin.Context) {
	mpID := c.Param("mpId")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	
	// 验证认证头
	authToken := c.GetHeader("Authorization")
	xid := c.GetHeader("xid")
	
	if authToken == "" || xid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "WeReadError401"})
		return
	}
	
	// 获取配置
	config := app.GetConfig()
	
	articles := []model.Article{}
	// 生成模拟文章数据
	for i := 0; i < config.DefaultArticles; i++ {
		articles = append(articles, model.Article{
			ID:          fmt.Sprintf("article_%s_%d_%d", mpID, page, i),
			Title:       fmt.Sprintf("文章标题 %d - 页码 %d", i, page),
			PicURL:      "https://example.com/pic.jpg",
			PublishTime: time.Now().AddDate(0, 0, -i).Unix(),
		})
	}
	
	c.JSON(http.StatusOK, articles)
}

// GetMpInfo 获取公众号信息
func GetMpInfo(c *gin.Context) {
	var request model.MPInfoRequest
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "WeReadError400"})
		return
	}
	
	// 验证认证头
	authToken := c.GetHeader("Authorization")
	xid := c.GetHeader("xid")
	
	if authToken == "" || xid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "WeReadError401"})
		return
	}
	
	// 验证URL是否是有效的微信文章URL
	if !strings.HasPrefix(request.URL, "https://mp.weixin.qq.com/s/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "WeReadError400 - Invalid URL format",
		})
		return
	}
	
	// 从URL中提取一些信息作为ID
	urlParts := strings.Split(request.URL, "/")
	mpID := "mp_" + urlParts[len(urlParts)-1]
	
	// 生成模拟公众号信息
	mpInfoList := []model.MPInfo{
		{
			ID:         mpID,
			Cover:      "https://example.com/cover.jpg",
			Name:       "示例公众号 - " + mpID[0:8],
			Intro:      "这是一个示例公众号的简介",
			UpdateTime: time.Now().Unix(),
		},
	}
	
	c.JSON(http.StatusOK, mpInfoList)
}

// 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
} 