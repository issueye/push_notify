# Push Notify Frontend

前端项目基于 Vue3 + Naive UI + TailwindCSS 开发。

## 技术栈

- **框架**: Vue 3
- **UI库**: Naive UI
- **样式**: TailwindCSS
- **HTTP**: Axios
- **路由**: Vue Router
- **状态管理**: Pinia
- **图表**: ECharts

## 快速开始

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

服务将在 `http://localhost:3000` 启动。

### 3. 构建生产版本

```bash
npm run build
```

## 项目结构

```
frontend/
├── public/                 # 静态资源
├── src/
│   ├── assets/             # 资源文件
│   ├── components/         # 公共组件
│   │   └── layout/         # 布局组件
│   ├── composables/        # 组合式函数
│   ├── router/             # 路由配置
│   ├── services/           # API服务层
│   ├── stores/             # Pinia状态管理
│   ├── types/              # 类型定义
│   ├── utils/              # 工具函数
│   ├── views/              # 页面视图
│   ├── App.vue
│   └── main.js
├── index.html
├── package.json
├── vite.config.js
└── tailwind.config.js
```

## 页面列表

| 页面 | 路由 | 说明 |
|------|------|------|
| 登录 | /login | 用户登录 |
| 工作台 | /dashboard | 统计概览 |
| 仓库管理 | /repos | 代码仓库配置 |
| 推送目标 | /targets | 钉钉/邮箱配置 |
| 推送记录 | /pushes | 推送历史查询 |
| 消息模板 | /templates | 模板管理 |
| 提示词 | /prompts | CODEVIEW提示词 |
| AI模型 | /models | AI服务配置 |
| 用户管理 | /users | 用户管理(管理员) |
| 系统日志 | /logs/system | 系统运行日志 |
| 操作日志 | /logs/operations | 用户操作日志 |
| AI调用日志 | /logs/ai-calls | AI调用记录 |
| 个人设置 | /settings | 个人设置 |

## API配置

在 `.env` 文件中配置API地址：

```env
VITE_API_BASE_URL=/api/v1
```

## 开发规范

- 使用 Vue 3 Composition API
- 组件命名采用 PascalCase
- 组合式函数命名采用 camelCase，以 use 开头
- 提交前运行 lint 检查
