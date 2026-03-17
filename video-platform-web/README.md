<div align="center">

# 🌐 VidFlow Web

**基于 SvelteKit 构建的现代化视频流媒体前端应用**

[!\[SvelteKit\](https://img.shields.io/badge/SvelteKit-2-FF3E00?style=flat\&logo=svelte null)](https://kit.svelte.dev/)
[!\[Svelte\](https://img.shields.io/badge/Svelte-4-FF3E00?style=flat\&logo=svelte null)](https://svelte.dev/)
[!\[TypeScript\](https://img.shields.io/badge/TypeScript-5-3178C6?style=flat\&logo=typescript null)](https://www.typescriptlang.org/)
[!\[Tailwind CSS\](https://img.shields.io/badge/Tailwind%20CSS-3-06B6D4?style=flat\&logo=tailwindcss null)](https://tailwindcss.com/)
[!\[License\](https://img.shields.io/badge/License-MIT-blue.svg null)](LICENSE)

[功能特性](#-功能特性) •
[技术架构](#-技术架构) •
[快速开始](#-快速开始) •
[组件文档](#-组件文档)

</div>

***

## 📖 项目简介

**VidFlow Web** 是 VidFlow 视频流媒体平台的前端应用，采用 SvelteKit 2 + Svelte 4 + Tailwind CSS 构建，提供流畅的用户体验和美观的界面设计，支持服务端渲染（SSR）和静态站点生成（SSG）。

### ✨ 核心亮点

| 特性           | 描述                      |
| ------------ | ----------------------- |
| ⚡ **极致性能**   | Svelte 编译型框架，无虚拟 DOM 开销 |
| 🎨 **精美 UI** | Bits UI 组件库，支持深色模式      |
| 📱 **响应式设计** | 完美适配桌面端、平板、移动端          |
| 🌙 **深色模式**  | 系统级深色模式自动切换             |
| 🔐 **类型安全**  | TypeScript + Zod 运行时验证  |

***

## 🎯 功能特性

### 📺 视频浏览

- 🎬 视频列表瀑布流展示
- 🔍 实时搜索与筛选
- 📂 分类标签导航
- 🔥 热门视频推荐

### 👤 用户中心

- 🔐 登录/注册/找回密码
- 👥 个人资料管理
- ❤️ 收藏夹管理
- 📜 观看历史记录

### 📤 视频上传

- 📁 拖拽上传支持
- 📊 上传进度显示
- 🎬 视频信息编辑
- 🖼️ 封面图选择

### 🎮 视频播放

- 📺 HLS 自适应播放
- 🎯 弹幕互动
- ⏯️ 播放控制
- 📊 播放进度记忆

***

## 🏗️ 技术架构

### 技术栈

<table>
<tr>
<td width="50%">

#### 核心框架

- **SvelteKit 2** - 全栈框架
- **Svelte 4** - 编译型 UI 框架
- **TypeScript 5** - 类型安全
- **Vite 5** - 构建工具

</td>
<td width="50%">

#### UI 组件

- **Tailwind CSS 3** - 原子化 CSS
- **Bits UI** - 无样式组件库
- **Lucide Svelte** - 图标库
- **Mode Watcher** - 深色模式

</td>
</tr>
<tr>
<td width="50%">

#### 状态管理

- **Svelte Stores** - 响应式状态
- **Zod** - 运行时验证

</td>
<td width="50%">

#### 工具链

- **ESLint** - 代码检查
- **Prettier** - 代码格式化
- **Svelte Check** - 类型检查

</td>
</tr>
</table>

### 项目结构

```
video-platform-web/
├── 📄 package.json              # 📦 依赖配置
├── 📄 svelte.config.js          # ⚙️ SvelteKit 配置
├── 📄 vite.config.ts            # ⚡ Vite 构建配置
├── 📄 tailwind.config.js        # 🎨 Tailwind 配置
├── 📄 Dockerfile                # 🐳 容器构建
└── 📂 src/
    ├── 📄 app.css               # 🎨 全局样式
    ├── 📄 app.html              # 📄 HTML 模板
    ├── 📂 lib/
    │   ├── 📂 api/              # 📡 API 客户端
    │   │   ├── auth.ts          #    认证 API
    │   │   └── video.ts         #    视频 API
    │   ├── 📂 components/       # 🧩 组件
    │   │   ├── 📂 ui/           #    UI 基础组件
    │   │   │   ├── button/      #    按钮组件
    │   │   │   ├── card/        #    卡片组件
    │   │   │   ├── input/       #    输入框组件
    │   │   │   ├── label/       #    标签组件
    │   │   │   └── toaster/     #    提示组件
    │   │   ├── Navigation.svelte  # 导航栏
    │   │   └── VideoCard.svelte   # 视频卡片
    │   ├── config.ts            # ⚙️ 应用配置
    │   ├── types.ts             # 📝 TypeScript 类型
    │   ├── utils.ts             # 🔧 工具函数
    │   └── index.ts             # 📦 导出入口
    └── 📂 routes/               # 🛣️ 路由页面
        ├── +layout.svelte       #    布局组件
        ├── +layout.server.ts    #    布局服务端
        ├── +page.svelte         #    首页
        ├── +page.server.ts      #    首页服务端
        ├── 📂 login/            #    登录页
        └── 📂 register/         #    注册页
```

***

## 🚀 快速开始

### 环境要求

| 依赖            | 版本  | 说明   |
| ------------- | --- | ---- |
| Node.js       | 18+ | 运行环境 |
| npm/pnpm/yarn | 最新  | 包管理器 |

### 本地开发

```bash
# 1️⃣ 克隆项目
git clone
cd video-platform-web

# 2️⃣ 安装依赖
npm install
# 或使用 pnpm
pnpm install

# 3️⃣ 配置环境变量
cp .env.example .env
# 编辑 .env 文件配置 API 地址

# 4️⃣ 启动开发服务器
npm run dev
```

服务启动后访问：

- 🌐 应用: <http://localhost:5173>
- 📱 移动端预览: <http://localhost:5173> (响应式)

### 可用脚本

| 命令                | 描述              |
| ----------------- | --------------- |
| `npm run dev`     | 启动开发服务器         |
| `npm run build`   | 构建生产版本          |
| `npm run preview` | 预览生产版本          |
| `npm run check`   | TypeScript 类型检查 |
| `npm run lint`    | 代码检查            |
| `npm run format`  | 代码格式化           |

### Docker 部署

```bash
# 构建镜像
docker build -t vidflow-web:latest .

# 运行容器
docker run -d -p 3000:3000 \
  -e VITE_API_BASE_URL=http://api.example.com \
  vidflow-web:latest
```

***

## 🧩 组件文档

### Button 按钮组件

```svelte
<script>
  import { Button } from '$lib/components/ui/button';
</script>

<Button variant="default">默认按钮</Button>
<Button variant="outline">轮廓按钮</Button>
<Button variant="ghost">幽灵按钮</Button>
<Button variant="destructive">危险按钮</Button>
<Button size="sm">小按钮</Button>
<Button size="lg">大按钮</Button>
<Button disabled>禁用按钮</Button>
```

### Card 卡片组件

```svelte
<script>
  import { 
    Card, CardHeader, CardTitle, 
    CardDescription, CardContent, CardFooter 
  } from '$lib/components/ui/card';
</script>

<Card>
  <CardHeader>
    <CardTitle>视频标题</CardTitle>
    <CardDescription>视频描述信息</CardDescription>
  </CardHeader>
  <CardContent>
    <p>视频内容区域</p>
  </CardContent>
  <CardFooter>
    <Button>观看视频</Button>
  </CardFooter>
</Card>
```

### Input 输入框组件

```svelte
<script>
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  
  let email = '';
</script>

<Label for="email">邮箱地址</Label>
<Input 
  id="email"
  type="email" 
  placeholder="请输入邮箱" 
  bind:value={email} 
/>
```

***

## 📡 API 集成

### 配置 API 地址

```typescript
// src/lib/config.ts
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
```

### 认证 API

```typescript
import { authApi } from '$lib/api/auth';

// 用户登录
const response = await authApi.login(email, password);

// 用户注册
const user = await authApi.register(username, email, password, displayName);

// 刷新令牌
const tokens = await authApi.refreshToken(refreshToken);

// 用户登出
await authApi.logout(token);
```

### 视频 API

```typescript
import { videoApi } from '$lib/api/video';

// 获取视频列表
const videos = await videoApi.list(token, page, pageSize);

// 获取视频详情
const video = await videoApi.get(videoId, token);

// 搜索视频
const results = await videoApi.search(query);

// 创建视频
const video = await videoApi.create(token, { 
  title, 
  description,
  category_id,
  visibility 
});

// 发布视频
await videoApi.publish(token, videoId);

// 删除视频
await videoApi.delete(token, videoId);
```

***

## 🛣️ 路由页面

| 路径            | 描述      |  认证 | SSR |
| ------------- | ------- | :-: | :-: |
| `/`           | 首页/视频列表 |  ❌  |  ✅  |
| `/login`      | 用户登录    |  ❌  |  ✅  |
| `/register`   | 用户注册    |  ❌  |  ✅  |
| `/video/[id]` | 视频详情/播放 |  ❌  |  ✅  |
| `/upload`     | 上传视频    |  ✅  |  ❌  |
| `/profile`    | 个人中心    |  ✅  |  ❌  |
| `/settings`   | 账户设置    |  ✅  |  ❌  |
| `/favorites`  | 我的收藏    |  ✅  |  ❌  |
| `/history`    | 观看历史    |  ✅  |  ❌  |

***

## 🎨 主题定制

### CSS 变量

```css
/* src/app.css */
:root {
  --background: 0 0% 100%;
  --foreground: 222.2 84% 4.9%;
  --primary: 262 83% 58%;
  --primary-foreground: 210 40% 98%;
  --secondary: 210 40% 96%;
  --muted: 210 40% 96%;
  --accent: 210 40% 96%;
  --destructive: 0 84% 60%;
  --border: 214.3 31.8% 91.4%;
  --ring: 262 83% 58%;
  --radius: 0.5rem;
}

.dark {
  --background: 222.2 84% 4.9%;
  --foreground: 210 40% 98%;
  --primary: 262 83% 58%;
  /* ... */
}
```

### Tailwind 配置

```javascript
// tailwind.config.js
export default {
  darkMode: 'class',
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))'
        },
        // ...
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)'
      }
    }
  }
}
```

***

## 📱 响应式设计

### 断点配置

| 断点  | 前缀     | 最小宽度   | 典型设备            |
| --- | ------ | ------ | --------------- |
| 手机  | `sm:`  | 640px  | iPhone, Android |
| 平板  | `md:`  | 768px  | iPad, Tablet    |
| 笔记本 | `lg:`  | 1024px | MacBook, Laptop |
| 桌面  | `xl:`  | 1280px | Desktop         |
| 大屏  | `2xl:` | 1536px | Large Monitor   |

### 响应式示例

```svelte
<!-- 视频卡片网格 -->
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
  {#each videos as video}
    <VideoCard {video} />
  {/each}
</div>

<!-- 导航栏 -->
<nav class="flex flex-col md:flex-row gap-4">
  <a href="/">首页</a>
  <a href="/upload">上传</a>
  <a href="/profile">我的</a>
</nav>
```

***

## ⚙️ 环境变量

| 变量                  | 描述        | 默认值                     |
| ------------------- | --------- | ----------------------- |
| `VITE_API_BASE_URL` | 后端 API 地址 | `http://localhost:8080` |

***

## 📝 开发规范

### 组件命名

- 组件文件使用 PascalCase: `VideoCard.svelte`
- 工具函数使用 camelCase: `formatDuration()`

### 目录结构

- UI 基础组件放在 `lib/components/ui/`
- 业务组件放在 `lib/components/`
- API 客户端放在 `lib/api/`

### 类型定义

- 接口使用 PascalCase: `VideoResponse`
- 类型导出放在 `lib/types.ts`

***

## 🗺️ 路线图

- [x] 基础 UI 组件库
- [x] 用户认证流程
- [x] 视频列表展示
- [ ] 视频播放器
- [ ] 弹幕系统
- [ ] 直播功能
- [ ] PWA 支持
- [ ] 国际化 (i18n)

***

## 📄 许可证

本项目采用 [MIT License](LICENSE) 许可证。

***

<div align="center">

[⬆ 返回顶部](#-vidflow-web)

</div>
