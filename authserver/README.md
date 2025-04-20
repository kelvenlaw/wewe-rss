# WeWe-RSS AuthServer

WeWe-RSS 项目的微信读书API转发服务。

## 功能

该服务模拟了微信读书API的基本功能，提供以下接口：

- 登录API
  - 创建登录URL：`GET /api/v2/login/platform`
  - 获取登录结果：`GET /api/v2/login/platform/:id`

- 公众号内容API
  - 获取公众号文章列表：`GET /api/v2/platform/mps/:mpId/articles`
  - 获取公众号信息：`POST /api/v2/platform/wxs2mp`

## 环境要求

- Go 1.20+

## 开发

1. 复制环境变量示例文件并根据需要修改：

```bash
cp .env.example .env
```

2. 运行开发服务：

```bash
go run cmd/server/main.go
```

## 构建

```bash
go build -o authserver ./cmd/server
```

## Docker部署

构建Docker镜像：

```bash
docker build -t wewe-rss/authserver .
```

运行Docker容器：

```bash
docker run -p 8080:8080 wewe-rss/authserver
```

## 配置

通过环境变量或`.env`文件配置服务参数：

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| PORT | 服务器端口 | 8080 |
| HOST | 服务器主机 | 0.0.0.0 |
| DEBUG | 调试模式 | false |
| DEFAULT_ARTICLES | 每页返回的默认文章数量 | 20 |
| WEREAD_PROXY_URL | 实际微信读书API代理URL（可选） | - |

## 接口详情

### 登录API

#### 创建登录URL

```
GET /api/v2/login/platform
```

内部实现流程：
1. 调用微信读书API `https://weread.qq.com/api/auth/getLoginUid` 获取UID
2. 使用UID构建确认URL `https://weread.qq.com/web/confirm?uid={uuid}`
3. 生成包含确认URL的二维码，并转换为Base64格式

返回示例：

```json
{
  "uuid": "e7637172-a2bd-45da-99af-ec1b3331d84b",
  "scanUrl": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgA..."
}
```

注意：
- `uuid` 是从微信读书API获取的UID，直接传递，不做修改
- `scanUrl` 是Base64编码的二维码图片数据，可直接用于前端显示
- 前端可以使用 `<img src="这里是base64数据">` 显示二维码

#### 获取登录结果

```
GET /api/v2/login/platform/:id
```

内部实现流程：
1. 调用微信读书API `https://weread.qq.com/api/auth/getLoginInfo?uid={uid}` 检查登录状态
2. 将登录结果返回给客户端

返回示例（成功）：

```json
{
  "vid": 12345,
  "token": "sample_token_login_1627123456",
  "username": "WeRead_User_login_16"
}
```

返回示例（失败）：

```json
{
  "message": "登录会话不存在或已过期"
}
```

### 公众号内容API

#### 获取公众号文章列表

```
GET /api/v2/platform/mps/:mpId/articles?page=1
```

请求头：

```
Authorization: Bearer {token}
xid: {vid}
```

返回示例：

```json
[
  {
    "id": "article_mp_12345_1_0",
    "title": "文章标题 0 - 页码 1",
    "picUrl": "https://example.com/pic.jpg",
    "publishTime": 1627123456
  },
  ...
]
```

#### 获取公众号信息

```
POST /api/v2/platform/wxs2mp
```

请求头：

```
Authorization: Bearer {token}
xid: {vid}
```

请求体：

```json
{
  "url": "https://mp.weixin.qq.com/s/abcdefg"
}
```

返回示例：

```json
[
  {
    "id": "mp_abcdefg",
    "cover": "https://example.com/cover.jpg",
    "name": "示例公众号 - mp_abcde",
    "intro": "这是一个示例公众号的简介",
    "updateTime": 1627123456
  }
]
```

## 扩展

1. 实现实际的微信读书API调用
2. 添加用户认证和权限控制
3. 实现文章内容缓存功能 