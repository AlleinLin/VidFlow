<div align="center">

# 🎬 VidFlow

**基于 Go 语言构建的高性能视频流媒体平台**

[!\[Go Version\](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat\&logo=go null)](https://golang.org/)
[!\[License\](https://img.shields.io/badge/License-MIT-blue.svg null)](LICENSE)
[!\[PostgreSQL\](https://img.shields.io/badge/PostgreSQL-16+-336791?style=flat\&logo=postgresql null)](https://www.postgresql.org/)
[!\[Redis\](https://img.shields.io/badge/Redis-7+-DC382D?style=flat\&logo=redis null)](https://redis.io/)

[功能特性](#-功能特性) •
[技术架构](#-技术架构) •
[快速开始](#-快速开始) •
[API 文档](#-api-文档) •
[部署指南](#-部署指南)

</div>

***

## 📖 项目简介

**VidFlow** 是一个采用 Go 语言从零构建的企业级视频流媒体平台后端服务。项目采用领域驱动设计（DDD）和清洁架构理念，具备高并发、高可用、易扩展的特点，适用于视频网站、在线教育、直播点播等多种场景。

### ✨ 核心亮点

| 特性            | 描述                            |
| ------------- | ----------------------------- |
| 🚀 **高性能**    | 基于 Go 原生并发，单机支持万级 QPS         |
| 🎥 **专业视频处理** | FFmpeg 驱动的多分辨率转码，HLS 自适应流媒体传输 |
| 🔐 **安全可靠**   | JWT 认证、RBAC 权限、内容审核机制         |
| 📊 **可观测性**   | Prometheus 指标 + Jaeger 链路追踪   |
| 🐳 **云原生**    | Docker/Kubernetes 友好，支持水平扩展   |

***

## 🎯 功能特性

### 👤 用户服务

- 🔐 用户注册/登录（密码加密存储）
- 🎫 JWT Token 认证与刷新机制
- 👥 用户资料管理与头像上传
- 🤝 关注/粉丝社交系统

### 📹 视频服务

- 📤 大文件分片上传
- 🎬 多分辨率自动转码（240p/480p/720p/1080p/4K）
- 📺 HLS 自适应流媒体分发
- 🖼️ 视频缩略图自动生成

### 💬 互动服务

- 💭 多级评论与回复系统
- ❤️ 点赞与收藏功能
- 🎯 实时弹幕系统
- 📜 观看历史记录

### 🔍 搜索与推荐

- 🔎 全文搜索（视频/用户）
- 🔥 热门视频推荐
- 🎯 个性化推荐算法
- 📺 相似视频推荐

### 💎 会员订阅

- 📦 多级会员套餐
- 🎁 权益管理系统
- 💳 支付集成框架

### 🛡️ 内容审核

- 🤖 自动内容审核
- 📋 审核规则配置
- 👨‍⚖️ 人工复审流程

***

## 🏗️ 技术架构

### 技术栈

<table>
<tr>
<td width="50%">

#### 核心框架

- **Go 1.23** - 主编程语言
- **Chi Router** - 轻量级 HTTP 路由
- **pgxpool** - PostgreSQL 连接池
- **go-redis** - Redis 客户端

</td>
<td width="50%">

#### 数据存储

- **PostgreSQL 16** - 主数据库
- **Redis 7** - 缓存与会话存储
- **ScyllaDB** - 海量数据存储（可选）

</td>
</tr>
<tr>
<td width="50%">

#### 视频处理

- **FFmpeg** - 视频转码引擎
- **HLS** - 自适应流媒体协议
- **CDN** - 内容分发网络

</td>
<td width="50%">

#### 基础设施

- **Docker** - 容器化部署
- **Kubernetes** - 容器编排
- **Prometheus** - 监控指标
- **Jaeger** - 分布式追踪
- **Kafka** - 消息队列

</td>
</tr>
</table>

### 项目结构

```
vidflow/
├── 📂 cmd/                        # 应用入口
│   └── 📂 api-gateway/            # API 网关服务
├── 📂 internal/                   # 内部代码（不可外部引用）
│   ├── 📂 config/                 # ⚙️ 配置管理
│   ├── 📂 domain/                 # 📦 领域模型
│   │   ├── 📂 user/               #    用户领域
│   │   ├── 📂 video/              #    视频领域
│   │   └── 📂 interaction/        #    互动领域
│   ├── 📂 handler/                # 🌐 HTTP 处理器
│   │   └── 📂 middleware/         #    中间件
│   ├── 📂 service/                # 💼 业务逻辑层
│   │   ├── 📂 user/               #    用户服务
│   │   ├── 📂 video/              #    视频服务
│   │   ├── 📂 interaction/        #    互动服务
│   │   ├── 📂 playback/           #    播放服务
│   │   ├── 📂 search/             #    搜索服务
│   │   ├── 📂 recommendation/     #    推荐服务
│   │   ├── 📂 transcode/          #    转码服务
│   │   ├── 📂 cdn/                #    CDN 服务
│   │   ├── 📂 subscription/       #    订阅服务
│   │   ├── 📂 audit/              #    审核服务
│   │   ├── 📂 notification/       #    通知服务
│   │   └── 📂 payment/            #    支付服务
│   ├── 📂 repository/             # 🗄️ 数据访问层
│   │   └── 📂 postgres/           #    PostgreSQL 实现
│   └── 📂 infrastructure/         # 🏗️ 基础设施
│       ├── 📂 cache/              #    Redis 缓存
│       ├── 📂 messaging/          #    Kafka 消息
│       ├── 📂 metrics/            #    Prometheus 指标
│       └── 📂 tracing/            #    Jaeger 追踪
├── 📂 pkg/                        # 公共包（可外部引用）
│   ├── 📂 jwt/                    # 🔑 JWT 工具
│   ├── 📂 hash/                   # 🔒 密码哈希
│   ├── 📂 logger/                 # 📝 日志工具
│   ├── 📂 errors/                 # ❌ 错误处理
│   └── 📂 response/               # 📤 响应格式
├── 📂 migrations/                 # 🗃️ 数据库迁移
└── 📂 deployments/                # 🚀 部署配置
    ├── 📂 docker/                 #    Docker 配置
    └── 📂 kubernetes/             #    K8s 配置
```

### 架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client (Web/App)                         │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                        API Gateway (Chi)                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │  Auth   │ │  Rate   │ │  CORS   │ │ Metrics │ │ Tracing │   │
│  │Middleware│ │  Limit  │ │         │ │         │ │         │   │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘   │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
    ┌─────────────────────────────┼─────────────────────────────┐
    │                             │                             │
    ▼                             ▼                             ▼
┌──────────┐               ┌──────────┐               ┌──────────┐
│   User   │               │  Video   │               │Interaction│
│ Service  │               │ Service  │               │ Service  │
└────┬─────┘               └────┬─────┘               └────┬─────┘
     │                          │                          │
     └──────────────────────────┼──────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Infrastructure Layer                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │PostgreSQL│  │  Redis   │  │  Kafka   │  │   CDN    │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
```

***

## 🚀 快速开始

### 环境要求

| 依赖         | 版本    | 说明        |
| ---------- | ----- | --------- |
| Go         | 1.23+ | 编程语言      |
| PostgreSQL | 16+   | 主数据库      |
| Redis      | 7+    | 缓存服务      |
| FFmpeg     | 6+    | 视频处理      |
| Docker     | 24+   | 容器化部署（可选） |

### 本地开发

```bash
# 1️⃣ 克隆项目
git clone 
cd vidflow

# 2️⃣ 安装依赖
go mod download

# 3️⃣ 创建数据库
createdb vidflow

# 4️⃣ 运行数据库迁移
psql -U postgres -d vidflow -f migrations/001_init_schema.up.sql

# 5️⃣ 配置环境变量
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=vidflow
export REDIS_ADDR=localhost:6379
export JWT_SECRET=your-secret-key-at-least-32-characters

# 6️⃣ 启动服务
go run ./cmd/api-gateway
```

服务启动后访问：

- 🌐 API: <http://localhost:8080>
- 📊 Metrics: <http://localhost:8080/metrics>

### Docker 部署

```bash
# 使用 Docker Compose 一键启动
docker-compose -f deployments/docker/docker-compose.yml up -d

# 查看服务状态
docker-compose -f deployments/docker/docker-compose.yml ps

# 查看日志
docker-compose -f deployments/docker/docker-compose.yml logs -f api-gateway
```

***

## 📚 API 文档

### 认证接口

| 方法     | 路径                      | 描述       |
| ------ | ----------------------- | -------- |
| `POST` | `/api/v1/auth/register` | 用户注册     |
| `POST` | `/api/v1/auth/login`    | 用户登录     |
| `POST` | `/api/v1/auth/refresh`  | 刷新 Token |

### 用户接口

| 方法       | 路径                            | 描述     |   认证   |
| -------- | ----------------------------- | ------ | :----: |
| `GET`    | `/api/v1/users/me`            | 获取当前用户 |    ✅   |
| `PUT`    | `/api/v1/users/me`            | 更新用户资料 |    ✅   |
| `DELETE` | `/api/v1/users/me`            | 删除用户   |    ✅   |
| `GET`    | `/api/v1/users/:username`     | 获取用户信息 | <br /> |
| `POST`   | `/api/v1/users/:id/follow`    | 关注用户   |    ✅   |
| `DELETE` | `/api/v1/users/:id/follow`    | 取消关注   |    ✅   |
| `GET`    | `/api/v1/users/:id/followers` | 获取粉丝列表 | <br /> |
| `GET`    | `/api/v1/users/:id/following` | 获取关注列表 | <br /> |

### 视频接口

| 方法       | 路径                           | 描述     |   认证   |
| -------- | ---------------------------- | ------ | :----: |
| `POST`   | `/api/v1/videos`             | 上传视频   |    ✅   |
| `GET`    | `/api/v1/videos`             | 获取视频列表 | <br /> |
| `GET`    | `/api/v1/videos/:id`         | 获取视频详情 | <br /> |
| `PUT`    | `/api/v1/videos/:id`         | 更新视频信息 |    ✅   |
| `DELETE` | `/api/v1/videos/:id`         | 删除视频   |    ✅   |
| `GET`    | `/api/v1/videos/:id/stream`  | 获取视频流  | <br /> |
| `POST`   | `/api/v1/videos/:id/publish` | 发布视频   |    ✅   |

### 互动接口

| 方法       | 路径                            | 描述     |   认证   |
| -------- | ----------------------------- | ------ | :----: |
| `POST`   | `/api/v1/comments`            | 发表评论   |    ✅   |
| `GET`    | `/api/v1/videos/:id/comments` | 获取评论列表 | <br /> |
| `DELETE` | `/api/v1/comments/:id`        | 删除评论   |    ✅   |
| `POST`   | `/api/v1/videos/:id/like`     | 点赞视频   |    ✅   |
| `DELETE` | `/api/v1/videos/:id/like`     | 取消点赞   |    ✅   |
| `POST`   | `/api/v1/videos/:id/favorite` | 收藏视频   |    ✅   |
| `DELETE` | `/api/v1/videos/:id/favorite` | 取消收藏   |    ✅   |
| `POST`   | `/api/v1/danmakus`            | 发送弹幕   |    ✅   |
| `GET`    | `/api/v1/videos/:id/danmakus` | 获取弹幕列表 | <br /> |

### 搜索接口

| 方法    | 路径                           | 描述   |
| ----- | ---------------------------- | ---- |
| `GET` | `/api/v1/search`             | 综合搜索 |
| `GET` | `/api/v1/search/videos`      | 搜索视频 |
| `GET` | `/api/v1/search/users`       | 搜索用户 |
| `GET` | `/api/v1/search/suggestions` | 搜索建议 |

***

## ⚙️ 配置说明

### 环境变量

| 变量名              | 说明            | 默认值              |
| ---------------- | ------------- | ---------------- |
| `SERVER_HOST`    | 服务地址          | `0.0.0.0`        |
| `SERVER_PORT`    | 服务端口          | `8080`           |
| `DB_HOST`        | 数据库地址         | `localhost`      |
| `DB_PORT`        | 数据库端口         | `5432`           |
| `DB_USER`        | 数据库用户         | `postgres`       |
| `DB_PASSWORD`    | 数据库密码         | -                |
| `DB_NAME`        | 数据库名称         | `vidflow`        |
| `REDIS_ADDR`     | Redis 地址      | `localhost:6379` |
| `REDIS_PASSWORD` | Redis 密码      | -                |
| `JWT_SECRET`     | JWT 密钥（≥32字符） | -                |
| `KAFKA_BROKERS`  | Kafka 集群地址    | `localhost:9092` |

### 配置文件

支持 YAML 格式配置文件，详见 [deployments/docker/config.yaml](deployments/docker/config.yaml)

***

## 📈 性能特性

| 特性           | 说明                           |
| ------------ | ---------------------------- |
| 🚀 **高并发处理** | 基于 Goroutine 的轻量级并发，单机万级 QPS |
| 🔗 **连接池优化** | 数据库和 Redis 连接池复用，减少连接开销      |
| ⚡ **异步处理**   | 视频转码、通知推送等耗时任务异步执行           |
| 💾 **多级缓存**  | Redis 缓存 + 本地缓存，减少数据库压力      |
| 📺 **流式传输**  | HLS 自适应码率，流畅播放体验             |

***

## 📊 监控与可观测性

### Prometheus 指标

服务暴露 `/metrics` 端点，提供以下指标：

```
# HTTP 请求
http_requests_total{method, path, status}
http_request_duration_seconds{method, path}

# 数据库
database_connections{state}
database_query_duration_seconds{query}

# 缓存
cache_hits_total
cache_misses_total

# 视频
video_uploads_total
video_transcode_duration_seconds{resolution}
video_views_total{video_id}

# 业务
user_registrations_total
user_logins_total
```

### Jaeger 链路追踪

集成 OpenTelemetry，支持分布式链路追踪，快速定位性能瓶颈。

***

## 🧪 测试

```bash
# 运行所有测试
go test ./... -v

# 运行带覆盖率的测试
go test ./... -cover

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

***

## 🗺️ 路线图

- [x] 用户认证与授权
- [x] 视频上传与转码
- [x] HLS 流媒体分发
- [x] 评论与互动系统
- [x] 搜索与推荐
- [ ] 直播功能
- [ ] AI 内容审核
- [ ] GraphQL API
- [ ] 移动端 SDK

***

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

***

## 📄 许可证

本项目采用 [MIT License](LICENSE) 许可证。

***

<div align="center">

[⬆ 返回顶部](#-vidflow)

</div>
