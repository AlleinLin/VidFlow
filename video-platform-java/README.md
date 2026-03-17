<div align="center">

# ☕ VidFlow Java

**基于 Spring Boot 3 构建的视频流媒体平台**

[!\[Java Version\](https://img.shields.io/badge/Java-21+-ED8B00?style=flat\&logo=openjdk null)](https://openjdk.org/)
[!\[Spring Boot\](https://img.shields.io/badge/Spring%20Boot-3.3-6DB33F?style=flat\&logo=springboot null)](https://spring.io/projects/spring-boot)
[!\[License\](https://img.shields.io/badge/License-MIT-blue.svg null)](LICENSE)

[功能特性](#-功能特性) •
[技术架构](#-技术架构) •
[快速开始](#-快速开始) •
[API 文档](#-api-文档)

</div>

***

## 📖 项目简介

**VidFlow Java** 是 VidFlow 视频流媒体平台的 Java 实现版本，采用 Spring Boot 3.3 + JDK 21 构建，充分利用 Virtual Threads 实现高并发处理，适用于企业级视频服务场景。

### ✨ 核心亮点

| 特性                     | 描述                               |
| ---------------------- | -------------------------------- |
| 🚀 **Virtual Threads** | JDK 21 虚拟线程，轻松应对高并发              |
| 🔐 **Spring Security** | 完整的安全认证授权体系                      |
| 📊 **微服务就绪**           | Spring Cloud 生态，支持分布式部署          |
| 🛡️ **Resilience4j**   | 熔断降级、限流保护                        |
| 📈 **可观测性**            | Micrometer + OpenTelemetry 全链路监控 |

***

## 🎯 功能特性

### 👤 用户服务

- 🔐 用户注册/登录（BCrypt 密码加密）
- 🎫 JWT Token 认证与 OAuth2 资源服务器
- 👥 用户资料管理与头像上传
- 🤝 关注/粉丝社交系统

### 📹 视频服务

- 📤 大文件分片上传
- 🎬 多分辨率转码调度
- 📺 HLS 自适应流媒体分发
- 🗂️ 视频分类与标签管理

### 💬 互动服务

- 💭 多级评论与回复系统
- ❤️ 点赞与收藏功能
- 🎯 实时弹幕系统
- 📜 观看历史记录

### 🔍 搜索与推荐

- 🔎 全文搜索（视频/用户）
- 🔥 热门视频推荐
- 🎯 个性化推荐算法

***

## 🏗️ 技术架构

### 技术栈

<table>
<tr>
<td width="50%">

#### 核心框架

- **Spring Boot 3.3** - 应用框架
- **Spring Security** - 安全框架
- **Spring Data JPA** - 数据访问
- **Spring Cloud** - 微服务生态

</td>
<td width="50%">

#### 数据存储

- **PostgreSQL** - 主数据库
- **Redis** - 缓存与会话
- **Cassandra** - 海量数据存储
- **MinIO** - 对象存储

</td>
</tr>
<tr>
<td width="50%">

#### 安全认证

- **JWT (jjwt)** - Token 生成
- **OAuth2** - 资源服务器
- **BCrypt** - 密码加密

</td>
<td width="50%">

#### 基础设施

- **Kafka** - 消息队列
- **Resilience4j** - 熔断限流
- **OpenTelemetry** - 链路追踪
- **Prometheus** - 指标监控

</td>
</tr>
</table>

### 项目结构

```
vidflow-java/
├── 📂 src/main/
│   ├── 📂 java/
│   │   └── 📂 com/vidflow/
│   │       ├── 📂 config/           # ⚙️ 配置类
│   │       ├── 📂 controller/       # 🌐 REST 控制器
│   │       ├── 📂 service/          # 💼 业务服务层
│   │       ├── 📂 repository/       # 🗄️ 数据访问层
│   │       ├── 📂 entity/           # 📦 实体类
│   │       ├── 📂 dto/              # 📋 数据传输对象
│   │       ├── 📂 security/         # 🔐 安全配置
│   │       └── 📂 exception/        # ❌ 异常处理
│   └── 📂 resources/
│       └── application.yml          # 📝 配置文件
├── 📄 pom.xml                       # Maven 依赖
└── 📄 Dockerfile                    # 🐳 容器构建
```

***

## 🚀 快速开始

### 环境要求

| 依赖         | 版本   | 说明                 |
| ---------- | ---- | ------------------ |
| JDK        | 21+  | 支持 Virtual Threads |
| Maven      | 3.9+ | 构建工具               |
| PostgreSQL | 16+  | 主数据库               |
| Redis      | 7+   | 缓存服务               |

### 本地开发

```bash
# 1️⃣ 克隆项目
git clone
cd vidflow-java

# 2️⃣ 安装依赖
mvn clean install

# 3️⃣ 配置环境变量
export SPRING_DATASOURCE_URL=jdbc:postgresql://localhost:5432/vidflow
export SPRING_DATASOURCE_USERNAME=postgres
export SPRING_DATASOURCE_PASSWORD=your_password
export SPRING_DATA_REDIS_HOST=localhost
export JWT_SECRET=your-secret-key-at-least-256-bits

# 4️⃣ 启动服务
mvn spring-boot:run
```

服务启动后访问：

- 🌐 API: <http://localhost:8081>
- 📚 Swagger UI: <http://localhost:8081/swagger-ui.html>
- 📊 Actuator: <http://localhost:8081/actuator>

### Docker 部署

```bash
# 构建镜像
docker build -t vidflow-java:latest .

# 运行容器
docker run -d -p 8081:8081 \
  -e SPRING_DATASOURCE_URL=jdbc:postgresql://host:5432/vidflow \
  -e SPRING_DATASOURCE_USERNAME=postgres \
  -e SPRING_DATASOURCE_PASSWORD=password \
  vidflow-java:latest
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

| 方法     | 路径                             | 描述     |   认证   |
| ------ | ------------------------------ | ------ | :----: |
| `GET`  | `/api/v1/users/me`             | 获取当前用户 |    ✅   |
| `PUT`  | `/api/v1/users/me`             | 更新用户资料 |    ✅   |
| `POST` | `/api/v1/users/{id}/follow`    | 关注用户   |    ✅   |
| `GET`  | `/api/v1/users/{id}/followers` | 获取粉丝列表 | <br /> |

### 视频接口

| 方法       | 路径                    | 描述     |   认证   |
| -------- | --------------------- | ------ | :----: |
| `POST`   | `/api/v1/videos`      | 上传视频   |    ✅   |
| `GET`    | `/api/v1/videos`      | 获取视频列表 | <br /> |
| `GET`    | `/api/v1/videos/{id}` | 获取视频详情 | <br /> |
| `PUT`    | `/api/v1/videos/{id}` | 更新视频信息 |    ✅   |
| `DELETE` | `/api/v1/videos/{id}` | 删除视频   |    ✅   |

***

## ⚙️ 配置说明

### application.yml 核心配置

```yaml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/vidflow
    username: postgres
    password: ${DB_PASSWORD}
  
  data:
    redis:
      host: localhost
      port: 6379
  
  jpa:
    hibernate:
      ddl-auto: validate
    show-sql: false

jwt:
  secret: ${JWT_SECRET}
  access-token-ttl: 15m
  refresh-token-ttl: 168h
```

***

## 📈 性能特性

| 特性                     | 说明                  |
| ---------------------- | ------------------- |
| 🧵 **Virtual Threads** | JDK 21 虚拟线程，轻量级并发   |
| 🔗 **连接池优化**           | HikariCP 高性能连接池     |
| 💾 **多级缓存**            | Redis + 本地缓存        |
| ⚡ **异步处理**             | @Async + Kafka 异步任务 |
| 🛡️ **熔断保护**           | Resilience4j 熔断降级   |

***

## 🧪 测试

```bash
# 运行所有测试
mvn test

# 运行集成测试
mvn verify -P integration-test

# 生成覆盖率报告
mvn jacoco:report
```

***

## 📄 许可证

本项目采用 [MIT License](LICENSE) 许可证。

***

<div align="center">

[⬆ 返回顶部](#-vidflow-java)

</div>
