#!/bin/bash

echo "=== 前端快速重启（仅重启容器）==="
echo ""

# 获取项目根目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

echo "📦 本地构建..."
cd "$PROJECT_ROOT/frontend"
npm run build
if [ $? -ne 0 ]; then
  echo "❌ 构建失败!"
  exit 1
fi

cd "$PROJECT_ROOT"

echo "🔄 重启容器..."
docker-compose restart frontend

sleep 2

STARTED_AT=$(docker inspect splendor-frontend | grep StartedAt | cut -d'"' -f4)
echo "✅ 容器重启成功: $STARTED_AT"
echo ""
echo "⚠️  注意: restart可能不会更新文件，如果没变化请用 rebuild_frontend.sh"
echo ""
echo "硬刷新浏览器: Cmd + Shift + R"
