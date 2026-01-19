# Push Notify Backend

后端服务基于 Golang + Gin 框架开发，提供 RESTful API 接口。

## 技术栈

- **语言**: Golang 1.20+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: SQLite (glebarez/sqlite)
- **缓存**: Redis
- **配置管理**: Viper
- **日志**: Zap
- **JWT认证**: gol-jwt

## 项目结构

```
backend/
├── cmd/server/           # 应用入口
├── internal/
│   ├── config/           # 配置模块
│   ├── middleware/       # 中间件
│   ├── models/           # 数据模型
│   ├── repository/       # 数据访问层
│   ├── services/         # 业务逻辑层
│   ├── handlers/         # API处理器
│   └── utils/            # 工具函数
├── pkg/                  # 外部包
│   ├── dingtalk/         # 钉钉推送
│   ├── email/            # 邮箱推送
│   └── ai/               # AI服务调用
├── router/               # 路由配置
├── database/             # 数据库相关
├── config.yaml           # 配置文件
├── go.mod
└── go.sum
```

## 快速开始

### 1. 安装依赖

```bash
cd backend
go mod tidy
```

### 2. 配置数据库

修改 `config.yaml` 文件，配置数据库路径：

```yaml
database:
  driver: "sqlite"
  path: "data.db"
```

### 3. 运行服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

## API文档

### 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/auth/login | 用户登录 |
| POST | /api/v1/auth/register | 用户注册 |
| POST | /api/v1/auth/refresh | 刷新Token |
| GET | /api/v1/auth/me | 获取当前用户 |
| PUT | /api/v1/auth/password | 修改密码 |

### 仓库管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/repos | 获取仓库列表 |
| POST | /api/v1/repos | 创建仓库 |
| GET | /api/v1/repos/:id | 获取仓库详情 |
| PUT | /api/v1/repos/:id | 更新仓库 |
| DELETE | /api/v1/repos/:id | 删除仓库 |
| POST | /api/v1/repos/:id/test | 测试Webhook |

### 推送目标

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/targets | 获取推送目标列表 |
| POST | /api/v1/targets | 创建推送目标 |
| GET | /api/v1/targets/:id | 获取目标详情 |
| PUT | /api/v1/targets/:id | 更新目标 |
| DELETE | /api/v1/targets/:id | 删除目标 |
| POST | /api/v1/targets/:id/test | 测试推送 |

### 消息模板

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/templates | 获取模板列表 |
| POST | /api/v1/templates | 创建模板 |
| GET | /api/v1/templates/:id | 获取模板详情 |
| PUT | /api/v1/templates/:id | 更新模板 |
| DELETE | /api/v1/templates/:id | 删除模板 |

### 提示词

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/prompts | 获取提示词列表 |
| POST | /api/v1/prompts | 创建提示词 |
| GET | /api/v1/prompts/:id | 获取提示词详情 |
| PUT | /api/v1/prompts/:id | 更新提示词 |
| DELETE | /api/v1/prompts/:id | 删除提示词 |
| POST | /api/v1/prompts/:id/test | 测试提示词 |

### AI模型

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/models | 获取模型列表 |
| POST | /api/v1/models | 创建模型 |
| GET | /api/v1/models/:id | 获取模型详情 |
| PUT | /api/v1/models/:id | 更新模型 |
| DELETE | /api/v1/models/:id | 删除模型 |
| POST | /api/v1/models/:id/default | 设置默认模型 |

### 推送记录

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/pushes | 获取推送记录 |
| GET | /api/v1/pushes/:id | 获取推送详情 |
| POST | /api/v1/pushes/:id/retry | 重试推送 |
| POST | /api/v1/pushes/batch-retry | 批量重试 |
| GET | /api/v1/pushes/stats | 获取统计 |

### 日志管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/logs/system | 获取系统日志 |
| GET | /api/v1/logs/operations | 获取操作日志 |
| GET | /api/v1/logs/ai-calls | 获取AI调用日志 |
| GET | /api/v1/logs/export | 导出日志 |

## 配置说明

### config.yaml

```yaml
# 应用配置
app:
  host: "0.0.0.0"
  port: 8080
  name: "Push Notify"
  env: "development"

# 数据库配置
database:
  driver: "sqlite"
  path: "data.db"
  log_mode: "info"

# Redis配置
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10

# JWT配置
jwt:
  secret: "your-secret-key-change-in-production"
  access_token_expire: 86400
  refresh_token_expire: 604800

# 日志配置
logging:
  level: "info"
  format: "json"
  output: "stdout"

# AI服务配置
ai:
  default_model_id: 1
  timeout: 60

# Webhook配置
webhook:
  signing_key: "webhook-secret-key"
```

## 开发规范

### 代码规范

- 遵循 Go 官方代码规范
- 使用 gofmt 格式化代码
- 注释覆盖率不低于 30%
- 复杂逻辑必须添加注释

### Git规范

- 分支命名: feature/*, bugfix/*, hotfix/*
- 提交信息: feat/fix/docs/refactor: 描述
- 代码合并前必须 Code Review

### 测试规范

- 核心业务逻辑单元测试覆盖率 >= 80%
- 集成测试覆盖所有 API 接口
- 提交前必须通过所有测试

## License

MIT
