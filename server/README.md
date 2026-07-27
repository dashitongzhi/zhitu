# 职途 后端服务

基于 Go + Gin + GORM + 纯 Go SQLite 驱动的求职全流程工具后端。

## 技术栈

| 层 | 选型 | 说明 |
|---|---|---|
| HTTP 框架 | `github.com/gin-gonic/gin` | 轻量高性能 |
| ORM | `gorm.io/gorm` | 主流 Go ORM |
| SQLite 驱动 | `github.com/glebarez/sqlite` | **纯 Go 实现**，底层 `modernc.org/sqlite`，无需 CGO |
| JWT | `github.com/golang-jwt/jwt/v5` | 标准 JWT 签发与校验 |
| 配置 | `gopkg.in/yaml.v3` | YAML 配置文件 |
| 密码哈希 | `golang.org/x/crypto/bcrypt` | 业界标准 |

> 不引入缓存层、不接入第三方 OAuth，符合方案要求。

## 目录结构

```
server/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
├── configs/
│   └── config.yaml              # 配置文件（管理员凭据在此）
├── internal/
│   ├── config/                  # 配置加载
│   │   └── config.go
│   ├── database/                # 数据库初始化 + 自动迁移
│   │   └── database.go
│   ├── models/                  # GORM 数据模型
│   │   └── user.go
│   ├── services/                # 业务逻辑层（无 HTTP 上下文）
│   │   ├── auth_service.go      # 注册/登录/改密/管理员登录
│   │   └── jwt_service.go       # JWT 签发与解析
│   ├── handlers/                # HTTP 处理器层
│   │   └── auth_handler.go
│   ├── middleware/              # 中间件
│   │   ├── jwt.go               # JWT 认证 + RequireAdmin
│   │   ├── cors.go              # 跨域
│   │   └── recovery.go          # panic 恢复
│   ├── routers/                 # 路由注册
│   │   └── router.go
│   └── utils/                   # 工具函数
│       ├── response.go          # 统一响应封装
│       └── password.go          # 密码哈希与强度校验
├── .gitignore
├── go.mod
└── README.md
```

## 分层架构

```
HTTP Request
    │
    ▼
[middleware]  Recovery → Logger → CORS → (JWTAuth → RequireAdmin)
    │
    ▼
[handlers]    解析请求 / 校验入参 / 调用 service / 组装响应
    │
    ▼
[services]    业务逻辑（邮箱规范化、密码校验、JWT 签发等）
    │
    ▼
[models]      GORM 模型
    │
    ▼
[database]    SQLite (glebarez/sqlite, 纯 Go)
```

- `handlers` 只负责 HTTP 协议层：参数绑定、错误码映射、调用 service。
- `services` 不依赖 gin，方便后续被定时任务 / RPC 等复用。
- `models` 仅描述数据结构，不含业务方法。
- `middleware` 提供横切关注点（认证、跨域、恢复）。

## 用户体系说明

| 角色 | 凭据来源 | 登录入口 | 权限 |
|---|---|---|---|
| 管理员 | `configs/config.yaml` 的 `admin` 字段 | `POST /api/auth/admin/login` | 可访问 `/api/admin/*` |
| 普通用户 | 数据库 `users` 表 | `POST /api/auth/login` | 仅能访问 `/api/v1/*` 业务接口 |

- 多用户，**无用户组**，**无角色细分**：除管理员外所有普通用户权限一致。
- 管理员**不入库**，凭据完全由配置文件维护；如需修改，编辑 `config.yaml` 后重启服务。
- `config.yaml` 中管理员密码支持**明文**或 **bcrypt 哈希**（以 `$2a$` / `$2b$` / `$2y$` 开头时自动按哈希校验）。

## 路由设计

### 公共路由（无需认证）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/health` | 健康检查 |
| POST | `/api/auth/register` | 普通用户注册（邮箱 + 密码） |
| POST | `/api/auth/login` | 普通用户登录 |
| POST | `/api/auth/admin/login` | 管理员登录 |

### 已登录路由（需 `Authorization: Bearer <token>`）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/auth/me` | 获取当前登录者信息 |
| POST | `/api/auth/change-password` | 修改当前用户密码 |
| GET | `/api/v1/_ping` | 业务模块占位（后续挂载简历/面试/投递） |

### 管理员路由（需管理员 Token）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/admin/ping` | 管理员身份探活（后续挂载管理功能） |

## 中间件设计

| 中间件 | 作用 |
|---|---|
| `Recovery` | 捕获 panic，输出 500 JSON，避免进程崩溃 |
| `gin.Logger` | 标准 access log |
| `CORS` | 按 `allow_origins` 配置放行跨域，OPTIONS 预检直接 204 |
| `JWTAuth` | 解析 `Authorization: Bearer`，校验签名/过期，注入 `user_id` / `email` / `is_admin` 到 gin.Context |
| `RequireAdmin` | 在 `JWTAuth` 之后使用，校验 `is_admin=true`，否则 403 |

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

业务码：`0` 成功；`4xxxx` 客户端错误；`5xxxx` 服务端错误（详见 `utils/response.go`）。

## 快速开始

### 1. 安装 Go

需要 Go 1.22+。

### 2. 安装依赖

```bash
cd server
go mod tidy
```

### 3. 修改配置

编辑 `configs/config.yaml`：

- `jwt.secret` 改为强随机字符串
- `admin.email` / `admin.password` 改为自定义管理员凭据
- 如需对接前端，把前端地址加入 `server.allow_origins`

### 4. 运行

```bash
go run ./cmd/server
# 或编译后运行
go build -o bin/server ./cmd/server && ./bin/server
```

默认监听 `:8080`，数据库文件 `./data/zhitu.db` 自动创建。

### 5. 接口验证

```bash
# 注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"abc12345"}'

# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"abc12345"}'

# 用返回的 token 调用 /me
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer <token>"

# 修改密码
curl -X POST http://localhost:8080/api/auth/change-password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"abc12345","new_password":"xyz67890"}'

# 管理员登录
curl -X POST http://localhost:8080/api/auth/admin/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@zhitu.com","password":"admin@zhitu.com"}'
```

## 业务模块扩展指引

后续接入简历实验室 / 面试训练场 / 投递看板时，按以下方式扩展：

1. **新增模型**：在 `internal/models/` 下新增 `resume.go` / `interview.go` / `application.go` 等。
2. **注册迁移**：在 `internal/database/database.go` 的 `autoMigrate` 中追加新模型。
3. **新增 service**：在 `internal/services/` 下新增对应业务 service。
4. **新增 handler**：在 `internal/handlers/` 下新增对应业务 handler。
5. **挂载路由**：在 `internal/routers/router.go` 的 `/api/v1` 组下新增业务路由（已自动应用 JWT 中间件）；管理类接口挂在 `/api/admin` 组下（自动应用 JWT + RequireAdmin）。
