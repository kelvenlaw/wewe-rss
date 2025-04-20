FROM golang:1.20-alpine AS builder

WORKDIR /app

# 设置Go Modules代理为国内镜像
ENV GOPROXY=https://goproxy.cn,direct

# 复制Go模块依赖文件
COPY go.mod go.sum* ./

# 安装依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/authserver ./cmd/server

# 使用轻量级的alpine镜像作为最终镜像
FROM alpine:latest

# 添加tzdata包以支持时区设置
RUN apk --no-cache add tzdata ca-certificates

# 设置东八区时间
ENV TZ=Asia/Shanghai

# 从构建阶段复制编译好的可执行文件
COPY --from=builder /app/authserver /usr/local/bin/authserver

# 创建工作目录
WORKDIR /app

# 应用运行时端口
EXPOSE 8080

# 确保端口绑定到0.0.0.0而不是localhost
ENV HOST=0.0.0.0
ENV PORT=8080

# 启动应用
ENTRYPOINT ["authserver"]