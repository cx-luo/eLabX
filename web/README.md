<h1 align="center">eLabX web</h1>
<p align="center"><strong>An AI-driven Electronic Laboratory Notebook (ELN)</strong></p>

This is the frontend of **eLabX**, an AI-powered Electronic Lab Notebook (ELN) built with:

- [Vue 3](https://vuejs.org/)
- [Vben Admin](https://github.com/vbenjs/vue-vben-admin)
- [Element Plus](https://element-plus.org/)
- [Vite](https://vitejs.dev/)

The frontend provides a modern, responsive interface for managing experiments, integrating AI features, and interacting with the backend.

---

## 🚀 Features

- 🧪 Dynamic table components for lab data
- 🧠 AI summarization UI (connects to OpenAI or LLM)
- 🔐 JWT-based login, RBAC permission control
- 🌙 Dark mode and multi-theme support (via Vben)
- 📄 Integrated form/table layouts with rich components

---

## ⚙️ Project Setup

### 1. Install dependencies

```bash
pnpm install
```

*(Also supports `npm install` or `yarn`, but `pnpm` is recommended by Vben-Admin)*

### 2. Start development server

```bash
pnpm dev
```

Open your browser at: [http://localhost:3000](http://localhost:3000)

> ⚠️ Ensure your backend (`server/`) is running at the configured API base URL

---

## 🔧 Environment Configuration

Modify `.env` or `.env.development` as needed:

```env
VITE_GLOB_API_URL=http://localhost:8080
VITE_GLOB_APP_TITLE=eLabX
VITE_USE_MOCK=false
```

You can configure proxy rules in `vite.config.ts` to redirect API requests.

---

## 📦 Build for Production

```bash
pnpm build
```

Output will be in `dist/`, ready to be deployed behind Nginx or any static file server.

---

## 📁 Frontend Directory Structure

```bash
web/
├── src/
│   ├── api/             # API definitions (Axios)
│   ├── components/      # Global/shared UI components
│   ├── hooks/           # Vue composition utilities
│   ├── layouts/         # App layout structure
│   ├── pages/           # Route views (experiments, login, etc.)
│   ├── router/          # Vue Router config
│   ├── store/           # Pinia store
│   ├── utils/           # Common helpers
│   ├── locales/         # i18n support
│   └── main.ts
├── public/              # Static assets
├── vite.config.ts       # Vite config
└── README.md            # ← 当前这个文件
```

---

## 🧪 API Integration

All API calls are made to the backend at `VITE_GLOB_API_URL`.

You can modify `src/api/modules/` to adjust endpoints, or use Axios interceptors for auth headers (JWT token injection).

---

## 🧩 Useful Commands

```bash
# Lint + format
pnpm lint
pnpm format

# Preview build
pnpm preview
```

---

## 🌐 i18n Support

This project supports multiple languages (English, 中文, etc.) via [Vue I18n](https://vue-i18n.intlify.dev/).
You can add translations in `src/locales/`.

---

## 🧬 Related Docs

* Backend README: [../server/README.md](../server/README.md)
* Project overview: [../README.md](../README.md)

---

## 📄 License

Licensed under the [MIT License](../LICENSE).

---

## 👤 Author

Developed by [@cx-luo](https://github.com/cx-luo)

> For issues and feature requests, please open an [issue](https://github.com/cx-luo/eLabX/issues).
