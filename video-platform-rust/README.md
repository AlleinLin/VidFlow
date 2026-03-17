<div align="center">

# 🦀 VidFlow Rust

**基于 Actix-web 构建的视频流媒体平台**

[!\[Rust Version\](https://img.shields.io/badge/Rust-1.78+-DEA584?style=flat\&logo=rust null)](https://www.rust-lang.org/)
[!\[Actix-web\](https://img.shields.io/badge/Actix--web-4.5-000000?style=flat null)](https://actix.rs/)
[!\[License\](https://img.shields.io/badge/License-MIT-blue.svg null)](LICENSE)

[功能特性](#-功能特性) •
[技术架构](#-技术架构) •
[快速开始](#-快速开始) •
[API 文档](#-api-文档)

</div>

***

## 📖 项目简介

**VidFlow Rust** 是 VidFlow 视频流媒体平台的 Rust 实现版本，采用 Actix-web 4 + Tokio 构建，实现极致性能和内存安全，适用于对性能要求极高的视频服务场景。

### ✨ 核心亮点

| 特性            | 描述                                   |
| ------------- | ------------------------------------ |
| ⚡ **极致性能**    | Actix-web 性能霸榜 TechEmpower           |
| 🛡️ **内存安全**  | Rust 编译期保证，无 GC 停顿                   |
| 🔥 **零成本抽象**  | 高级抽象不牺牲运行效率                          |
| 🧵 **异步运行时**  | Tokio 支持百万级并发连接                      |
| 📊 **完整可观测性** | Tracing + Prometheus + OpenTelemetry |

***

## 🎯 功能特性

### 👤 用户服务

- 🔐 用户注册/登录（Argon2 密码加密）
- 🎫 JWT Token 认证
- 👥 用户资料管理
- 🤝 关注/粉丝社交系统

### 📹 视频服务

- 📤 视频上传
- 🎬 多分辨率转码调度
- 📺 HLS 流媒体分发
- 🗂️ 视频分类管理

### 💬 互动服务

- 💭 多级评论系统
- ❤️ 点赞与收藏
- 🎯 实时弹幕
- 📜 观看历史

### 🔍 搜索与推荐

- 🔎 全文搜索
- 🔥 热门推荐
- 🎯 相关视频推荐

***

## 🏗️ 技术架构

### 技术栈

<table>
<tr>
<td width="50%">

#### 核心框架

- **Actix-web 4** - 高性能 Web 框架
- **Tokio** - 异步运行时
- **SQLx** - 编译期 SQL 验证
- **Serde** - 序列化框架

</td>
<td width="50%">

#### 数据存储

- **PostgreSQL** - 主数据库
- **Redis** - 缓存层
- **ScyllaDB** - 弹幕存储
- **S3** - 对象存储

</td>
</tr>
<tr>
<td width="50%">

#### 安全认证

- **JWT** - jsonwebtoken
- **Argon2** - 密码哈希
- **CORS** - 跨域配置

</td>
<td width="50%">

#### 基础设施

- **Kafka (rdkafka)** - 消息队列
- **Tracing** - 日志追踪
- **Prometheus** - 指标监控
- **OpenTelemetry** - 分布式追踪

</td>
</tr>
</table>

### 项目结构

```
vidflow-rust/
├── 📄 Cargo.toml               # 📦 项目依赖
├── 📄 .env.example             # 🔑 环境变量模板
├── 📄 Dockerfile               # 🐳 容器构建
├── 📂 migrations/              # 🗃️ 数据库迁移
└── 📂 src/
    ├── 📄 main.rs              # 🚀 应用入口
    ├── 📄 config.rs            # ⚙️ 配置管理
    ├── 📄 db.rs                # 🗄️ 数据库连接
    ├── 📄 redis.rs             # 📮 Redis 客户端
    ├── 📄 error.rs             # ❌ 错误处理
    ├── 📄 utils.rs             # 🔧 工具函数
    ├── 📂 models/              # 📦 数据模型
    │   ├── mod.rs
    │   ├── user.rs             #    用户模型
    │   ├── video.rs            #    视频模型
    │   └── playback.rs         #    播放模型
    ├── 📂 handlers/            # 🌐 HTTP 处理器
    │   ├── mod.rs
    │   ├── auth.rs             #    认证接口
    │   ├── users.rs            #    用户接口
    │   ├── videos.rs           #    视频接口
    │   ├── comments.rs         #    评论接口
    │   ├── interactions.rs     #    互动接口
    │   ├── danmakus.rs         #    弹幕接口
    │   ├── recommend.rs        #    推荐接口
    │   ├── search.rs           #    搜索接口
    │   ├── playback.rs         #    播放接口
    │   ├── health.rs           #    健康检查
    │   └── metrics.rs          #    指标暴露
    ├── 📂 middleware/          # 🔒 中间件
    │   ├── mod.rs
    │   ├── auth.rs             #    JWT 认证
    │   └── metrics.rs          #    指标采集
    ├── 📂 services/            # 💼 业务服务
    │   ├── mod.rs
    │   ├── auth.rs             #    认证服务
    │   └── video.rs            #    视频服务
    └── 📂 utils/               # 🔧 工具函数
        └── helpers.rs
```

***

## 🚀 快速开始

### 环境要求

| 依赖         | 版本    | 说明   |
| ---------- | ----- | ---- |
| Rust       | 1.78+ | 编程语言 |
| PostgreSQL | 16+   | 主数据库 |
| Redis      | 7+    | 缓存服务 |

### 本地开发

```bash
# 1️⃣ 克隆项目
git clone
cd vidflow-rust

# 2️⃣ 复制环境变量
cp .env.example .env

# 3️⃣ 编辑配置
# 修改 .env 文件中的数据库连接信息

# 4️⃣ 构建项目
cargo build --release

# 5️⃣ 运行服务
cargo run --release
```

服务启动后访问：

- 🌐 API: <http://localhost:8082>
- 📊 Metrics: <http://localhost:8082/metrics>
- ❤️ Health: <http://localhost:8082/health>

### Docker 部署

```bash
# 构建镜像
docker build -t vidflow-rust:latest .

# 运行容器
docker run -d -p 8082:8082 \
  -e DATABASE__URL=postgresql://user:pass@host:5432/vidflow \
  -e REDIS__URL=redis://host:6379 \
  -e JWT__SECRET=your-secret-key \
  vidflow-rust:latest
```

***

## 📚 API 文档

### 认证接口 `/api/v1/auth`

| 方法     | 路径          | 描述   |  认证 |
| ------ | ----------- | ---- | :-: |
| `POST` | `/register` | 用户注册 |  ❌  |
| `POST` | `/login`    | 用户登录 |  ❌  |
| `POST` | `/refresh`  | 刷新令牌 |  ❌  |
| `POST` | `/logout`   | 用户登出 |  ✅  |

### 用户接口 `/api/v1/users`

| 方法     | 路径                | 描述     |   认证   |
| ------ | ----------------- | ------ | :----: |
| `GET`  | `/{id}`           | 获取用户资料 | <br /> |
| `PUT`  | `/{id}`           | 更新用户资料 |    ✅   |
| `POST` | `/{id}/follow`    | 关注用户   |    ✅   |
| `POST` | `/{id}/unfollow`  | 取消关注   |    ✅   |
| `GET`  | `/{id}/followers` | 获取粉丝列表 | <br /> |
| `GET`  | `/{id}/following` | 获取关注列表 | <br /> |

### 视频接口 `/api/v1/videos`

| 方法       | 路径              | 描述     |   认证   |
| -------- | --------------- | ------ | :----: |
| `GET`    | `/`             | 视频列表   | <br /> |
| `POST`   | `/`             | 创建视频   |    ✅   |
| `GET`    | `/{id}`         | 获取视频详情 | <br /> |
| `PUT`    | `/{id}`         | 更新视频   |    ✅   |
| `DELETE` | `/{id}`         | 删除视频   |    ✅   |
| `POST`   | `/{id}/publish` | 发布视频   |    ✅   |

### 互动接口 `/api/v1/interactions`

| 方法     | 路径                       | 描述     |  认证 |
| ------ | ------------------------ | ------ | :-: |
| `POST` | `/like/{video_id}`       | 点赞视频   |  ✅  |
| `POST` | `/unlike/{video_id}`     | 取消点赞   |  ✅  |
| `POST` | `/favorite/{video_id}`   | 收藏视频   |  ✅  |
| `POST` | `/unfavorite/{video_id}` | 取消收藏   |  ✅  |
| `GET`  | `/status/{video_id}`     | 获取互动状态 |  ✅  |

### 推荐接口 `/api/v1/recommend`

| 方法    | 路径                    | 描述   |   认证   |
| ----- | --------------------- | ---- | :----: |
| `GET` | `/home`               | 首页推荐 | <br /> |
| `GET` | `/related/{video_id}` | 相关推荐 | <br /> |
| `GET` | `/search`             | 搜索视频 | <br /> |

***

## ⚙️ 配置说明

### 环境变量

| 变量                       | 描述             | 默认值       |
| ------------------------ | -------------- | --------- |
| `DATABASE__URL`          | PostgreSQL 连接串 | -         |
| `REDIS__URL`             | Redis 连接串      | -         |
| `SERVER__HOST`           | 服务监听地址         | `0.0.0.0` |
| `SERVER__PORT`           | 服务监听端口         | `8082`    |
| `SERVER__WORKERS`        | 工作线程数          | `4`       |
| `JWT__SECRET`            | JWT 密钥         | -         |
| `JWT__ACCESS_TOKEN_TTL`  | 访问令牌有效期(秒)     | `900`     |
| `JWT__REFRESH_TOKEN_TTL` | 刷新令牌有效期(秒)     | `604800`  |

***

## 📈 性能特性

| 特性            | 说明                         |
| ------------- | -------------------------- |
| ⚡ **极致性能**    | Actix-web 性能霸榜 TechEmpower |
| 🧵 **异步 I/O** | Tokio 支持百万级并发              |
| 💾 **零拷贝**    | Serde 零拷贝序列化               |
| 🔒 **内存安全**   | 编译期保证，无数据竞争                |
| 🎯 **智能背压**   | 自动流量控制                     |

***

## 🧪 测试

```bash
# 运行所有测试
cargo test

# 运行特定测试
cargo test test_user_registration

# 测试覆盖率
cargo tarpaulin --out Html
```

***

## 📊 监控指标

访问 `/metrics` 端点获取 Prometheus 指标：

```
# HTTP 请求
http_requests_total{method, path, status}
http_request_duration_seconds{method, path}

# 数据库
db_connections_active
db_queries_total

# 缓存
cache_hits_total
cache_misses_total
```

***

## 📄 许可证

本项目采用 [MIT License](LICENSE) 许可证。

***

<div align="center">

[⬆ 返回顶部](#-vidflow-rust)

</div>
