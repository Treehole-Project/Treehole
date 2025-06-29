# 🕳️ Treehole Project

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.0+-green.svg)](https://vuejs.org/)
[![Website](https://img.shields.io/badge/Website-treehole.club-orange.svg)](https://treehole.club/)

一个完全开源的匿名化交流平台，由某不知名 BIT 学生创建。树洞项目旨在为用户提供一个安全、匿名的空间来分享想法、倾听他人的声音。

🌐 **在线访问**: https://treehole.club/

## ✨ 功能特性

### 🔐 匿名交流
- 完全匿名发帖和回复
- 无需注册即可使用
- 保护用户隐私安全

### 💬 丰富的交流方式
- 支持文字、图片发布
- 多层级回复系统
- 点赞和互动功能

### 🔍 强大的搜索功能
- 关键词搜索
- 高级搜索筛选
- 实时搜索建议

### 🏷️ 分类管理
- 多种帖子标签
- 校区分组功能  
- 话题分类组织

### 📊 数据统计
- 浏览量统计
- 点赞数追踪
- 回复数显示

## 🏗️ 技术架构

### 后端 (Go)
- **Web 框架**: Gin - 高性能的 HTTP Web 框架
- **数据库 ORM**: GORM - 功能强大的 Go ORM 库
- **数据库**: SQLite/MySQL - 支持多种数据库
- **定时任务**: Cron - 自动化数据同步
- **配置管理**: Godotenv - 环境变量管理

### 前端 (Vue 3)
- **框架**: Vue 3 + Composition API
- **路由**: Vue Router 4
- **HTTP 客户端**: Axios
- **UI 框架**: Tailwind CSS
- **构建工具**: Vite

### 部署
- **容器化**: Docker + Docker Compose
- **反向代理**: 支持 Nginx 配置
- **跨平台**: 适配 Webapp

## 📁 项目结构

```
Treehole/
├── 📁 GO/                          # 后端 Go 代码
│   ├── main.go                     # 应用程序入口
│   ├── go.mod                      # Go 模块依赖
│   └── internal/                   # 内部包
│       ├── api/                    # API 路由和处理器
│       ├── config/                 # 配置管理
│       ├── database/               # 数据库连接和迁移
│       ├── models/                 # 数据模型
│       ├── scheduler/              # 定时任务调度
│       └── scraper/                # 数据爬取服务
├── 📁 WEBSITE/                     # 前端 Vue 项目
│   ├── src/                        # 源代码
│   │   ├── components/             # Vue 组件
│   │   ├── views/                  # 页面视图
│   │   ├── router/                 # 路由配置
│   │   ├── services/               # API 服务
│   │   └── utils/                  # 工具函数
│   ├── public/                     # 静态资源
│   └── package.json                # 项目依赖
├── docker-compose.yml              # Docker 编排配置
├── Dockerfile                      # Docker 镜像构建
└── README.md                       # 项目说明文档
```

## 🚀 快速开始

### 使用 Docker (推荐)

1. **克隆仓库**
   ```bash
   git clone https://github.com/yourusername/treehole.git
   cd treehole
   ```

2. **启动服务**
   ```bash
   docker-compose up -d
   ```

3. **访问应用**
   - 前端: http://localhost:8081
   - API: http://localhost:8081/api/v1

### 本地开发

#### 后端开发

1. **环境要求**
   - Go 1.23+
   - SQLite (或 MySQL)

2. **启动后端**
   ```bash
   cd GO
   go mod download
   go run main.go
   ```

#### 前端开发

1. **环境要求**
   - Node.js 16+
   - npm 或 yarn

2. **启动前端**
   ```bash
   cd WEBSITE
   npm install
   npm run dev
   ```

## 🔧 配置说明

### 环境变量

在 `GO/` 目录下创建 `.env` 文件：

```env
# 数据库配置
DATABASE_URL=./data/treehole.db

# 服务端口
PORT=8080

# 时区设置
TZ=Asia/Shanghai

# 爬虫配置
SCRAPER_INTERVAL=300  # 5分钟同步一次
```

### API 端点

| 方法 | 端点 | 描述 |
|------|------|------|
| GET | `/api/v1/posts` | 获取帖子列表 |
| GET | `/api/v1/posts/:id` | 获取单个帖子详情 |
| GET | `/api/v1/posts/:id/replies` | 获取帖子回复 |
| GET | `/api/v1/search` | 搜索帖子 |
| GET | `/api/v1/stats` | 获取统计信息 |
| GET | `/api/v1/tags` | 获取标签列表 |

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. **Fork 项目**
2. **创建特性分支** (`git checkout -b feature/AmazingFeature`)
3. **提交更改** (`git commit -m 'Add some AmazingFeature'`)
4. **推送到分支** (`git push origin feature/AmazingFeature`)
5. **打开 Pull Request**

### 开发规范

- 遵循 Go 和 Vue.js 的最佳实践
- 编写清晰的提交信息
- 添加必要的注释和文档
- 确保代码通过所有测试

## 📸 项目截图

### 主页
展示最新的匿名帖子，支持实时搜索和分类浏览。

### 帖子详情
完整的帖子阅读体验，包含多层级评论系统。

### 搜索功能
强大的搜索功能，支持多条件筛选。

## 🛡️ 隐私与安全

- **匿名性保证**: 不收集用户个人信息
- **内容自由**: 没有审查机制
- **开源透明**: 所有代码公开可审计

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。


## 📞 联系我们

- **网站**: https://treehole.club/
- **Issues**: [GitHub Issues](https://github.com/Treehole-Project/treehole/issues)
- **讨论**: [GitHub Discussions](https://github.com/Treehole-Project/treehole/discussions)

---

<div align="center">

** © 2025 Treehole Project. 完全开源的匿名交流平台 **

*Made with ❤️ by 某不知名 BIT 学生*

</div>
