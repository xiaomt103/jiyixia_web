# 记一下会员管理系统

一个基于 Go + Vue + PostgreSQL + Redis + Docker 的会员管理系统脚手架，包含用户端、后台管理端和可扩展的后台 API。

## 功能

- 用户端：查看会员套餐、提交会员注册、查询会员状态。
- 后台管理端：管理员登录、会员管理、状态控制、套餐管理。
- API：`/api/v1/admin/*` 暴露后台控制接口，使用 Bearer Token 认证。
- 基础设施：Docker Compose 一键启动 Go API、Vue 前端、PostgreSQL、Redis。

## 本地启动

```bash
docker compose up --build
```

服务地址：

- 前端：http://localhost:8080
- API：http://localhost:8080/api/health
- PostgreSQL：localhost:5432，数据库 `members`
- Redis：localhost:6379

默认后台账号由环境变量配置，可参考 `.env.example` 覆盖；Docker Compose 默认：

- 用户名：`admin`
- 密码：`admin123`（仅用于本地开发，请在真实环境覆盖）

## 常用 API

完整 OpenAPI 描述见 [`docs/openapi.yaml`](docs/openapi.yaml)。

### 用户端

- `GET /api/health` 健康检查
- `GET /api/plans` 套餐列表
- `POST /api/members` 创建会员申请
- `GET /api/members/{id}` 查询会员信息

### 后台管理端

- `POST /api/admin/login` 登录并返回 token
- `GET /api/v1/admin/members` 会员列表
- `PATCH /api/v1/admin/members/{id}/status` 更新会员状态
- `GET /api/v1/admin/plans` 套餐列表
- `POST /api/v1/admin/plans` 创建套餐

## 开发命令

```bash
# 后端测试
cd backend && go test ./...

# 前端构建
cd frontend && npm install && npm run build
```
