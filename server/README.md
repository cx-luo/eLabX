<h1 align="center">eLabX Server</h1>
This is the backend service of **eLabX** – an AI-driven Electronic Laboratory Notebook (ELN).  
Built with [Gin](https://github.com/gin-gonic/gin) and [GORM](https://gorm.io/), it provides a RESTful API for managing experiments, users, and AI-enhanced lab data.

---

## 🚀 Features

- 🧪 CRUD APIs for lab experiments, notes, and projects
- 🧠 Optional AI summarization via OpenAI API
- 🔐 JWT-based user authentication
- 🗃️ GORM-based ORM layer for MySQL/PostgreSQL
- 📊 Dynamic table querying with filtering, sorting, and pagination
- 📦 Modular routing and service structure

---

## 📁 Directory Overview

```bash
server/
├── api/             # HTTP handlers (Gin controllers)
├── dao/             # Data access layer (GORM logic)
├── middleware/      # JWT, CORS, logging, etc.
├── models/          # GORM models
├── router/          # Route definitions
├── service/         # Business logic
├── utils/           # Common helpers
├── config/          # App and DB config
├── main.go          # Application entrypoint
└── go.mod
```

---

## ⚙️ Configuration

App config is stored in `config/config.yaml` (or .env).

Example:

```yaml
server:
  port: 8080
  mode: release

database:
  type: mysql
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: elabx

auth:
  jwt_secret: your_jwt_secret_key
```

---

## 🛠️ Getting Started

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Run the Server

```bash
go run main.go
```

By default, it starts on: `http://localhost:8080`

### 3. Test API

Use tools like Postman or `curl`:

```bash
curl http://localhost:8080/api/v1/ping
```

---

## 🧪 API Docs

> You can enable Swagger in development for easier API browsing.
> (Optional integration: [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger))

---

## 🧹 Lint & Format

```bash
go fmt ./...
golangci-lint run
```

---

## 🧩 Environment Variables (Optional)

Alternatively, you can use `.env`:

```env
GIN_MODE=release
DB_HOST=localhost
DB_USER=root
DB_PASS=pass
JWT_SECRET=your_jwt
```

---

## 🐳 Docker Support

> *(Optional – if Dockerfile exists)*

```bash
docker build -t elabx-server .
docker run -p 8080:8080 elabx-server
```

---

## 🧭 Contribution Guide

Please refer to the root-level [CONTRIBUTING.md](../docs/CONTRIBUTING.md).

---

## 📄 License

See [LICENSE](../LICENSE) for licensing info (MIT).

---

## 📬 Contact

Maintainer: [@cx-luo](https://github.com/cx-luo)
