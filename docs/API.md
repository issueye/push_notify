# Push Notify 前后端对接文档

> **文档版本**: V1.0  
> **创建日期**: 2026年1月19日  
> **文档状态**: 初稿

---

## 目录

1. [文档概述](#1-文档概述)
2. [通用规范](#2-通用规范)
3. [认证授权](#3-认证授权)
4. [仓库管理模块](#4-仓库管理模块)
5. [推送目标管理模块](#5-推送目标管理模块)
6. [推送内容查询模块](#6-推送内容查询模块)
7. [AI模型管理模块](#7-ai模型管理模块)
8. [用户管理模块](#8-用户管理模块)
9. [日志管理模块](#9-日志管理模块)
10. [消息模板管理模块](#10-消息模板管理模块)
11. [提示词管理模块](#11-提示词管理模块)
12. [Webhook接口](#12-webhook接口)
13. [错误码说明](#13-错误码说明)

---

## 1. 文档概述

本文档描述 Push Notify 系统的 RESTful API 接口规范，是前端与后端开发对接的技术指南。文档涵盖所有业务模块的接口定义、请求参数、响应格式及错误处理说明。

### 1.1 接口基础信息

| 项目 | 内容 |
|------|------|
| Base URL | `/api/v1` |
| 数据格式 | JSON |
| 字符编码 | UTF-8 |
| 时区 | UTC（请求参数中的时间字段需使用UTC时间） |

### 1.2 接口分类

| 分类 | 前缀 | 说明 |
|------|------|------|
| 认证接口 | `/api/v1/auth/*` | 登录、注册、Token刷新等 |
| 用户接口 | `/api/v1/users/*` | 用户管理相关 |
| 仓库接口 | `/api/v1/repos/*` | 仓库管理相关 |
| 推送目标接口 | `/api/v1/targets/*` | 推送目标管理相关 |
| 推送记录接口 | `/api/v1/pushes/*` | 推送历史查询相关 |
| AI模型接口 | `/api/v1/models/*` | AI模型配置相关 |
| 模板接口 | `/api/v1/templates/*` | 消息模板相关 |
| 提示词接口 | `/api/v1/prompts/*` | 提示词管理相关 |
| 日志接口 | `/api/v1/logs/*` | 系统日志相关 |
| Webhook | `/webhook/*` | 代码仓库Webhook回调 |

---

## 2. 通用规范

### 2.1 请求格式

#### 2.1.1 请求头（Headers）

```http
Content-Type: application/json
Accept: application/json
Authorization: Bearer <token>
X-Request-Id: <uuid>
Accept-Language: zh-CN
```

| Header | 必填 | 说明 |
|--------|------|------|
| Content-Type | 是 | 固定为 `application/json` |
| Accept | 否 | 期望的响应内容类型，默认 `application/json` |
| Authorization | 是 | JWT Token，格式为 `Bearer <token>` |
| X-Request-Id | 否 | 请求唯一标识，用于链路追踪 |
| Accept-Language | 否 | 语言偏好，如 `zh-CN`、`en-US` |

#### 2.1.2 请求参数类型

| 类型 | 说明 | 示例 |
|------|------|------|
| Path | URL路径参数 | `/api/v1/users/:id` |
| Query | URL查询参数 | `/api/v1/users?page=1&size=10` |
| Body | 请求体参数 | POST请求的JSON数据 |

### 2.2 响应格式

#### 2.2.1 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": { },
  "request_id": "uuid-string"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码，200表示成功 |
| message | string | 状态描述信息 |
| data | object | 响应数据主体 |
| request_id | string | 请求ID，用于问题排查 |

#### 2.2.2 分页响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [ ],
    "pagination": {
      "page": 1,
      "size": 10,
      "total": 100,
      "total_pages": 10
    }
  },
  "request_id": "uuid-string"
}
```

#### 2.2.3 错误响应

```json
{
  "code": 400,
  "message": "参数错误",
  "details": [ "字段name不能为空" ],
  "request_id": "uuid-string"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 错误状态码 |
| message | string | 错误描述 |
| details | array | 详细错误信息（可选） |

### 2.3 分页参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 当前页码，从1开始 |
| size | int | 否 | 10 | 每页条数，最大100 |
| sort | string | 否 | created_at | 排序字段 |
| order | string | 否 | desc | 排序方向：asc/desc |

### 2.4 日期时间格式

| 类型 | 格式 | 示例 |
|------|------|------|
| 日期 | YYYY-MM-DD | 2026-01-19 |
| 时间 | HH:mm:ss | 14:30:00 |
| 日期时间 | RFC3339 | 2026-01-19T14:30:00Z |
| 时间戳 | Unix秒 | 1705677000 |

---

## 3. 认证授权

### 3.1 JWT Token机制

系统采用 JWT (JSON Web Token) 进行身份认证。用户登录成功后，后端返回 access_token，前端后续请求需在 Header 中携带该 Token。

#### Token结构

```
Header: { "alg": "HS256", "typ": "JWT" }
Payload: { "user_id": 1, "role": "admin", "exp": 1705677000 }
Signature: HMAC-SHA256(secret_key, base64url(header) + "." + base64url(payload))
```

#### Token有效期

| Token类型 | 有效期 | 说明 |
|----------|--------|------|
| access_token | 24小时 | 访问API的凭证 |
| refresh_token | 7天 | 刷新access_token的凭证 |

### 3.2 登录

**接口说明**: 用户登录获取Token

```http
POST /api/v1/auth/login
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |
| captcha | string | 否 | 验证码（连续失败5次后必填） |
| captcha_key | string | 否 | 验证码key |

**请求示例**

```json
{
  "username": "admin",
  "password": "password123"
}
```

**响应示例**

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "token_type": "Bearer",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin"
    }
  }
}
```

### 3.3 注册

**接口说明**: 新用户注册

```http
POST /api/v1/auth/register
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名，3-20位字母数字下划线 |
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码，最少8位包含大小写字母和数字 |
| confirm_password | string | 是 | 确认密码 |

**响应示例**

```json
{
  "code": 200,
  "message": "注册成功，请前往邮箱激活账号",
  "data": {
    "user_id": 2
  }
}
```

### 3.4 获取当前用户信息

**接口说明**: 获取已登录用户的详细信息

```http
GET /api/v1/auth/me
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin",
    "status": "active",
    "last_login_at": "2026-01-19T10:30:00Z",
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

### 3.5 修改密码

**接口说明**: 修改当前用户密码

```http
PUT /api/v1/auth/password
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| old_password | string | 是 | 原密码 |
| new_password | string | 是 | 新密码 |
| confirm_password | string | 是 | 确认新密码 |

### 3.6 刷新Token

**接口说明**: 使用refresh_token获取新的access_token

```http
POST /api/v1/auth/refresh
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| refresh_token | string | 是 | 刷新令牌 |

### 3.7 登出

**接口说明**: 用户登出，使当前Token失效

```http
POST /api/v1/auth/logout
```

---

## 4. 仓库管理模块

### 4.1 获取仓库列表

**接口说明**: 获取已配置的所有代码仓库

```http
GET /api/v1/repos
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| size | int | 否 | 每页条数，默认10 |
| keyword | string | 否 | 搜索关键词（仓库名称） |
| sort | string | 否 | 排序字段，默认created_at |
| order | string | 否 | 排序方向：asc/desc |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "后端服务",
        "url": "https://github.com/company/backend-service",
        "type": "github",
        "webhook_url": "https://api.push-notify.com/webhook/abc123",
        "target_count": 2,
        "status": "active",
        "created_at": "2026-01-01T00:00:00Z",
        "updated_at": "2026-01-19T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 10,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

### 4.2 获取仓库详情

**接口说明**: 获取单个仓库的详细信息

```http
GET /api/v1/repos/:id
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 仓库ID |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "后端服务",
    "url": "https://github.com/company/backend-service",
    "type": "github",
    "access_token": "******",
    "webhook_url": "https://api.push-notify.com/webhook/abc123",
    "webhook_secret": "secret123",
    "webhook_events": ["push", "merge_request"],
    "model_id": 1,
    "model_name": "GPT-4",
    "targets": [
      {
        "id": 1,
        "name": "开发群",
        "type": "dingtalk"
      }
    ],
    "push_stats": {
      "total": 150,
      "success": 145,
      "failed": 5
    },
    "status": "active",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-19T10:00:00Z"
  }
}
```

### 4.3 创建仓库

**接口说明**: 添加新的代码仓库配置

```http
POST /api/v1/repos
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 仓库名称，唯一 |
| url | string | 是 | 仓库地址（Git HTTPS/SSH地址） |
| type | string | 是 | 仓库类型：github/gitlab/gitee |
| access_token | string | 否 | 私有仓库访问令牌 |
| webhook_secret | string | 否 | Webhook密钥，用于签名验证 |
| model_id | int | 否 | 关联的AI模型ID |

**请求示例**

```json
{
  "name": "后端服务",
  "url": "https://github.com/company/backend-service",
  "type": "github",
  "access_token": "ghp_xxxxxxxxxxxx",
  "webhook_secret": "mysecret",
  "model_id": 1
}
```

**响应示例**

```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 1,
    "name": "后端服务",
    "webhook_url": "https://api.push-notify.com/webhook/abc123",
    "webhook_secret": "mysecret"
  }
}
```

### 4.4 更新仓库

**接口说明**: 更新已配置的仓库信息

```http
PUT /api/v1/repos/:id
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 仓库ID |

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 仓库名称 |
| url | string | 否 | 仓库地址 |
| type | string | 否 | 仓库类型 |
| access_token | string | 否 | 访问令牌 |
| webhook_secret | string | 否 | Webhook密钥 |
| model_id | int | 否 | 关联AI模型ID |
| status | string | 否 | 状态：active/inactive |

### 4.5 删除仓库

**接口说明**: 删除仓库及其关联数据

```http
DELETE /api/v1/repos/:id
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 仓库ID |

**响应示例**

```json
{
  "code": 200,
  "message": "删除成功"
}
```

### 4.6 测试Webhook

**接口说明**: 触发测试Webhook请求验证配置

```http
POST /api/v1/repos/:id/test
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 仓库ID |

**响应示例**

```json
{
  "code": 200,
  "message": "测试成功",
  "data": {
    "request_id": "test-uuid-123",
    "status": "success",
    "response": {
      "code": 200,
      "message": "Webhook测试成功"
    }
  }
}
```

### 4.7 获取仓库关联的推送目标

**接口说明**: 获取指定仓库已关联的推送目标列表

```http
GET /api/v1/repos/:id/targets
```

### 4.8 仓库关联推送目标

**接口说明**: 将推送目标关联到仓库

```http
POST /api/v1/repos/:id/targets
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_ids | array | 是 | 推送目标ID列表 |

---

## 5. 推送目标管理模块

### 5.1 获取推送目标列表

**接口说明**: 获取所有已配置的推送目标

```http
GET /api/v1/targets
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页条数 |
| keyword | string | 否 | 搜索关键词 |
| type | string | 否 | 筛选类型：dingtalk/email |
| scope | string | 否 | 筛选范围：global/repo |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "开发群",
        "type": "dingtalk",
        "scope": "repo",
        "status": "active",
        "push_count": 150,
        "success_rate": 98.5,
        "created_at": "2026-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 10,
      "total": 3,
      "total_pages": 1
    }
  }
}
```

### 5.2 获取推送目标详情

**接口说明**: 获取单个推送目标的详细信息

```http
GET /api/v1/targets/:id
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "开发群",
    "type": "dingtalk",
    "scope": "repo",
    "config": {
      "access_token": "dingtalk-token-xxx",
      "secret": "******"
    },
    "repos": [
      {
        "id": 1,
        "name": "后端服务"
      }
    ],
    "push_stats": {
      "total": 150,
      "success": 147,
      "failed": 3
    },
    "status": "active",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-19T10:00:00Z"
  }
}
```

### 5.3 创建钉钉推送目标

**接口说明**: 添加钉钉群机器人作为推送目标

```http
POST /api/v1/targets
```

**请求参数（钉钉类型）**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 目标名称，同类型下唯一 |
| type | string | 是 | 固定值：dingtalk |
| config | object | 是 | 钉钉配置 |
| config.access_token | string | 是 | 钉钉群机器人AccessToken |
| config.secret | string | 否 | 钉钉机器人Secret（签名验证用） |
| scope | string | 否 | 范围：global/repo，默认global |
| repo_ids | array | 否 | 关联的仓库ID列表（scope为repo时必填） |

**请求示例（钉钉）**

```json
{
  "name": "开发群",
  "type": "dingtalk",
  "config": {
    "access_token": "dingtalk-access-token-xxx",
    "secret": "SECxxxxxxxx"
  },
  "scope": "repo",
  "repo_ids": [1, 2]
}
```

### 5.4 创建邮箱推送目标

**接口说明**: 添加邮箱作为推送目标

```http
POST /api/v1/targets
```

**请求参数（邮箱类型）**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 目标名称 |
| type | string | 是 | 固定值：email |
| config | object | 是 | 邮箱配置 |
| config.smtp_host | string | 是 | SMTP服务器地址 |
| config.smtp_port | int | 是 | SMTP端口，如587 |
| config.from | string | 是 | 发件人邮箱 |
| config.password | string | 是 | 发件人密码或授权码 |
| config.to | array | 是 | 收件人邮箱列表 |
| scope | string | 否 | 范围：global/repo |
| repo_ids | array | 否 | 关联的仓库ID列表 |

**请求示例（邮箱）**

```json
{
  "name": "技术团队邮箱",
  "type": "email",
  "config": {
    "smtp_host": "smtp.example.com",
    "smtp_port": 587,
    "from": "notify@company.com",
    "password": "smtp-password",
    "to": ["dev@company.com", "tech-lead@company.com"]
  },
  "scope": "global"
}
```

### 5.5 更新推送目标

**接口说明**: 更新推送目标配置

```http
PUT /api/v1/targets/:id
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 目标名称 |
| config | object | 否 | 配置信息 |
| status | string | 否 | 状态：active/inactive |

### 5.6 删除推送目标

**接口说明**: 删除推送目标

```http
DELETE /api/v1/targets/:id
```

### 5.7 测试推送

**接口说明**: 向推送目标发送测试消息

```http
POST /api/v1/targets/:id/test
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| template_id | int | 否 | 使用的模板ID，不传则使用默认模板 |

**响应示例**

```json
{
  "code": 200,
  "message": "测试消息已发送",
  "data": {
    "push_id": 1001,
    "status": "success"
  }
}
```

### 5.8 关联仓库

**接口说明**: 将仓库关联到推送目标

```http
POST /api/v1/targets/:id/repos
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| repo_ids | array | 是 | 仓库ID列表 |

### 5.9 取消仓库关联

**接口说明**: 取消仓库与推送目标的关联

```http
DELETE /api/v1/targets/:id/repos/:repoId
```

---

## 6. 推送内容查询模块

### 6.1 获取推送记录列表

**接口说明**: 查询推送历史记录

```http
GET /api/v1/pushes
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页条数 |
| repo_id | int | 否 | 筛选仓库ID |
| target_id | int | 否 | 筛选推送目标ID |
| status | string | 否 | 筛选状态：success/failed/pending |
| start_time | string | 否 | 开始时间，RFC3339格式 |
| end_time | string | 否 | 结束时间，RFC3339格式 |
| keyword | string | 否 | 搜索关键词 |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1001,
        "repo_id": 1,
        "repo_name": "后端服务",
        "target_id": 1,
        "target_name": "开发群",
        "target_type": "dingtalk",
        "commit_id": "abc123def",
        "commit_msg": "feat: 新增用户登录功能",
        "status": "success",
        "pushed_at": "2026-01-19T14:30:00Z",
        "created_at": "2026-01-19T14:30:05Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 10,
      "total": 500,
      "total_pages": 50
    }
  }
}
```

### 6.2 获取推送详情

**接口说明**: 获取单条推送记录的详细信息

```http
GET /api/v1/pushes/:id
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1001,
    "repo": {
      "id": 1,
      "name": "后端服务",
      "url": "https://github.com/company/backend-service"
    },
    "target": {
      "id": 1,
      "name": "开发群",
      "type": "dingtalk"
    },
    "template": {
      "id": 1,
      "name": "代码提交通知"
    },
    "commit": {
      "id": "abc123def",
      "msg": "feat: 新增用户登录功能",
      "author": "zhangsan",
      "branch": "main",
      "changed_files": ["login.go", "auth.go"],
      "file_count": 2
    },
    "codeview": {
      "result": "通过",
      "issues": [],
      "summary": "代码审查通过，未发现问题"
    },
    "content": {
      "title": "代码提交通知",
      "body": "## 提交信息\n\n..."
    },
    "status": "success",
    "retry_count": 0,
    "error_msg": null,
    "pushed_at": "2026-01-19T14:30:00Z",
    "created_at": "2026-01-19T14:30:05Z"
  }
}
```

### 6.3 重试推送

**接口说明**: 重新发送失败的推送

```http
POST /api/v1/pushes/:id/retry
```

**响应示例**

```json
{
  "code": 200,
  "message": "重试任务已创建",
  "data": {
    "new_push_id": 1002,
    "original_push_id": 1001
  }
}
```

### 6.4 批量重试

**接口说明**: 批量重试失败的推送记录

```http
POST /api/v1/pushes/batch-retry
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| push_ids | array | 是 | 要重试的推送记录ID列表 |

**响应示例**

```json
{
  "code": 200,
  "message": "批量重试任务已创建",
  "data": {
    "total": 5,
    "success": 5,
    "failed": 0
  }
}
```

### 6.5 批量删除

**接口说明**: 批量删除推送记录

```http
DELETE /api/v1/pushes/batch-delete
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| push_ids | array | 是 | 要删除的推送记录ID列表 |
| before_date | string | 否 | 删除此日期之前的记录 |

### 6.6 获取推送统计

**接口说明**: 获取推送数据统计

```http
GET /api/v1/pushes/stats
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | string | 否 | 开始日期 YYYY-MM-DD |
| end_date | string | 否 | 结束日期 YYYY-MM-DD |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "today": {
      "total": 15,
      "success": 14,
      "failed": 1,
      "rate": 93.3
    },
    "this_week": {
      "total": 85,
      "success": 82,
      "failed": 3,
      "rate": 96.5
    },
    "this_month": {
      "total": 320,
      "success": 310,
      "failed": 10,
      "rate": 96.9
    },
    "trend": [
      { "date": "2026-01-13", "total": 12, "success": 12, "failed": 0 },
      { "date": "2026-01-14", "total": 18, "success": 17, "failed": 1 }
    ]
  }
}
```

---

## 7. AI模型管理模块

### 7.1 获取模型列表

**接口说明**: 获取所有已配置的AI模型

```http
GET /api/v1/models
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页条数 |
| keyword | string | 否 | 搜索关键词 |
| provider | string | 否 | 按提供商筛选 |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "GPT-4",
        "type": "gpt-4",
        "provider": "OpenAI",
        "api_url": "https://api.openai.com/v1/chat/completions",
        "is_default": true,
        "call_count": 1500,
        "status": "active",
        "created_at": "2026-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 10,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

### 7.2 获取模型详情

**接口说明**: 获取AI模型的详细信息

```http
GET /api/v1/models/:id
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "GPT-4",
    "type": "gpt-4",
    "provider": "OpenAI",
    "api_url": "https://api.openai.com/v1/chat/completions",
    "api_key": "******",
    "params": {
      "temperature": 0.3,
      "max_tokens": 4000,
      "top_p": 0.9,
      "presence_penalty": 0,
      "frequency_penalty": 0
    },
    "is_default": true,
    "call_count": 1500,
    "success_count": 1495,
    "failed_count": 5,
    "avg_response_time": 2.5,
    "status": "active",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-19T10:00:00Z"
  }
}
```

### 7.3 创建模型

**接口说明**: 添加新的AI模型配置

```http
POST /api/v1/models
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 模型名称 |
| type | string | 是 | 模型类型 |
| provider | string | 否 | 提供商名称 |
| api_url | string | 是 | API地址 |
| api_key | string | 是 | API密钥 |
| timeout | int | 否 | 超时时间（秒），默认60 |
| params | object | 否 | 模型参数配置 |

**请求示例**

```json
{
  "name": "GPT-4",
  "type": "gpt-4",
  "provider": "OpenAI",
  "api_url": "https://api.openai.com/v1/chat/completions",
  "api_key": "sk-xxxxxxxxxxxx",
  "timeout": 60,
  "params": {
    "temperature": 0.3,
    "max_tokens": 4000,
    "top_p": 0.9
  }
}
```

### 7.4 更新模型

**接口说明**: 更新AI模型配置

```http
PUT /api/v1/models/:id
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 模型名称 |
| api_url | string | 否 | API地址 |
| api_key | string | 否 | API密钥 |
| timeout | int | 否 | 超时时间 |
| params | object | 否 | 模型参数 |
| status | string | 否 | 状态 |

### 7.5 删除模型

**接口说明**: 删除AI模型配置

```http
DELETE /api/v1/models/:id
```

### 7.6 设置默认模型

**接口说明**: 设置指定模型为默认模型

```http
POST /api/v1/models/:id/default
```

### 7.7 获取调用日志

**接口说明**: 获取模型的调用日志

```http
GET /api/v1/models/:id/logs
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页条数 |
| status | string | 否 | 筛选状态：success/failed |
| start_time | string | 否 | 开始时间 |
| end_time | string | 否 | 结束时间 |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 5001,
        "model_id": 1,
        "repo_id": 1,
        "input_summary": "CODEVIEW请求：main.go",
        "output_summary": "代码审查通过...",
        "status": "success",
        "duration_ms": 2500,
        "tokens_used": 1500,
        "created_at": "2026-01-19T14:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 10,
      "total": 1500,
      "total_pages": 150
    }
  }
}
```

### 7.8 验证模型配置

**接口说明**: 验证模型API配置是否有效

```http
POST /api/v1/models/:id/verify
```

**响应示例**

```json
{
  "code": 200,
  "message": "验证成功",
  "data": {
    "valid": true,
    "response_time_ms": 520,
    "model_info": {
      "id": "gpt-4",
      "object": "model",
      "created": 1688815741,
      "owned_by": "openai"
    }
  }
}
```

---

## 8. 用户管理模块

### 8.1 获取用户列表

**接口说明**: 获取所有用户（管理员功能）

```http
GET /api/v1/users
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页条数 |
| keyword | string | 否 | 搜索用户名或邮箱 |
| role | string | 否 | 筛选角色：admin/user |
| status | string | 否 | 筛选状态：active/locked |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "role": "admin",
        "status": "active",
        "last_login_at": "2026-01-19T10:30:00Z",
        "created_at": "2026-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 10,
      "total": 20,
      "total_pages": 2
    }
  }
}
```

### 8.2 获取用户详情

**接口说明**: 获取指定用户的详细信息

```http
GET /api/v1/users/:id
```

### 8.3 创建用户

**接口说明**: 管理员创建新用户

```http
POST /api/v1/users
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| email | string | 是 | 邮箱 |
| password | string | 是 | 初始密码 |
| role | string | 是 | 角色：admin/user |

### 8.4 更新用户

**接口说明**: 更新用户信息

```http
PUT /api/v1/users/:id
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 否 | 邮箱 |
| role | string | 否 | 角色 |
| status | string | 否 | 状态 |

### 8.5 删除用户

**接口说明**: 删除用户（软删除）

```http
DELETE /api/v1/users/:id
```

### 8.6 重置用户密码

**接口说明**: 管理员重置用户密码

```http
POST /api/v1/users/:id/reset-password
```

**响应示例**

```json
{
  "code": 200,
  "message": "密码已重置",
  "data": {
    "new_password": "TempP@ss123"
  }
}
```

### 8.7 锁定/解锁用户

**接口说明**: 管理员锁定或解锁用户账户

```http
PUT /api/v1/users/:id/lock
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| action | string | 是 | 操作：lock/unlock |

### 8.8 获取个人设置

**接口说明**: 获取当前用户的个人设置

```http
GET /api/v1/settings
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "notify": {
      "channels": ["dingtalk", "email"],
      "quiet_hours": {
        "enabled": true,
        "start": "22:00",
        "end": "08:00"
      }
    },
    "preferences": {
      "language": "zh-CN",
      "timezone": "Asia/Shanghai",
      "theme": "light"
    }
  }
}
```

### 8.9 更新个人设置

**接口说明**: 更新当前用户的个人设置

```http
PUT /api/v1/settings
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| notify | object | 否 | 通知偏好设置 |
| preferences | object | 否 | 界面偏好设置 |

---

## 9. 日志管理模块

### 9.1 获取系统运行日志

**接口说明**: 获取系统运行日志

```http
GET /api/v1/logs/system
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页条数 |
| level | string | 否 | 日志级别：debug/info/warn/error |
| keyword | string | 否 | 搜索关键词 |
| start_time | string | 否 | 开始时间 |
| end_time | string | 否 | 结束时间 |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 10001,
        "level": "INFO",
        "module": "repository",
        "message": "Webhook请求处理成功",
        "request_id": "req-uuid-123",
        "client_ip": "192.168.1.100",
        "created_at": "2026-01-19T14:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 50,
      "total": 1000,
      "total_pages": 20
    }
  }
}
```

### 9.2 获取用户操作日志

**接口说明**: 获取用户操作日志

```http
GET /api/v1/logs/operations
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页条数 |
| user_id | int | 否 | 筛选用户ID |
| action | string | 否 | 筛选操作类型 |
| start_time | string | 否 | 开始时间 |
| end_time | string | 否 | 结束时间 |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 2001,
        "user_id": 1,
        "username": "admin",
        "action": "repo.create",
        "object_type": "repo",
        "object_id": 1,
        "object_name": "后端服务",
        "detail": "创建仓库：后端服务",
        "client_ip": "192.168.1.100",
        "user_agent": "Mozilla/5.0...",
        "created_at": "2026-01-19T14:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 50,
      "total": 500,
      "total_pages": 10
    }
  }
}
```

### 9.3 获取AI调用日志

**接口说明**: 获取AI模型调用日志

```http
GET /api/v1/logs/ai-calls
```

### 9.4 综合搜索日志

**接口说明**: 多维度搜索日志

```http
GET /api/v1/logs/search
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| type | string | 否 | 日志类型：system/operation/ai-call |
| page | int | 否 | 页码 |
| size | int | 否 | 每页条数 |

### 9.5 导出日志

**接口说明**: 导出日志文件

```http
GET /api/v1/logs/export
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 日志类型：system/operation/ai-call |
| start_time | string | 是 | 开始时间 |
| end_time | string | 是 | 结束时间 |
| format | string | 否 | 导出格式：csv/json，默认csv |

**响应**: 文件下载

### 9.6 获取日志统计

**接口说明**: 获取日志统计数据

```http
GET /api/v1/logs/stats
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 否 | 日志类型 |
| start_date | string | 否 | 开始日期 |
| end_date | string | 否 | 结束日期 |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "level_distribution": {
      "debug": 100,
      "info": 500,
      "warn": 20,
      "error": 5
    },
    "daily_trend": [
      { "date": "2026-01-13", "count": 85 },
      { "date": "2026-01-14", "count": 92 }
    ],
    "top_errors": [
      { "message": "连接AI服务失败", "count": 3 },
      { "message": "推送钉钉失败", "count": 2 }
    ]
  }
}
```

---

## 10. 消息模板管理模块

### 10.1 获取模板列表

**接口说明**: 获取所有消息模板

```http
GET /api/v1/templates
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页条数 |
| keyword | string | 否 | 搜索关键词 |
| type | string | 否 | 模板类型：dingtalk/email |
| scene | string | 否 | 模板场景 |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "代码提交通知",
        "type": "dingtalk",
        "scene": "commit_notify",
        "title": "代码提交通知",
        "is_default": true,
        "status": "active",
        "version": 3,
        "created_at": "2026-01-01T00:00:00Z",
        "updated_at": "2026-01-19T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 10,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

### 10.2 获取模板详情

**接口说明**: 获取模板详细信息

```http
GET /api/v1/templates/:id
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "代码提交通知",
    "type": "dingtalk",
    "scene": "commit_notify",
    "title": "代码提交通知 - {{.RepoName}}",
    "content": "## {{.RepoName}} 提交了代码\n\n**提交信息**: {{.CommitMsg}}\n\n**提交者**: {{.Author}}\n\n**分支**: {{.Branch}}\n\n...",
    "is_default": true,
    "status": "active",
    "version": 3,
    "variables": ["RepoName", "CommitMsg", "Author", "Branch"],
    "push_count": 200,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-19T10:00:00Z"
  }
}
```

### 10.3 创建模板

**接口说明**: 创建新的消息模板

```http
POST /api/v1/templates
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 模板名称 |
| type | string | 是 | 模板类型：dingtalk/email |
| scene | string | 是 | 模板场景 |
| title | string | 是 | 模板标题 |
| content | string | 是 | 模板内容（支持变量替换） |
| remark | string | 否 | 备注说明 |

**请求示例**

```json
{
  "name": "代码提交通知",
  "type": "dingtalk",
  "scene": "commit_notify",
  "title": "代码提交通知 - {{.RepoName}}",
  "content": "## {{.RepoName}} 提交了代码\n\n**提交信息**: {{.CommitMsg}}\n**提交者**: {{.Author}}\n**分支**: {{.Branch}}\n\n### 变更文件\n{{range .ChangedFiles}}- {{.}}\n{{end}}",
  "remark": "用于代码提交时的通知"
}
```

### 10.4 更新模板

**接口说明**: 更新模板配置

```http
PUT /api/v1/templates/:id
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 模板名称 |
| title | string | 否 | 模板标题 |
| content | string | 否 | 模板内容 |
| remark | string | 否 | 备注说明 |
| status | string | 否 | 状态 |

### 10.5 删除模板

**接口说明**: 删除模板

```http
DELETE /api/v1/templates/:id
```

### 10.6 获取版本历史

**接口说明**: 获取模板的版本历史

```http
GET /api/v1/templates/:id/versions
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "version": 3,
        "content": "模板内容v3...",
        "created_by": "admin",
        "created_at": "2026-01-19T10:00:00Z"
      },
      {
        "version": 2,
        "content": "模板内容v2...",
        "created_by": "admin",
        "created_at": "2026-01-15T08:00:00Z"
      }
    ]
  }
}
```

### 10.7 回滚版本

**接口说明**: 回滚到指定版本

```http
POST /api/v1/templates/:id/rollback
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| version | int | 是 | 要回滚的目标版本号 |

### 10.8 模板预览

**接口说明**: 预览模板渲染效果

```http
POST /api/v1/templates/preview
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| template_id | int | 是 | 模板ID |
| data | object | 是 | 模拟数据 |

**请求示例**

```json
{
  "template_id": 1,
  "data": {
    "RepoName": "后端服务",
    "RepoUrl": "https://github.com/company/backend",
    "CommitId": "abc123",
    "CommitMsg": "feat: 新增登录功能",
    "Author": "zhangsan",
    "Branch": "main",
    "ChangedFiles": ["login.go", "auth.go"],
    "FileCount": 2
  }
}
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "title": "代码提交通知 - 后端服务",
    "body": "## 后端服务 提交了代码\n\n**提交信息**: feat: 新增登录功能\n..."
  }
}
```

### 10.9 模板测试

**接口说明**: 向指定目标发送测试消息

```http
POST /api/v1/templates/:id/test
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_id | int | 是 | 推送目标ID |
| test_data | object | 否 | 测试数据 |

### 10.10 启用/禁用模板

**接口说明**: 更新模板启用状态

```http
PUT /api/v1/templates/:id/status
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 是 | 状态：active/inactive |

---

## 11. 提示词管理模块

### 11.1 获取提示词列表

**接口说明**: 获取所有提示词

```http
GET /api/v1/prompts
```

**Query参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页条数 |
| keyword | string | 否 | 搜索关键词 |
| type | string | 否 | 类型：codeview/message |
| scene | string | 否 | 场景 |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "Go代码规范检查",
        "type": "codeview",
        "scene": "code_style",
        "language": "go",
        "model_id": 1,
        "model_name": "GPT-4",
        "version": 2,
        "created_at": "2026-01-01T00:00:00Z",
        "updated_at": "2026-01-19T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 10,
      "total": 8,
      "total_pages": 1
    }
  }
}
```

### 11.2 获取提示词详情

**接口说明**: 获取提示词详细信息

```http
GET /api/v1/prompts/:id
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Go代码规范检查",
    "type": "codeview",
    "scene": "code_style",
    "language": "go",
    "content": "你是一位资深Go语言代码审查专家。请审查以下代码，重点关注：\n1. 代码规范和最佳实践\n2. 潜在的bug\n3. 性能问题\n4. 安全性问题\n\n文件名：{{.FileName}}\n代码内容：\n{{.FileContent}}",
    "model_id": 1,
    "model_name": "GPT-4",
    "variables": ["FileName", "FileContent"],
    "version": 2,
    "usage_count": 150,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-19T10:00:00Z"
  }
}
```

### 11.3 创建CODEVIEW提示词

**接口说明**: 创建用于代码审查的提示词

```http
POST /api/v1/prompts
```

**请求参数（CODEVIEW类型）**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 提示词名称 |
| type | string | 是 | 固定值：codeview |
| scene | string | 否 | 使用场景 |
| language | string | 否 | 适用编程语言 |
| content | string | 是 | 提示词内容（需包含变量） |
| model_id | int | 否 | 关联的AI模型ID |

**请求示例**

```json
{
  "name": "Go代码规范检查",
  "type": "codeview",
  "scene": "code_style",
  "language": "go",
  "content": "你是一位资深Go语言代码审查专家。请审查以下代码，重点关注：\n1. 代码规范和最佳实践\n2. 潜在的bug\n3. 性能问题\n4. 安全性问题\n\n文件名：{{.FileName}}\n代码内容：\n{{.FileContent}}",
  "model_id": 1
}
```

### 11.4 创建推送消息提示词

**接口说明**: 创建用于生成推送消息的提示词

```http
POST /api/v1/prompts
```

**请求参数（推送消息类型）**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 提示词名称 |
| type | string | 是 | 固定值：message |
| target_type | string | 是 | 目标类型：dingtalk/email |
| content | string | 是 | 提示词内容 |

### 11.5 更新提示词

**接口说明**: 更新提示词配置

```http
PUT /api/v1/prompts/:id
```

### 11.6 删除提示词

**接口说明**: 删除提示词

```http
DELETE /api/v1/prompts/:id
```

### 11.7 获取版本历史

**接口说明**: 获取提示词的版本历史

```http
GET /api/v1/prompts/:id/versions
```

### 11.8 回滚版本

**接口说明**: 回滚到指定版本

```http
POST /api/v1/prompts/:id/rollback
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| version | int | 是 | 目标版本号 |

### 11.9 测试提示词

**接口说明**: 测试提示词效果

```http
POST /api/v1/prompts/:id/test
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| test_data | object | 是 | 测试数据 |

**请求示例**

```json
{
  "test_data": {
    "FileName": "main.go",
    "FileContent": "package main\n\nfunc main() {\n    fmt.Println(\"Hello\")\n}",
    "Language": "go",
    "RepoName": "backend-service",
    "Branch": "main",
    "CommitMsg": "initial commit"
  }
}
```

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "result": "代码审查结果：\n1. 缺少错误处理\n2. 建议使用 fmt.Println() 替代打印",
    "model_name": "GPT-4",
    "duration_ms": 2500,
    "tokens_used": 500
  }
}
```

### 11.10 导出提示词

**接口说明**: 导出提示词配置

```http
GET /api/v1/prompts/:id/export
```

**响应**: JSON文件下载

### 11.11 导入提示词

**接口说明**: 从文件导入提示词配置

```http
POST /api/v1/prompts/import
```

**Content-Type**: multipart/form-data

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 提示词配置文件（JSON格式） |

---

## 12. Webhook接口

### 12.1 Webhook回调

**接口说明**: 接收代码仓库的Webhook回调

```http
POST /webhook/:webhookId
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| webhookId | string | 是 | Webhook唯一标识符 |

**请求头（部分仓库平台）**

```http
X-GitHub-Event: push
X-Hub-Signature-256: sha256=...
X-GitHub-Delivery: uuid-string
```

**GitHub Push事件请求体示例**

```json
{
  "ref": "refs/heads/main",
  "before": "abc123",
  "after": "def456",
  "repository": {
    "id": 12345,
    "name": "backend-service",
    "full_name": "company/backend-service",
    "html_url": "https://github.com/company/backend-service",
    "clone_url": "https://github.com/company/backend-service.git"
  },
  "pusher": {
    "name": "zhangsan",
    "email": "zhangsan@company.com"
  },
  "sender": {
    "login": "zhangsan"
  },
  "commits": [
    {
      "id": "def456",
      "message": "feat: 新增用户登录功能",
      "timestamp": "2026-01-19T14:00:00Z",
      "author": {
        "name": "zhangsan",
        "email": "zhangsan@company.com"
      },
      "added": ["login.go"],
      "modified": [],
      "removed": []
    }
  ],
  "head_commit": {
    "id": "def456",
    "message": "feat: 新增用户登录功能"
  }
}
```

**GitLab Push事件请求体示例**

```json
{
  "object_kind": "push",
  "event_name": "push",
  "before": "abc123",
  "after": "def456",
  "ref": "refs/heads/main",
  "checkout_sha": "def456",
  "user_id": 100,
  "user_name": "zhangsan",
  "user_email": "zhangsan@company.com",
  "project": {
    "id": 123,
    "name": "backend-service",
    "web_url": "https://gitlab.com/company/backend-service",
    "git_ssh_url": "git@gitlab.com:company/backend-service.git",
    "git_http_url": "https://gitlab.com/company/backend-service.git"
  },
  "commits": [
    {
      "id": "def456",
      "message": "feat: 新增用户登录功能",
      "timestamp": "2026-01-19T14:00:00Z",
      "author": {
        "name": "zhangsan",
        "email": "zhangsan@company.com"
      }
    }
  ]
}
```

**Gitee Push事件请求体示例**

```json
{
  "type": "push",
  "before": "abc123",
  "after": "def456",
  "ref": "refs/heads/main",
  "repository": {
    "id": 12345,
    "name": "backend-service",
    "full_name": "company/backend-service",
    "html_url": "https://gitee.com/company/backend-service"
  },
  "sender": {
    "login": "zhangsan"
  },
  "commits": [
    {
      "id": "def456",
      "message": "feat: 新增用户登录功能",
      "timestamp": "2026-01-19T14:00:00Z",
      "author": {
        "name": "zhangsan",
        "email": "zhangsan@company.com"
      }
    }
  ]
}
```

**响应示例**

```json
{
  "code": 200,
  "message": "Webhook处理成功",
  "data": {
    "request_id": "webhook-uuid-123",
    "repo_id": 1,
    "commits_count": 1,
    "push_id": 1001,
    "status": "processing"
  }
}
```

### 12.2 Webhook安全验证

#### 12.2.1 签名验证（钉钉）

钉钉机器人使用签名验证机制。服务端需要验证请求签名：

```go
func verifyDingTalkSignature(timestamp, secret, signature string) bool {
    stringToSign := timestamp + "\n" + secret
    hmac256 := hmac.New(sha256.New, []byte(secret))
    hmac256.Write([]byte(stringToSign))
    sign := base64.StdEncoding.EncodeToString(hmac256.Sum(nil))
    return sign == signature
}
```

#### 12.2.2 Secret验证

Webhook Secret 用于验证请求来源：

```http
X-Webhook-Secret: <配置的secret值>
```

服务端会比较请求头中的 Secret 与配置是否一致。

### 12.3 Webhook触发的事件类型

| 事件类型 | 说明 | GitHub | GitLab | Gitee |
|---------|------|--------|--------|-------|
| push | 代码推送 | push | push | push |
| merge_request | 合并请求 | pull_request | merge_request | merge_request |
| tag_push | 标签推送 | create | tag_push | push |

---

## 13. 错误码说明

### 13.1 错误码列表

| 错误码 | 说明 | 处理建议 |
|--------|------|----------|
| 200 | 成功 | - |
| 400 | 参数错误 | 检查请求参数是否正确 |
| 401 | 未认证或Token过期 | 重新登录获取Token |
| 403 | 无权限访问 | 联系管理员开通权限 |
| 404 | 资源不存在 | 检查资源ID是否正确 |
| 409 | 资源冲突 | 检查唯一性约束字段 |
| 422 | 数据验证失败 | 检查数据格式和约束 |
| 429 | 请求过于频繁 | 稍后重试 |
| 500 | 服务器内部错误 | 联系技术支持 |
| 501 | 功能未实现 | 等待版本更新 |
| 503 | 服务不可用 | 检查服务状态 |

### 13.2 业务错误码

| 错误码 | 说明 |
|--------|------|
| 10001 | 用户名已存在 |
| 10002 | 邮箱已被注册 |
| 10003 | 密码强度不足 |
| 10004 | 账户已被锁定 |
| 10005 | 登录密码错误 |
| 20001 | 仓库名称已存在 |
| 20002 | 仓库URL格式无效 |
| 20003 | Webhook配置无效 |
| 20004 | 仓库不存在 |
| 30001 | 推送目标名称已存在 |
| 30002 | 钉钉Token无效 |
| 30003 | 邮箱SMTP配置无效 |
| 30004 | 推送目标不存在 |
| 40001 | AI模型配置无效 |
| 40002 | AI服务调用失败 |
| 40003 | 模型不存在 |
| 50001 | 模板名称已存在 |
| 50002 | 模板内容不能为空 |
| 50003 | 默认模板不能删除 |
| 60001 | 提示词名称已存在 |
| 60002 | 提示词必须包含必要变量 |
| 70001 | 推送记录不存在 |
| 70002 | 推送重试次数已用尽 |

### 13.3 错误响应格式

```json
{
  "code": 400,
  "message": "参数错误",
  "details": [
    "字段name不能为空",
    "字段url格式不正确"
  ],
  "request_id": "uuid-string"
}
```

---

## 附录

### 附录A：变量列表

#### A.1 消息模板变量

| 变量名 | 说明 | 示例 |
|-------|------|------|
| {{.RepoName}} | 仓库名称 | backend-service |
| {{.RepoUrl}} | 仓库地址 | https://github.com/company/backend |
| {{.CommitId}} | 提交ID | abc123def |
| {{.CommitMsg}} | 提交信息 | feat: 新增登录功能 |
| {{.Author}} | 提交者 | zhangsan |
| {{.Branch}} | 分支名称 | main |
| {{.ChangedFiles}} | 变更文件列表 | login.go, auth.go |
| {{.FileCount}} | 变更文件数量 | 2 |
| {{.CodeViewResult}} | 代码审查结果 | 通过/有建议 |
| {{.CodeViewIssues}} | 审查问题列表 | 问题描述列表 |
| {{.ReviewTime}} | 审查时间 | 2026-01-19 14:30 |

#### A.2 CODEVIEW提示词变量

| 变量名 | 说明 | 示例 |
|-------|------|------|
| {{.FileName}} | 文件名 | main.go |
| {{.FileContent}} | 完整文件内容 | package main... |
| {{.DiffContent}} | 代码差异内容 | diff内容 |
| {{.Language}} | 编程语言 | Go |
| {{.RepoName}} | 仓库名称 | backend-service |
| {{.Branch}} | 分支名称 | main |
| {{.CommitMsg}} | 提交信息 | feat: 新增登录功能 |

### 附录B：支持的仓库类型

| 类型值 | 仓库平台 | Webhook文档 |
|--------|---------|------------|
| github | GitHub | https://docs.github.com/en/webhooks |
| gitlab | GitLab | https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html |
| gitee | Gitee | https://gitee.com/help/articles/4330 |

### 附录C：支持的推送类型

| 类型值 | 说明 | 配置项 |
|--------|------|--------|
| dingtalk | 钉钉群机器人 | access_token, secret |
| email | 邮箱 | smtp_host, smtp_port, from, password, to |

### 附录D：支持的模板场景

| 场景值 | 说明 |
|--------|------|
| commit_notify | 代码提交通知 |
| review_notify | 审查结果通知 |

---

**文档结束**
