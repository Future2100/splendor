# Splendor UI 问题检查清单

## 每次修改后必须验证的问题（按优先级排序）

### 🔴 Critical Issues（必须修复，否则不可用）

- [ ] **1. 按钮位置固定在卡片底部**
  - Buy/Reserve按钮必须在所有卡片的同一高度
  - 不能因为cost项目数量变化而移动位置
  - 验证：检查1个cost和4个cost的卡片，按钮应在相同位置

- [ ] **2. 所有cost项目完全可见**
  - 4个cost项目的卡片，最后一个（Ruby/Onyx等）必须完全显示
  - 不能被截断，不能需要滚动才能看到
  - 验证：找到有4个cost的卡片，检查第4个项目是否完整显示

- [ ] **3. Buy/Reserve按钮完整显示**
  - 按钮不能被卡片边缘截断
  - 按钮文字和边框必须完整可见
  - 验证：所有卡片的按钮都应该完整显示在卡片内

### 🟡 High Priority（影响美观和一致性）

- [ ] **4. 卡片高度完全一致**
  - 所有卡片（1个cost ~ 4个cost）高度必须一致
  - 验证：目测所有卡片顶部和底部对齐

- [ ] **5. Gem icon使用emoji（和Gem Bank一致）**
  - 使用：💎🔷💚❤️⚫
  - 不使用CSS旋转方块
  - 验证：左上角gem icon和cost列表中的icon应该是emoji

- [ ] **6. Gem icon背景干净**
  - 左上角gem icon不应该有背景色、边框、或badge
  - 只显示emoji本身（无"T1"文字）
  - 验证：左上角应该只有💎，没有背景，没有"T1"文字

- [ ] **7. 卡片背景色匹配gem类型**
  - Diamond卡 → 蓝色系
  - Sapphire卡 → 深蓝色系
  - Emerald卡 → 绿色系
  - Ruby卡 → 红色系
  - Onyx卡 → 灰黑色系
  - 验证：卡片边框颜色和背景颜色应该对应gem类型

### 🟢 Medium Priority（用户体验优化）

- [ ] **8. 去掉"COST"标题文字**
  - Cost section不显示"Cost"标题
  - 直接显示cost列表
  - 验证：卡片上不应该有"COST"/"Cost"文字

- [ ] **9. 没有多余的gem显示**
  - 只在左上角显示一个gem emoji
  - 卡片中央区域只有cost列表
  - 验证：整个卡片应该只有一个gem icon（左上角）

- [ ] **10. 布局紧凑**
  - 一屏能看到更多内容
  - 间距、padding、margin都要小
  - 验证：一屏应该能看到所有3个tier

- [ ] **11. 左侧sidebar高度匹配右侧tiers**
  - Gem Bank + Nobles总高度 ≈ 3个tier总高度
  - Nobles应该填充剩余空间
  - 验证：左右两侧顶部和底部应该大致对齐

## 验证流程（每次部署后执行）

### Step 1: 硬刷新浏览器
```bash
Mac: Cmd + Shift + R
Windows/Linux: Ctrl + Shift + R
```

### Step 2: 视觉检查
1. 打开 http://localhost:3000/game/4
2. 逐项检查上述清单
3. 截图记录任何问题

### Step 3: 测试检查
```bash
npm test -- --run
```
- 所有测试必须通过
- 特别注意卡片高度测试

### Step 4: 不同卡片检查
- [ ] 检查1个cost的卡片
- [ ] 检查2个cost的卡片
- [ ] 检查3个cost的卡片
- [ ] 检查4个cost的卡片（重点！）

## 常见错误及原因

### 按钮被截断或位置移动
**原因**: Cost section使用flex-1或自适应高度
**解决**: Cost section必须用固定高度（如100px、110px等）

### Ruby/Onyx等第4个cost项目看不见
**原因**: Cost section高度不够，或overflow设置错误
**解决**:
- Cost section高度必须 ≥ 100px
- 必须有 overflow-y-auto
- 必须有 h-full

### 按钮位置不固定
**原因**: 使用justify-center或flex-1在cost section
**解决**: Cost section用flex-none + 固定高度，按钮用flex-none

### Gem icon有背景方块
**原因**: 使用CSS div元素创建icon
**解决**: 直接用<span>包裹emoji，不加任何背景

## 当前正确的布局结构（200px总高度）

```
卡片结构：
├─ Top badges (absolute positioned)
│  ├─ 左上: 💎 (只有gem emoji, 无背景，无文字)
│  └─ 右上: [5] (Victory points badge)
│
├─ Card Content (p-2 pt-7)
│  ├─ Cost section (flex-none h-[130px])
│  │  └─ Cost items (overflow-y-auto)
│  │     └─ 4个cost项目都要完全可见
│  │
│  └─ Buttons (flex-none mt-1)
│     ├─ Buy
│     └─ Reserve
```

## 空间计算（200px卡片）

- Gem icon: absolute position top-1 left-1/2 (不占用layout空间)
  - 位置: 顶部4px处
  - 大小: text-2xl (~24px高)
  - z-index: 20（显示在最上层）
- 顶部padding: 28px (pt-7)
- Cost section: 130px（确保4个cost项目可见）
  - 4个cost项: 4 × 24px = 96px
  - 3个间隙: 3 × 2px = 6px
  - 容器padding: ~12px
  - 小计: ~114px（130px留buffer）
- 按钮区域: ~30px (py-1.5 + margin)
- 底部padding: 8px (p-2)
- 间距: 4px (mt-1)

**总计**: 28 + 130 + 30 + 8 + 4 = 200px
**Gem icon**: absolute定位，不占用flex空间，显示在蓝色背景区域顶部中间

## 部署后验证命令

```bash
# 检查部署的JS文件
docker exec splendor-frontend ls -lh /usr/share/nginx/html/assets/ | grep "index-.*\.js"

# 检查是否有emoji
docker exec splendor-frontend grep -o '💎' /usr/share/nginx/html/assets/index-*.js | head -1

# 检查cost section高度
docker exec splendor-frontend grep -o 'h-\[100px\]' /usr/share/nginx/html/assets/index-*.js | head -1

# 检查卡片高度
docker exec splendor-frontend grep -o 'minHeight:"200px"' /usr/share/nginx/html/assets/index-*.js | head -1
```

## 本次部署检查结果（2026-01-25 - Latest）

⚠️ **待用户验证** - 需要硬刷新浏览器（Cmd+Shift+R / Ctrl+Shift+R）

🔧 **最新修改**:
1. 修复NobleDisplay组件的gem icon从CSS方块改为emoji
2. 修复PlayerPanels中gem总数计算逻辑
3. 修复数据库中game 4的gold数量（7→5）

### 代码级验证（已完成）
✅ **1. 按钮位置固定** - Cost section h-[130px] flex-none
✅ **2. 所有cost项目可见** - 130px高度（计算：4×24px + 3×2px + 12px = 114px，留16px buffer）
✅ **3. 按钮完整显示** - 固定在底部，mt-1间距
✅ **4. 卡片高度一致** - minHeight/maxHeight: 200px
✅ **5. Gem icon用emoji** - 💎🔷💚❤️⚫ (包括Nobles)
✅ **6. Gem icon无背景** - 使用<span>包裹emoji，无div背景
✅ **7. 卡片背景色匹配gem** - Diamond蓝色、Ruby红色、Emerald绿色等
✅ **8. 无"COST"标题** - Cost section直接显示列表
✅ **9. 无中间gem方块** - 只在tier badge显示小emoji
✅ **10. 布局紧凑** - 200px卡片高度
✅ **11. Sidebar高度匹配** - flex布局，Nobles填充剩余空间
✅ **12. Gem计算正确** - PlayerPanels使用正确的reduce逻辑
✅ **13. Gold初始化正确** - 2人游戏gold永远5个（已修复旧游戏bug）
✅ **14. Nobles用emoji icon** - NobleDisplay的GemIcon改为emoji

**部署版本**: index-Dle57CiV.js

### ⚠️ 关键验证项（必须用户目测确认）
- [ ] **Ruby/Onyx完全可见** - 在有4个cost的卡片上，第4个cost项目（Ruby/Onyx）必须完整显示，不能被截断
- [ ] **按钮在卡片底部** - 所有卡片的Buy/Reserve按钮都在相同高度
- [ ] **左上角gem icon清晰** - 左上角只有💎（或其他gem emoji），没有背景，没有"T1"文字

## 修改历史

- 2026-01-25 17:45: **添加当前用户显示** - 在LobbyPage和GamePage顶部显示"Logged in as/Playing as [username]"，让用户知道当前登录身份
- 2026-01-25 17:35: **添加永久显示的游戏规则** - 创建独立的GameRules组件在页面底部永久显示，包含完整的游戏规则（拿宝石、购买卡片、保留卡片、胜利条件）；从ActionPanel移除重复的规则文本
- 2026-01-25 17:25: **改进游戏卡片按钮** - GameCard根据用户是否是玩家显示不同按钮：玩家看到"Continue Game"/"Return to Lobby"，非玩家看到"View Game"；已完成游戏显示"View Results"
- 2026-01-25 17:15: **添加分享功能** - WaitingRoom添加一键复制房间码和加入链接按钮；LobbyPage支持?join=房间码自动加入
- 2026-01-25 17:00: **实现游戏结束逻辑** - 后端：达到15分时自动结束游戏，计算获胜者；前端：添加游戏结束界面显示最终排名
- 2026-01-25 16:45: **修复Nobles gem icon** - NobleDisplay的GemIcon从CSS旋转方块改为emoji（💎🔷💚❤️⚫）
- 2026-01-25 16:40: **修复gold数量** - 修复数据库中game 4的gold从7改为5（2人游戏应该是5个gold）
- 2026-01-25 16:35: **修复gem计算** - PlayerPanels中totalGems计算从`(a||0)+(b||0)`改为正确的`sum+(count||0)`
- 2026-01-25 15:35: **缩小Nobles** - text-lg→text-sm, p-2.5→p-1.5, text-xs→text-[10px], space-y-2→space-y-1.5，让所有nobles显示完全
- 2026-01-25 15:30: **调整左右栏宽度** - 左边lg:col-span-2→lg:col-span-3 (25%)，右边lg:col-span-10→lg:col-span-9 (75%)，卡片变窄，左边有更多空间
- 2026-01-25 15:25: **缩小Gem Bank** - 图标text-2xl→text-xl，padding p-2→p-1.5，gap-2→gap-1.5，标题text-base→text-sm，给Nobles更多空间
- 2026-01-25 15:20: **正确对齐** - 左边栏整体固定高度664px（3×216px + 2×8px），Nobles用flex-1填充剩余空间
- 2026-01-25 15:15: ❌ 错误尝试 - Nobles设置maxHeight: 480px，没有正确对齐
- 2026-01-25 15:10: **布局优化** - 去掉"TIER 1/2/3"标题，deck card缩小到120px固定宽度
- 2026-01-25 15:05: gem icon移到top-0（最顶端），不遮挡白色区域
- 2026-01-25 15:00: **正确修复** - gem icon移到top-1（往上移emoji），恢复pt-7和h-[130px]（Ruby可见）
- 2026-01-25 14:50: ❌ 错误方案 - pt-12给gem icon留空间，但导致cost section变为h-[114px]，Ruby被遮挡
- 2026-01-25 14:40: 恢复pt-7和h-[130px]，gem icon位置top-3，添加drop-shadow增强可见性（但gem icon被白色区域覆盖）
- 2026-01-25 14:30: 提高gem icon的z-index到z-20，确保完全显示在蓝色背景区域，不被白色cost区域覆盖
- 2026-01-25 14:25: 将gem icon移到卡片顶部中间（left-1/2 -translate-x-1/2水平居中）
- 2026-01-25 14:20: 调整gem icon高度到top-1.5，和victory points badge对齐
- 2026-01-25 14:15: 调整gem icon位置从top-2 left-2到top-3 left-3（离边缘更远）
- 2026-01-25 14:10: 去掉tier badge背景，只保留gem emoji图标（无"T1"文字）
- 2026-01-25 14:05: 增加cost section到130px，等待用户视觉验证
- 2026-01-25 13:32: 修复所有问题并验证通过（120px版本）
- 2026-01-25 13:30: 创建checklist

## 测试方法说明

### 我的测试限制
我只能验证**部署的代码文件**中是否包含正确的代码模式（如`h-[130px]`、`minHeight:"200px"`等），但**无法看到实际的视觉结果**。

### 你需要做的测试（每次部署后）
1. **硬刷新浏览器**（清除缓存）
   - Mac: `Cmd + Shift + R`
   - Windows/Linux: `Ctrl + Shift + R`

2. **视觉检查（最重要！）**
   - 找一个有4个cost的卡片（例如：Diamond 1, Sapphire 1, Emerald 1, Ruby 3）
   - 检查第4个cost项目（Ruby）是否**完整显示**（不被截断）
   - 检查Buy/Reserve按钮是否在**卡片底部固定位置**（不是中间）
   - 检查Tier badge（💎 T1）中的emoji是否**没有背景方块**

3. **如果有问题**
   - 截图给我看
   - 告诉我具体哪个项目有问题
