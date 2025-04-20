#!/bin/bash

# 停止并删除现有容器（如果存在）
echo "Stopping existing containers..."
docker stop wewe-rss-authserver 2>/dev/null || true
docker rm wewe-rss-authserver 2>/dev/null || true

# 构建镜像
echo "Building Docker image..."
docker build -t wewe-rss/authserver .

# 运行容器
echo "Starting container..."
docker run -d --name wewe-rss-authserver -p 8080:8080 wewe-rss/authserver

# 检查容器是否启动成功
echo "Checking if container is running..."
sleep 2
docker ps | grep wewe-rss-authserver

echo "Container logs:"
docker logs wewe-rss-authserver

echo "You can access the service at http://localhost:8080/health"
echo "Full logs: docker logs wewe-rss-authserver -f" 