# 梅花易数占卜系统

该项目实现了一个完整的梅花易数占卜系统，包括基于 Gin + Gorm 的 Go 后端以及基于 Vue 3 + Vite 的前端界面。

## 项目结构

```
backend/   Go 后端服务
frontend/  Vue 3 前端应用
```

## 后端（Go）

### 运行步骤

1. 进入 `backend` 目录：
   ```bash
   cd backend
   ```
2. 安装依赖（首次需要可以访问公共 Go 模块代理）：
   ```bash
   go mod tidy
   ```
3. 启动服务：
   ```bash
   go run .
   ```

服务默认监听 `http://localhost:8080`，并提供以下接口：

- `POST /api/divination`：提交占卜请求，返回卦象信息并写入历史记录。
- `GET /api/divination`：获取历史占卜记录，按时间倒序返回。

后端默认使用 SQLite 数据库存储结果，文件名为 `divinations.db`。

## 前端（Vue 3）

### 安装与运行

1. 进入 `frontend` 目录并安装依赖：
   ```bash
   cd frontend
   npm install
   ```
2. 启动开发服务器：
   ```bash
   npm run dev
   ```

Vite 开发服务器默认运行在 `http://localhost:5173`，并通过代理将 `/api` 请求转发至本地的 Go 服务。

### 构建生产包

```bash
npm run build
```

构建完成后，输出位于 `frontend/dist` 目录。

## 功能简介

- 支持输入事件与三个起卦数字，自动计算对应的卦名、动爻和解读摘要。
- 历史记录自动保存，可随时回顾过往卦象。
- 前端界面采用 Vue 3 `<script setup>` 写法，样式使用 SCSS 编写，强调国风质感的视觉风格。

## 许可证

本项目仅用于演示用途，可按需修改与扩展。
