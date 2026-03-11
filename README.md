# Online Agent (Frontend + Golang Backend + DeepSeek)

一个前后端分离的在线 Agent 项目：
- 前端：React + Vite（部署到 Vercel）
- 后端：Golang API（部署到 Vercel）
- 模型接口：DeepSeek Chat Completions

## 项目结构

```txt
.
├── frontend/          # React 前端
├── backend/           # Golang 后端（/api/chat）
└── .github/workflows/ # CI 自动化自测
```

## 1) 本地运行

### 启动后端

```bash
cd backend
export DEEPSEEK_API_KEY="你的 DeepSeek API Key"
go run .
```

后端默认运行在 `http://localhost:8080`。

### 启动前端

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

默认访问 `http://localhost:5173`。

## 2) 部署到 Vercel（前后端分离）

建议在同一个 GitHub 仓库里创建两个 Vercel Project：

### Project A: `online-agent-backend`
- Root Directory: `backend`
- Framework Preset: Other
- Environment Variables:
  - `DEEPSEEK_API_KEY` = 你的真实 DeepSeek API Key
- 部署后得到后端地址，例如：
  - `https://online-agent-backend.vercel.app`

### Project B: `online-agent-frontend`
- Root Directory: `frontend`
- Framework Preset: Vite
- Build Command: `npm run build`
- Output Directory: `dist`
- Environment Variables:
  - `VITE_API_BASE_URL` = 后端地址，例如 `https://online-agent-backend.vercel.app`

## 3) API 说明

### POST `/api/chat`

Request:

```json
{
  "message": "你好，帮我做一份本周学习计划"
}
```

Response:

```json
{
  "reply": "..."
}
```

## 4) 自动化自测（CI）

仓库已配置 GitHub Actions，在 Push / Pull Request 时自动执行：
- 后端：`gofmt` 检查 + `go test ./...`
- 前端：`npm install` + `npm run build`

这样可以实现“自动自测通过后再合并/提 PR”的流程。

## 注意事项

- 不要把 API Key 直接写进前端代码或提交到仓库。
- 通过 Vercel / GitHub Secrets 注入 `DEEPSEEK_API_KEY` 更安全。
