# 树洞镜像站后端

这是一个使用 Go 语言开发的树洞镜像站后端项目，提供 Web API 接口和数据同步功能。

## 功能特性

- 🌐 **RESTful API 接口** - 完整的帖子和回复查询接口
- 📊 **完整数据存储** - 支持帖子和回复的所有字段（IP、openid、点赞数、图片等）
- 🕷️ **智能数据同步** - 支持增量同步和新回复检测
- ⏰ **定时任务调度** - 自动定时同步数据
- 🔍 **搜索功能** - 支持关键词搜索帖子
- 🏷️ **标签系统** - 分类管理帖子
- 📈 **统计信息** - 同步状态和数据统计
- 💾 **多数据库支持** - 支持 SQLite 和 MySQL
- 🔄 **递归评论处理** - 完整处理嵌套评论结构
- 📱 **完整字段映射** - 所有原始 API 字段都能正确存储

## 技术栈

- **Web 框架**: Gin
- **数据库 ORM**: GORM
- **爬虫**: Colly
- **定时任务**: Cron
- **配置管理**: Godotenv

## 项目结构

```
.
├── main.go                 # 程序入口
├── go.mod                  # Go 模块依赖
├── .env                    # 环境变量配置
├── internal/
│   ├── api/                # API 路由和处理器
│   │   └── handlers.go
│   ├── config/             # 配置管理
│   │   └── config.go
│   ├── database/           # 数据库连接
│   │   └── database.go
│   ├── models/             # 数据模型
│   │   └── models.go
│   ├── scraper/            # 爬虫服务
│   │   └── scraper.go
│   └── scheduler/          # 定时任务
│       └── scheduler.go
└── .github/
    └── copilot-instructions.md
```

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- Git

### 安装依赖

```bash
go mod tidy
```

### 配置环境变量

创建 `.env` 文件：

```env
# 数据库配置
DATABASE_URL=data.db

# 源站点配置
SOURCE_URL=https://example-tree-hole.com

# 爬虫配置
SCRAPE_INTERVAL=30m
MAX_RETRIES=3
REQUEST_TIMEOUT=30s
USER_AGENT=TreeHoleMirror/1.0
RATE_LIMIT_DELAY=1s

# 定时任务配置
SYNC_CRON=0 */30 * * * *

# 服务器配置
PORT=8080
```

### 运行项目

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

## API 接口

### 帖子相关

- `GET /api/v1/posts` - 获取帖子列表
- `GET /api/v1/posts/:id` - 获取单个帖子
- `GET /api/v1/posts/:id/replies` - 获取帖子回复

### 搜索

- `GET /api/v1/search?q=关键词` - 基础搜索帖子（搜索标题和内容）
- `GET /api/v1/search/advanced` - 高级搜索（支持多字段和逻辑关系）
- `GET /api/v1/search/users?q=用户名` - 搜索用户
- `GET /api/v1/search/comments?q=关键词` - 搜索评论

#### 高级搜索参数

- `title` - 搜索标题
- `content` - 搜索内容
- `author` - 搜索作者用户名
- `author_id` - 搜索作者 ID (openid)
- `post_id` - 搜索帖子 ID
- `original_id` - 搜索原始 ID
- `comment` - 搜索评论内容（返回包含该评论的帖子）
- `tag` - 搜索标签
- `state` - 搜索状态 (normal, deleted, complaint, chosen, hot)
- `radio_group` - 搜索分组
- `logic` - 逻辑关系：`and`（与）或 `or`（或），默认为 `and`

**示例：**

```
# 搜索标题包含"树洞"且作者为"张三"的帖子
GET /api/v1/search/advanced?title=树洞&author=张三&logic=and

# 搜索标题包含"树洞"或内容包含"表白"的帖子
GET /api/v1/search/advanced?title=树洞&content=表白&logic=or

# 通过评论内容搜索帖子
GET /api/v1/search/advanced?comment=好棒
```

### 用户相关

- `GET /api/v1/users/:user_id/posts` - 获取指定用户的帖子
- `GET /api/v1/users/:user_id/replies` - 获取指定用户的回复

_注：`user_id` 可以是用户名(author)或用户 ID(author_id)_

### 标签

- `GET /api/v1/tags` - 获取所有标签
- `GET /api/v1/tags/:name/posts` - 根据标签获取帖子

### 统计

- `GET /api/v1/stats` - 获取统计信息

### 同步

- `POST /api/v1/sync?source_url=URL` - 手动触发同步
- `GET /api/v1/sync/status` - 获取同步状态

### 健康检查

- `GET /health` - 健康检查

## 数据模型

### Post (帖子)

包含完整的帖子信息，字段如下：

- **基础信息**: ID、原始 ID、标题、内容、作者、作者 ID(openid)
- **统计信息**: 点赞数(likeNum)、回复数、浏览数、评论数
- **时间信息**: 创建时间、更新时间
- **分类信息**: 分组(radioGroup)、校区分组(campusGroup)、地区(region)、标签(tag)
- **扩展信息**: 价格、微信号、图片(images)、封面(cover)
- **状态信息**: 统一状态(state) - normal, deleted, complaint, chosen, hot
- **网络信息**: IP 地址
- **关联**: 回复关联

### Reply (回复)

包含完整的回复信息，字段如下：

- **基础信息**: ID、帖子 ID、原始 ID、内容、作者、作者 ID(openid)
- **回复关系**: 回复对象(applyTo)、层级(level)、父评论 ID(parent_id)
- **互动信息**: 点赞数(like_num)、图片(images)
- **分类信息**: 标签(tag)
- **时间信息**: 创建时间、更新时间
- **关联**: 与帖子的外键关联

### SyncStatus (同步状态)

- **同步信息**: 同步时间、最后帖子 ID、状态
- **统计信息**: 总帖子数、总回复数
- **错误信息**: 错误消息（如有）

## 开发说明

### 自定义爬虫规则

修改 `internal/scraper/scraper.go` 中的 `parsePost` 方法，根据目标网站的 HTML 结构调整选择器：

```go
// 示例：修改帖子解析规则
s.collector.OnHTML(".your-post-selector", func(e *colly.HTMLElement) {
    post := s.parsePost(e)
    if post != nil {
        posts = append(posts, *post)
    }
})
```

### 添加新的 API 端点

在 `internal/api/handlers.go` 中添加新的处理函数，并在 `SetupRouter` 中注册路由。

### 数据库迁移

当修改数据模型后，GORM 会自动处理数据库迁移。如需手动控制，可以修改 `internal/database/database.go` 中的 `Migrate` 函数。

## 部署

### 构建

```bash
go build -o tree-hole-mirror main.go
```

### Docker 部署

创建 `Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

