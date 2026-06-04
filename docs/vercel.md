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
VITE_API_BASE_URL=https://你的-api-域名
```

如果不设置 `VITE_API_BASE_URL`，前端会继续请求同域 `/api`，这只适合本地 Docker/Nginx 或你在同域下另行配置 API 代理的场景。

## 通过 Vercel CLI 部署

```bash
npm i -g vercel
vercel --prod
```

首次执行时按提示关联到当前仓库对应的 Vercel Project。之后每次 push 到 GitHub 后，Vercel Git 集成也会自动重新部署。

## 完整后端能力

Vercel 静态前端部署不会运行本项目的 Go API，也不会启动 PostgreSQL/Redis。完整会员管理功能需要先部署后端，然后把 `VITE_API_BASE_URL` 指向后端域名。
