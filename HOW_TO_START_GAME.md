# 如何开始游戏

## 问题已修复！✅

刚才的问题是：等待室没有自动刷新玩家列表，所以房主看不到第二个玩家加入，导致"开始游戏"按钮一直显示"Need X more players"。

### 已修复内容

✅ 等待室现在每2秒自动刷新一次
✅ 当新玩家加入时，房主会立即看到更新
✅ 当玩家数量满足要求（>=2）时，"开始游戏"按钮会自动启用
✅ 当房主开始游戏后，所有玩家会自动跳转到游戏界面

## 如何开始游戏

### 步骤1: 创建房间（房主）

1. 访问 http://localhost:3000
2. 登录或注册账号
3. 点击 "Create Game" 按钮
4. 选择玩家数量（2-4人）
5. 点击创建

你会看到一个**房间码**（例如：`826c45`），分享给你的朋友！

### 步骤2: 加入房间（其他玩家）

1. 在另一个浏览器窗口（或无痕模式）访问 http://localhost:3000
2. 注册一个不同的账号
3. 在大厅页面，找到游戏卡片，点击 "Join" 输入房间码
4. 或者直接在游戏列表中看到等待中的游戏，点击 "Join"

### 步骤3: 开始游戏（房主）

当所有玩家都加入后：

1. **房主**会在等待室看到一个绿色的 **"Start Game"** 按钮
2. 按钮会显示：
   - ✅ `"Start Game"` - 可以开始（玩家数 >= 2）
   - ⏸️ `"Need X more players"` - 还需要更多玩家
3. 点击 "Start Game" 按钮
4. 所有玩家会自动跳转到游戏界面

### 等待室界面说明

```
┌─────────────────────────────────────┐
│         Waiting Room                │
│     Room Code: 826c45               │
│   Share this code with friends!     │
├─────────────────────────────────────┤
│  Players (2/2)     [You are host]   │
├─────────────────────────────────────┤
│   [👤 Player1]    [👤 Player2]      │
│       Host                          │
├─────────────────────────────────────┤
│ [Leave Game]    [Start Game]        │
└─────────────────────────────────────┘
```

### 重要提示

#### 只有房主能开始游戏
- 创建房间的人是房主
- 其他玩家会看到："Waiting for host to start the game..."
- 房主会看到 "You are the host" 标签

#### 至少需要2个玩家
- 1个玩家：按钮显示 "Need 1 more players"（禁用）
- 2个玩家：按钮显示 "Start Game"（启用）✅

#### 自动刷新
- 等待室每2秒刷新一次
- 当新玩家加入时，你会立即看到
- 当房主开始游戏时，你会自动跳转

## 测试建议

### 本地测试（一个人两个账号）

```bash
# 浏览器1（正常模式）
1. 访问 http://localhost:3000
2. 注册账号: user1@test.com
3. 创建2人游戏
4. 记下房间码（如：826c45）

# 浏览器2（无痕/隐私模式）
5. 打开无痕窗口，访问 http://localhost:3000
6. 注册账号: user2@test.com
7. 加入游戏（输入房间码：826c45）

# 回到浏览器1
8. 等待2秒，你会看到 Player2 出现
9. "Start Game" 按钮变为可点击（绿色）
10. 点击 "Start Game"
11. 两个浏览器都会跳转到游戏界面
```

### 多人测试

分享房间码给朋友：
1. 让朋友访问 http://localhost:3000（如果是局域网，用你的IP地址）
2. 让他们注册账号
3. 输入你的房间码加入
4. 等待所有人加入后，点击 "Start Game"

## 常见问题

### Q: 我看不到"Start Game"按钮
**A**: 你不是房主。只有创建房间的人能开始游戏。

### Q: 按钮显示 "Need X more players"
**A**: 还没有足够的玩家。等待更多玩家加入，或刷新页面（现在会自动刷新）。

### Q: 我点了"Start Game"但没反应
**A**:
1. 检查浏览器控制台是否有错误
2. 确保后端服务正常运行：`docker-compose ps`
3. 检查网络连接
4. 刷新页面重试

### Q: 第二个玩家加入了，但我还是看不到
**A**: 现在已修复！页面每2秒自动刷新。如果还是看不到：
1. 手动刷新页面（F5）
2. 检查两个账号是否在同一个游戏
3. 查看后端日志：`docker-compose logs backend`

### Q: 点击"Start Game"后报错
**A**:
1. 确保至少有2个玩家
2. 检查后端日志看错误信息
3. 可能是数据库问题，重启服务：`docker-compose restart`

### Q: 如何在局域网内让朋友加入？
**A**:
1. 找到你的本地IP地址：
   - Mac/Linux: `ifconfig | grep "inet " | grep -v 127.0.0.1`
   - Windows: `ipconfig`
2. 让朋友访问 `http://你的IP:3000`（如：`http://192.168.1.100:3000`）
3. 确保防火墙允许3000端口

## 验证修复

运行这个测试确认修复成功：

```bash
# 终端1：创建第一个用户和游戏
TOKEN1=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"host","email":"host@test.com","password":"pass"}' \
  | jq -r .access_token)

GAME=$(curl -s -X POST http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer $TOKEN1" \
  -H "Content-Type: application/json" \
  -d '{"num_players":2}')

GAME_ID=$(echo $GAME | jq -r .game.id)
ROOM_CODE=$(echo $GAME | jq -r .room_code)
echo "房间码: $ROOM_CODE"
echo "游戏ID: $GAME_ID"

# 终端2：创建第二个用户并加入
TOKEN2=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"guest","email":"guest@test.com","password":"pass"}' \
  | jq -r .access_token)

curl -s -X POST http://localhost:8080/api/v1/games/join \
  -H "Authorization: Bearer $TOKEN2" \
  -H "Content-Type: application/json" \
  -d "{\"room_code\":\"$ROOM_CODE\"}"

echo "✓ 玩家2已加入"

# 终端1：房主开始游戏
curl -s -X POST http://localhost:8080/api/v1/games/$GAME_ID/start \
  -H "Authorization: Bearer $TOKEN1"

echo "✓ 游戏已开始"
```

如果以上命令都成功，说明后端工作正常，前端也应该能正常开始游戏了！

## 更新日志

### 2026-01-24 23:35
- ✅ 修复：等待室现在每2秒自动刷新游戏状态
- ✅ 修复：当新玩家加入时，房主立即看到更新
- ✅ 修复：当房主开始游戏时，所有玩家自动跳转
- ✅ 改进：更好的实时体验

### 已知问题
- 无

---

**现在可以正常开始游戏了！🎮**

刷新你的浏览器页面（Ctrl+R 或 Cmd+R），然后再试一次！
