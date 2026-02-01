# Splendor Production Deployment Guide for AWS EC2

## 前置要求

在你的AWS EC2实例上需要安装：
- Docker
- Docker Compose
- Git（可选，用于拉取代码）

## 部署步骤

### 1. 准备EC2实例

确保EC2安全组开放以下端口：
- **80** (HTTP)
- **443** (HTTPS，如果使用SSL)
- **22** (SSH，用于管理)

```bash
# 检查Docker是否运行
sudo systemctl status docker

# 如果Docker未启动
sudo systemctl start docker
sudo systemctl enable docker
```

### 2. 上传代码到EC2

**方式A：使用Git（推荐）**
```bash
# 在EC2上
cd /home/ec2-user
git clone YOUR_REPO_URL splendor
cd splendor
```

**方式B：使用SCP上传**
```bash
# 在本地机器上
scp -r /Users/shanks/go/src/splendor ec2-user@YOUR_EC2_IP:/home/ec2-user/
```

### 3. 配置环境变量

```bash
cd /home/ec2-user/splendor

# 复制并编辑环境变量文件
cp .env.production .env.production.local
nano .env.production.local
```

**重要：必须修改以下配置：**

```bash
# 强密码
POSTGRES_PASSWORD=your_strong_password_here

# JWT密钥（至少32位随机字符串）
JWT_SECRET=your_very_long_random_secret_key_minimum_32_chars

# 你的域名或EC2公网IP
FRONTEND_URL=http://your-domain.com
# 或使用IP：FRONTEND_URL=http://YOUR_EC2_PUBLIC_IP

VITE_API_URL=http://your-domain.com/api
VITE_WS_URL=ws://your-domain.com/api
```

生成安全的JWT密钥：
```bash
openssl rand -base64 32
```

### 4. 部署应用

```bash
# 重命名配置文件
mv .env.production.local .env.production

# 运行部署脚本
./deploy.sh
```

### 5. 验证部署

```bash
# 查看所有容器状态
docker ps

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f

# 检查特定服务
docker-compose -f docker-compose.prod.yml logs backend
docker-compose -f docker-compose.prod.yml logs frontend
```

访问你的应用：
- **前端**: `http://YOUR_EC2_PUBLIC_IP`
- **API**: `http://YOUR_EC2_PUBLIC_IP/api`
- **WebSocket**: `ws://YOUR_EC2_PUBLIC_IP/ws`

### 6. 设置域名（可选但推荐）

如果你有域名：

1. 在域名提供商处添加A记录指向你的EC2公网IP
2. 更新`.env.production`中的URL配置
3. 重新部署

## 配置HTTPS（强烈推荐）

### 使用Let's Encrypt免费SSL证书

```bash
# 安装certbot
sudo yum install -y certbot python3-certbot-nginx  # Amazon Linux
# 或
sudo apt-get install -y certbot python3-certbot-nginx  # Ubuntu

# 创建SSL目录
mkdir -p nginx/ssl

# 获取证书
sudo certbot certonly --standalone -d your-domain.com

# 复制证书到项目
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem nginx/ssl/
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem nginx/ssl/
sudo chown ec2-user:ec2-user nginx/ssl/*

# 编辑nginx配置文件，取消HTTPS部分的注释
nano nginx/nginx.conf

# 重启nginx
docker-compose -f docker-compose.prod.yml restart nginx
```

## 日常运维命令

### 查看日志
```bash
# 所有服务
docker-compose -f docker-compose.prod.yml logs -f

# 特定服务
docker-compose -f docker-compose.prod.yml logs -f backend
docker-compose -f docker-compose.prod.yml logs -f frontend
docker-compose -f docker-compose.prod.yml logs -f postgres

# 最近100行
docker-compose -f docker-compose.prod.yml logs --tail=100
```

### 重启服务
```bash
# 重启所有服务
docker-compose -f docker-compose.prod.yml restart

# 重启特定服务
docker-compose -f docker-compose.prod.yml restart backend
docker-compose -f docker-compose.prod.yml restart frontend
```

### 停止服务
```bash
docker-compose -f docker-compose.prod.yml down
```

### 更新代码并重新部署
```bash
# 如果使用Git
git pull
./deploy.sh

# 或者手动
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --build
```

### 数据库备份
```bash
# 备份数据库
docker exec splendor-db pg_dump -U splendor splendor > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据库
docker exec -i splendor-db psql -U splendor splendor < backup_file.sql
```

### 清理Docker资源
```bash
# 清理未使用的镜像
docker image prune -f

# 清理所有未使用的资源
docker system prune -a
```

## 监控和健康检查

### 检查服务健康
```bash
# Backend健康检查
curl http://localhost/health

# 查看容器状态
docker ps

# 查看资源使用
docker stats
```

### 设置自动重启
容器已配置为`restart: always`，系统重启后会自动启动。

## 性能优化建议

1. **启用Gzip压缩**: 已在nginx配置中启用
2. **配置CDN**: 可以使用CloudFront
3. **数据库优化**:
   ```bash
   # 进入数据库容器
   docker exec -it splendor-db psql -U splendor

   # 创建索引等优化
   ```

## 安全建议

1. ✅ 使用强密码
2. ✅ 启用HTTPS
3. ✅ 定期更新Docker镜像
4. ✅ 配置防火墙规则
5. ✅ 定期备份数据库
6. ✅ 监控日志查看异常活动

## 故障排查

### 服务无法启动
```bash
# 查看详细日志
docker-compose -f docker-compose.prod.yml logs

# 检查端口占用
sudo netstat -tulpn | grep :80
sudo netstat -tulpn | grep :443
```

### 数据库连接失败
```bash
# 检查数据库容器
docker logs splendor-db

# 测试数据库连接
docker exec -it splendor-db psql -U splendor -d splendor
```

### WebSocket连接失败
- 检查nginx配置中的WebSocket设置
- 确保EC2安全组允许WebSocket连接
- 查看浏览器控制台错误

## 成本优化

1. 选择合适的EC2实例类型（t3.small或t3.medium适合小型应用）
2. 使用Elastic IP避免IP变更
3. 考虑使用RDS代替容器数据库（生产环境推荐）
4. 设置CloudWatch告警监控资源使用

## 获取帮助

如有问题，检查：
1. Docker日志: `docker-compose -f docker-compose.prod.yml logs`
2. 系统日志: `sudo journalctl -u docker`
3. Nginx日志: `docker exec splendor-nginx cat /var/log/nginx/error.log`
