package model

// 账号状态常量
const (
	StatusInvalid = 0
	StatusEnable  = 1
	StatusDisable = 2
)

// LoginResult 登录结果
type LoginResult struct {
	Message  string `json:"message,omitempty"`
	Vid      int    `json:"vid,omitempty"`
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
}

// LoginURL 登录扫码URL
type LoginURL struct {
	UUID    string `json:"uuid"`
	ScanURL string `json:"scanUrl"`
}

// Article 公众号文章
type Article struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	PicURL      string `json:"picUrl"`
	PublishTime int64  `json:"publishTime"`
}

// MPInfo 公众号信息
type MPInfo struct {
	ID         string `json:"id"`
	Cover      string `json:"cover"`
	Name       string `json:"name"`
	Intro      string `json:"intro"`
	UpdateTime int64  `json:"updateTime"`
}

// MPInfoRequest 获取公众号信息的请求
type MPInfoRequest struct {
	URL string `json:"url" binding:"required"`
} 