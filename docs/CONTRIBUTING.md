# 🤝 Contributing to eLabX

First off, thanks for taking the time to contribute! 🎉  
Your help makes **eLabX** better for everyone — whether it's fixing bugs, improving documentation, or proposing new features.

---

## 📦 Project Structure Overview

```bash
eLabX/
├── web/             # Frontend: Vue 3 + Vben Admin + Element Plus
├── server/          # Backend: Gin + GORM + MySQL/PostgreSQL
├── docs/            # Documentation
├── LICENSE
└── README.md
```

---

## 🧩 How to Contribute

### 🐛 1. Report Bugs

If you find a bug, please [create an issue](https://github.com/cx-luo/eLabX/issues) with:

* A clear title
* Steps to reproduce
* Expected behavior
* Actual behavior
* (Optional) Screenshots or error logs

### ✨ 2. Request Features

We welcome suggestions! If your idea benefits the community, open a feature request issue and describe:

* The use case
* Why it's valuable
* Alternatives considered

### 💻 3. Submit Pull Requests (PR)

Before submitting a PR:

1. Fork the repository
2. Create a feature branch:

   ```bash
   git checkout -b feature/my-new-feature
   ```
3. Commit your changes with clear messages
4. Push and submit a Pull Request (PR)

> ✅ Ensure your code passes linting and builds successfully

---

## 🧹 Coding Standards

### Frontend (Vue + TypeScript)

* Follow Vue 3 composition API practices
* Use ESLint/Prettier (configured in `web/.eslintrc`)
* Use i18n for UI text (avoid hard-coded strings)

### Backend (Golang)

* Follow idiomatic Go style (`go fmt`)
* Group routes logically
* Use consistent error handling

---

## 📘 Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard:

Examples:

```bash
feat: add AI summarization support
fix: correct sorting issue in experiment table
docs: update API usage guide
refactor: restructure database schema
```

---

## ✅ Pull Request Checklist

Before you submit:

* [ ] Code compiles and runs
* [ ] Tests pass (if applicable)
* [ ] Documentation is updated (README/API/docs)
* [ ] No breaking changes without discussion
* [ ] PR description explains what and why

---

## 👥 Code of Conduct

All contributors are expected to follow our [Code of Conduct](./CODE_OF_CONDUCT.md) (optional). We are committed to providing a welcoming and respectful environment.

---

## 🧪 Test Locally

You can run the full stack locally with:

```bash
# Backend
cd server
go run main.go

# Frontend
cd web
pnpm install
pnpm dev
```

---

## 📮 Questions?

If you have questions before contributing, feel free to:

* Open an issue
* Start a discussion
* Email the maintainer: [chengxiang.luo@foxmail.com](chengxiang.luo@foxmail.com)

Thanks again! ❤️
