# Push Notify 前端开发计划

> **文档版本**: V1.0  
> **创建日期**: 2026年1月19日  
> **技术栈**: Vue3 + Naive UI + TailwindCSS + Axios + JavaScript

---

## 一、技术架构

### 1.1 技术选型

| 组件 | 技术 | 版本要求 | 说明 |
|------|------|---------|------|
| 框架 | Vue3 | >= 3.4 | 主要开发框架 |
| UI库 | Naive UI | latest | Vue3组件库 |
| 样式 | TailwindCSS | latest | 原子化CSS |
| HTTP | Axios | latest | 网络请求 |
| 路由 | Vue Router | latest | 路由管理 |
| 状态管理 | Pinia | latest | 状态管理 |
| 持久化 | pinia-plugin-persistedstate | latest | 状态持久化 |
| 构建工具 | Vite | latest | 开发构建 |
| 代码规范 | ESLint + Prettier | latest | 代码格式化 |

### 1.2 浏览器兼容性

| 浏览器 | 最低版本 | 说明 |
|--------|---------|------|
| Chrome | 90 | 推荐 |
| Firefox | 88 | 支持 |
| Safari | 15 | 支持 |
| Edge | 90 | 支持 |

### 1.3 目录结构

```
frontend/
├── public/                     # 静态资源
│   └── favicon.ico
├── src/
│   ├── assets/                 # 资源文件
│   │   └── logo.png
│   ├── components/             # 公共组件
│   │   ├── common/             # 通用组件
│   │   │   ├── PageHeader.vue
│   │   │   ├── PageContent.vue
│   │   │   ├── SearchForm.vue
│   │   │   ├── DataTable.vue
│   │   │   └── EmptyData.vue
│   │   ├── form/               # 表单组件
│   │   │   ├── BaseForm.vue
│   │   │   ├── BaseInput.vue
│   │   │   ├── BaseSelect.vue
│   │   │   ├── BaseSwitch.vue
│   │   │   └── BaseUpload.vue
│   │   └── layout/             # 布局组件
│   │       ├── MainLayout.vue
│   │       ├── Sidebar.vue
│   │       ├── Header.vue
│   │       └── Footer.vue
│   ├── composables/            # 组合式函数
│   │   ├── useTable.js
│   │   ├── useForm.js
│   │   ├── useModal.js
│   │   └── useMessage.js
│   ├── directives/             # 自定义指令
│   │   ├── permission.js
│   │   └── auth.js
│   ├── hooks/                  # 自定义Hook
│   │   ├── useAuth.js
│   │   └── usePermission.js
│   ├── router/                 # 路由配置
│   │   ├── index.js
│   │   ├── routes.js
│   │   └── guard.js
│   ├── stores/                 # Pinia状态管理
│   │   ├── index.js
│   │   ├── user.js
│   │   ├── app.js
│   │   └── permission.js
│   ├── services/               # API服务层
│   │   ├── index.js
│   │   ├── auth.js
│   │   ├── user.js
│   │   ├── repo.js
│   │   ├── target.js
│   │   ├── push.js
│   │   ├── template.js
│   │   ├── prompt.js
│   │   ├── model.js
│   │   └── log.js
│   ├── utils/                  # 工具函数
│   │   ├── request.js
│   │   ├── auth.js
│   │   ├── constants.js
│   │   ├── helpers.js
│   │   └── format.js
│   ├── views/                  # 页面视图
│   │   ├── login/              # 登录页
│   │   │   ├── index.vue
│   │   │   └── LoginForm.vue
│   │   ├── dashboard/          # 工作台
│   │   │   └── index.vue
│   │   ├── repos/              # 仓库管理
│   │   │   ├── index.vue
│   │   │   ├── RepoList.vue
│   │   │   ├── RepoForm.vue
│   │   │   └── RepoDetail.vue
│   │   ├── targets/            # 推送目标
│   │   │   ├── index.vue
│   │   │   ├── TargetList.vue
│   │   │   ├── TargetForm.vue
│   │   │   └── TargetDetail.vue
│   │   ├── pushes/             # 推送记录
│   │   │   ├── index.vue
│   │   │   ├── PushList.vue
│   │   │   └── PushDetail.vue
│   │   ├── templates/          # 消息模板
│   │   │   ├── index.vue
│   │   │   ├── TemplateList.vue
│   │   │   ├── TemplateForm.vue
│   │   │   ├── TemplatePreview.vue
│   │   │   └── TemplateVersion.vue
│   │   ├── prompts/            # 提示词
│   │   │   ├── index.vue
│   │   │   ├── PromptList.vue
│   │   │   ├── PromptForm.vue
│   │   │   └── PromptTest.vue
│   │   ├── models/             # AI模型
│   │   │   ├── index.vue
│   │   │   ├── ModelList.vue
│   │   │   └── ModelForm.vue
│   │   ├── users/              # 用户管理
│   │   │   ├── index.vue
│   │   │   ├── UserList.vue
│   │   │   └── UserForm.vue
│   │   ├── logs/               # 日志管理
│   │   │   ├── index.vue
│   │   │   ├── SystemLog.vue
│   │   │   ├── OperationLog.vue
│   │   │   └── AICallLog.vue
│   │   ├── settings/           # 个人设置
│   │   │   └── index.vue
│   │   └── error/              # 错误页
│   │       ├── 403.vue
│   │       ├── 404.vue
│   │       └── 500.vue
│   ├── App.vue
│   ├── main.js
│   └── style.css               # 全局样式
├── .env                        # 环境变量
├── .env.development            # 开发环境
├── .env.production             # 生产环境
├── index.html
├── package.json
├── vite.config.js
├── tailwind.config.js
└── postcss.config.js
```

---

## 二、页面设计

### 2.1 页面结构

```
┌─────────────────────────────────────────────────────────────┐
│                        Header                                │
│  [Logo] [系统名称]           [消息] [个人] [退出]           │
├──────────┬──────────────────────────────────────────────────┤
│          │                                                  │
│  Sidebar │                   Main Content                   │
│          │                                                  │
│  ┌──────┐│   ┌──────────────────────────────────────────┐  │
│  │Dashboard│   │                                          │  │
│  ├──────┤│   │           页面内容区域                     │  │
│  │仓库管理││   │                                          │  │
│  ├──────┤│   │                                          │  │
│  │推送目标││   │                                          │  │
│  ├──────┤│   │                                          │  │
│  │推送记录││   │                                          │  │
│  ├──────┤│   │                                          │  │
│  │模板管理││   │                                          │  │
│  ├──────┤│   │                                          │  │
│  │提示词  ││   │                                          │  │
│  ├──────┤│   │                                          │  │
│  │AI模型  ││   │                                          │  │
│  ├──────┤│   │                                          │  │
│  │用户管理││   │                                          │  │
│  ├──────┤│   │                                          │  │
│  │日志管理││   │                                          │  │
│  └──────┘│   └──────────────────────────────────────────┘  │
│          │                                                  │
└──────────┴──────────────────────────────────────────────────┘
```

### 2.2 页面清单

| 页面 | 路由 | 说明 | 权限 |
|------|------|------|------|
| 登录页 | /login | 用户登录 | 公开 |
| 工作台 | /dashboard | 统计概览 | 普通用户 |
| 仓库列表 | /repos | 仓库管理列表 | 普通用户 |
| 仓库详情 | /repos/:id | 仓库详情 | 普通用户 |
| 仓库表单 | /repos/create, /repos/:id/edit | 创建/编辑仓库 | 管理员 |
| 推送目标列表 | /targets | 推送目标列表 | 普通用户 |
| 推送目标详情 | /targets/:id | 推送目标详情 | 普通用户 |
| 推送目标表单 | /targets/create, /targets/:id/edit | 创建/编辑 | 管理员 |
| 推送记录列表 | /pushes | 推送历史记录 | 普通用户 |
| 推送详情 | /pushes/:id | 推送详情 | 普通用户 |
| 模板列表 | /templates | 消息模板列表 | 普通用户 |
| 模板表单 | /templates/create, /templates/:id/edit | 创建/编辑 | 管理员 |
| 模板预览 | /templates/:id/preview | 模板预览 | 普通用户 |
| 提示词列表 | /prompts | 提示词列表 | 普通用户 |
| 提示词表单 | /prompts/create, /prompts/:id/edit | 创建/编辑 | 管理员 |
| 提示词测试 | /prompts/:id/test | 提示词测试 | 管理员 |
| AI模型列表 | /models | AI模型列表 | 普通用户 |
| 模型表单 | /models/create, /models/:id/edit | 创建/编辑 | 管理员 |
| 用户列表 | /users | 用户管理列表 | 管理员 |
| 用户表单 | /users/create, /users/:id/edit | 创建/编辑 | 管理员 |
| 系统日志 | /logs/system | 系统运行日志 | 管理员 |
| 操作日志 | /logs/operations | 用户操作日志 | 管理员 |
| AI调用日志 | /logs/ai-calls | AI调用日志 | 管理员 |
| 个人设置 | /settings | 个人设置 | 普通用户 |
| 403 | /403 | 无权限 | 公开 |
| 404 | /404 | 页面不存在 | 公开 |
| 500 | /500 | 服务器错误 | 公开 |

---

## 三、开发计划

### 第一阶段：项目初始化（第1周）

#### 3.1.1 项目搭建

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 初始化项目 | 使用Vite创建Vue3项目 | 0.5天 | - |
| 安装依赖 | 安装Naive UI、TailwindCSS等 | 0.5天 | - |
| 配置Vite | 开发/生产环境配置 | 0.5天 | - |
| 配置TailwindCSS | 样式配置 | 0.5天 | - |
| 配置ESLint | 代码规范配置 | 0.5天 | - |

#### 3.1.2 基础架构

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 路由配置 | Vue Router基础配置 | 0.5天 | - |
| Axios封装 | 请求拦截和响应拦截 | 0.5天 | - |
| Pinia配置 | 状态管理基础配置 | 0.5天 | - |
| 持久化配置 | Token和状态持久化 | 0.5天 | - |

#### 3.1.3 布局组件

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 主布局 | MainLayout组件 | 0.5天 | - |
| 侧边栏 | Sidebar组件 | 0.5天 | 主布局 |
| 顶部栏 | Header组件 | 0.5天 | 主布局 |
| 路由守卫 | 登录验证和权限控制 | 0.5天 | 路由配置 |

#### 3.1.4 交付物

- [ ] Vue3项目骨架
- [ ] 完整依赖安装
- [ ] 路由配置
- [ ] Axios封装
- [ ] Pinia状态管理
- [ ] 基础布局组件

### 第二阶段：登录认证（第2周）

#### 3.2.1 登录页面

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 登录页布局 | Login页面 | 0.5天 | - |
| 登录表单 | 包含用户名密码和验证码 | 1天 | Axios封装 |
| 登录逻辑 | 调用登录API | 0.5天 | 登录表单 |
| Token存储 | 保存Token到本地 | 0.5天 | Pinia配置 |

#### 3.2.2 注册页面

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 注册页布局 | Register页面 | 0.5天 | - |
| 注册表单 | 包含用户名、邮箱、密码 | 1天 | - |
| 注册逻辑 | 调用注册API | 0.5天 | 注册表单 |

#### 3.2.3 权限控制

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 路由守卫 | 未登录跳转登录页 | 0.5天 | 路由配置 |
| 权限指令 | v-permission指令 | 0.5天 | - |
| 无权限页面 | 403页面 | 0.5天 | - |

#### 3.2.4 交付物

- [ ] 登录页面（/login）
- [ ] 注册页面（/register）
- [ ] Token管理
- [ ] 路由守卫
- [ ] 权限控制

### 第三阶段：工作台和用户模块（第2-3周）

#### 3.3.1 工作台

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 工作台布局 | Dashboard页面 | 0.5天 | - |
| 统计卡片 | 推送统计卡片 | 0.5天 | API服务 |
| 近期推送 | 最近推送列表 | 0.5天 | API服务 |
| 快捷入口 | 常用功能入口 | 0.5天 | - |

#### 3.3.2 个人设置

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 设置页面 | Settings页面 | 0.5天 | - |
| 基本信息 | 邮箱、密码修改 | 1天 | API服务 |
| 通知设置 | 通知偏好配置 | 0.5天 | API服务 |
| 外观设置 | 主题、语言切换 | 0.5天 | Pinia |

#### 3.3.3 交付物

- [ ] 工作台（/dashboard）
- [ ] 个人设置（/settings）
- [ ] 统计展示

### 第四阶段：仓库管理模块（第3-4周）

#### 3.4.1 仓库列表

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 列表页面 | Repos页面布局 | 0.5天 | - |
| 搜索表单 | 搜索和筛选 | 0.5天 | 公共组件 |
| 数据表格 | 仓库列表展示 | 1天 | DataTable |
| 分页组件 | 分页控制 | 0.5天 | - |

#### 3.4.2 仓库表单

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 创建表单 | 添加仓库表单 | 1天 | BaseForm |
| 编辑表单 | 编辑仓库表单 | 0.5天 | 创建表单 |
| 仓库类型选择 | GitHub/GitLab/Gitee | 0.5天 | - |
| Webhook配置展示 | 显示生成的Webhook URL | 0.5天 | - |

#### 3.4.3 仓库详情

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 详情页面 | RepoDetail页面 | 1天 | - |
| 基本信息展示 | 仓库信息展示 | 0.5天 | - |
| Webhook配置 | Webhook URL和密钥 | 0.5天 | - |
| 关联目标 | 推送目标列表 | 0.5天 | - |
| 测试按钮 | 测试Webhook | 0.5天 | API服务 |

#### 3.4.4 交付物

- [ ] 仓库列表页面（/repos）
- [ ] 创建/编辑仓库（/repos/create, /repos/:id/edit）
- [ ] 仓库详情页面（/repos/:id）
- [ ] Webhook测试功能

### 第五阶段：推送目标模块（第4-5周）

#### 3.5.1 推送目标列表

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 列表页面 | Targets页面布局 | 0.5天 | - |
| 类型筛选 | 钉钉/邮箱筛选 | 0.5天 | - |
| 数据表格 | 目标列表展示 | 1天 | DataTable |

#### 3.5.2 钉钉目标表单

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 钉钉表单 | 钉钉配置表单 | 1天 | BaseForm |
| Token验证 | 验证AccessToken | 0.5天 | API服务 |
| Secret配置 | 密钥配置（可选） | 0.5天 | - |

#### 3.5.3 邮箱目标表单

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 邮箱表单 | 邮箱配置表单 | 1天 | BaseForm |
| SMTP配置 | SMTP服务器配置 | 0.5天 | - |
| 收件人配置 | 多个收件人 | 0.5天 | - |
| SMTP验证 | 验证SMTP配置 | 0.5天 | API服务 |

#### 3.5.4 推送目标详情

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 详情页面 | TargetDetail页面 | 1天 | - |
| 配置信息展示 | 展示配置详情 | 0.5天 | - |
| 关联仓库 | 仓库列表 | 0.5天 | - |
| 统计信息 | 推送统计 | 0.5天 | API服务 |
| 测试推送 | 发送测试消息 | 0.5天 | API服务 |

#### 3.5.5 交付物

- [ ] 推送目标列表（/targets）
- [ ] 创建/编辑推送目标（/targets/create, /targets/:id/edit）
- [ ] 推送目标详情（/targets/:id）
- [ ] 钉钉推送配置
- [ ] 邮箱推送配置
- [ ] 测试推送功能

### 第六阶段：推送记录模块（第5-6周）

#### 3.6.1 推送记录列表

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 列表页面 | Pushes页面布局 | 0.5天 | - |
| 高级搜索 | 时间范围、状态筛选 | 1天 | SearchForm |
| 数据表格 | 推送记录列表 | 1天 | DataTable |

#### 3.6.2 推送详情

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 详情页面 | PushDetail页面 | 1天 | - |
| 基本信息 | 推送基本信息展示 | 0.5天 | - |
| 提交信息 | 提交ID、信息、分支 | 0.5天 | - |
| CODEVIEW结果 | 审查结果展示 | 0.5天 | - |
| 推送内容 | 消息内容展示 | 0.5天 | - |
| 失败原因 | 错误信息展示 | 0.5天 | - |

#### 3.6.3 推送操作

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 重试按钮 | 单条重试 | 0.5天 | API服务 |
| 批量重试 | 批量重试功能 | 0.5天 | API服务 |
| 批量删除 | 批量删除历史记录 | 0.5天 | API服务 |

#### 3.6.4 推送统计

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 统计卡片 | 今日/本周/本月统计 | 0.5天 | API服务 |
| 趋势图表 | 推送趋势图 | 1天 | ECharts |
| 成功率图表 | 成功/失败占比 | 0.5天 | ECharts |

#### 3.6.5 交付物

- [ ] 推送记录列表（/pushes）
- [ ] 推送详情（/pushes/:id）
- [ ] 重试功能
- [ ] 批量操作
- [ ] 统计数据展示
- [ ] 图表展示

### 第七阶段：消息模板模块（第6-7周）

#### 3.7.1 模板列表

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 列表页面 | Templates页面布局 | 0.5天 | - |
| 类型筛选 | 钉钉/邮箱筛选 | 0.5天 | - |
| 场景筛选 | 提交通知/审查通知 | 0.5天 | - |
| 数据表格 | 模板列表 | 1天 | DataTable |

#### 3.7.2 模板表单

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 模板表单 | 创建/编辑模板表单 | 1天 | BaseForm |
| 模板类型选择 | 钉钉/邮件 | 0.5天 | - |
| 场景选择 | 通知场景选择 | 0.5天 | - |
| Markdown编辑器 | 钉钉模板编辑 | 1天 | - |
| HTML编辑器 | 邮件模板编辑 | 1天 | - |
| 变量帮助 | 显示可用变量 | 0.5天 | - |

#### 3.7.3 模板预览

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 预览页面 | TemplatePreview页面 | 1天 | - |
| 实时预览 | 输入时预览效果 | 1天 | - |
| 模拟数据 | 默认模拟数据 | 0.5天 | - |
| 自定义数据 | 自定义预览数据 | 0.5天 | - |

#### 3.7.4 模板版本

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 版本历史 | 版本历史列表 | 0.5天 | API服务 |
| 版本回滚 | 回滚到指定版本 | 0.5天 | API服务 |
| 版本对比 | 版本差异对比 | 1天 | - |

#### 3.7.5 交付物

- [ ] 模板列表（/templates）
- [ ] 创建/编辑模板（/templates/create, /templates/:id/edit）
- [ ] 模板预览（/templates/:id/preview）
- [ ] 模板版本管理
- [ ] Markdown/HTML编辑器

### 第八阶段：提示词模块（第7-8周）

#### 3.8.1 提示词列表

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 列表页面 | Prompts页面布局 | 0.5天 | - |
| 类型筛选 | CODEVIEW/消息提示词 | 0.5天 | - |
| 数据表格 | 提示词列表 | 1天 | DataTable |

#### 3.8.2 提示词表单

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| CODEVIEW表单 | CODEVIEW提示词表单 | 1天 | BaseForm |
| 消息表单 | 推送消息提示词表单 | 0.5天 | - |
| 变量帮助 | 显示可用变量 | 0.5天 | - |
| 模型选择 | 关联AI模型选择 | 0.5天 | - |

#### 3.8.3 提示词测试

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 测试页面 | PromptTest页面 | 1天 | - |
| 测试数据输入 | 输入测试数据 | 0.5天 | - |
| 测试结果展示 | 显示AI返回结果 | 1天 | - |

#### 3.8.4 提示词导入导出

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 导出功能 | 导出为JSON | 0.5天 | API服务 |
| 导入功能 | 从JSON导入 | 0.5天 | API服务 |
| 导入验证 | 格式验证 | 0.5天 | - |

#### 3.8.5 交付物

- [ ] 提示词列表（/prompts）
- [ ] 创建/编辑提示词（/prompts/create, /prompts/:id/edit）
- [ ] 提示词测试（/prompts/:id/test）
- [ ] 导入导出功能

### 第九阶段：AI模型模块（第8周）

#### 3.9.1 模型列表

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 列表页面 | Models页面布局 | 0.5天 | - |
| 数据表格 | 模型列表 | 1天 | DataTable |
| 调用统计 | 调用次数展示 | 0.5天 | - |

#### 3.9.2 模型表单

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 模型表单 | 创建/编辑模型表单 | 1天 | BaseForm |
| 模型类型选择 | GPT-4/Claude等 | 0.5天 | - |
| API配置 | URL和Key输入 | 0.5天 | - |
| 参数配置 | 温度、token等参数 | 1天 | - |
| 配置验证 | 验证API配置 | 0.5天 | API服务 |

#### 3.9.3 交付物

- [ ] AI模型列表（/models）
- [ ] 创建/编辑模型（/models/create, /models/:id/edit）
- [ ] 参数配置
- [ ] 配置验证

### 第十阶段：用户管理模块（第8-9周）

#### 3.10.1 用户列表

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 列表页面 | Users页面布局 | 0.5天 | - |
| 角色筛选 | 管理员/普通用户 | 0.5天 | - |
| 状态筛选 | 激活/锁定 | 0.5天 | - |
| 数据表格 | 用户列表 | 1天 | DataTable |

#### 3.10.2 用户表单

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 创建表单 | 创建用户表单 | 0.5天 | BaseForm |
| 编辑表单 | 编辑用户信息 | 0.5天 | - |
| 角色选择 | 分配用户角色 | 0.5天 | - |

#### 3.10.3 用户操作

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 重置密码 | 重置用户密码 | 0.5天 | API服务 |
| 锁定/解锁 | 账户状态管理 | 0.5天 | API服务 |
| 删除用户 | 删除用户（软删除） | 0.5天 | API服务 |

#### 3.10.4 交付物

- [ ] 用户列表（/users）
- [ ] 创建/编辑用户（/users/create, /users/:id/edit）
- [ ] 用户操作功能

### 第十一阶段：日志管理模块（第9周）

#### 3.11.1 日志列表

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 系统日志 | SystemLog页面 | 1天 | - |
| 操作日志 | OperationLog页面 | 1天 | - |
| AI调用日志 | AICallLog页面 | 1天 | - |

#### 3.11.2 日志功能

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 日志筛选 | 时间、级别筛选 | 1天 | SearchForm |
| 关键词搜索 | 搜索日志内容 | 0.5天 | - |
| 日志详情 | 查看日志详情 | 0.5天 | - |
| 日志导出 | 导出CSV/JSON | 1天 | API服务 |

#### 3.11.3 交付物

- [ ] 系统日志页面（/logs/system）
- [ ] 操作日志页面（/logs/operations）
- [ ] AI调用日志页面（/logs/ai-calls）
- [ ] 搜索和导出功能

### 第十二阶段：测试和优化（第10周）

#### 3.12.1 测试

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 单元测试 | 组件单元测试 | 2天 | - |
| E2E测试 | 关键流程测试 | 1天 | - |
| 兼容性测试 | 多浏览器测试 | 0.5天 | - |

#### 3.12.2 优化

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 性能优化 | 懒加载优化 | 0.5天 | - |
| 首屏优化 | 骨架屏加载 | 0.5天 | - |
| 错误处理 | 全局错误捕获 | 0.5天 | - |

#### 3.12.3 文档

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| README | 项目说明文档 | 0.5天 | - |
| 组件文档 | 重要组件说明 | 0.5天 | - |

---

## 四、API服务层设计

### 4.1 服务模块结构

```typescript
// src/services/index.js - 导出所有服务
export * from './auth'
export * from './user'
export * from './repo'
export * from './target'
export * from './push'
export * from './template'
export * from './prompt'
export * from './model'
export * from './log'
```

### 4.2 Axios封装示例

```typescript
// src/utils/request.js
import axios, { AxiosRequestConfig, AxiosResponse } from 'axios'
import { useUserStore } from '@/stores/user'
import { useMessage } from '@/composables/useMessage'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000,
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message, data } = response.data
    if (code === 200) {
      return data
    }
    useMessage().error(message || '请求失败')
    return Promise.reject(new Error(message))
  },
  (error) => {
    const { response } = error
    if (response) {
      switch (response.status) {
        case 401:
          useMessage().error('登录已过期，请重新登录')
          useUserStore().logout()
          break
        case 403:
          useMessage().error('没有权限访问')
          break
        case 404:
          useMessage().error('请求的资源不存在')
          break
        case 500:
          useMessage().error('服务器错误')
          break
        default:
          useMessage().error(response.data?.message || '请求失败')
      }
    } else {
      useMessage().error('网络连接失败')
    }
    return Promise.reject(error)
  }
)

export default request
```

### 4.3 API服务示例

```typescript
// src/services/repo.js
import request from '@/utils/request'

// 获取仓库列表
export function getRepoList(params) {
  return request({
    url: '/api/v1/repos',
    method: 'get',
    params,
  })
}

// 获取仓库详情
export function getRepoDetail(id) {
  return request({
    url: `/api/v1/repos/${id}`,
    method: 'get',
  })
}

// 创建仓库
export function createRepo(data) {
  return request({
    url: '/api/v1/repos',
    method: 'post',
    data,
  })
}

// 更新仓库
export function updateRepo(id, data) {
  return request({
    url: `/api/v1/repos/${id}`,
    method: 'put',
    data,
  })
}

// 删除仓库
export function deleteRepo(id) {
  return request({
    url: `/api/v1/repos/${id}`,
    method: 'delete',
  })
}

// 测试Webhook
export function testWebhook(id) {
  return request({
    url: `/api/v1/repos/${id}/test`,
    method: 'post',
  })
}
```

---

## 五、组件设计

### 5.1 公共组件

| 组件名 | 说明 | 参数 |
|--------|------|------|
| PageHeader | 页面标题栏 | title, description, actions |
| PageContent | 页面内容容器 | - |
| SearchForm | 搜索表单 | fields, onSearch |
| DataTable | 数据表格 | columns, data, pagination, loading |
| EmptyData | 空数据展示 | message, image |
| BaseForm | 基础表单 | schema, model, rules |
| BaseInput | 输入框 | v-model, placeholder, disabled |
| BaseSelect | 选择器 | v-model, options, multiple |
| BaseSwitch | 开关 | v-model, label |
| BaseUpload | 上传组件 | accept, maxSize, maxCount |

### 5.2 布局组件

| 组件名 | 说明 | 参数 |
|--------|------|------|
| MainLayout | 主布局 | - |
| Sidebar | 侧边导航 | menu, collapsed |
| Header | 顶部栏 | userInfo, onLogout |
| Footer | 底部 | - |

---

## 六、状态管理

### 6.1 Pinia Store

#### 6.1.1 User Store

```typescript
// src/stores/user.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, logout, getUserInfo } from '@/services/auth'

export const useUserStore = defineStore('user', () => {
  // State
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)
  const roles = ref<string[]>([])

  // Getters
  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => roles.value.includes('admin'))

  // Actions
  async function doLogin(username, password) {
    const data = await login({ username, password })
    token.value = data.access_token
    localStorage.setItem('token', data.access_token)
    await getUserProfile()
  }

  async function getUserProfile() {
    const info = await getUserInfo()
    userInfo.value = info
    roles.value = [info.role]
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    roles.value = []
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    roles,
    isLoggedIn,
    isAdmin,
    doLogin,
    getUserProfile,
    logout,
  }
}, {
  persist: true,
})
```

#### 6.1.2 App Store

```typescript
// src/stores/app.js
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false)
  const theme = ref(localStorage.getItem('theme') || 'light')
  const language = ref(localStorage.getItem('language') || 'zh-CN')

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setTheme(newTheme) {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    document.documentElement.className = newTheme
  }

  function setLanguage(lang) {
    language.value = lang
    localStorage.setItem('language', lang)
  }

  return {
    sidebarCollapsed,
    theme,
    language,
    toggleSidebar,
    setTheme,
    setLanguage,
  }
}, {
  persist: true,
})
```

---

## 七、路由配置

### 7.1 路由定义

```typescript
// src/router/routes.js
import { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', public: true },
  },
  {
    path: '/',
    component: () => import('@/components/layout/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '工作台', icon: 'DashboardOutlined' },
      },
      {
        path: 'repos',
        name: 'Repos',
        component: () => import('@/views/repos/index.vue'),
        meta: { title: '仓库管理', icon: 'GitBranchOutlined' },
      },
      {
        path: 'targets',
        name: 'Targets',
        component: () => import('@/views/targets/index.vue'),
        meta: { title: '推送目标', icon: 'NotificationOutlined' },
      },
      {
        path: 'pushes',
        name: 'Pushes',
        component: () => import('@/views/pushes/index.vue'),
        meta: { title: '推送记录', icon: 'SendOutlined' },
      },
      {
        path: 'templates',
        name: 'Templates',
        component: () => import('@/views/templates/index.vue'),
        meta: { title: '消息模板', icon: 'FileTextOutlined' },
      },
      {
        path: 'prompts',
        name: 'Prompts',
        component: () => import('@/views/prompts/index.vue'),
        meta: { title: '提示词', icon: 'BulbOutlined' },
      },
      {
        path: 'models',
        name: 'Models',
        component: () => import('@/views/models/index.vue'),
        meta: { title: 'AI模型', icon: 'RobotOutlined' },
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/users/index.vue'),
        meta: { title: '用户管理', icon: 'UserOutlined', roles: ['admin'] },
      },
      {
        path: 'logs',
        name: 'Logs',
        redirect: '/logs/system',
        meta: { title: '日志管理', icon: 'FileSearchOutlined' },
        children: [
          {
            path: 'system',
            name: 'SystemLog',
            component: () => import('@/views/logs/SystemLog.vue'),
            meta: { title: '系统日志' },
          },
          {
            path: 'operations',
            name: 'OperationLog',
            component: () => import('@/views/logs/OperationLog.vue'),
            meta: { title: '操作日志' },
          },
          {
            path: 'ai-calls',
            name: 'AICallLog',
            component: () => import('@/views/logs/AICallLog.vue'),
            meta: { title: 'AI调用日志' },
          },
        ],
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/settings/index.vue'),
        meta: { title: '个人设置' },
      },
    ],
  },
  {
    path: '/403',
    name: '403',
    component: () => import('@/views/error/403.vue'),
    meta: { title: '无权限', public: true },
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '页面不存在', public: true },
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404',
  },
]

export default routes
```

### 7.2 路由守卫

```typescript
// src/router/guard.js
import { Router } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { usePermissionStore } from '@/stores/permission'

export function createRouterGuard(router) {
  router.beforeEach(async (to, from, next) => {
    const userStore = useUserStore()
    const permissionStore = usePermissionStore()

    // 设置页面标题
    document.title = to.meta.title
      ? `${to.meta.title} - Push Notify`
      : 'Push Notify'

    // 公开页面直接放行
    if (to.meta.public) {
      next()
      return
    }

    // 未登录跳转登录页
    if (!userStore.isLoggedIn) {
      next({ name: 'Login', query: { redirect: to.fullPath } })
      return
    }

    // 获取用户信息
    if (!userStore.userInfo) {
      await userStore.getUserProfile()
    }

    // 权限检查
    if (to.meta.roles && !to.meta.roles.includes(userStore.roles[0])) {
      next({ name: '403' })
      return
    }

    next()
  })

  router.afterEach(() => {
    // 滚动到顶部
    window.scrollTo(0, 0)
  })
}
```

---

## 八、里程碑

| 阶段 | 内容 | 工期 | 预计完成 |
|------|------|------|----------|
| 第一阶段 | 项目初始化 | 1周 | 第1周末 |
| 第二阶段 | 登录认证 | 1周 | 第2周末 |
| 第三阶段 | 工作台和用户模块 | 1周 | 第3周末 |
| 第四阶段 | 仓库管理模块 | 1周 | 第4周末 |
| 第五阶段 | 推送目标模块 | 1周 | 第5周末 |
| 第六阶段 | 推送记录模块 | 1周 | 第6周末 |
| 第七阶段 | 消息模板模块 | 1周 | 第7周末 |
| 第八阶段 | 提示词模块 | 1周 | 第8周末 |
| 第九阶段 | AI模型模块 | 0.5周 | 第8周末 |
| 第十阶段 | 用户管理模块 | 0.5周 | 第9周末 |
| 第十一阶段 | 日志管理模块 | 1周 | 第9周末 |
| 第十二阶段 | 测试和优化 | 1周 | 第10周末 |

**预计总工期**: 10周

---

## 九、开发规范

### 9.1 代码规范

- 遵循Vue 3 Composition API规范
- 使用TypeScript进行类型检查
- 组件命名采用PascalCase
- 组合式函数命名采用camelCase，以use开头

### 9.2 Git规范

- 分支命名：feature/*, bugfix/*, hotfix/*
- 提交信息：feat/fix/docs/refactor: 描述
- 代码合并前必须Code Review

### 9.3 组件规范

- 组件文件大小不超过200行
- 复杂逻辑提取为组合式函数
- 通用组件放入components/common目录

---

**文档结束**
