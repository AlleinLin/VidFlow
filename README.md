<div align="center">

# 🎬 VidFlow

**新一代高性能视频流媒体平台 - 多语言实现**

[!\[Go\](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat\&logo=go null)](https://golang.org/)
[!\[Java\](https://img.shields.io/badge/Java-21+-ED8B00?style=flat\&logo=openjdk null)](https://openjdk.org/)
[!\[Rust\](https://img.shields.io/badge/Rust-1.78+-DEA584?style=flat\&logo=rust null)](https://www.rust-lang.org/)
[!\[License\](https://img.shields.io/badge/License-MIT-blue.svg null)](LICENSE)

**企业级视频流媒体平台，提供 Go、Java、Rust 三种语言实现**

[项目简介](#-项目简介) •
[功能特性](#-功能特性) •
[技术架构](#-技术架构) •
[快速开始](#-快速开始) •
[部署指南](#-部署指南)

</div>

***

## 📖 项目简介

**VidFlow** 是一个采用现代技术栈从零构建的企业级视频流媒体平台，提供 **Go**、**Java**、**Rust** 三种语言的独立实现。项目采用领域驱动设计（DDD）和清洁架构理念，具备高并发、高可用、易扩展的特点，适用于视频网站、在线教育、直播点播等多种场景。

### 🎯 为什么选择 VidFlow？

| 特性            | 说明                            |
| ------------- | ----------------------------- |
| 🌍 **多语言实现**  | Go、Java、Rust 三种语言，满足不同团队技术栈需求 |
| 🏗️ **清洁架构**  | 领域驱动设计，业务逻辑与技术实现分离            |
| 🚀 **高性能**    | 支持万级 QPS，毫秒级响应                |
| 🔐 **安全可靠**   | JWT 认证、内容审核、数据加密              |
| 📊 **完整可观测性** | Prometheus 指标 + Jaeger 链路追踪   |
| 🐳 **云原生**    | Docker/Kubernetes 友好，支持水平扩展   |

***

## 🏗️ 项目结构

```
vidflow/
├── 📂 video-platform-go/      # 🐹 Go 语言实现 (端口: 8080)
│   ├── cmd/                   #    应用入口
│   ├── internal/              #    内部代码
│   ├── pkg/                   #    公共包
│   └── migrations/            #    数据库迁移
│
├── 📂 video-platform-java/    # ☕ Java 语言实现 (端口: 8081)
│   ├── src/                   #    源代码
│   ├── pom.xml                #    Maven 配置
│   └── Dockerfile             #    容器构建
│
├── 📂 video-platform-rust/    # 🦀 Rust 语言实现 (端口: 8082)
│   ├── src/                   #    源代码
│   ├── Cargo.toml             #    Cargo 配置
│   └── Dockerfile             #    容器构建
│
├── 📂 deploy/                 # 🚀 部署配置
│   ├── k8s/                   #    Kubernetes 配置
│   └── prometheus.yml         #    Prometheus 配置
│
├── 📄 docker-compose.yml      # 🐳 开发环境编排
└── 📄 README.md               # 📝 项目说明
```

***

## 🎯 功能特性

### 核心功能模块

<table>
<tr>
<td width="33%" valign="top">

#### 👤 用户服务

- 🔐 用户注册/登录
- 🎫 JWT Token 认证
- 👥 用户资料管理
- 🤝 关注/粉丝系统
- 📧 邮箱验证

</td>
<td width="33%" valign="top">

#### 📹 视频服务

- 📤 大文件分片上传
- 🎬 多分辨率转码
- 📺 HLS 自适应流媒体
- 🖼️ 缩略图自动生成
- 🗂️ 分类标签管理

</td>
<td width="33%" valign="top">

#### 💬 互动服务

- 💭 多级评论系统
- ❤️ 点赞与收藏
- 🎯 实时弹幕
- 📜 观看历史
- 📤 分享功能

</td>
</tr>
<tr>
<td width="33%" valign="top">

#### 🔍 搜索与推荐

- 🔎 全文搜索
- 🔥 热门推荐
- 🎯 个性化推荐
- 📺 相关视频推荐
- 💡 搜索建议

</td>
<td width="33%" valign="top">

#### 💎 会员订阅

- 📦 多级会员套餐
- 🎁 权益管理系统
- 💳 支付集成
- 📊 消费记录
- 🔄 自动续费

</td>
<td width="33%" valign="top">

#### 🛡️ 内容审核

- 🤖 自动内容审核
- 📋 审核规则配置
- 👨‍⚖️ 人工复审流程
- 📝 审核记录
- ⚠️ 违规处理

</td>
</tr>
</table>

***

## 🛠️ 技术架构

### 三种语言实现对比

| 特性        | Go 版本      | Java 版本                  | Rust 版本     |
| --------- | ---------- | ------------------------ | ----------- |
| **框架**    | Chi Router | Spring Boot 3.3          | Actix-web 4 |
| **运行时**   | Goroutine  | Virtual Threads (JDK 21) | Tokio       |
| **ORM**   | pgxpool    | Spring Data JPA          | SQLx        |
| **性能**    | ⭐⭐⭐⭐⭐      | ⭐⭐⭐⭐                     | ⭐⭐⭐⭐⭐       |
| **开发效率**  | ⭐⭐⭐⭐⭐      | ⭐⭐⭐⭐⭐                    | ⭐⭐⭐⭐        |
| **生态成熟度** | ⭐⭐⭐⭐       | ⭐⭐⭐⭐⭐                    | ⭐⭐⭐         |
| **内存安全**  | ⭐⭐⭐⭐       | ⭐⭐⭐⭐                     | ⭐⭐⭐⭐⭐       |

### 技术栈详情

<table>
<tr>
<th>类别</th>
<th>技术选型</th>
</tr>
<tr>
<td><strong>数据库</strong></td>
<td>PostgreSQL 16 / Redis 7 / ScyllaDB</td>
</tr>
<tr>
<td><strong>消息队列</strong></td>
<td>Kafka / RabbitMQ</td>
</tr>
<tr>
<td><strong>对象存储</strong></td>
<td>MinIO / AWS S3 / 阿里云 OSS</td>
</tr>
<tr>
<td><strong>视频处理</strong></td>
<td>FFmpeg / HLS</td>
</tr>
<tr>
<td><strong>监控</strong></td>
<td>Prometheus / Grafana</td>
</tr>
<tr>
<td><strong>追踪</strong></td>
<td>Jaeger / OpenTelemetry</td>
</tr>
<tr>
<td><strong>容器化</strong></td>
<td>Docker / Kubernetes</td>
</tr>
</table>

### 系统架构图

```
┌─────────────────────────────────────────────────────────────────────────┐
│                            Load Balancer (Nginx)                         │
└─────────────────────────────────────┬───────────────────────────────────┘
                                      │
          ┌───────────────────────────┼───────────────────────────┐
          │                           │                           │
          ▼                           ▼                           ▼
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│   VidFlow Go    │         │  VidFlow Java   │         │  VidFlow Rust   │
│   (Port 8080)   │         │   (Port 8081)   │         │   (Port 8082)   │
│                 │         │                 │         │                 │
│  • Chi Router   │         │  • Spring Boot  │         │  • Actix-web    │
│  • pgxpool      │         │  • JPA/Hibernate│         │  • SQLx         │
│  • go-redis     │         │  • Spring Data  │         │  • redis-rs     │
└────────┬────────┘         └────────┬────────┘         └────────┬────────┘
         │                           │                           │
         └───────────────────────────┼───────────────────────────┘
                                     │
         ┌───────────────────────────┼───────────────────────────┐
         │                           │                           │
         ▼                           ▼                           ▼
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│   PostgreSQL    │         │     Redis       │         │     Kafka       │
│   (主数据库)     │         │   (缓存层)      │         │   (消息队列)    │
└─────────────────┘         └─────────────────┘         └─────────────────┘
         │
         ▼
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│     MinIO       │         │    Jaeger       │         │   Prometheus    │
│   (对象存储)     │         │   (链路追踪)    │         │    (监控)       │
└─────────────────┘         └─────────────────┘         └─────────────────┘
```

***

## 🚀 快速开始

### 环境要求

| 依赖             | 版本    | 说明        |
| -------------- | ----- | --------- |
| Docker         | 24+   | 容器运行时     |
| Docker Compose | 2.20+ | 容器编排      |
| Go             | 1.23+ | Go 版本开发   |
| JDK            | 21+   | Java 版本开发 |
| Rust           | 1.78+ | Rust 版本开发 |

### 一键启动开发环境

```bash
# 1️⃣ 克隆项目
git clone
cd vidflow

# 2️⃣ 启动基础设施服务
docker-compose up -d postgres redis kafka minio jaeger prometheus grafana

# 3️⃣ 初始化数据库
psql -U postgres -d vidflow -f video-platform-go/migrations/001_init_schema.up.sql
```

### 启动各语言服务

```bash
# 🐹 启动 Go 服务 (端口 8080)
cd video-platform-go
go mod download
go run ./cmd/api-gateway

# ☕ 启动 Java 服务 (端口 8081)
cd video-platform-java
mvn spring-boot:run

# 🦀 启动 Rust 服务 (端口 8082)
cd video-platform-rust
cargo run --release
```

### 访问服务

| 服务         | 地址                                      | 说明          |
| ---------- | --------------------------------------- | ----------- |
| Go API     | <http://localhost:8080>                 | Go 版本 API   |
| Java API   | <http://localhost:8081>                 | Java 版本 API |
| Rust API   | <http://localhost:8082>                 | Rust 版本 API |
| Swagger UI | <http://localhost:8081/swagger-ui.html> | API 文档      |
| MinIO 控制台  | <http://localhost:9001>                 | 对象存储管理      |
| Jaeger     | <http://localhost:16686>                | 链路追踪        |
| Prometheus | <http://localhost:9090>                 | 指标查询        |
| Grafana    | <http://localhost:3000>                 | 监控面板        |

***

## 📚 子项目文档

| 项目                  | 描述              | 链接                                      |
| ------------------- | --------------- | --------------------------------------- |
| 🐹 **VidFlow Go**   | Go 语言实现，轻量高效    | [README](video-platform-go/README.md)   |
| ☕ **VidFlow Java**  | Java 语言实现，企业级生态 | [README](video-platform-java/README.md) |
| 🦀 **VidFlow Rust** | Rust 语言实现，极致性能  | [README](video-platform-rust/README.md) |

***

## 📡 API 文档

### 认证接口

| 方法     | 路径                      | 描述       |
| ------ | ----------------------- | -------- |
| `POST` | `/api/v1/auth/register` | 用户注册     |
| `POST` | `/api/v1/auth/login`    | 用户登录     |
| `POST` | `/api/v1/auth/refresh`  | 刷新 Token |

### 用户接口

| 方法     | 路径                          | 描述     |  认证 |
| ------ | --------------------------- | ------ | :-: |
| `GET`  | `/api/v1/users/me`          | 获取当前用户 |  ✅  |
| `PUT`  | `/api/v1/users/me`          | 更新用户资料 |  ✅  |
| `POST` | `/api/v1/users/{id}/follow` | 关注用户   |  ✅  |

### 视频接口

| 方法     | 路径                           | 描述     |   认证   |
| ------ | ---------------------------- | ------ | :----: |
| `POST` | `/api/v1/videos`             | 上传视频   |    ✅   |
| `GET`  | `/api/v1/videos`             | 获取视频列表 | <br /> |
| `GET`  | `/api/v1/videos/{id}`        | 获取视频详情 | <br /> |
| `GET`  | `/api/v1/videos/{id}/stream` | 获取视频流  | <br /> |

### 互动接口

| 方法     | 路径                             | 描述   |  认证 |
| ------ | ------------------------------ | ---- | :-: |
| `POST` | `/api/v1/videos/{id}/like`     | 点赞视频 |  ✅  |
| `POST` | `/api/v1/videos/{id}/favorite` | 收藏视频 |  ✅  |
| `POST` | `/api/v1/comments`             | 发表评论 |  ✅  |

***

## 🚢 部署指南

### Docker Compose 部署

```bash
# 构建所有服务镜像
docker-compose build

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### Kubernetes 部署

```bash
# 创建命名空间和配置
kubectl apply -f deploy/k8s/deployment.yaml

# 查看部署状态
kubectl get pods -n vidflow

# 查看服务
kubectl get svc -n vidflow
```

***

## 📊 监控与可观测性

### Prometheus 指标

所有服务暴露 `/metrics` 端点，提供以下指标：

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

# 业务
video_uploads_total
user_registrations_total
```

### Jaeger 链路追踪

集成 OpenTelemetry，支持分布式链路追踪，快速定位性能瓶颈。

***

## 🧪 测试

```bash
# Go 测试
cd video-platform-go && go test ./... -v

# Java 测试
cd video-platform-java && mvn test

# Rust 测试
cd video-platform-rust && cargo test
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
- [ ] WebRTC 实时通信

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
