# 重启服务指南

## 快速重启命令

### 方式1: 只重启前端（最快）
```bash
cd /Users/shanks/go/src/splendor
docker-compose restart frontend
```
**用途**: 当只修改了前端代码时
**时间**: ~3秒

### 方式2: 只重启后端
```bash
docker-compose restart backend
```
**用途**: 当只修改了后端代码时
**时间**: ~3秒

### 方式3: 重启所有服务（推荐）
```bash
docker-compose restart
```
**用途**: 重启前端、后端、数据库
**时间**: ~10秒

### 方式4: 完全重新构建和启动
```bash
docker-compose down
docker-compose up -d --build
```
**用途**: 当有依赖变化或需要完全重置时
**时间**: ~1-2分钟

## 详细说明

### 1️⃣ 重启前端
```bash
docker-compose restart frontend
```

**何时使用**:
- 修改了 `frontend/src` 下的任何文件
- 前端出现问题需要重启
- 应用新的前端代码

**验证**:
```bash
curl http://localhost:3000
```

### 2️⃣ 重启后端
```bash
docker-compose restart backend
```

**何时使用**:
- 修改了 `backend` 下的 Go 代码
- 后端出现问题需要重启
- 环境变量有变化

**验证**:
```bash
curl http://localhost:8080/health
# 应返回: {"message":"Splendor API is running","status":"ok"}
```

### 3️⃣ 重启数据库
```bash
docker-compose restart postgres
```

**何时使用**:
- 数据库连接出现问题
- 很少需要，除非数据库卡住

**注意**: 重启数据库不会丢失数据（数据存储在 volume 中）

### 4️⃣ 重启所有服务
```bash
docker-compose restart
```

**等同于**:
```bash
docker-compose restart postgres
docker-compose restart backend
docker-compose restart frontend
```

**何时使用**:
- 不确定问题在哪个服务
- 想要全部刷新
- 最安全的重启方式

### 5️⃣ 停止所有服务
```bash
docker-compose down
```

**作用**:
- 停止所有容器
- 删除容器（但保留数据）
- 删除网络

**验证停止成功**:
```bash
docker-compose ps
# 应该显示空列表
```

### 6️⃣ 启动所有服务
```bash
docker-compose up -d
```

**参数说明**:
- `-d`: 后台运行（detached mode）
- 不加 `-d` 会在前台显示所有日志

**验证启动成功**:
```bash
docker-compose ps
# 所有服务状态应该是 Up
```

### 7️⃣ 重新构建并启动
```bash
docker-compose up -d --build
```

**作用**:
- 重新构建 Docker 镜像
- 应用代码更改
- 重新安装依赖

**何时使用**:
- 修改了 `Dockerfile`
- 修改了 `go.mod` 或 `package.json`
- 代码有重大更新
- **刚才修改前端代码就是用的这个**

### 8️⃣ 完全重置（包括数据）
```bash
docker-compose down -v
docker-compose up -d --build
```

**警告**: ⚠️ 这会**删除所有数据**！

**作用**:
- 停止并删除所有容器
- **删除所有数据卷**（包括数据库数据）
- 重新构建
- 从头开始

**何时使用**:
- 数据库损坏
- 想要完全重新开始
- 测试初始化流程

## 常用组合命令

### 应用前端代码更改
```bash
cd /Users/shanks/go/src/splendor/frontend
npm run build
cd ..
docker-compose restart frontend
```

### 应用后端代码更改（如果直接修改容器内代码）
```bash
docker-compose restart backend
```

### 应用后端代码更改（推荐 - 重新构建）
```bash
docker-compose up -d --build backend
```

### 查看实时日志
```bash
# 所有服务
docker-compose logs -f

# 只看后端
docker-compose logs -f backend

# 只看前端
docker-compose logs -f frontend

# 只看最近50行
docker-compose logs --tail=50 backend
```

### 停止查看日志
按 `Ctrl + C`

## 检查服务状态

### 查看所有容器
```bash
docker-compose ps
```

### 查看后端日志
```bash
docker-compose logs backend --tail=20
```

### 查看前端日志
```bash
docker-compose logs frontend --tail=20
```

### 测试后端健康
```bash
curl http://localhost:8080/health
```

### 测试前端
```bash
curl http://localhost:3000
```

## 刚才的修复流程回顾

我刚才做的事情：

1. **修改了前端代码** (`frontend/src/pages/LobbyPage.tsx`)
2. **重新构建前端**:
   ```bash
   cd frontend && npm run build
   ```
3. **重新构建并重启前端容器**:
   ```bash
   docker-compose up -d --build frontend
   ```

所以现在前端已经是最新的了！你只需要：

### ✅ 刷新浏览器页面

按 `F5` 或 `Ctrl+R` (Mac: `Cmd+R`) 刷新页面即可！

不需要重启任何东西，因为我已经帮你重启了。

## 快速参考

| 操作 | 命令 | 时间 |
|------|------|------|
| 刷新浏览器 | `F5` 或 `Ctrl+R` | 即时 |
| 重启前端 | `docker-compose restart frontend` | 3秒 |
| 重启后端 | `docker-compose restart backend` | 3秒 |
| 重启所有 | `docker-compose restart` | 10秒 |
| 重新构建 | `docker-compose up -d --build` | 1-2分钟 |
| 查看日志 | `docker-compose logs -f` | - |
| 查看状态 | `docker-compose ps` | 即时 |
| 停止服务 | `docker-compose down` | 5秒 |
| 启动服务 | `docker-compose up -d` | 20秒 |
| 完全重置 | `docker-compose down -v && docker-compose up -d --build` | 2分钟 |

## 故障排除

### 如果重启后还是有问题

1. **查看日志找错误**:
   ```bash
   docker-compose logs backend
   docker-compose logs frontend
   ```

2. **完全重新构建**:
   ```bash
   docker-compose down
   docker-compose up -d --build
   ```

3. **清除浏览器缓存**:
   - Chrome: `Ctrl+Shift+Delete`
   - 选择"缓存图像和文件"
   - 点击"清除数据"

4. **使用无痕模式测试**:
   - Chrome: `Ctrl+Shift+N`
   - Firefox: `Ctrl+Shift+P`

### 如果端口被占用

```bash
# 查看谁在使用端口
lsof -i :3000  # 前端
lsof -i :8080  # 后端

# 强制停止
docker-compose down
```

### 如果容器无法启动

```bash
# 查看详细错误
docker-compose up

# 不加 -d，可以看到完整日志
```

## 总结

**现在你只需要做一件事**:

### 🔄 刷新浏览器页面！

按 `F5` 或 `Ctrl+R`，前端已经更新了，新代码已经生效！

---

**快速命令速查**:
```bash
# 最常用的3个命令
docker-compose ps              # 查看状态
docker-compose logs -f         # 查看日志
docker-compose restart         # 重启所有服务
```
