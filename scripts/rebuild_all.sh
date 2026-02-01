#!/bin/bash

echo "=== 前后端完全重建脚本 ==="
echo ""

# 获取项目根目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

cd "$PROJECT_ROOT"

echo "📦 Step 1: 本地构建前端..."
cd frontend
npm run build
if [ $? -ne 0 ]; then
  echo "❌ 前端构建失败!"
  exit 1
fi
cd ..
echo "✅ 前端构建成功"
echo ""

echo "🐳 Step 2: 停止所有容器..."
docker-compose down
echo "✅ 所有容器已停止"
echo ""

echo "🗑️  Step 3: 删除旧镜像..."
docker rmi splendor-frontend splendor-backend 2>/dev/null || echo "镜像已不存在"
echo "✅ 旧镜像已删除"
echo ""

echo "🔨 Step 4: 重新构建所有镜像..."
docker-compose build --no-cache
if [ $? -ne 0 ]; then
  echo "❌ 构建失败!"
  exit 1
fi
echo "✅ 镜像构建成功"
echo ""

echo "🚀 Step 5: 启动所有容器..."
docker-compose up -d
if [ $? -ne 0 ]; then
  echo "❌ 容器启动失败!"
  exit 1
fi
echo "✅ 容器启动成功"
echo ""

sleep 3

echo "🔍 Step 6: 验证部署..."
docker-compose ps

echo ""
echo "=== 部署完成 ==="
echo ""
echo "前端: http://localhost:3000/game/4"
echo "后端: http://localhost:8080/api/v1/games/4/state"
echo ""
echo "记得硬刷新浏览器: Cmd + Shift + R"
