#!/bin/bash

echo "=== 前端完全重建脚本 ==="
echo ""

# 获取项目根目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# 进入前端目录
cd "$PROJECT_ROOT/frontend"

echo "📦 Step 1: 本地构建前端..."
npm run build
if [ $? -ne 0 ]; then
  echo "❌ npm build 失败!"
  exit 1
fi
echo "✅ 前端构建成功"
echo ""

# 获取新的JS文件名
NEW_JS=$(ls -t dist/assets/index-*.js | head -1 | xargs basename)
echo "新JS文件: $NEW_JS"
echo ""

# 进入项目根目录
cd "$PROJECT_ROOT"

echo "🐳 Step 2: 停止并删除旧容器..."
docker-compose stop frontend
docker-compose rm -f frontend
echo "✅ 旧容器已删除"
echo ""

echo "🗑️  Step 3: 删除旧镜像..."
docker rmi splendor-frontend 2>/dev/null || echo "镜像已不存在"
echo "✅ 旧镜像已删除"
echo ""

echo "🔨 Step 4: 重新构建镜像 (--no-cache)..."
docker-compose build --no-cache frontend
if [ $? -ne 0 ]; then
  echo "❌ Docker build 失败!"
  exit 1
fi
echo "✅ 镜像构建成功"
echo ""

echo "🚀 Step 5: 启动新容器..."
docker-compose up -d frontend
if [ $? -ne 0 ]; then
  echo "❌ 容器启动失败!"
  exit 1
fi
echo "✅ 容器启动成功"
echo ""

# 等待容器完全启动
sleep 2

echo "🔍 Step 6: 验证部署..."
CONTAINER_JS=$(docker exec splendor-frontend ls /usr/share/nginx/html/assets/*.js 2>/dev/null | xargs -n1 basename)
SERVED_JS=$(curl -s http://localhost:3000/ | grep -o 'index-[^"]*\.js')
STARTED_AT=$(docker inspect splendor-frontend | grep StartedAt | cut -d'"' -f4)

echo "容器内JS文件: $CONTAINER_JS"
echo "Served JS文件: $SERVED_JS"
echo "容器启动时间: $STARTED_AT"
echo ""

if [ "$CONTAINER_JS" = "$NEW_JS" ] && [ "$SERVED_JS" = "$NEW_JS" ]; then
  echo "✅ 验证成功! 前端已更新到最新版本"
else
  echo "⚠️  警告: 文件名不匹配"
  echo "   本地: $NEW_JS"
  echo "   容器: $CONTAINER_JS"
  echo "   Served: $SERVED_JS"
fi
echo ""

echo "=== 部署完成 ==="
echo ""
echo "📝 下一步:"
echo "1. 硬刷新浏览器: Cmd + Shift + R"
echo "2. 打开Console查看是否有错误"
echo "3. 测试新功能"
echo ""
echo "访问: http://localhost:3000/game/4"
