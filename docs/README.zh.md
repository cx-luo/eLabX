# eLabX

**AI 驱动的电子实验室记录本（ELN）**

---

## 🌟 项目简介

eLabX 是一个使用 [Vben-Admin](https://github.com/vbenjs/vue-vben-admin)、[Element Plus](https://element-plus.org/)、[Gin](https://github.com/gin-gonic/gin) 和 [GORM](https://gorm.io/) 构建的 AI 驱动 ELN 系统，支持：

- ✅ 实验记录管理（实验、协议、样本等）
- 🧠 AI 自动摘要与建议（可接入 OpenAI 或本地 LLM）
- 🔐 用户认证、角色权限控制
- 🧩 表格排序 & 过滤（前端交互）
- 🔁 导出/导入记录（CSV/JSON）
- 🚀 REST API 后端，模块化设计

---

## 📸 截图预览

---

## 🛠 技术栈

| 层级       | 技术               |
|----------|------------------|
| 前端       | Vue 3 + Vben‑Admin + Element Plus |
| 后端       | Golang + Gin + GORM + MySQL/Postgres |
| AI（可选） | OpenAI API / 自己部署 LLM |
| 认证       | JWT (JSON Web Token) |

---

## 📦 快速开始

### 后端

```bash
cd server
go mod tidy
go run main.go
```

### 前端

```bash
cd web
pnpm install
pnpm dev
```

> ⚙️ 可通过 `.env` 配置数据库连接、JWT 密钥、AI 接口 KEY 等

---

## 📁 项目结构

```
eLabX/
├── server/        # 后端服务（Gin + GORM + REST API）
├── web/           # 前端代码（Vue + Vben‑Admin）
├── docs/          # 文档（架构、接口、部署等）
├── LICENSE        # 开源协议
└── README.md      # 本说明文档
```

---

## 🔭 功能与 Roadmap

* ✅ CRUD 实验记录
* ✅ AI 自动摘要功能
* 🔲 多用户协作 & 权限控制
* 🔲 审计日志（修改记录跟踪）
* 🔲 移动端响应式支持

---

## 🤝 贡献指南

我们真诚欢迎你的贡献！

1. 🍴 Fork 本仓库
2. 🧩 创建分支 (`git checkout -b feature/xxx`)
3. 🧪 完成功能或修复后提交 PR
4. 📬 我们会在 PR 中讨论，并进行合并

请阅读 `docs/CONTRIBUTING.md` 了解详细流程。

---

## 📄 开源协议

本项目采用 **MIT 许可证**。详见 [`LICENSE`](../LICENSE) 文件：

---

## 📬 联系方式

创建者：cx-Luo
Email: `chengxiang.luo@foxmail.com`
GitHub：[@cx-luo](https://github.com/cx-luo)
---

## 🌐 多语言支持

* 🇨🇳 中文文档：(this README)
* 🇺🇸 English [README.md](../README.md)（可选）