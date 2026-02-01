# Splendor 部署脚本

这个目录包含了快速部署和重建前后端的自动化脚本。

## 脚本列表

### 1. `rebuild_frontend.sh` - 前端完全重建 ⭐ 推荐

**用途**: 修改前端代码后，完全重建并部署

**步骤**:
- 本地构建前端 (`npm run build`)
- 停止并删除旧容器
- 删除旧Docker镜像
- 使用 `--no-cache` 重新构建镜像
- 启动新容器
- 验证部署

**使用场景**:
- 添加/修改前端组件
- 修改React代码
- 修改CSS样式
- **重要**: 当 `restart_frontend.sh` 无效时使用

**运行**:
```bash
./scripts/rebuild_frontend.sh
```

### 2. `rebuild_backend.sh` - 后端完全重建

**用途**: 修改后端代码后，完全重建并部署

**步骤**:
- 停止并删除旧容器
- 删除旧Docker镜像
- 使用 `--no-cache` 重新构建镜像
- 启动新容器
- 测试API响应

**使用场景**:
- 修改Go代码
- 修改API逻辑
- 修改数据库查询

**运行**:
```bash
./scripts/rebuild_backend.sh
```

### 3. `restart_frontend.sh` - 前端快速重启

**用途**: 快速重启前端容器（不完全重建）

**步骤**:
- 本地构建前端
- 重启容器

**使用场景**:
- 小改动
- 快速测试

**注意**:
- ⚠️  可能不会更新文件（Docker缓存问题）
- 如果没有变化，请使用 `rebuild_frontend.sh`

**运行**:
```bash
./scripts/restart_frontend.sh
```

### 4. `rebuild_all.sh` - 前后端完全重建

**用途**: 同时重建前端和后端

**步骤**:
- 本地构建前端
- 停止所有容器
- 删除所有镜像
- 重新构建所有镜像
- 启动所有容器

**使用场景**:
- 同时修改了前后端
- 重大更新
- 清理环境

**运行**:
```bash
./scripts/rebuild_all.sh
```

## 使用流程

### 场景1: 修改前端代码

```bash
# 1. 修改代码
vim frontend/src/components/game/GameBoard.tsx

# 2. 重建前端
./scripts/rebuild_frontend.sh

# 3. 硬刷新浏览器
# Mac: Cmd + Shift + R
# Windows: Ctrl + Shift + R
```

### 场景2: 修改后端代码

```bash
# 1. 修改代码
vim backend/internal/gamelogic/validator.go

# 2. 重建后端
./scripts/rebuild_backend.sh

# 3. 测试API
curl http://localhost:8080/api/v1/games/4/state
```

### 场景3: 同时修改前后端

```bash
# 1. 修改代码
vim frontend/src/...
vim backend/internal/...

# 2. 重建全部
./scripts/rebuild_all.sh

# 3. 测试
```

## 为什么需要 `--no-cache`？

Docker有时会使用缓存的旧文件，导致代码更新后容器内文件没有变化。使用 `--no-cache` 强制从头重建，确保使用最新代码。

## 验证部署

### 前端验证

1. **检查容器内文件**:
```bash
docker exec splendor-frontend ls -lh /usr/share/nginx/html/assets/
```

2. **检查served文件**:
```bash
curl -s http://localhost:3000/ | grep -o 'index-[^"]*\.js'
```

3. **浏览器验证**:
- 打开 http://localhost:3000/game/4
- 打开Console (Cmd+Option+J)
- 硬刷新 (Cmd+Shift+R)
- 查看是否有错误

### 后端验证

1. **检查API响应**:
```bash
curl -s http://localhost:8080/api/v1/games/4/state | jq .
```

2. **检查日志**:
```bash
docker-compose logs backend | tail -20
```

## 常见问题

### Q: 浏览器没有变化？
**A**: 硬刷新浏览器 (Cmd+Shift+R) 或清空缓存

### Q: `restart_frontend.sh` 没有效果？
**A**: 使用 `rebuild_frontend.sh` 完全重建

### Q: 脚本执行失败？
**A**: 检查错误信息，可能需要：
- 检查Docker是否运行
- 检查端口是否被占用
- 检查是否有语法错误

### Q: 如何查看容器启动时间？
**A**:
```bash
docker inspect splendor-frontend | grep StartedAt
docker inspect splendor-backend | grep StartedAt
```

## 脚本特点

✅ 自动化部署流程
✅ 包含验证步骤
✅ 清晰的错误提示
✅ 相对路径，可以在任何位置运行
✅ 完整的状态反馈

## 项目目录结构

```
splendor/
├── scripts/              # 部署脚本
│   ├── rebuild_frontend.sh
│   ├── rebuild_backend.sh
│   ├── restart_frontend.sh
│   ├── rebuild_all.sh
│   └── README.md
├── frontend/            # React前端
├── backend/             # Go后端
└── docker-compose.yml   # Docker配置
```
