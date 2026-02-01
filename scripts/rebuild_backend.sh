#!/bin/bash

echo "=== 后端完全重建脚本 ==="
echo ""

# 获取项目根目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

cd "$PROJECT_ROOT"

echo "🐳 Step 1: 停止并删除旧容器..."
docker-compose stop backend
docker-compose rm -f backend
echo "✅ 旧容器已删除"
echo ""

echo "🗑️  Step 2: 删除旧镜像..."
docker rmi splendor-backend 2>/dev/null || echo "镜像已不存在"
echo "✅ 旧镜像已删除"
echo ""

echo "🔨 Step 3: 重新构建镜像 (--no-cache)..."
docker-compose build --no-cache backend
if [ $? -ne 0 ]; then
  echo "❌ Docker build 失败!"
  exit 1
fi
echo "✅ 镜像构建成功"
echo ""

echo "🚀 Step 4: 启动新容器..."
docker-compose up -d backend
if [ $? -ne 0 ]; then
  echo "❌ 容器启动失败!"
  exit 1
fi
echo "✅ 容器启动成功"
echo ""

# 等待容器完全启动
sleep 2

STARTED_AT=$(docker inspect splendor-backend | grep StartedAt | cut -d'"' -f4)
echo "容器启动时间: $STARTED_AT"
echo ""

echo "🔍 测试后端API..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/games/4/state)
if [ "$HTTP_CODE" = "200" ]; then
  echo "✅ 后端API正常响应 (HTTP $HTTP_CODE)"
else
  echo "⚠️  后端API响应异常 (HTTP $HTTP_CODE)"
fi
echo ""

echo "=== 部署完成 ==="
echo ""
echo "访问: http://localhost:8080/api/v1/games/4/state"
