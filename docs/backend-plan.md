# Push Notify 后端开发计划

> **文档版本**: V1.0  
> **创建日期**: 2026年1月19日  
> **技术栈**: Golang + Gin + GORM + SQLite + Redis

---

## 一、技术架构

### 1.1 技术选型

| 组件 | 技术 | 版本要求 | 说明 |
|------|------|---------|------|
| 开发语言 | Golang | >= 1.20 | 主要开发语言 |
| Web框架 | Gin | latest | HTTP路由和中间件 |
| ORM | GORM | latest | 数据库操作 |
| 数据库 | SQLite | - | 使用 glebarez/sqlite 驱动 |
| 缓存 | Redis | >= 6.0 | 会话缓存和热点数据 |
| 配置管理 | Viper | latest | 配置文件管理 |
| 日志 | Zap | latest | 结构化日志 |
| JWT | gol-jwt | latest | Token认证 |
| API文档 | Swag | latest | Swagger文档生成 |

### 1.2 项目结构

```
backend/
├── cmd/                        # 应用入口
│   └── server/                 # 服务启动入口
│       └── main.go
├── internal/                   # 内部包
│   ├── config/                 # 配置模块
│   │   └── config.go
│   ├── middleware/             # 中间件
│   │   ├── auth.go             # 认证中间件
│   │   ├── logger.go           # 日志中间件
│   │   ├── rate_limit.go       # 限流中间件
│   │   └── cors.go             # 跨域中间件
│   ├── models/                 # 数据模型
│   │   ├── user.go
│   │   ├── repo.go
│   │   ├── target.go
│   │   ├── push.go
│   │   ├── template.go
│   │   ├── prompt.go
│   │   ├── model.go
│   │   └── log.go
│   ├── repository/             # 数据访问层
│   │   ├── user_repo.go
│   │   ├── repo_repo.go
│   │   ├── target_repo.go
│   │   ├── push_repo.go
│   │   ├── template_repo.go
│   │   ├── prompt_repo.go
│   │   ├── model_repo.go
│   │   └── log_repo.go
│   ├── services/               # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── repo_service.go
│   │   ├── target_service.go
│   │   ├── push_service.go
│   │   ├── template_service.go
│   │   ├── prompt_service.go
│   │   ├── model_service.go
│   │   ├── log_service.go
│   │   └── webhook_service.go
│   ├── handlers/               # 处理器（API层）
│   │   ├── auth_handler.go
│   │   ├── repo_handler.go
│   │   ├── target_handler.go
│   │   ├── push_handler.go
│   │   ├── template_handler.go
│   │   ├── prompt_handler.go
│   │   ├── model_handler.go
│   │   ├── log_handler.go
│   │   └── webhook_handler.go
│   └── utils/                  # 工具函数
│       ├── jwt.go
│       ├── password.go
│       ├── encrypt.go
│       └── response.go
├── pkg/                        # 外部包
│   ├── dingtalk/               # 钉钉推送
│   ├── email/                  # 邮箱推送
│   └── ai/                     # AI服务调用
├── router/                     # 路由配置
│   └── router.go
├── database/                   # 数据库相关
│   └── database.go
├── migrations/                 # 数据库迁移
│   └── migrations.go
├── docs/                       # API文档
├── config.yaml                 # 配置文件
├── go.mod
└── go.sum
```

### 1.3 数据库设计

#### 1.3.1 ER图

```
users (用户表)
  ├── id (PK)
  ├── username (唯一)
  ├── email (唯一)
  ├── password
  ├── role
  ├── status
  └── timestamps

repos (仓库表)
  ├── id (PK)
  ├── name (唯一)
  ├── url
  ├── type
  ├── access_token
  ├── webhook_url
  ├── webhook_secret
  ├── model_id (FK)
  └── timestamps

targets (推送目标表)
  ├── id (PK)
  ├── name
  ├── type
  ├── config (JSON)
  ├── scope
  └── timestamps

repos_targets (仓库-推送目标关联表)
  ├── id (PK)
  ├── repo_id (FK)
  ├── target_id (FK)

pushes (推送记录表)
  ├── id (PK)
  ├── repo_id (FK)
  ├── target_id (FK)
  ├── template_id (FK)
  ├── commit_id
  ├── commit_msg
  ├── status
  ├── content
  ├── error_msg
  ├── retry_count
  └── timestamps

templates (消息模板表)
  ├── id (PK)
  ├── name
  ├── type
  ├── scene
  ├── title
  ├── content
  ├── is_default
  ├── version
  └── timestamps

prompts (提示词表)
  ├── id (PK)
  ├── name
  ├── type
  ├── scene
  ├── language
  ├── content
  ├── model_id (FK)
  ├── version
  └── timestamps

models (AI模型表)
  ├── id (PK)
  ├── name
  ├── type
  ├── api_url
  ├── api_key
  ├── params (JSON)
  ├── is_default
  ├── call_count
  └── timestamps

logs (日志表)
  ├── id (PK)
  ├── type
  ├── level
  ├── module
  ├── message
  ├── detail (JSON)
  ├── user_id (FK)
  ├── request_id
  └── timestamps
```

#### 1.3.2 数据库初始化脚本

```sql
-- 用户表
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    status VARCHAR(20) DEFAULT 'active',
    last_login_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 仓库表
CREATE TABLE repos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) UNIQUE NOT NULL,
    url VARCHAR(500) NOT NULL,
    type VARCHAR(50) NOT NULL,
    access_token VARCHAR(255),
    webhook_url VARCHAR(500) NOT NULL,
    webhook_secret VARCHAR(100),
    model_id INTEGER,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (model_id) REFERENCES models(id)
);

-- 推送目标表
CREATE TABLE targets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    config TEXT NOT NULL,
    scope VARCHAR(20) DEFAULT 'global',
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 仓库-推送目标关联表
CREATE TABLE repo_targets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_id INTEGER NOT NULL,
    target_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(repo_id, target_id),
    FOREIGN KEY (repo_id) REFERENCES repos(id) ON DELETE CASCADE,
    FOREIGN KEY (target_id) REFERENCES targets(id) ON DELETE CASCADE
);

-- 推送记录表
CREATE TABLE pushes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_id INTEGER NOT NULL,
    target_id INTEGER NOT NULL,
    template_id INTEGER,
    commit_id VARCHAR(50) NOT NULL,
    commit_msg VARCHAR(500) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    content TEXT NOT NULL,
    error_msg TEXT,
    retry_count INTEGER DEFAULT 0,
    pushed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (repo_id) REFERENCES repos(id),
    FOREIGN KEY (target_id) REFERENCES targets(id),
    FOREIGN KEY (template_id) REFERENCES templates(id)
);

-- 消息模板表
CREATE TABLE templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    scene VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'active',
    version INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 提示词表
CREATE TABLE prompts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    scene VARCHAR(50),
    language VARCHAR(50),
    content TEXT NOT NULL,
    model_id INTEGER,
    version INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (model_id) REFERENCES models(id)
);

-- AI模型表
CREATE TABLE models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    api_url VARCHAR(500) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    params TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    call_count INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 日志表
CREATE TABLE logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type VARCHAR(20) NOT NULL,
    level VARCHAR(20) NOT NULL,
    module VARCHAR(50),
    message TEXT NOT NULL,
    detail TEXT,
    user_id INTEGER,
    request_id VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 索引
CREATE INDEX idx_repos_status ON repos(status);
CREATE INDEX idx_targets_type ON targets(type);
CREATE INDEX idx_pushes_status ON pushes(status);
CREATE INDEX idx_pushes_repo_id ON pushes(repo_id);
CREATE INDEX idx_pushes_created_at ON pushes(created_at);
CREATE INDEX idx_logs_type ON logs(type);
CREATE INDEX idx_logs_created_at ON logs(created_at);
```

---

## 二、开发计划

### 第一阶段：项目基础架构（第1-2周）

#### 2.1.1 项目初始化

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 创建项目结构 | 初始化Gin项目，配置Go Module | 0.5天 | - |
| 配置管理 | 使用Viper实现配置文件管理 | 0.5天 | - |
| 数据库连接 | 配置SQLite和GORM连接 | 0.5天 | - |
| Redis连接 | 实现Redis客户端封装 | 0.5天 | 项目结构 |
| 日志系统 | 集成Zap日志库 | 0.5天 | - |

#### 2.1.2 基础中间件

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 跨域中间件 | CORS配置 | 0.5天 | - |
| 日志中间件 | 请求日志记录 | 0.5天 | 日志系统 |
| 认证中间件 | JWT Token验证 | 1天 | JWT工具 |
| 限流中间件 | 请求频率限制 | 0.5天 | Redis连接 |
| 错误处理 | 统一错误响应 | 0.5天 | - |

#### 2.1.3 交付物

- [ ] 完整的项目骨架
- [ ] 配置文件（config.yaml）
- [ ] 数据库连接和初始化
- [ ] 基础中间件
- [ ] 单元测试覆盖核心工具函数

### 第二阶段：用户认证模块（第3周）

#### 2.2.1 用户注册

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 数据模型 | User模型定义 | 0.5天 | - |
| 密码加密 | bcrypt密码加密工具 | 0.5天 | - |
| Repository层 | User数据访问实现 | 0.5天 | 数据模型 |
| Service层 | 注册业务逻辑 | 0.5天 | Repository |
| Handler层 | 注册API接口 | 0.5天 | Service |
| 邮箱验证 | 发送验证邮件 | 1天 | Service |

#### 2.2.2 用户登录

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| JWT工具 | Token生成和验证 | 0.5天 | - |
| Repository层 | 用户查询实现 | 0.5天 | 数据模型 |
| Service层 | 登录业务逻辑 | 0.5天 | Repository |
| Handler层 | 登录API接口 | 0.5天 | Service |
| 登录锁定 | 失败次数限制 | 0.5天 | Redis连接 |
| Token刷新 | 刷新Token接口 | 0.5天 | JWT工具 |

#### 2.2.3 用户管理

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 用户列表 | 分页查询接口 | 1天 | Repository |
| 用户创建 | 管理员创建用户 | 0.5天 | Service |
| 用户更新 | 信息编辑接口 | 0.5天 | Service |
| 用户删除 | 软删除实现 | 0.5天 | Service |
| 密码管理 | 修改/重置密码 | 1天 | 密码加密 |
| 个人设置 | 获取/更新个人设置 | 0.5天 | Service |

#### 2.2.4 交付物

- [ ] 用户注册/登录接口（/api/v1/auth/*）
- [ ] 用户管理接口（/api/v1/users/*）
- [ ] JWT认证流程
- [ ] 单元测试覆盖

### 第三阶段：仓库管理模块（第4周）

#### 2.3.1 仓库CRUD

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 数据模型 | Repo模型定义 | 0.5天 | - |
| Repository层 | Repo数据访问 | 0.5天 | 数据模型 |
| Service层 | 仓库业务逻辑 | 1天 | Repository |
| Handler层 | 仓库API接口 | 1天 | Service |
| Webhook生成 | 生成唯一Webhook URL | 0.5天 | - |
| URL验证 | Git仓库地址格式验证 | 0.5天 | - |

#### 2.3.2 Webhook管理

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| Webhook测试 | 测试接口实现 | 0.5天 | Service |
| 仓库关联目标 | 关联/取消关联接口 | 1天 | Repository |
| 仓库详情 | 获取仓库完整信息 | 0.5天 | Service |

#### 2.3.3 交付物

- [ ] 仓库CRUD接口（/api/v1/repos/*）
- [ ] Webhook配置生成
- [ ] 仓库-推送目标关联
- [ ] 单元测试覆盖

### 第四阶段：推送目标模块（第5周）

#### 2.4.1 钉钉推送

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 数据模型 | Target模型定义 | 0.5天 | - |
| 钉钉SDK | 钉钉消息推送封装 | 1天 | - |
| Token验证 | 钉钉AccessToken验证 | 0.5天 | 钉钉SDK |
| Repository层 | Target数据访问 | 0.5天 | 数据模型 |
| Service层 | 钉钉推送逻辑 | 1天 | 钉钉SDK |
| Handler层 | 钉钉目标API | 1天 | Service |

#### 2.4.2 邮箱推送

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 邮箱SDK | SMTP邮件发送封装 | 1天 | - |
| SMTP验证 | SMTP配置验证 | 0.5天 | 邮箱SDK |
| Service层 | 邮箱推送逻辑 | 1天 | 邮箱SDK |
| Handler层 | 邮箱目标API | 1天 | Service |

#### 2.4.3 推送目标管理

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 测试推送 | 发送测试消息 | 0.5天 | Service |
| 目标列表 | 分页查询接口 | 0.5天 | Repository |
| 目标详情 | 获取完整信息 | 0.5天 | Service |

#### 2.4.4 交付物

- [ ] 钉钉推送集成
- [ ] 邮箱推送集成
- [ ] 推送目标管理接口（/api/v1/targets/*）
- [ ] 单元测试覆盖

### 第五阶段：消息模板模块（第6周）

#### 2.5.1 模板管理

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 数据模型 | Template模型定义 | 0.5天 | - |
| 模板变量 | 变量替换引擎 | 1天 | - |
| Repository层 | Template数据访问 | 0.5天 | 数据模型 |
| Service层 | 模板业务逻辑 | 1天 | Repository |
| Handler层 | 模板API接口 | 1天 | Service |
| 版本管理 | 模板版本控制 | 1天 | Service |

#### 2.5.2 模板功能

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 模板预览 | 渲染预览接口 | 0.5天 | 模板变量 |
| 模板测试 | 发送测试消息 | 0.5天 | Service |
| 启用禁用 | 状态管理 | 0.5天 | Service |

#### 2.5.3 交付物

- [ ] 消息模板CRUD接口（/api/v1/templates/*）
- [ ] 模板变量引擎
- [ ] 版本管理
- [ ] 单元测试覆盖

### 第六阶段：提示词管理模块（第6-7周）

#### 2.6.1 提示词管理

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 数据模型 | Prompt模型定义 | 0.5天 | - |
| Repository层 | Prompt数据访问 | 0.5天 | 数据模型 |
| Service层 | 提示词业务逻辑 | 1天 | Repository |
| Handler层 | 提示词API接口 | 1天 | Service |
| 版本管理 | 提示词版本控制 | 1天 | Service |

#### 2.6.2 提示词功能

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 提示词测试 | 测试接口 | 0.5天 | Service |
| 导入导出 | JSON格式导入导出 | 1天 | Service |

#### 2.6.3 交付物

- [ ] 提示词CRUD接口（/api/v1/prompts/*）
- [ ] 版本管理
- [ ] 导入导出功能
- [ ] 单元测试覆盖

### 第七阶段：AI模型管理模块（第7周）

#### 2.7.1 模型配置

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 数据模型 | Model模型定义 | 0.5天 | - |
| Repository层 | Model数据访问 | 0.5天 | 数据模型 |
| Service层 | 模型业务逻辑 | 1天 | Repository |
| Handler层 | 模型API接口 | 1天 | Service |

#### 2.7.2 AI服务集成

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| AI SDK | OpenAI兼容API封装 | 1.5天 | - |
| 模型验证 | API连接测试 | 0.5天 | AI SDK |
| 调用日志 | 记录AI调用日志 | 1天 | Service |
| 成本统计 | 调用次数统计 | 0.5天 | Service |

#### 2.7.3 交付物

- [ ] AI模型管理接口（/api/v1/models/*）
- [ ] OpenAI兼容API调用
- [ ] 调用日志记录
- [ ] 单元测试覆盖

### 第八阶段：推送记录模块（第8周）

#### 2.8.1 推送记录

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 数据模型 | Push模型定义 | 0.5天 | - |
| Repository层 | Push数据访问 | 0.5天 | 数据模型 |
| Service层 | 推送业务逻辑 | 1天 | Repository |
| Handler层 | 推送记录API | 1天 | Service |

#### 2.8.2 推送功能

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 推送重试 | 重试失败推送 | 1天 | Service |
| 批量操作 | 批量重试/删除 | 1天 | Service |
| 统计接口 | 推送数据统计 | 0.5天 | Repository |

#### 2.8.3 交付物

- [ ] 推送记录查询接口（/api/v1/pushes/*）
- [ ] 重试机制
- [ ] 批量操作
- [ ] 统计数据接口

### 第九阶段：Webhook和CODEVIEW模块（第8-9周）

#### 2.9.1 Webhook接收

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| Webhook接口 | 接收仓库回调 | 1天 | - |
| 事件解析 | GitHub/GitLab/Gitee事件解析 | 1.5天 | Webhook接口 |
| 签名验证 | Webhook安全验证 | 0.5天 | - |
| 代码获取 | 获取提交代码差异 | 1天 | Repo配置 |

#### 2.9.2 CODEVIEW处理

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 提示词选择 | 选择对应提示词 | 0.5天 | Prompt服务 |
| AI调用 | 调用AI模型审查代码 | 1天 | AI SDK |
| 结果处理 | 解析AI返回结果 | 0.5天 | AI调用 |
| 模板渲染 | 生成推送消息 | 1天 | 模板服务 |

#### 2.9.3 推送执行

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 目标查找 | 查找关联推送目标 | 0.5天 | Target服务 |
| 消息推送 | 执行推送操作 | 1天 | 钉钉/邮箱SDK |
| 状态更新 | 更新推送状态 | 0.5天 | Push服务 |
| 重试机制 | 失败自动重试 | 1天 | Push服务 |

#### 2.9.4 交付物

- [ ] Webhook接收接口（/webhook/*）
- [ ] 多平台事件解析
- [ ] CODEVIEW自动执行
- [ ] 推送消息生成和发送
- [ ] 重试机制

### 第十阶段：日志管理模块（第9周）

#### 2.10.1 日志记录

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 数据模型 | Log模型定义 | 0.5天 | - |
| 系统日志 | 运行日志记录中间件 | 0.5天 | - |
| 操作日志 | 用户操作日志记录 | 1天 | Auth中间件 |
| AI调用日志 | AI调用记录 | 0.5天 | AI SDK |

#### 2.10.2 日志查询

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 日志列表 | 分页查询接口 | 1天 | Repository |
| 日志搜索 | 关键词搜索 | 0.5天 | Repository |
| 日志导出 | CSV/JSON导出 | 1天 | Service |
| 统计接口 | 日志统计分析 | 0.5天 | Service |

#### 2.10.3 交付物

- [ ] 日志管理接口（/api/v1/logs/*）
- [ ] 多类型日志记录
- [ ] 搜索和导出功能
- [ ] 统计接口

### 第十一阶段：测试和优化（第10周）

#### 2.11.1 测试

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 单元测试 | 核心业务逻辑测试 | 2天 | 各模块 |
| 集成测试 | API端到端测试 | 2天 | 单元测试 |
| 压力测试 | 并发性能测试 | 1天 | 集成测试 |

#### 2.11.2 优化

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| 数据库优化 | 查询优化和索引 | 1天 | - |
| 缓存优化 | Redis缓存策略 | 1天 | - |
| 日志优化 | 日志轮转和压缩 | 0.5天 | - |

#### 2.11.3 文档

| 任务 | 描述 | 工期 | 依赖 |
|------|------|------|------|
| API文档 | Swagger文档生成 | 0.5天 | - |
| 部署文档 | 部署和运维手册 | 0.5天 | - |

---

## 三、接口清单

### 3.1 认证接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 登录 | POST | /api/v1/auth/login | 用户登录 |
| 注册 | POST | /api/v1/auth/register | 用户注册 |
| 获取当前用户 | GET | /api/v1/auth/me | 获取登录用户信息 |
| 修改密码 | PUT | /api/v1/auth/password | 修改密码 |
| 刷新Token | POST | /api/v1/auth/refresh | 刷新访问令牌 |
| 登出 | POST | /api/v1/auth/logout | 用户登出 |

### 3.2 用户接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 用户列表 | GET | /api/v1/users | 获取用户列表 |
| 创建用户 | POST | /api/v1/users | 创建新用户 |
| 用户详情 | GET | /api/v1/users/:id | 获取用户详情 |
| 更新用户 | PUT | /api/v1/users/:id | 更新用户信息 |
| 删除用户 | DELETE | /api/v1/users/:id | 删除用户 |
| 重置密码 | POST | /api/v1/users/:id/reset-password | 重置用户密码 |
| 锁定用户 | PUT | /api/v1/users/:id/lock | 锁定/解锁用户 |
| 个人设置 | GET/PUT | /api/v1/settings | 获取/更新个人设置 |

### 3.3 仓库接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 仓库列表 | GET | /api/v1/repos | 获取仓库列表 |
| 创建仓库 | POST | /api/v1/repos | 创建仓库 |
| 仓库详情 | GET | /api/v1/repos/:id | 获取仓库详情 |
| 更新仓库 | PUT | /api/v1/repos/:id | 更新仓库 |
| 删除仓库 | DELETE | /api/v1/repos/:id | 删除仓库 |
| 测试Webhook | POST | /api/v1/repos/:id/test | 测试Webhook |
| 仓库目标 | GET/POST | /api/v1/repos/:id/targets | 获取/关联推送目标 |

### 3.4 推送目标接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 目标列表 | GET | /api/v1/targets | 获取推送目标列表 |
| 创建目标 | POST | /api/v1/targets | 创建推送目标 |
| 目标详情 | GET | /api/v1/targets/:id | 获取目标详情 |
| 更新目标 | PUT | /api/v1/targets/:id | 更新目标 |
| 删除目标 | DELETE | /api/v1/targets/:id | 删除目标 |
| 测试推送 | POST | /api/v1/targets/:id/test | 发送测试消息 |
| 关联仓库 | POST | /api/v1/targets/:id/repos | 关联仓库 |
| 取消关联 | DELETE | /api/v1/targets/:id/repos/:repoId | 取消关联 |

### 3.5 推送记录接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 推送列表 | GET | /api/v1/pushes | 获取推送记录列表 |
| 推送详情 | GET | /api/v1/pushes/:id | 获取推送详情 |
| 重试推送 | POST | /api/v1/pushes/:id/retry | 重试推送 |
| 批量重试 | POST | /api/v1/pushes/batch-retry | 批量重试 |
| 批量删除 | DELETE | /api/v1/pushes/batch-delete | 批量删除 |
| 推送统计 | GET | /api/v1/pushes/stats | 获取统计数据 |

### 3.6 消息模板接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 模板列表 | GET | /api/v1/templates | 获取模板列表 |
| 创建模板 | POST | /api/v1/templates | 创建模板 |
| 模板详情 | GET | /api/v1/templates/:id | 获取模板详情 |
| 更新模板 | PUT | /api/v1/templates/:id | 更新模板 |
| 删除模板 | DELETE | /api/v1/templates/:id | 删除模板 |
| 版本历史 | GET | /api/v1/templates/:id/versions | 获取版本历史 |
| 回滚版本 | POST | /api/v1/templates/:id/rollback | 回滚版本 |
| 模板预览 | POST | /api/v1/templates/preview | 预览渲染 |
| 模板测试 | POST | /api/v1/templates/:id/test | 测试模板 |
| 启用禁用 | PUT | /api/v1/templates/:id/status | 更新状态 |

### 3.7 提示词接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 提示词列表 | GET | /api/v1/prompts | 获取提示词列表 |
| 创建提示词 | POST | /api/v1/prompts | 创建提示词 |
| 提示词详情 | GET | /api/v1/prompts/:id | 获取详情 |
| 更新提示词 | PUT | /api/v1/prompts/:id | 更新提示词 |
| 删除提示词 | DELETE | /api/v1/prompts/:id | 删除提示词 |
| 版本历史 | GET | /api/v1/prompts/:id/versions | 获取版本历史 |
| 回滚版本 | POST | /api/v1/prompts/:id/rollback | 回滚版本 |
| 测试提示词 | POST | /api/v1/prompts/:id/test | 测试提示词 |
| 导出提示词 | GET | /api/v1/prompts/:id/export | 导出 |
| 导入提示词 | POST | /api/v1/prompts/import | 导入 |

### 3.8 AI模型接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 模型列表 | GET | /api/v1/models | 获取模型列表 |
| 创建模型 | POST | /api/v1/models | 创建模型 |
| 模型详情 | GET | /api/v1/models/:id | 获取模型详情 |
| 更新模型 | PUT | /api/v1/models/:id | 更新模型 |
| 删除模型 | DELETE | /api/v1/models/:id | 删除模型 |
| 默认模型 | POST | /api/v1/models/:id/default | 设置默认 |
| 调用日志 | GET | /api/v1/models/:id/logs | 获取日志 |
| 验证模型 | POST | /api/v1/models/:id/verify | 验证配置 |

### 3.9 日志接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 系统日志 | GET | /api/v1/logs/system | 获取系统日志 |
| 操作日志 | GET | /api/v1/logs/operations | 获取操作日志 |
| AI调用日志 | GET | /api/v1/logs/ai-calls | 获取AI日志 |
| 日志搜索 | GET | /api/v1/logs/search | 搜索日志 |
| 日志导出 | GET | /api/v1/logs/export | 导出日志 |
| 日志统计 | GET | /api/v1/logs/stats | 获取统计 |

### 3.10 Webhook接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| Webhook回调 | POST | /webhook/:webhookId | 接收仓库回调 |

---

## 四、开发规范

### 4.1 代码规范

- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 注释覆盖率不低于30%
- 复杂逻辑必须添加注释

### 4.2 Git规范

- 分支命名：feature/*, bugfix/*, hotfix/*
- 提交信息：feat/fix/docs/refactor: 描述
- 代码合并前必须Code Review

### 4.3 测试规范

- 核心业务逻辑单元测试覆盖率 >= 80%
- 集成测试覆盖所有API接口
- 提交前必须通过所有测试

---

## 五、里程碑

| 阶段 | 内容 | 工期 | 预计完成 |
|------|------|------|----------|
| 第一阶段 | 项目基础架构 | 2周 | 第2周末 |
| 第二阶段 | 用户认证模块 | 1周 | 第3周末 |
| 第三阶段 | 仓库管理模块 | 1周 | 第4周末 |
| 第四阶段 | 推送目标模块 | 1周 | 第5周末 |
| 第五阶段 | 消息模板模块 | 1周 | 第6周末 |
| 第六阶段 | 提示词管理模块 | 1周 | 第7周末 |
| 第七阶段 | AI模型管理模块 | 1周 | 第7周末 |
| 第八阶段 | 推送记录模块 | 1周 | 第8周末 |
| 第九阶段 | Webhook和CODEVIEW | 2周 | 第9周末 |
| 第十阶段 | 日志管理模块 | 1周 | 第9周末 |
| 第十一阶段 | 测试和优化 | 1周 | 第10周末 |

**预计总工期**: 10周

---

**文档结束**
