# Vercel 部署说明

本项目当前的 Vercel 配置用于部署 Vue 前端。Go API、PostgreSQL 和 Redis 仍建议部署在支持 Docker Compose 的服务器或容器平台上。

## 通过 Vercel Dashboard 部署

1. 登录 Vercel，选择 **Add New Project**。
2. 导入 GitHub 仓库 `xiaomt103/jiyixia_web`。
3. 保持仓库根目录作为 Root Directory。
4. Vercel 会读取根目录的 `vercel.json`：
   - Install Command: `cd frontend && npm install`
   - Build Command: `cd frontend && npm run build`
   - Output Directory: `frontend/dist`
5. 如果已经有线上 API，请在 Vercel 项目的 Environment Variables 中设置：

```text
API_BASE_URL=https://你的-api-域名
```

Vercel 会通过根目录 `api/` 下的 Serverless Functions 处理 `/api/*` 请求，并把请求代理到 `API_BASE_URL`，例如 `/api/admin/login` 会转发到 `https://你的-api-域名/api/admin/login`。这样浏览器仍访问同域 `/api`，不需要额外处理 CORS。

请不要把 `API_BASE_URL` 设置成当前 Vercel 站点域名（例如 `https://jiyixia-web.vercel.app`），否则代理会请求自己并形成循环。本项目已内置保护：当 `API_BASE_URL` 未配置或指向当前 Vercel 域名时，`/api/*` 会进入临时演示模式，返回内置套餐并支持默认后台登录，便于先查看页面效果。项目同时提供显式函数文件（例如 `api/admin/login.js`、`api/plans.js`），避免 catch-all API 在部分 Vercel 配置下没有被部署而导致 404。

如果你希望浏览器直接请求后端域名，也可以设置 `VITE_API_BASE_URL=https://你的-api-域名`，但这要求后端正确允许 Vercel 域名的 CORS。

## 通过 Vercel CLI 部署

```bash
npm i -g vercel
vercel --prod
```

首次执行时按提示关联到当前仓库对应的 Vercel Project。之后每次 push 到 GitHub 后，Vercel Git 集成也会自动重新部署。

## 完整后端能力

Vercel 前端部署不会运行本项目的 Go API，也不会启动 PostgreSQL/Redis，因此默认没有可持久化的 psql 数据库。内置演示模式仅适合验证页面和登录流程：Serverless Function 内存可能随实例变化而丢失，前端会把当前浏览器提交的 demo 会员/套餐同步到 `localStorage`，让同一浏览器里的后台管理端可以看到刚提交的用户数据。不同浏览器、不同设备或重新部署后不能依赖这些 demo 数据。完整会员管理功能需要先部署后端和 PostgreSQL/Redis，然后把 Vercel 环境变量 `API_BASE_URL` 指向后端域名。
